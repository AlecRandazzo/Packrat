// Copyright (c) 2022 Alec Randazzo

//go:build windows

package windows

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
	"io"
	"os"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/windows/vbr"
)

type handler interface {
	GetHandle() error
	VolumeLetter() string
	Handle() *os.File
	Reader() io.Reader
	UpdateReader(io.Reader)
	Vbr() vbr.VolumeBootRecord
	LastOffset() int64
	UpdateLastOffset(int64)
}

// VolumeHandler contains everything needed for basic collection functionality
type VolumeHandler struct {
	handle       *os.File
	volumeLetter string
	vbr          vbr.VolumeBootRecord
	reader       io.Reader
	lastOffset   int64
}

func NewVolumeHandler(driveLetter string) *VolumeHandler {
	return &VolumeHandler{
		volumeLetter: driveLetter,
	}
}

func (volumeHandler VolumeHandler) VolumeLetter() string {
	return volumeHandler.volumeLetter
}

func (volumeHandler VolumeHandler) Handle() *os.File {
	return volumeHandler.handle
}

func (volumeHandler *VolumeHandler) UpdateReader(newReader io.Reader) {
	volumeHandler.reader = newReader
	return
}

func (volumeHandler VolumeHandler) Vbr() vbr.VolumeBootRecord {
	return volumeHandler.vbr
}

func (volumeHandler VolumeHandler) Reader() io.Reader {
	return volumeHandler.reader
}

func (volumeHandler VolumeHandler) LastOffset() int64 {
	return volumeHandler.lastOffset
}

func (volumeHandler *VolumeHandler) UpdateLastOffset(newOffset int64) {
	volumeHandler.lastOffset = newOffset
	return
}

// GetHandle will get a file handle to the underlying NTFS volume. We need this in order to bypass file locks.
func (volumeHandler *VolumeHandler) GetHandle() error {
	dwDesiredAccess := uint32(0x80000000) //0x80 FILE_READ_ATTRIBUTES
	dwShareMode := uint32(0x02 | 0x01)
	dwCreationDisposition := uint32(0x03)
	dwFlagsAndAttributes := uint32(0x00)

	// We check to see if volumeHandler is set to account for testing. If we weren't doing tests, then the code in this
	// if statement could live outside the if statement
	if volumeHandler.handle == nil {
		volumePath, _ := windows.UTF16PtrFromString(fmt.Sprintf("\\\\.\\%s:", volumeHandler.volumeLetter))
		syscallHandle, err := windows.CreateFile(volumePath, dwDesiredAccess, dwShareMode, nil, dwCreationDisposition, dwFlagsAndAttributes, 0)
		if err != nil {
			return fmt.Errorf("getHandle() failed to get volumeHandle to volumeHandler %s: %w", volumeHandler.volumeLetter, err)
		}
		volumeHandler.handle = os.NewFile(uintptr(syscallHandle), "")
	}

	// Parse the VBR to get details we need about the volume.
	volumeBootRecord := make([]byte, 512)
	_, err := volumeHandler.handle.Read(volumeBootRecord)
	if err != nil {
		return fmt.Errorf("GetHandle() failed to read the volume boot record on volume %v: %w", volumeHandler.volumeLetter, err)
	}
	volumeHandler.vbr, err = vbr.Parse(volumeBootRecord)
	if err != nil {
		return fmt.Errorf("NewOldVolumeHandler() failed to parse vbr from volume letter %s: %w", volumeHandler.volumeLetter, err)
	}
	log.Debugf("Successfully got a file handle to volume %v and read its volume boot record.", volumeHandler.volumeLetter)

	return nil
}
