package main

import "github.com/AlecRandazzo/Packrat/internal/parse"

type ParseCmd struct {
	Target string `arg:"" short:"t" help:"Target file, archive, or folder"`
	Output string `arg:"" short:"o" help:"Output file"`
}

func (c *ParseCmd) Run() error {
	writer, err := parse.NewCsvWriter(CLI.Parse.Output)
	if err != nil {
		return err
	}

	err = parse.Parse(CLI.Parse.Target, writer, defaultBytesPerSector, defaultBytesPerCluster)
	if err != nil {
		return err
	}
	return nil
}
