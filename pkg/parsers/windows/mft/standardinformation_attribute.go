// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"errors"
	"time"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/timestamp"
)

// RawStandardInformationAttribute is a []byte alias for raw standard information attribute. Used with the Parse() method.
type RawStandardInformationAttribute []byte

// StandardInformationAttribute contains information from a parsed standard information attribute.
type StandardInformationAttribute struct {
	SiCreated    time.Time
	SiModified   time.Time
	SiAccessed   time.Time
	SiChanged    time.Time
	FlagResident bool
}

// Parse parses the raw standard information attribute receiver and returns a parsed standard information attribute.
func (rawStandardInformationAttribute RawStandardInformationAttribute) Parse() (standardInformationAttribute StandardInformationAttribute, err error) {
	const offsetResidentFlag = 0x08

	const offsetSiCreated = 0x18
	const lengthSiCreated = 0x08

	const offsetSiModified = 0x20
	const lengthSiModified = 0x08

	const offsetSiChanged = 0x28
	const lengthSiChanged = 0x08

	const offsetSiAccessed = 0x30
	const lengthSiAccessed = 0x08

	// The standard information Attribute has a minimum length of 0x30
	if len(rawStandardInformationAttribute) < 0x30 {
		err = errors.New("StandardInformationAttributes.parse() received invalid bytes")
		return
	}
	// Check to see if the standard information Attribute is resident to the MFT or not
	rawResidencyFlag := RawResidencyFlag(rawStandardInformationAttribute[offsetResidentFlag])
	standardInformationAttribute.FlagResident = rawResidencyFlag.Parse()
	if standardInformationAttribute.FlagResident == false {
		err = errors.New("non resident standard information flag found")
		return
	}

	// parse timestamps
	standardInformationAttribute.SiCreated, _ = timestamp.Parse(rawStandardInformationAttribute[offsetSiCreated : offsetSiCreated+lengthSiCreated])
	standardInformationAttribute.SiModified, _ = timestamp.Parse(rawStandardInformationAttribute[offsetSiModified : offsetSiModified+lengthSiModified])
	standardInformationAttribute.SiChanged, _ = timestamp.Parse(rawStandardInformationAttribute[offsetSiChanged : offsetSiChanged+lengthSiChanged])
	standardInformationAttribute.SiAccessed, _ = timestamp.Parse(rawStandardInformationAttribute[offsetSiAccessed : offsetSiAccessed+lengthSiAccessed])
	return
}
