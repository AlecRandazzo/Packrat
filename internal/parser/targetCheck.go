package parser

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/windows/evtx"
	"github.com/AlecRandazzo/Packrat/pkg/parsers/windows/mft"
	"github.com/AlecRandazzo/Packrat/pkg/parsers/windows/registry"
)

type targetType int

const (
	unknown targetType = iota
	directoryType
	zipType
	mftType
	registryType
	eventLogType
)

var (
	zipMagicNumber = []byte{0x50, 0x4B, 0x03, 0x04} // PK
)

func getTargetType(target string) (targetType, error) {
	// Figure out if the target is a directory
	stats, err := os.Stat(target)
	if err != nil {
		return unknown, fmt.Errorf("could not determine what the target is: %w", err)
	}
	if stats.IsDir() {
		return directoryType, nil
	}

	// Check what kind of file it is
	var f *os.File
	f, err = os.Open(target)
	if err != nil {
		return unknown, fmt.Errorf("could not determine what the target is: %w", err)
	}
	defer f.Close()

	buf := make([]byte, 10)
	_, err = f.Read(buf)
	if err != nil {
		return unknown, fmt.Errorf("failed to read file: %w", err)
	}

	if bytes.Compare(buf[:len(zipMagicNumber)], zipMagicNumber) == 0 {
		return zipType, nil
	}

	if bytes.Compare(buf[:len(mft.MagicNumber)], mft.MagicNumber) == 0 {
		return mftType, nil
	}

	if bytes.Compare(buf[:len(registry.MagicNumber)], registry.MagicNumber) == 0 {
		return registryType, nil
	}

	if bytes.Compare(buf[:len(evtx.MagicNumber)], evtx.MagicNumber) == 0 {
		return eventLogType, nil
	}

	return unknown, errors.New("unknown target type")
}
