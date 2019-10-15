package windowscollector

import (
	"errors"
	"fmt"
	"github.com/AlecRandazzo/GoFor-MFT-Parser"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

func parseMFTRecord0(volume VolumeHandler) (mftRecord0 mft.MasterFileTableRecord, err error) {
	// Move handle pointer back to beginning of volume
	_, err = windows.Seek(volume.Handle, 0x00, 0)
	if err != nil {
		err = fmt.Errorf("failed to see back to volume offset 0x00: %w", err)
		return
	}

	// Seek to the offset where the MFT starts. If it errors, bomb.
	_, err = windows.Seek(volume.Handle, volume.Vbr.MftByteOffset, 0)
	if err != nil {
		err = fmt.Errorf("failed to seek to mft: %w", err)
		return
	}

	// Read the first entry in the MFT. The first record in the MFT always is for the MFT itself. If it errors, bomb.
	buffer := make([]byte, volume.Vbr.MftRecordSize)
	_, err = windows.Read(volume.Handle, buffer)
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
	logrus.Debugf("Identified the following data runs for the MFT itself: %+v", mftRecord0.DataAttribute.NonResidentDataAttribute.DataRuns)

	return
}

func buildDirectoryTree(volumeHandler VolumeHandler, mftRecord0 mft.MasterFileTableRecord) (directoryTree mft.DirectoryTree, err error) {
	log.Debugf("Building directory tree for volume %v", volumeHandler.VolumeLetter)
	unresolvedDirectoryTree := mft.UnresolvedDirectoryTree{}
	dataRunNumber := 0
	numberOfDataRuns := len(mftRecord0.DataAttribute.NonResidentDataAttribute.DataRuns)

	for dataRunNumber < numberOfDataRuns {
		dataRunReader := DataRunsReader{
			VolumeHandler:          volumeHandler,
			DataRuns:               mftRecord0.DataAttribute.NonResidentDataAttribute.DataRuns,
			fileName:               "$MFT",
			dataRunTracker:         0,
			bytesLeftToReadTracker: 0,
			initialized:            false,
		}
		dataRunNumber++
		tempUnresolvedDirectoryTree := mft.UnresolvedDirectoryTree{}
		tempUnresolvedDirectoryTree, err = mft.BuildUnresolvedDirectoryTree(&dataRunReader)
		if err != nil {
			err = fmt.Errorf("failed to build an unresolved directory tree for mft data starting at datarun %+v", dataRunReader.DataRuns)
			return
		}
		log.Debugf("Found %v directories that need resolution in the MFT datarun %+v. Adding these to the master unresolved directory tracker.", len(tempUnresolvedDirectoryTree), dataRunReader.DataRuns)

		// Merge temporary directory tree with the master tree
		for recordNumber, directory := range tempUnresolvedDirectoryTree {
			unresolvedDirectoryTree[recordNumber] = directory
		}
	}

	directoryTree, _ = unresolvedDirectoryTree.Resolve(volumeHandler.VolumeLetter)
	log.Debugf("Resolved %v directories.", len(directoryTree))

	return
}
