// Copyright (c) 2022 Alec Randazzo

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
	chunkMagicNumber, _            = hex.DecodeString("456c6643686e6b") //0x456c6643686e6b
	chunkMagicNumberMeta           = byteshelper.NewDataLocation(0, 0x08)
	firstEventRecordNumberMeta     = byteshelper.NewDataLocation(0x08, 0x08)
	lastEventRecordNumberMeta      = byteshelper.NewDataLocation(0x10, 0x08)
	firstEventRecordIdentifierMeta = byteshelper.NewDataLocation(0x18, 0x08)
	lastEventRecordIdentifierMeta  = byteshelper.NewDataLocation(0x20, 0x08)
	chunkHeaderSizeMeta            = byteshelper.NewDataLocation(0x28, 0x04)
	lastEventRecordDataOffsetMeta  = byteshelper.NewDataLocation(0x2C, 0x04)
	freeSpaceOffsetMeta            = byteshelper.NewDataLocation(0x30, 0x04)
	eventRecordsChecksumMeta       = byteshelper.NewDataLocation(0x34, 0x04)
	chunkChecksumMeta              = byteshelper.NewDataLocation(0x7C, 0x04)
)

type chunk struct {
	firstEventRecordNumber     int64
	lastEventRecordNumber      int64
	firstEventRecordIdentifier int64
	lastEventRecordIdentifier  int64
	headerSize                 int32
	lastEventRecordDataOffset  int32
	freeSpaceOffset            int32
	eventRecordsChecksum       int32
	checksum                   int32
}

func (chunk chunk) parse(inBytes []byte) error {
	// Sanity checking
	err := sanitycheck.Bytes(inBytes, 512)
	if err != nil {
		return fmt.Errorf("chunk should be 512 inBytes in size, received %d: %w", len(inBytes), err)
	}

	buffer, _ := byteshelper.GetValue(inBytes, chunkMagicNumberMeta)
	if bytes.Compare(buffer, chunkMagicNumber) != 0 {
		return fmt.Errorf("this is not a chunk since magic number is %x", buffer)
	}

	buffer, _ = byteshelper.GetValue(inBytes, firstEventRecordNumberMeta)
	chunk.firstEventRecordNumber, _ = byteshelper.LittleEndianBinaryToInt64(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, lastEventRecordNumberMeta)
	chunk.lastEventRecordNumber, _ = byteshelper.LittleEndianBinaryToInt64(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, firstEventRecordIdentifierMeta)
	chunk.firstEventRecordIdentifier, _ = byteshelper.LittleEndianBinaryToInt64(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, lastEventRecordIdentifierMeta)
	chunk.lastEventRecordIdentifier, _ = byteshelper.LittleEndianBinaryToInt64(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, chunkHeaderSizeMeta)
	chunk.headerSize, _ = byteshelper.LittleEndianBinaryToInt32(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, lastEventRecordDataOffsetMeta)
	chunk.lastEventRecordDataOffset, _ = byteshelper.LittleEndianBinaryToInt32(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, freeSpaceOffsetMeta)
	chunk.freeSpaceOffset, _ = byteshelper.LittleEndianBinaryToInt32(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, eventRecordsChecksumMeta)
	chunk.eventRecordsChecksum, _ = byteshelper.LittleEndianBinaryToInt32(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, chunkChecksumMeta)
	chunk.checksum, _ = byteshelper.LittleEndianBinaryToInt32(buffer)

	return nil
}
