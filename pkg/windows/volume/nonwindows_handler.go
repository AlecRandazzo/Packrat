// Copyright (c) 2022 Alec Randazzo

// This is a dummy file that exists so we can run tests on mac and linux systems for go files that don't have a Windows dependency

//go:build !windows

package volume

import (
	"io"
	"os"

	"github.com/AlecRandazzo/Packrat/pkg/windows/vbr"
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
}

// GetHandle will get a file handle to the underlying NTFS volume. We need this in order to bypass file locks.
func (handler *Volume) GetHandle() error {
	return nil
}
