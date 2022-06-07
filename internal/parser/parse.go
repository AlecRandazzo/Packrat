// Copyright (c) 2022 Alec Randazzo

package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecRandazzo/Packrat/pkg/windows/mft"
)

type Format int

const (
	Csv Format = iota
)

const outputDirectory = "results"

func Parse(target string, format Format, bytesPerSector, bytesPerCluster uint) error {
	targetType, err := getTargetType(target)
	if err != nil {
		return fmt.Errorf("target not supported: %w", err)
	}

	err = os.Mkdir(outputDirectory, os.ModePerm)
	if err != nil {
		if err.Error() != "mkdir results: file exists" {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	switch targetType {
	case mftType:
		err = handleMft(target, format, bytesPerSector, bytesPerCluster)
		if err != nil {
			return fmt.Errorf("failed to parser mft file: %w", err)
		}
	case zipType:
		_ = ""
	case directoryType:
		_ = ""
	case registryType:
		return errors.New("registry hives not support yet")
	case eventLogType:
		return errors.New("evtx not support yet")
	case unknown:
		return fmt.Errorf("target not supported: %w", err)
	}

	return nil
}

func handleMft(target string, format Format, bytesPerSector, bytesPerCluster uint) error {
	reader, err := os.Open(target)
	if err != nil {
		return fmt.Errorf("failed to open mft file: %w", err)
	}
	defer reader.Close()

	var outFile *os.File

	switch format {
	case Csv:
		fp := filepath.Join(outputDirectory, filepath.Base(target)+".csv")
		outFile, err = os.Create(fp)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		var writer mft.CsvWriter
		writer, err = mft.NewCsvWriter(outFile)
		if err != nil {
			return fmt.Errorf("failed to create csvWriter: %w", err)
		}
		err = mft.ParseFile(reader, writer, bytesPerSector, bytesPerCluster)
		if err != nil {
			return fmt.Errorf("failed to parser mft: %w", err)
		}
	default:
		return fmt.Errorf("unknown format enum: %d", format)
	}

	return nil
}
