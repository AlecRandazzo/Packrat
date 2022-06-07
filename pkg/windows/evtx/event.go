// Copyright (c) 2022 Alec Randazzo

package evtx

import (
	"bytes"
	"fmt"
	"time"

	"github.com/AlecRandazzo/Packrat/pkg/general/byteshelper"
	"github.com/AlecRandazzo/Packrat/pkg/general/sanitycheck"
	"github.com/AlecRandazzo/Packrat/pkg/general/timestamp"
)

// These are effectively constants
var (
	eventMagicNumber          = []byte{0x2a, 0x2a, 0x00, 0x00}
	eventMagicNumberMeta      = byteshelper.NewDataLocation(0x00, 0x04)
	eventSizeMeta             = byteshelper.NewDataLocation(0x04, 0x04)
	eventRecordIdentifierMeta = byteshelper.NewDataLocation(0x08, 0x04)
	eventWrittenTimestampMeta = byteshelper.NewDataLocation(0x10, 0x08)
	eventBinXmlMeta           = byteshelper.DataLocation{
		Offset: 0x18, // TODO get length and change this to NewDataLocation()
	}
)

type event struct {
	eventSize                int32
	evenRecordIdentifier     int32
	writtenTimestamp         time.Time
	eventGeneratedTime       time.Time
	eventDetailKeyValuePairs map[int]map[string]string
}

func (event *event) parse(inBytes []byte) error {
	// Sanity checking
	err := sanitycheck.Bytes(inBytes, 0x18)
	if err != nil {
		return fmt.Errorf("event record size too small: %w", err)
	}

	buffer, _ := byteshelper.GetValue(inBytes, eventMagicNumberMeta)
	if !bytes.Equal(buffer, eventMagicNumber) {
		return fmt.Errorf("this is not an event since magic number is %x", buffer)
	}

	buffer, _ = byteshelper.GetValue(inBytes, eventSizeMeta)
	event.eventSize, _ = byteshelper.LittleEndianBinaryToInt32(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, eventRecordIdentifierMeta)
	event.evenRecordIdentifier, _ = byteshelper.LittleEndianBinaryToInt32(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, eventWrittenTimestampMeta)
	event.writtenTimestamp, _ = timestamp.Parse(buffer)

	buffer, _ = byteshelper.GetValue(inBytes, eventBinXmlMeta)
	_ = buffer // TODO parser the byteshelper XML

	return nil
}
