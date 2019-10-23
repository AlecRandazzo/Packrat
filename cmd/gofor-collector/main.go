/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package main

import (
	collector "github.com/AlecRandazzo/GoFor-Windows-Collector"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Options struct {
	Debug string `short:"d" long:"debug" default:"" description:"Log debug information to output file."`
	//SendTo             string   `short:"s" long:"sendto" required:"true" description:"Where to send collected files to." choice:"zip"`
	ZipName            string `short:"z" long:"zipname" description:"Output file name for the zip." required:"true"`
	DataTypesToCollect string `short:"g" long:"gather" default:"a" description:"Types of data to collect. Concatenate the abbreviation characters together for what you want. The order doesn't matter. Valid values are 'a' for all, 'm' for $MFT, 'r' for system registries, 'u' for user registries, 'e' for event logs. Examples: '/g mrue', '/g a'"`
}

func init() {
	// Log configuration
	log.SetFormatter(&log.JSONFormatter{})
	// runtime.GOMAXPROCS(1)
}

func main() {
	opts := new(Options)
	parsedOpts := flags.NewParser(opts, flags.Default)
	_, err := parsedOpts.Parse()
	if err != nil {
		os.Exit(-1)
	}

	log.SetFormatter(&log.JSONFormatter{})
	if opts.Debug == "" {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.ErrorLevel)
	} else {
		debugLog, _ := os.Create(opts.Debug)
		defer debugLog.Close()
		log.SetOutput(debugLog)
		log.SetLevel(log.DebugLevel)
	}

	var exportList collector.ListOfFilesToExport
	if strings.Contains(opts.DataTypesToCollect, "a") {
		exportList = collector.ListOfFilesToExport{
			{
				FullPath:        `C:\$MFT`,
				IsFullPathRegex: false,
				FileName:        `$MFT`,
				IsFileNameRegex: false,
			},
			{
				FullPath:        `C:\Windows\System32\config\SYSTEM`,
				IsFullPathRegex: false,
				FileName:        `SYSTEM`,
				IsFileNameRegex: false,
			},
			{
				FullPath:        `C:\Windows\System32\config\SOFTWARE`,
				IsFullPathRegex: false,
				FileName:        `SOFTWARE`,
				IsFileNameRegex: false,
			},
			{
				FullPath:        `C:\\Windows\\System32\\winevt\\Logs\\.*\\.evtx$`,
				IsFullPathRegex: true,
				FileName:        `.*\\.evtx$`,
				IsFileNameRegex: true,
			},
			{
				FullPath:        `C:\\users\\([^\\]+)\\ntuser.dat`,
				IsFullPathRegex: true,
				FileName:        `ntuser.dat`,
				IsFileNameRegex: false,
			},
			{
				FullPath:        `C:\\Users\\([^\\]+)\\AppData\\Local\\Microsoft\\Windows\\usrclass.dat`,
				IsFullPathRegex: true,
				FileName:        `usrclass.dat`,
				IsFileNameRegex: false,
			},
		}
	} else {
		if strings.Contains(opts.DataTypesToCollect, "m") {
			exportList = append(exportList, collector.FileToExport{
				FullPath:        `C:\$MFT`,
				IsFullPathRegex: false,
				FileName:        `$MFT`,
				IsFileNameRegex: false,
			})
		}
		if strings.Contains(opts.DataTypesToCollect, "r") {
			exportList = append(exportList, collector.FileToExport{
				FullPath:        `C:\Windows\System32\config\SYSTEM`,
				IsFullPathRegex: false,
				FileName:        `SYSTEM`,
				IsFileNameRegex: false,
			})
			exportList = append(exportList, collector.FileToExport{
				FullPath:        `C:\Windows\System32\config\SOFTWARE`,
				IsFullPathRegex: false,
				FileName:        `SOFTWARE`,
				IsFileNameRegex: false,
			})
		}
		if strings.Contains(opts.DataTypesToCollect, "u") {
			exportList = append(exportList, collector.FileToExport{
				FullPath:        `C:\\users\\([^\\]+)\\ntuser.dat`,
				IsFullPathRegex: true,
				FileName:        `ntuser.dat`,
				IsFileNameRegex: false,
			})
			exportList = append(exportList, collector.FileToExport{
				FullPath:        `C:\\Users\\([^\\]+)\\AppData\\Local\\Microsoft\\Windows\\usrclass.dat`,
				IsFullPathRegex: true,
				FileName:        `usrclass.dat`,
				IsFileNameRegex: false,
			})
		}
		if strings.Contains(opts.DataTypesToCollect, "e") {
			exportList = append(exportList, collector.FileToExport{
				FullPath:        `C:\\Windows\\System32\\winevt\\Logs\\.*\\.evtx$`,
				IsFullPathRegex: true,
				FileName:        `.*\\.evtx$`,
				IsFileNameRegex: true,
			})
		}
	}

	resultWriter := collector.ZipResultWriter{ZipFileName: opts.ZipName}

	err = collector.Collect("C", exportList, resultWriter)
	if err != nil {
		log.Panic(err)
	}
}
