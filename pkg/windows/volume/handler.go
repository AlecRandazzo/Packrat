// Copyright (c) 2022 Alec Randazzo

//go:build windows

package volume

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/windows"

	"github.com/AlecRandazzo/Packrat/pkg/windows/vbr"
	log "github.com/sirupsen/logrus"
)

// Handler handles the file handle to an NTFS volume.
type Handler interface {
	GetHandle() error
	Letter() string
	Handle() *os.File
	Reader() io.Reader
	UpdateReader(io.Reader)
	Vbr() vbr.VolumeBootRecord
	LastOffset() int64
	UpdateLastOffset(int64)
}

// Volume contains everything needed to support an ntfs volume handle.
type Volume struct {
	handle     *os.File
	letter     string
	vbr        vbr.VolumeBootRecord
	reader     io.Reader
	lastOffset int64
}

// NewHandler creates a new NTFS volume handler.
func NewHandler(driveLetter string) (*Volume, error) {
	v := Volume{letter: driveLetter}
	err := v.GetHandle()
	return &v, err
}

// Letter returns the volume's assigned letter.
func (handler Volume) Letter() string {
	return handler.letter
}

// Handle returns the file handle to an NTFS volume
func (handler Volume) Handle() *os.File {
	return handler.handle
}

// UpdateReader changes the io.Reader
func (handler *Volume) UpdateReader(newReader io.Reader) {
	handler.reader = newReader
	return
}

// Vbr returns information about the volume boot record.
func (handler Volume) Vbr() vbr.VolumeBootRecord {
	return handler.vbr
}

// Reader reads an NTFS volume.
func (handler Volume) Reader() io.Reader {
	return handler.reader
}

// LastOffset returns the last read offset.
func (handler Volume) LastOffset() int64 {
	return handler.lastOffset
}

// UpdateLastOffset updates what the last offset read was.
func (handler *Volume) UpdateLastOffset(newOffset int64) {
	handler.lastOffset = newOffset
	return
}

// GetHandle will get a file handle to the underlying NTFS volume. We need this in order to bypass file locks.
func (handler *Volume) GetHandle() error {
	dwDesiredAccess := uint32(0x80000000) //0x80 FILE_READ_ATTRIBUTES
	dwShareMode := uint32(0x02 | 0x01)
	dwCreationDisposition := uint32(0x03)
	dwFlagsAndAttributes := uint32(0x00)

	// We check to see if handler is set to account for testing.
	if handler.handle == nil {
		volumePath, _ := windows.UTF16PtrFromString(fmt.Sprintf("\\\\.\\%s:", handler.letter))
		syscallHandle, err := windows.CreateFile(volumePath, dwDesiredAccess, dwShareMode, nil, dwCreationDisposition, dwFlagsAndAttributes, 0)
		if err != nil {
			return fmt.Errorf("getHandle() failed to get volumeHandle to handler %s: %w", handler.letter, err)
		}
		handler.handle = os.NewFile(uintptr(syscallHandle), "")
	}

	// Parse the VBR to get details we need about the volume.
	volumeBootRecord := make([]byte, 512)
	_, err := handler.handle.Read(volumeBootRecord)
	if err != nil {
		return fmt.Errorf("GetHandle() failed to read the volume boot record on volume %v: %w", handler.letter, err)
	}
	handler.vbr, err = vbr.Parse(volumeBootRecord)
	if err != nil {
		return fmt.Errorf("NewOldVolumeHandler() failed to parser vbr from volume letter %s: %w", handler.letter, err)
	}
	log.Debugf("Successfully got a file handle to volume %v and read its volume boot record.", handler.letter)

	return nil
}
