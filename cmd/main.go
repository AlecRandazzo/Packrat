// Copyright (c) 2022 Alec Randazzo

package main

import (
	"archive/zip"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"

	"github.com/AlecRandazzo/Packrat/internal/collect/windows"
	"github.com/AlecRandazzo/Packrat/internal/parse"
)

var CLI struct {
	Collect struct {
		Throttle bool   `short:"t" optional:"" help:"Throttle the process to a single thread."`
		Output   string `short:"o" optional:"" help:"Output file. If not specified, the file name defaults to the host name and a timestamp."`
		Debug    bool   `short:"d" optional:"" help:"Debug mode"`
	} `cmd help:"Collect forensic data."`
	Parse struct {
		Mft    string `arg:"" short:"m" help:"Mft File"`
		Output string `arg:"" short:"o" help:"Output file"`
	} `cmd help:"Parse forensic data."`
}

const (
	defaultBytesPerSector  = 512
	defaultBytesPerCluster = 4096
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)

	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "collect":
		if CLI.Collect.Throttle {
			runtime.GOMAXPROCS(1)
		}
		if CLI.Collect.Debug {
			debugLog, err := os.Create("debug.log")
			if err != nil {
				panic(err)
			}
			log.SetOutput(debugLog)
			log.SetLevel(log.DebugLevel)
		}

		systemDrive := os.Getenv("SYSTEMDRIVE")
		exportList := windows.FileExportList{
			{
				FullPath:      fmt.Sprintf(`%s\$MFT`, systemDrive),
				FullPathRegex: false,
				FileName:      `$MFT`,
				FileNameRegex: false,
			},
			{
				FullPath:      fmt.Sprintf(`%s\\Windows\\System32\\winevt\\Logs\\.*\.evtx$`, systemDrive),
				FullPathRegex: true,
				FileName:      `.*\.evtx$`,
				FileNameRegex: true,
			},
			{
				FullPath:      fmt.Sprintf(`%s\Windows\System32\config\SYSTEM`, systemDrive),
				FullPathRegex: false,
				FileName:      `SYSTEM`,
				FileNameRegex: false,
			},
			{
				FullPath:      fmt.Sprintf(`%s\Windows\System32\config\SOFTWARE`, systemDrive),
				FullPathRegex: false,
				FileName:      `SOFTWARE`,
				FileNameRegex: false,
			},
			{
				FullPath:      fmt.Sprintf(`%s\\users\\([^\\]+)\\ntuser.dat`, systemDrive),
				FullPathRegex: true,
				FileName:      `ntuser.dat`,
				FileNameRegex: false,
			},
			{
				FullPath:      fmt.Sprintf(`%s\\Users\\([^\\]+)\\AppData\\Local\\Microsoft\\Windows\\usrclass.dat`, systemDrive),
				FullPathRegex: true,
				FileName:      `usrclass.dat`,
				FileNameRegex: false,
			},
			{
				FullPath:      fmt.Sprintf(`%s\\Users\\([^\\]+)\\AppData\\Local\\Microsoft\\Windows\\WebCache\\WebCacheV01.dat`, systemDrive),
				FullPathRegex: true,
				FileName:      `WebCacheV01.dat`,
				FileNameRegex: false,
			},
		}

		var zipName string
		if CLI.Collect.Output != "" {
			zipName = CLI.Collect.Output
		} else {
			hostName, err := os.Hostname()
			if err != nil {
				panic(err)
			}

			zipName = fmt.Sprintf("%s_%s.zip", hostName, time.Now().Format("2006-01-02T15.04.05Z"))
		}
		fileHandle, err := os.Create(zipName)
		if err != nil {
			err = fmt.Errorf("failed to create zip file %s", zipName)
		}
		defer fileHandle.Close()

		zipWriter := zip.NewWriter(fileHandle)
		//resultWriter := collect.ZipResultWriter{
		//	ZipWriter:  zipWriter,
		//	FileHandle: fileHandle,
		//}
		defer zipWriter.Close()

		err = windows.Collect(exportList, zipWriter)
		if err != nil {
			log.Panic(err)
		}
	case "parse <mft>":
		writer, err := parse.NewCsvWriter(CLI.Parse.Output)
		if err != nil {
			panic(err)
		}

		err = parse.Parse(CLI.Parse.Mft, writer, defaultBytesPerSector, defaultBytesPerCluster)
		if err != nil {
			panic(err)
		}
	default:
		ctx.Command()
	}
}
