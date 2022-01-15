// Copyright (c) 2020 Alec Randazzo

package evtx

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/byteshelper"
	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/sanitycheck"
)

// These are effectively constants
var (
	fileHeaderMagicNumber, _  = hex.DecodeString("456c6646696c6500") // ElfFile\x00
	fileHeaderMagicNumberMeta = byteshelper.NewDataLocation(0x00, 0x08)
	firstChunkNumberMeta      = byteshelper.NewDataLocation(0x08, 0x08)
	lastChunkNumberMeta       = byteshelper.NewDataLocation(0x10, 0x08)
	nextRecordIdentifierMeta  = byteshelper.NewDataLocation(0x18, 0x08)
	fileHeaderSizeMeta        = byteshelper.NewDataLocation(0x20, 0x04)
	minorVersionMeta          = byteshelper.NewDataLocation(0x24, 0x02)
	majorVersionMeta          = byteshelper.NewDataLocation(0x26, 0x02)
	headerBlockSizeMeta       = byteshelper.NewDataLocation(0x28, 0x02)
	numberOfChunksMeta        = byteshelper.NewDataLocation(0x2A, 0x02)
	fileFlagsMeta             = byteshelper.NewDataLocation(0x78, 0x04)
	fileFlagDirty             = 0x0001
	fileFlagFull              = 0x0002
	fileHeaderCheckSumMeta    = byteshelper.NewDataLocation(0x7C, 0x04)
)

type fileHeader struct {
	firstChunkNumber     int64
	lastChunkNumber      int64
	nextRecordIdentifier int64
	headerSize           int32
	minorVersion         int16
	majorVersion         int16
	headerBlockSize      int16
	numberOfChunks       int16
	fileFlags            fileFlags
	checksum             int32
}

type fileFlags struct{}

func (fileHeader fileHeader) parse(inBytes []byte) error {
	// Sanity checking
	err := sanitycheck.Bytes(inBytes, 4096)
	if err != nil {
		return fmt.Errorf("file header should be 4096 inBytes in size, received %d: %w", len(inBytes), err)
	}

	buffer, _ := byteshelper.GetValue(inBytes, fileHeaderMagicNumberMeta)
	if bytes.Compare(buffer, fileHeaderMagicNumber) != 0 {
		return fmt.Errorf("this is not an evtx file header since magic number is %x", buffer)
	}

	buffer, _ = byteshelper.GetValue(inBytes, firstChunkNumberMeta)
	fileHeader.firstChunkNumber, _ = byteshelper.LittleEndianBinaryToInt64(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, lastChunkNumberMeta)
	fileHeader.lastChunkNumber, _ = byteshelper.LittleEndianBinaryToInt64(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, nextRecordIdentifierMeta)
	fileHeader.nextRecordIdentifier, _ = byteshelper.LittleEndianBinaryToInt64(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, fileHeaderSizeMeta)
	fileHeader.headerSize, _ = byteshelper.LittleEndianBinaryToInt32(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, minorVersionMeta)
	fileHeader.minorVersion, _ = byteshelper.LittleEndianBinaryToInt16(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, majorVersionMeta)
	fileHeader.majorVersion, _ = byteshelper.LittleEndianBinaryToInt16(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, headerBlockSizeMeta)
	fileHeader.headerBlockSize, _ = byteshelper.LittleEndianBinaryToInt16(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, numberOfChunksMeta)
	fileHeader.numberOfChunks, _ = byteshelper.LittleEndianBinaryToInt16(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, fileFlagsMeta)
	_ = buffer // TODO

	buffer, _ = byteshelper.GetValue(inBytes, fileHeaderCheckSumMeta)
	fileHeader.checksum, _ = byteshelper.LittleEndianBinaryToInt32(buffer)

	return nil
}
