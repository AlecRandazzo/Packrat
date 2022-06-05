// Copyright (c) 2022 Alec Randazzo

package parse

import (
	"encoding/csv"
	"fmt"
	"github.com/AlecRandazzo/Packrat/pkg/parsers/windows/mft"
	"os"
	"strconv"
	"strings"
)

type CsvWriter struct {
	writer *csv.Writer
}

func NewCsvWriter(file string) (CsvWriter, error) {
	handle, err := os.Create(file)
	if err != nil {
		return CsvWriter{}, err
	}
	writer := csv.NewWriter(handle)
	writer.Comma = []rune("|")[0]
	err = writer.Write([]string{
		"RecordNumber",
		"MftOffset",
		"FileName",
		"FullPath",
		"LogicalFileSize",
		"PhysicalFileSize",
		"SiCreated",
		"SiModified",
		"SiAccessed",
		"SiChanged",
		"FnCreated",
		"FnModified",
		"FnAccessed",
		"FnChanged",
		"Directory",
		"Deleted",
		"Hidden",
		"SystemFile",
	})
	if err != nil {
		return CsvWriter{}, err
	}

	return CsvWriter{writer: writer}, nil
}

func boolToIntString(b bool) string {
	result := "0"
	if b {
		result = "1"
	}
	return result
}

func (csvWriter *CsvWriter) write(record mft.Record) error {
	row := []string{
		strconv.FormatUint(uint64(record.Header.RecordNumber), 10),
		strconv.FormatUint(record.Metadata.MftOffset, 10),
		record.Metadata.ChosenFileNameAttribute.FileName,
		record.Metadata.FullPath,
		strconv.FormatUint(record.Metadata.ChosenFileNameAttribute.LogicalFileSize, 10),
		strconv.FormatUint(record.Metadata.ChosenFileNameAttribute.PhysicalFileSize, 10),
		record.StandardInformationAttributes.Created.Format("01/02/2006 15:04:05"),
		record.StandardInformationAttributes.Modified.Format("01/02/2006 15:04:05"),
		record.StandardInformationAttributes.Accessed.Format("01/02/2006 15:04:05"),
		record.StandardInformationAttributes.Changed.Format("01/02/2006 15:04:05"),
		record.Metadata.ChosenFileNameAttribute.Created.Format("01/02/2006 15:04:05"),
		record.Metadata.ChosenFileNameAttribute.Modified.Format("01/02/2006 15:04:05"),
		record.Metadata.ChosenFileNameAttribute.Accessed.Format("01/02/2006 15:04:05"),
		record.Metadata.ChosenFileNameAttribute.Changed.Format("01/02/2006 15:04:05"),
		boolToIntString(record.Header.Flags.Directory),
		boolToIntString(record.Header.Flags.Deleted),
		boolToIntString(record.Metadata.ChosenFileNameAttribute.FileNameFlags.Hidden),
		boolToIntString(record.Metadata.ChosenFileNameAttribute.FileNameFlags.System),
	}
	return csvWriter.writer.Write(row)
}

func Parse(target string, writer CsvWriter, bytesPerSector, bytesPerCluster uint) error {
	handle, err := os.Open(target)
	if err != nil {
		return err
	}
	defer handle.Close()

	// Make directory tree
	dirTree := make(mft.DirectoryTree, 0)
	dirTree, err = mft.BuildDirectoryTree(handle, "", bytesPerSector)

	// Reset offset pointer back to the beginning
	_, err = handle.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek back to the beginning of the mft for second pass: %w", err)
	}

	var offset uint64
	for {
		buf := make([]byte, 1024)
		var n int
		n, err = handle.Read(buf)
		if err != nil {
			break
		}
		offset += uint64(n)
		record, _ := mft.ParseRecord(buf, bytesPerSector, bytesPerCluster)
		record.Metadata.MftOffset = offset

		// Pick an FN attribute
		for _, attribute := range record.FileNameAttributes {
			if strings.Contains(string(attribute.FileNamespace), "WIN32") ||
				strings.Contains(string(attribute.FileNamespace), "WIN32 & DOS") ||
				strings.Contains(string(attribute.FileNamespace), "POSIX") {
				record.Metadata.ChosenFileNameAttribute = attribute
				break
			}
		}
		record.Metadata.FullPath = fmt.Sprintf(`%s\%s`, dirTree[record.Header.RecordNumber], record.Metadata.ChosenFileNameAttribute.FileName)
		_ = writer.write(record)
	}
	return nil
}
