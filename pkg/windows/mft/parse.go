// Copyright (c) 2022 Alec Randazzo

package mft

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/AlecRandazzo/Packrat/pkg/general/byteshelper"
	log "github.com/sirupsen/logrus"
)

// Record contains information on a parsed MFT record
type Record struct {
	Header                        RecordHeader
	StandardInformationAttributes StandardInformationAttribute
	FileNameAttributes            []FileNameAttribute
	DataAttribute                 DataAttribute
	AttributeList                 AttributeListAttributes
	Metadata                      struct {
		MftOffset               uint64
		FullPath                string
		ChosenFileNameAttribute FileNameAttribute
	}
}

const defaultRecordSize = 1024

func ParseFile(reader *os.File, writer Writer, bytesPerSector, bytesPerCluster uint) error {
	// Make directory tree
	dirTree, err := BuildDirectoryTree(reader, "", bytesPerSector, bytesPerCluster)
	if err != nil {
		return fmt.Errorf("failed to build directory tree: %w", err)
	}

	// Reset offset pointer back to the beginning
	_, err = reader.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek back to the beginning of the mft for second pass: %w", err)
	}

	var offset uint64
	for {
		buf := make([]byte, defaultRecordSize)
		var n int
		n, err = reader.Read(buf)
		if err != nil {
			break
		}
		offset += uint64(n)
		record, _ := ParseRecord(buf, bytesPerSector, bytesPerCluster)
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
		err = writer.Write(record)
		if err != nil && err != invalidRecord {
			log.Errorf("failed to parser record: %s", err.Error())
		}
	}

	return nil
}

// ParseRecord the raw MFT record and returns a parsed mft record.
func ParseRecord(input []byte, bytesPerSector, bytesPerCluster uint) (Record, error) {
	// Sanity checks
	size := len(input)
	if size == 0 {
		return Record{}, errors.New("received nil input")
	}
	if bytesPerSector == 0 {
		return Record{}, errors.New("input of 0 for bytesPerSector is not valid, typically this value is 512")
	}
	if bytesPerCluster == 0 {
		return Record{}, errors.New("input of 0 for bytesPerCluster is not valid, typically this value is 4096")
	}

	// init return variables
	var mft Record

	err := ValidateMftRecordBytes(input)
	if err != nil {
		return Record{}, fmt.Errorf("this is not an mft record: %w", err)
	}
	err = fixup(input, bytesPerSector)
	if err != nil {
		return Record{}, fmt.Errorf("fixup failed: %w", err)
	}

	input = trimSlackSpace(input)

	mft.Header, err = GetRecordHeaders(input)
	if err != nil {
		return Record{}, fmt.Errorf("failed to get record headers: %w", err)
	}

	var rawAttributes [][]byte
	rawAttributes, err = GetRawAttributes(input, mft.Header)
	if err != nil {
		return Record{}, fmt.Errorf("failed to get raw data attributes: %w", err)
	}

	mft.FileNameAttributes, mft.StandardInformationAttributes, mft.DataAttribute, mft.AttributeList, _ = GetAttributes(rawAttributes, bytesPerCluster)
	return mft, nil
}

// trimSlackSpace trims off slack space after end sequence 0xffffffff
func trimSlackSpace(input []byte) []byte {
	lenMftRecordBytes := len(input)
	mftRecordEndByteSequence := []byte{0xff, 0xff, 0xff, 0xff}
	for i := 0; i < (lenMftRecordBytes - 4); i++ {
		if bytes.Equal(input[i:i+0x04], mftRecordEndByteSequence) {
			input = input[:i]
			break
		}
	}

	return input
}

var (
	updateSequenceOffsetLocation     = byteshelper.NewDataLocation(0x04, 0x02)
	updateSequenceBufferSizeLocation = byteshelper.NewDataLocation(0x06, 0x02)
)

// fixup MFT record
func fixup(input []byte, bytesPerSector uint) error {
	// Sanity checks
	inputSize := uint(len(input))
	if inputSize == 0 {
		return errors.New("nil input bytes received by fixup()")
	}
	if inputSize < bytesPerSector {
		return errors.New("input is smaller than sector size")
	}
	if bytesPerSector == 0 {
		return errors.New("bytesPerSector is 0")
	}

	buffer, _ := byteshelper.GetValue(input, updateSequenceOffsetLocation)
	updateSequenceOffset := binary.LittleEndian.Uint16(buffer)
	updateSequence := input[updateSequenceOffset : updateSequenceOffset+2]

	buffer, _ = byteshelper.GetValue(input, updateSequenceBufferSizeLocation)

	updateSequenceBufferSize := binary.LittleEndian.Uint16(buffer)
	updateSequenceBuffer := input[updateSequenceOffset+2 : updateSequenceOffset+(updateSequenceBufferSize*2)]

	i := uint(512)
	bufferIndex := 0
	for i <= inputSize {
		if reflect.DeepEqual(input[i-2:i], updateSequence) {
			input[i-2] = updateSequenceBuffer[bufferIndex]
			input[i-1] = updateSequenceBuffer[bufferIndex+1]
		}
		i += bytesPerSector
		bufferIndex += 2
	}

	return nil
}