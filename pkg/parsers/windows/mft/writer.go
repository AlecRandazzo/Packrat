package mft

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
)

type Writer interface {
	Write(Record) error
}

type CsvWriter struct {
	writer *csv.Writer
}

var invalidRecord = errors.New("invalid record")

func (writer CsvWriter) Write(record Record) error {
	if record.Header.RecordNumber == 0 {
		return invalidRecord
	}
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
	return writer.writer.Write(row)
}

func NewCsvWriter(writer io.Writer) (CsvWriter, error) {
	csvWriter := csv.NewWriter(writer)
	csvWriter.Comma = []rune("|")[0]
	err := csvWriter.Write([]string{
		"RecordNumber",
		"MftOffset",
		"FileName",
		"fullPath",
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

	return CsvWriter{writer: csvWriter}, nil
}

func boolToIntString(b bool) string {
	result := "0"
	if b {
		result = "1"
	}
	return result
}
