package main

import (
	"github.com/AlecRandazzo/Packrat/internal/parser"
)

const (
	defaultBytesPerSector  = 512
	defaultBytesPerCluster = 4096
)

type ParseCmd struct {
	Target      string `arg:"" short:"t" help:"Target file, archive, or folder"`
	SectorSize  uint   `arg:"" optional:"" help:"NTFS sector size in bytes. If not provided, the parser will default to 512"`
	ClusterSize uint   `arg:"" optional:"" help:"NTFS cluster size in bytes. If not provided, the parser will default to 4096"`
}

func (c *ParseCmd) Run() error {
	var sectorSize, clusterSize uint
	if c.SectorSize == 0 {
		sectorSize = defaultBytesPerSector
	} else {
		sectorSize = c.SectorSize
	}

	if c.ClusterSize == 0 {
		clusterSize = defaultBytesPerCluster
	} else {
		clusterSize = c.ClusterSize
	}

	err := parser.Parse(CLI.Parse.Target, parser.Csv, sectorSize, clusterSize)
	if err != nil {
		return err
	}
	return nil
}
