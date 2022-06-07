// Copyright (c) 2022 Alec Randazzo

package windows

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/windows/mft"
)

func parseMFTRecord0(handler handler) (mft.Record, error) {
	// Move handle pointer back to beginning of handler
	_, err := handler.Handle().Seek(0x00, 0)
	if err != nil {
		return mft.Record{}, fmt.Errorf("failed to seek back to handler offset 0x00: %w", err)
	}

	// Seek to the offset where the MFT starts. If it errors, bomb.
	_, err = handler.Handle().Seek(handler.Vbr().MftOffset, 0)
	if err != nil {
		return mft.Record{}, fmt.Errorf("failed to seek to mft: %w", err)
	}

	// Read the first entry in the MFT. The first record in the MFT always is for the MFT itself. If it errors, bomb.
	buffer := make([]byte, handler.Vbr().MftRecordSize)
	_, err = handler.Handle().Read(buffer)
	if err != nil {
		return mft.Record{}, fmt.Errorf("failed to read the mft: %w", err)
	}

	// Sanity check that this is indeed an mft record
	err = mft.ValidateMftRecordBytes(buffer)
	if err != nil {
		return mft.Record{}, fmt.Errorf("invalid mft record: %w", err)
	}

	// Parse the MFT record
	var mft0 mft.Record
	mft0, err = mft.ParseRecord(buffer, handler.Vbr().BytesPerSector, handler.Vbr().BytesPerCluster)
	if err != nil {
		return mft.Record{}, fmt.Errorf("Handler.parseMFTRecord0() failed to parser the mft's mft record: %w", err)
	}
	log.Debugf("Identified the following data runs for the MFT itself: %+v", mft0.DataAttribute.NonResidentDataAttribute.DataRuns)

	return mft0, nil
}
