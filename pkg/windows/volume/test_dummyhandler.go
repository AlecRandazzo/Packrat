package volume

import (
	"fmt"
	"io"
	"os"

	"github.com/AlecRandazzo/Packrat/pkg/windows/vbr"
	log "github.com/sirupsen/logrus"
)

type Dummy struct {
	handle       *os.File
	volumeLetter string
	vbr          vbr.VolumeBootRecord
	reader       io.Reader
	lastOffset   int64
	FilePath     string
	err          error
}

func (handler *Dummy) GetHandle() error {
	if handler.err != nil {
		return handler.err
	}
	var err error
	handler.handle, err = os.Open(handler.FilePath)
	if err != nil {
		return fmt.Errorf("failed to get handle to file: %w", err)
	}
	// Parse the VBR to get details we need about the volume.
	volumeBootRecord := make([]byte, 512)
	_, err = handler.handle.Read(volumeBootRecord)
	if err != nil {
		return fmt.Errorf("GetHandle() failed to read the volume boot record on volume %v: %w", handler.volumeLetter, err)
	}
	handler.vbr, err = vbr.Parse(volumeBootRecord)
	if err != nil {
		return fmt.Errorf("NewOldVolumeHandler() failed to parser vbr from volume letter %s: %w", handler.volumeLetter, err)
	}
	log.Debugf("Successfully got a file handle to volume %v and read its volume boot record.", handler.volumeLetter)

	return err
}

func (handler Dummy) Letter() string {
	return handler.volumeLetter
}

func (handler Dummy) Handle() *os.File {
	return handler.handle
}

func (handler *Dummy) UpdateReader(newReader io.Reader) {
	handler.reader = newReader
}

func (handler Dummy) Vbr() vbr.VolumeBootRecord {
	return handler.vbr
}

func (handler Dummy) Reader() io.Reader {
	return handler.reader
}

func (handler Dummy) LastOffset() int64 {
	return handler.lastOffset
}

func (handler *Dummy) UpdateLastOffset(newOffset int64) {
	handler.lastOffset = newOffset
}
