package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/AlecRandazzo/Packrat/internal/collect/windows"
	"os"
	"runtime"
	"time"
)

type CollectCmd struct {
	Throttle bool   `short:"t" optional:"" help:"Throttle the process to a single thread."`
	Output   string `short:"o" optional:"" help:"Output file. If not specified, the file name defaults to the host name and a timestamp."`
}

func (c *CollectCmd) Run() error {
	if CLI.Collect.Throttle {
		runtime.GOMAXPROCS(1)
	}

	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case "windows":
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
			return fmt.Errorf("failed to create zip file %s: %w", zipName, err)
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
			return fmt.Errorf("failed to collect forensic data: %w", err)
		}
	case "darwin":
		return errors.New("mac forensic collection not implemented yet")
	case "linux":
		return errors.New("linux forensic collection not implemented yet")
	default:
		return fmt.Errorf("unsupported operating system: %s", operatingSystem)
	}

	return nil
}
