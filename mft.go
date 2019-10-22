package windowscollector

import (
	"errors"
	"fmt"
	"github.com/AlecRandazzo/GoFor-MFT-Parser"
	log "github.com/sirupsen/logrus"
	syscall "golang.org/x/sys/windows"
)

func parseMFTRecord0(volume VolumeHandler) (mftRecord0 mft.MasterFileTableRecord, err error) {
	// Move handle pointer back to beginning of volume
	_, err = syscall.Seek(volume.Handle, 0x00, 0)
	if err != nil {
		err = fmt.Errorf("failed to see back to volume offset 0x00: %w", err)
		return
	}

	// Seek to the offset where the MFT starts. If it errors, bomb.
	_, err = syscall.Seek(volume.Handle, volume.Vbr.MftByteOffset, 0)
	if err != nil {
		err = fmt.Errorf("failed to seek to mft: %w", err)
		return
	}

	// Read the first entry in the MFT. The first record in the MFT always is for the MFT itself. If it errors, bomb.
	buffer := make([]byte, volume.Vbr.MftRecordSize)
	_, err = syscall.Read(volume.Handle, buffer)
	if err != nil {
		err = fmt.Errorf("failed to read the mft: %w", err)
		return
	}

	// Sanity check that this is indeed an mft record
	result, err := mft.RawMasterFileTableRecord(buffer).IsThisAnMftRecord()
	if err != nil {
		err = fmt.Errorf("IsThisAnMftRecord() returned an error: %v", err)
	} else if result == false {
		err = errors.New("VolumeHandler.parseMFTRecord0() received an invalid mft record")
		return
	}

	// Parse the MFT record

	mftRecord0, err = mft.RawMasterFileTableRecord(buffer).Parse(volume.Vbr.BytesPerCluster)
	if err != nil {
		err = fmt.Errorf("VolumeHandler.parseMFTRecord0() failed to parse the mft's mft record: %w", err)
		return
	}
	log.Debugf("Identified the following data runs for the MFT itself: %+v", mftRecord0.DataAttribute.NonResidentDataAttribute.DataRuns)

	return
}
