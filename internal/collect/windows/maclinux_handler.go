// Copyright (c) 2020 Alec Randazzo

// This is a dummy file that exist purely so we can run tests on mac and linux systems

//go:build darwin || linux

package windows

import (
	"github.com/AlecRandazzo/Packrat/pkg/parsers/windows/vbr"
	"io"
	"os"
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
	return nil
}
