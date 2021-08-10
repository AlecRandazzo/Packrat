// Copyright (c) 2020 Alec Randazzo

package main

import (
	"archive/zip"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/disk"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2" //TODO Replace with kong

	"github.com/AlecRandazzo/Packrat/internal/collector"
)

var (
	debug      = kingpin.Flag("debug", "Enable debug mode.").Bool()
	all        = kingpin.Flag("all", "Collect all forensic artifacts.").Bool()
	mft        = kingpin.Flag("mft", "Collect the system drive MFT.").Bool()
	mftAll     = kingpin.Flag("mft-all", "Collect all attached volume MFTs.").Bool()
	mftLetters = kingpin.Flag("mft-letters", "Collect all attached volume MFTs.").Strings()
	reg        = kingpin.Flag("reg", "Collect all registry hives, both system and user hives.").Bool()
	events     = kingpin.Flag("events", "Collect all event logs.").Bool()
	browser    = kingpin.Flag("browser", "Collect browser history").Bool()
	config     = kingpin.Flag("custom-config", "Custom configuration file that will overwrite built in config.").File()
	throttle   = kingpin.Flag("throttle", "This setting will limit the process to a single thread. This will reduce the CPU load.").Bool()
	output     = kingpin.Flag("output", "Specify the name of the output file. If not specified, the file name defaults to the host name and a timestamp.").String()
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	kingpin.Parse()
	if *throttle {
		runtime.GOMAXPROCS(1)
	}

	if *debug {
		debugLog, _ := os.Create("debug.log")
		log.SetOutput(debugLog)
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.ErrorLevel)
	}

	exportList := make(collector.ListOfFilesToExport, 0)
	if *config != nil {
		configData, err := ioutil.ReadAll(*config)
		if err != nil {
			log.Panic(err)
		}

		err = yaml.Unmarshal(configData, &exportList)
		if err != nil {
			log.Panic(err)
		}
	} else {
		if *all {
			*mftAll = true
			*events = true
			*reg = true
			*browser = true
		}
		if *mftAll {
			volumes, _ := disk.Partitions(true)
			for _, volume := range volumes {
				file := collector.FileToExport{
					FullPath:        fmt.Sprintf(`%s\$MFT`, volume.Mountpoint),
					IsFullPathRegex: false,
					FileName:        `$MFT`,
					IsFileNameRegex: false,
				}
				exportList = append(exportList, file)
			}
		} else if len(*mftLetters) > 0 {
			for _, v := range *mftLetters {
				file := collector.FileToExport{
					FullPath:        fmt.Sprintf(`%s:\$MFT`, v),
					IsFullPathRegex: false,
					FileName:        `$MFT`,
					IsFileNameRegex: false,
				}
				exportList = append(exportList, file)
			}
		} else if *mft {
			file := collector.FileToExport{
				FullPath:        `%SYSTEMDRIVE%:\$MFT`,
				IsFullPathRegex: false,
				FileName:        `$MFT`,
				IsFileNameRegex: false,
			}
			exportList = append(exportList, file)
		}

		if *events {
			file := collector.FileToExport{
				FullPath:        `%SYSTEMDRIVE%:\\Windows\\System32\\winevt\\Logs\\.*\.evtx$`,
				IsFullPathRegex: true,
				FileName:        `.*\.evtx$`,
				IsFileNameRegex: true,
			}
			exportList = append(exportList, file)
		}

		if *reg {
			files := collector.ListOfFilesToExport{
				{
					FullPath:        `%SYSTEMDRIVE%:\Windows\System32\config\SYSTEM`,
					IsFullPathRegex: false,
					FileName:        `SYSTEM`,
					IsFileNameRegex: false,
				},
				{
					FullPath:        `%SYSTEMDRIVE%:\Windows\System32\config\SOFTWARE`,
					IsFullPathRegex: false,
					FileName:        `SOFTWARE`,
					IsFileNameRegex: false,
				},
				{
					FullPath:        `%SYSTEMDRIVE%:\\users\\([^\\]+)\\ntuser.dat`,
					IsFullPathRegex: true,
					FileName:        `ntuser.dat`,
					IsFileNameRegex: false,
				},
				{
					FullPath:        `%SYSTEMDRIVE%:\\Users\\([^\\]+)\\AppData\\Local\\Microsoft\\Windows\\usrclass.dat`,
					IsFullPathRegex: true,
					FileName:        `usrclass.dat`,
					IsFileNameRegex: false,
				},
			}
			for _, v := range files {
				exportList = append(exportList, v)
			}
		}

		if *browser {
			file := collector.FileToExport{
				FullPath:        `%SYSTEMDRIVE%:\\Users\\([^\\]+)\\AppData\\Local\\Microsoft\\Windows\\WebCache\\WebCacheV01.dat`,
				IsFullPathRegex: true,
				FileName:        `WebCacheV01.dat`,
				IsFileNameRegex: false,
			}
			exportList = append(exportList, file)
		}
	}

	var zipName string
	if *output != "" {
		zipName = *output
	} else {
		hostName, _ := os.Hostname()
		zipName = fmt.Sprintf("%s_%s.zip", hostName, time.Now().Format("2006-01-02T15.04.05Z"))
	}
	fileHandle, err := os.Create(zipName)
	if err != nil {
		err = fmt.Errorf("failed to create zip file %s", zipName)
	}
	defer fileHandle.Close()

	zipWriter := zip.NewWriter(fileHandle)
	resultWriter := collector.ZipResultWriter{
		ZipWriter:  zipWriter,
		FileHandle: fileHandle,
	}
	defer zipWriter.Close()
	var volume collector.VolumeHandler
	err = collector.Collect(volume, exportList, &resultWriter)
	if err != nil {
		log.Panic(err)
	}
}
