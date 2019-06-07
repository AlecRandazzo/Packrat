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
	"flag"
	"github.com/AlecRandazzo/GoFor/pkg/gofor"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Log configuration
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)
}

func main() {
	hostname, _ := os.Hostname()
	outZip := flag.String("zip", hostname, "")
	flag.Parse()

	client := gofor.CollectorClient{}
	exportList := gofor.ExportList{
		{FullPath: "C:\\\\$MFT", Type: "equal"},
		{FullPath: "C:\\Windows\\System32\\config\\SYSTEM", Type: "equal"},
		{FullPath: "C:\\Windows\\System32\\config\\SOFTWARE", Type: "equal"},
		{FullPath: "C:\\\\Windows\\\\System32\\\\winevt\\\\Logs\\\\.*\\.evtx$", Type: "regex"},
		{FullPath: "C:\\\\users\\\\([^\\\\]+)\\\\ntuser.dat$", Type: "regex"},
		{FullPath: "C:\\\\Users\\\\([^\\\\]+)\\\\AppData\\\\Local\\\\Microsoft\\\\Windows\\\\usrclass.dat$", Type: "regex"},
	}
	client.ExportToZip(exportList, *outZip)
}
