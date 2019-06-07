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
	"bytes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

type MasterFileTableRecord struct {
	bytesPerCluster               int64
	RecordHeader                  RecordHeader
	StandardInformationAttributes StandardInformationAttributes
	FileNameAttributes            []FileNameAttributes
	DataAttributes                DataAttributes
	MftRecordBytes                []byte
	AttributeInfo                 []AttributeInfo
}

type MftFile struct {
	FileHandle        *os.File
	MappedDirectories map[uint64]string
	outputChannel     chan MasterFileTableRecord
}

// Parse an already extracted MFT and write the results to a file.
func ParseMFT(mftFilePath, outFileName string) (err error) {
	file := MftFile{}
	file.FileHandle, err = os.Open(mftFilePath)
	if err != nil {
		err = errors.Wrapf(err, "failed to open MFT file %s", mftFilePath)
		return
	}
	defer file.FileHandle.Close()

	err = file.buildDirectoryTree()
	if err != nil {
		return
	}

	file.outputChannel = make(chan MasterFileTableRecord, 100)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go file.mftToCSV(outFileName, &waitGroup)

	var offset int64 = 0
	for {
		buffer := make([]byte, 1024)
		_, err = file.FileHandle.ReadAt(buffer, offset)
		if err == io.EOF {
			err = nil
			break
		}
		mftRecord := MasterFileTableRecord{}
		mftRecord.MftRecordBytes = buffer
		err = mftRecord.ParseMFTRecord()
		if err != nil {
			log.WithFields(log.Fields{
				"mft_offset":   offset,
				"deleted_flag": mftRecord.RecordHeader.FlagDeleted,
			}).Debug(err)
			offset += 1024
			continue
		}
		file.outputChannel <- mftRecord
		offset += 1024
		if len(mftRecord.FileNameAttributes) == 0 {
			continue
		}

	}
	close(file.outputChannel)
	waitGroup.Wait()
	return
}

// Parse the bytes of an MFT record
func (mftRecord *MasterFileTableRecord) ParseMFTRecord() (err error) {

	recordHeaderPresent := mftRecord.checkForRecordHeader()
	if recordHeaderPresent == false {
		return
	}

	mftRecord.extractValidPortionOfMFTRecord()

	// Get attributes from the MFT record
	mftRecord.getRecordHeader()

	err = mftRecord.getAttributeList()
	if err != nil {
		err = errors.Wrap(err, "failed to get attribute list")
		return
	}

	err = mftRecord.getStandardInformationAttribute()
	if err != nil {
		err = errors.Wrap(err, "failed to get standard information attribute")
		return
	}

	err = mftRecord.getFileNameAttributes()
	if err != nil {
		err = errors.Wrap(err, "failed to get file name attributes")
		return
	}
	err = mftRecord.getDataAttribute()
	if err != nil {
		err = errors.Wrap(err, "failed to get data attribute")
		return
	}
	return
}

// Extract everything before the end sequence 0xffffffff
func (mftRecord *MasterFileTableRecord) extractValidPortionOfMFTRecord() {
	lenMftRecordBytes := len(mftRecord.MftRecordBytes)
	mftRecordEndByteSequence := []byte{0xff, 0xff, 0xff, 0xff}
	for i := 0; i < (lenMftRecordBytes - 4); i++ {
		if bytes.Equal(mftRecord.MftRecordBytes[i:i+0x04], mftRecordEndByteSequence) {
			mftRecord.MftRecordBytes = mftRecord.MftRecordBytes[:i]
			break
		}
	}
}

// Verifies that the bytes receives is actually an MFT record. All MFT records start with "FILE0".
func (mftRecord *MasterFileTableRecord) checkForRecordHeader() (recordHeaderPresent bool) {
	const offsetRecordMagicNumber = 0x00
	const lengthRecordMagicNumber = 0x05
	valueRecordHeader := string(mftRecord.MftRecordBytes[offsetRecordMagicNumber : offsetRecordMagicNumber+lengthRecordMagicNumber])
	if valueRecordHeader == "FILE0" {
		recordHeaderPresent = true
	} else {
		recordHeaderPresent = false
	}
	return
}
