/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package GoFor_Collector

import (
	"encoding/binary"
	"errors"
	"fmt"
	mft "github.com/AlecRandazzo/Gofor-MFT-Parser"
	"math"
)

type VolumeBootRecord struct {
	VolumeLetter           string
	BytesPerSector         int64
	SectorsPerCluster      int64
	BytesPerCluster        int64
	MftByteOffset          int64
	MftRecordSize          int64
	ClustersPerIndexRecord int64
}

// Parses a byte slice containing an NTFS volume boot record (VBR)
func ParseVolumeBootRecord(volumeBootRecordBytes []byte) (vbr VolumeBootRecord, err error) {
	const codeNTFSMagicNumber = "NTFS"
	const offsetNTFSMagicNumber = 0x03
	const lengthNTFSMagicNumber = 0x04
	const offsetBytesPerSector = 0x0b
	const lengthBytesPerSector = 0x02
	const offsetSectorsPerCluster = 0x0d
	const offsetClustersPerMFTRecord = 0x40
	const offsetMftClusterOffset = 0x30
	const lengthMftClusterOffset = 0x08
	const offsetClustersPerIndexRecord = 0x44

	// Sanity check to verify that the function actually received a VBR. Bomb if we didn't.
	valueNTFSMagicNumber := string(volumeBootRecordBytes[offsetNTFSMagicNumber : offsetNTFSMagicNumber+lengthNTFSMagicNumber])
	if valueNTFSMagicNumber != codeNTFSMagicNumber {
		err = errors.New("vbr is invalid")
		return
	}

	// Start pulling out data based on pre-defined offsets in the VBR
	valueBytesPerSector := volumeBootRecordBytes[offsetBytesPerSector : offsetBytesPerSector+lengthBytesPerSector]
	vbr.BytesPerSector = int64(binary.LittleEndian.Uint16(valueBytesPerSector))
	vbr.SectorsPerCluster = int64(volumeBootRecordBytes[offsetSectorsPerCluster])
	clustersPerMFTRecord := int(volumeBootRecordBytes[offsetClustersPerMFTRecord])
	if clustersPerMFTRecord < 128 {
		err = fmt.Errorf("clusters per MFT record is %d, which is less than 128", clustersPerMFTRecord)
		return
	}
	signedTwosComplement := int8(volumeBootRecordBytes[0x40]) * -1
	vbr.MftRecordSize = int64(math.Pow(2, float64(signedTwosComplement)))
	vbr.BytesPerCluster = vbr.SectorsPerCluster * vbr.BytesPerSector
	valueMftClusterOffset := volumeBootRecordBytes[offsetMftClusterOffset : offsetMftClusterOffset+lengthMftClusterOffset]
	mftClusterOffset := mft.ConvertLittleEndianByteSliceToInt64(valueMftClusterOffset)
	if mftClusterOffset == 0 {
		err = fmt.Errorf("failed to get mft offset clusters: %w", err)
		return
	}
	vbr.MftByteOffset = mftClusterOffset * vbr.BytesPerCluster
	vbr.ClustersPerIndexRecord = int64(volumeBootRecordBytes[offsetClustersPerIndexRecord])
	return
}
