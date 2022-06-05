// Copyright (c) 2022 Alec Randazzo

package mft

import (
	"errors"
	"fmt"
	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/byteshelper"
	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/sanitycheck"
	"time"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/timestamp"
)

// StandardInformationAttribute contains information from a parsed standard information attribute.
type StandardInformationAttribute struct {
	Created      time.Time
	Modified     time.Time
	Accessed     time.Time
	Changed      time.Time
	FlagResident bool
}

var (
	siResidentFlagLocation = byteshelper.NewDataLocation(0x08, 0x01)
	siCreatedLocation      = byteshelper.NewDataLocation(0x18, 0x08)
	siModifiedLocation     = byteshelper.NewDataLocation(0x20, 0x08)
	siChangedLocation      = byteshelper.NewDataLocation(0x28, 0x08)
	siAccessedLocation     = byteshelper.NewDataLocation(0x30, 0x08)
)

// getStandardInformationAttribute parses a raw standard information attribute.
func getStandardInformationAttribute(input []byte) (StandardInformationAttribute, error) {
	err := sanitycheck.Bytes(input, 0x30)
	if err != nil {
		return StandardInformationAttribute{}, fmt.Errorf("invalid input: %w", err)
	}

	// init return values
	var si StandardInformationAttribute

	// Check to see if the standard information Attribute is resident to the MFT or not
	var buffer []byte
	buffer, err = byteshelper.GetValue(input, siResidentFlagLocation)
	if err != nil {
		return StandardInformationAttribute{}, fmt.Errorf("failed to get FlagResidency from standard information attribute, %w", err)
	}
	si.FlagResident = checkResidency(buffer[0])
	if !si.FlagResident {
		return StandardInformationAttribute{}, errors.New("non resident standard information attribute")
	}

	// parse timestamps
	buffer, _ = byteshelper.GetValue(input, siCreatedLocation)
	si.Created, _ = timestamp.Parse(buffer)

	buffer, _ = byteshelper.GetValue(input, siModifiedLocation)
	si.Modified, _ = timestamp.Parse(buffer)

	buffer, _ = byteshelper.GetValue(input, siChangedLocation)
	si.Changed, _ = timestamp.Parse(buffer)

	buffer, _ = byteshelper.GetValue(input, siAccessedLocation)
	si.Accessed, _ = timestamp.Parse(buffer)

	return si, nil
}
