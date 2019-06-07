/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package gofor

import (
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func (file MftFile) mftToCSV(outFileName string, waitgroup *sync.WaitGroup) (err error) {
	outFile, err := os.Create(outFileName)
	if err != nil {
		err = errors.Wrapf(err, "ParseMFT(): failed to create output file %s", outFileName)
	}
	defer outFile.Close()
	csvWriter := csv.NewWriter(outFile)
	csvWriter.Comma = '|'
	csvHeader := []string{
		"Record Number",
		"Directory Flag",
		"System File Flag",
		"Hidden Flag",
		"Read-only Flag",
		"Deleted Flag",
		"File Path",
		"File Name",
		"File Size",
		"File Created",
		"File Modified",
		"File Accessed",
		"File Entry Modified",
		"FileName Created ",
		"FileName Modified ",
		"Filename Accessed ",
		"Filename Entry Modified ",
	}
	err = csvWriter.Write(csvHeader)
	if err != nil {
		log.Fatal(err)
	}

	openChannel := true
	for openChannel != false {
		var csvRow []string
		var mftRecord MasterFileTableRecord
		mftRecord, openChannel = <-file.outputChannel
		for _, record := range mftRecord.FileNameAttributes {
			if strings.Contains(record.FileNamespace, "WIN32") == true || strings.Contains(record.FileNamespace, "POSIX") {
				var fileDirectory string
				if value, ok := file.MappedDirectories[record.ParentDirRecordNumber]; ok {
					fileDirectory = value
				} else {
					fileDirectory = "$ORPHANFILE"
				}
				recordNumber := fmt.Sprint(mftRecord.RecordHeader.RecordNumber)
				physicalFileSize := strconv.FormatUint(record.PhysicalFileSize, 10)
				csvRow = []string{
					recordNumber, //Record Number
					strconv.FormatBool(mftRecord.RecordHeader.FlagDirectory), //Directory Flag
					strconv.FormatBool(record.FileNameFlags.System),          //System file flag
					strconv.FormatBool(record.FileNameFlags.Hidden),          //Hidden flag
					strconv.FormatBool(record.FileNameFlags.ReadOnly),        //Read only flag
					strconv.FormatBool(mftRecord.RecordHeader.FlagDeleted),   //Deleted Flag
					fileDirectory,    //File directory
					record.FileName,  //File Name
					physicalFileSize, // File Size
					mftRecord.StandardInformationAttributes.SiCreated,  //File Created
					mftRecord.StandardInformationAttributes.SiModified, //File Modified
					mftRecord.StandardInformationAttributes.SiAccessed, //File Accessed
					mftRecord.StandardInformationAttributes.SiChanged,  //File entry Modified
					record.FnCreated,  //FileName Created
					record.FnModified, //FileName Modified
					record.FnAccessed, //FileName Accessed
					record.FnChanged,  //FileName Entry Modified
				}
				err = csvWriter.Write(csvRow)
				if err != nil {
					log.Fatal(err)
				}
				break
			}
		}
	}
	waitgroup.Done()
	return
}
