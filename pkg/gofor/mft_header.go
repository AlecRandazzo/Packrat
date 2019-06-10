/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package gofor

import "encoding/binary"

type recordHeader struct {
	AttributesOffset uint16
	RecordNumber     uint32
	FlagDeleted      bool
	FlagDirectory    bool
}

func (mftRecord *masterFileTableRecord) getRecordHeader() {
	const offsetAttributesOffset = 0x14

	const offsetRecordNumber = 0x2c
	const lengthRecordNumber = 0x04

	mftRecord.RecordHeader.AttributesOffset = uint16(mftRecord.MftRecordBytes[offsetAttributesOffset])

	mftRecord.getHeaderFlags()

	valueRecordNumber := mftRecord.MftRecordBytes[offsetRecordNumber : offsetRecordNumber+lengthRecordNumber]
	mftRecord.RecordHeader.RecordNumber = binary.LittleEndian.Uint32(valueRecordNumber)
	return
}

func (mftRecord *masterFileTableRecord) getHeaderFlags() {
	const offsetRecordFlag = 0x16
	const codeDeletedFile = 0x00
	//const codeActiveFile = 0x01
	//const codeDeletedDirectory = 0x02
	const codeDirectory = 0x03
	recordFlag := mftRecord.MftRecordBytes[offsetRecordFlag]
	if recordFlag == codeDeletedFile {
		mftRecord.RecordHeader.FlagDeleted = true
		mftRecord.RecordHeader.FlagDirectory = false
	} else if recordFlag == codeDirectory {
		mftRecord.RecordHeader.FlagDirectory = true
		mftRecord.RecordHeader.FlagDeleted = false
	} else {
		mftRecord.RecordHeader.FlagDeleted = false
		mftRecord.RecordHeader.FlagDirectory = false
	}
	return
}
