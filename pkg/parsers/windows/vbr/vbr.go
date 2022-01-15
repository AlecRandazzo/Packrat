// Copyright (c) 2020 Alec Randazzo

package vbr

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/byteshelper"
	"math"
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

var (
	magicNumber                  = []byte("NTFS")
	magicNumberLocation          = byteshelper.NewDataLocation(0x03, 0x04)
	bytesPerSectorLocation       = byteshelper.NewDataLocation(0x0B, 0x02)
	sectorsPerClusterLocation    = byteshelper.NewDataLocation(0x0D, 0x01)
	clustersPerMftRecordLocation = byteshelper.NewDataLocation(0x40, 0x01)
	clustersPerOffsetLocation    = byteshelper.NewDataLocation(0x30, 0x08)
	clustersPerIndexRecord       = byteshelper.NewDataLocation(0x44, 0x01)
)

// Parse parses a byte slice containing an NTFS volume boot record (VBR)
func Parse(input []byte) (vbr VolumeBootRecord, err error) {
	// Sanity check that we have the right data
	vbrSize := len(input)
	if vbrSize == 0 {
		return VolumeBootRecord{}, errors.New("received nil bytes")
	} else if vbrSize < 512 {
		return VolumeBootRecord{}, errors.New("received less than 512 bytes")
	}

	var buffer []byte

	// Sanity check to verify that the function actually received a VBR. Bomb if we didn't.
	buffer, err = byteshelper.GetValue(input, magicNumberLocation)
	if bytes.Compare(magicNumber, buffer) != 0 {
		return VolumeBootRecord{}, errors.New("received byte slice that does not start with 'NTFS' magic number")
	}

	// Pull out data based on pre-defined offsets in the VBR
	buffer, err = byteshelper.GetValue(input, bytesPerSectorLocation)
	vbr.BytesPerSector = int64(binary.LittleEndian.Uint16(buffer))

	buffer, err = byteshelper.GetValue(input, sectorsPerClusterLocation)
	vbr.SectorsPerCluster = int64(buffer[0])

	buffer, err = byteshelper.GetValue(input, clustersPerMftRecordLocation)
	clustersPerMFTRecord := int(buffer[0])
	if clustersPerMFTRecord < 128 {
		return VolumeBootRecord{}, fmt.Errorf("found the clusters per MFT record is %d, which is less than 128", clustersPerMFTRecord)
	}
	signedTwosComplement := int8(buffer[0]) * -1
	vbr.MftRecordSize = int64(math.Pow(2, float64(signedTwosComplement)))
	vbr.BytesPerCluster = vbr.SectorsPerCluster * vbr.BytesPerSector

	buffer, err = byteshelper.GetValue(input, clustersPerOffsetLocation)
	var mftClusterOffset int64
	mftClusterOffset, err = byteshelper.LittleEndianBinaryToInt64(buffer)
	if mftClusterOffset == 0 {
		return VolumeBootRecord{}, fmt.Errorf("failed to get mft offset clusters: %w", err)
	}
	vbr.MftByteOffset = mftClusterOffset * vbr.BytesPerCluster

	buffer, err = byteshelper.GetValue(input, clustersPerIndexRecord)
	vbr.ClustersPerIndexRecord = int64(buffer[0])

	return vbr, nil
}
