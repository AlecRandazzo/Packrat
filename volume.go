/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package windowscollector

import (
	"fmt"
	//mft "github.com/AlecRandazzo/GoFor-MFT-Parser"
	vbr "github.com/AlecRandazzo/GoFor-VBR-Parser"
	log "github.com/sirupsen/logrus"
	syscall "golang.org/x/sys/windows"
)

type VolumeHandler struct {
	Handle       syscall.Handle
	VolumeLetter string
	Vbr          vbr.VolumeBootRecord
	//MappedDirectories map[uint64]string
	//MftRecord0        mft.MasterFileTableRecord
}

func getHandle(volumeLetter string) (handle syscall.Handle, err error) {
	dwDesiredAccess := uint32(0x80000000) //0x80 FILE_READ_ATTRIBUTES
	dwShareMode := uint32(0x02 | 0x01)
	dwCreationDisposition := uint32(0x03)
	dwFlagsAndAttributes := uint32(0x00)

	volumePath, _ := syscall.UTF16PtrFromString(fmt.Sprintf("\\\\.\\%s:", volumeLetter))
	handle, err = syscall.CreateFile(volumePath, dwDesiredAccess, dwShareMode, nil, dwCreationDisposition, dwFlagsAndAttributes, 0)
	if err != nil {
		err = fmt.Errorf("getHandle() failed to get handle to volume %s: %w", volumeLetter, err)
		return
	}
	return
}

// GetVolumeHandler gets a file handle to the specified volume and parses its volume boot record.
func GetVolumeHandler(volumeLetter string) (volume VolumeHandler, err error) {
	const volumeBootRecordSize = 512
	volume.VolumeLetter = volumeLetter
	volume.Handle, err = getHandle(volumeLetter)
	if err != nil {
		err = fmt.Errorf("GetVolumeHandler() failed to get handle to volume %s: %w", volumeLetter, err)
		return
	}

	// Parse the VBR to get details we need about the volume.
	volumeBootRecord := make([]byte, volumeBootRecordSize)
	_, err = syscall.Read(volume.Handle, volumeBootRecord)
	if err != nil {
		err = fmt.Errorf("GetVolumeHandler() failed to read the volume boot record on volume %v: %w", volumeLetter, err)
		return
	}
	volume.Vbr, err = vbr.RawVolumeBootRecord(volumeBootRecord).Parse()
	if err != nil {
		err = fmt.Errorf("GetVolumeHandler() failed to parse vbr from volume letter %s: %w", volumeLetter, err)
		return
	}
	log.Debugf("Successfully got a file handle to volume %v and read its volume boot record.", volumeLetter)
	return
}
