/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package gofor

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
)

type AttributeInfo struct {
	AttributeType  byte
	AttributeBytes []byte
}

// Get MFT record attributes list.
func (mftRecord *MasterFileTableRecord) getAttributeList() (err error) {
	const offsetAttributeSize = 0x04
	const lengthAttributeSize = 0x04

	//const offsetResidentFlag = 0x08
	//const offsetHeaderNameLength = 0x09
	//
	//const offsetAttributeType = 0x18
	//const lengthAttributeType = 0x04
	//
	//const offsetRecordLength = 0x1c
	//const lengthRecordLength = 0x02
	//
	//const offsetNameLength = 0x1e
	//const offsetNameOffset = 0x1f
	//
	//const offsetStartingVCN = 0x20
	//const lengthStartingVCN = 0x08
	//
	//const offsetBaseFileReference = 0x28
	//const lengthBaseFileReference = 0x08
	//
	//const offsetAttributeId = 0x30
	//const lengthAttributeId = 0x02
	//
	//const offsetName = 0x32

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic %s, hex dump: %s", fmt.Sprint(r), hex.EncodeToString(mftRecord.MftRecordBytes))
		}
	}()

	// Init variable that tracks how far to the next attribute
	var distanceToNext uint16 = 0

	// Init a byte slice that tracks all possible valid MFT attribute types.
	// We'll be used this to verify if what we are looking at is actually an MFT attribute.
	const codeStandardInformation = 0x10
	const codeAttributeList = 0x20
	const codeFileName = 0x30
	const codeVolumeVersion = 0x40
	const codeSecurityDescriptor = 0x50
	const codeVolumeName = 0x60
	const codeVolumeInformation = 0x70
	const codeData = 0x80
	const codeIndexRoot = 0x90
	const codeIndexAllocation = 0xA0
	const codeBitmap = 0xB0
	const codeSymbolicLink = 0xC0
	const codeReparsePoint = 0xD0
	const codeEaInformation = 0xE0
	const codePropertySet = 0xF0

	validAttributeTypes := []byte{
		codeStandardInformation,
		codeAttributeList,
		codeFileName,
		codeVolumeVersion,
		codeSecurityDescriptor,
		codeVolumeName,
		codeVolumeInformation,
		codeData,
		codeIndexRoot,
		codeIndexAllocation,
		codeBitmap,
		codeSymbolicLink,
		codeReparsePoint,
		codeEaInformation,
		codePropertySet,
	}

	for {
		// Calculate offset to next attribute
		mftRecord.RecordHeader.AttributesOffset = mftRecord.RecordHeader.AttributesOffset + distanceToNext
		lenBytesIn := len(mftRecord.MftRecordBytes)
		if mftRecord.RecordHeader.AttributesOffset > uint16(lenBytesIn) || mftRecord.RecordHeader.AttributesOffset+0x04 > uint16(lenBytesIn) {
			break
		}

		// Verify if the byte slice is actually an MFT attribute
		var validAttribute bool
		for _, validType := range validAttributeTypes {
			if mftRecord.MftRecordBytes[mftRecord.RecordHeader.AttributesOffset] == validType {
				validAttribute = true
				break
			} else {
				validAttribute = false
			}
		}
		if validAttribute == false {
			break
		}

		// Pull out information describing the attribute and the attribute bytes
		attributeInfoSlice := AttributeInfo{}
		attributeInfoSlice.AttributeType = mftRecord.MftRecordBytes[mftRecord.RecordHeader.AttributesOffset]
		attributeSize := binary.LittleEndian.Uint16(mftRecord.MftRecordBytes[mftRecord.RecordHeader.AttributesOffset+offsetAttributeSize : mftRecord.RecordHeader.AttributesOffset+offsetAttributeSize+lengthAttributeSize])
		//AttributeSize := binary.LittleEndian.Uint16(MftRecordBytes[attributeOffset+4:attributeOffset+5])
		end := uint16(mftRecord.RecordHeader.AttributesOffset) + attributeSize
		attributeInfoSlice.AttributeBytes = mftRecord.MftRecordBytes[mftRecord.RecordHeader.AttributesOffset:end]

		// Append the attribute to the attribute struct
		mftRecord.AttributeInfo = append(mftRecord.AttributeInfo, attributeInfoSlice)

		// Track the distance to the next attribute based on the size of the current attribute
		distanceToNext = attributeSize
		if distanceToNext == 0 {
			break
		}
	}
	return
}
