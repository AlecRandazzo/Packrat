// Copyright (c) 2020 Alec Randazzo

package vbr

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/byteshelper"
)

// VolumeBootRecord contains relevant data about an NTFS volume
type VolumeBootRecord struct {
	VolumeLetter           string
	BytesPerSector         int64
	SectorsPerCluster      int64
	BytesPerCluster        int64
	MftByteOffset          int64
	MftRecordSize          int64
	ClustersPerIndexRecord int64
}

// RawVolumeBootRecord is a []byte alias containing bytes of a raw volume boot record. Used as a receiver for Parse().
type RawVolumeBootRecord []byte

// Parse parses a byte slice containing an NTFS volume boot record (VBR)
func (rawVolumeBootRecord RawVolumeBootRecord) Parse() (vbr VolumeBootRecord, err error) {
	// Sanity check that we have the right data
	sizeOfRawVolumeBootRecord := len(rawVolumeBootRecord)
	if sizeOfRawVolumeBootRecord == 0 {
		err = errors.New("RawVolumeBootRecord.Parse() received nil bytes")
		return
	} else if sizeOfRawVolumeBootRecord < 512 {
		err = errors.New("RawVolumeBootRecord.Parse() received less than 512 bytes")
		return
	}

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
	valueNTFSMagicNumber := string(rawVolumeBootRecord[offsetNTFSMagicNumber : offsetNTFSMagicNumber+lengthNTFSMagicNumber])
	if valueNTFSMagicNumber != codeNTFSMagicNumber {
		err = errors.New("RawVolumeBootRecord.Parse() received byte slice that does not start with 'NTFS' magic number")
		return
	}

	// Start pulling out data based on pre-defined offsets in the VBR
	valueBytesPerSector := rawVolumeBootRecord[offsetBytesPerSector : offsetBytesPerSector+lengthBytesPerSector]
	vbr.BytesPerSector = int64(binary.LittleEndian.Uint16(valueBytesPerSector))
	vbr.SectorsPerCluster = int64(rawVolumeBootRecord[offsetSectorsPerCluster])
	clustersPerMFTRecord := int(rawVolumeBootRecord[offsetClustersPerMFTRecord])
	if clustersPerMFTRecord < 128 {
		err = fmt.Errorf("RawVolumeBootRecord.Parse() found the clusters per MFT record is %d, which is less than 128", clustersPerMFTRecord)
		vbr = VolumeBootRecord{}
		return
	}
	signedTwosComplement := int8(rawVolumeBootRecord[0x40]) * -1
	vbr.MftRecordSize = int64(math.Pow(2, float64(signedTwosComplement)))
	vbr.BytesPerCluster = vbr.SectorsPerCluster * vbr.BytesPerSector
	valueMftClusterOffset := rawVolumeBootRecord[offsetMftClusterOffset : offsetMftClusterOffset+lengthMftClusterOffset]
	mftClusterOffset, err := byteshelper.LittleEndianBinaryToInt64(valueMftClusterOffset)
	if mftClusterOffset == 0 {
		err = fmt.Errorf("RawVolumeBootRecord.Parse() failed to get mft offset clusters: %w", err)
		vbr = VolumeBootRecord{}
		return
	}
	vbr.MftByteOffset = mftClusterOffset * vbr.BytesPerCluster
	vbr.ClustersPerIndexRecord = int64(rawVolumeBootRecord[offsetClustersPerIndexRecord])
	return
}
