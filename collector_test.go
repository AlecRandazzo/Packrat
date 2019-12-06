/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package windowscollector

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	vbr "github.com/Go-Forensics/VBR-Parser"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"testing"
)

func TestCollect(t *testing.T) {
	type args struct {
		exportList   ListOfFilesToExport
		resultWriter ResultWriter
		handler      Handler
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Collect(tt.args.handler, tt.args.exportList, tt.args.resultWriter); (err != nil) != tt.wantErr {
				t.Errorf("Collect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getFiles(t *testing.T) {
	type args struct {
		volumeHandler        *VolumeHandler
		resultWriter         ZipResultWriter
		listOfSearchKeywords listOfSearchTerms
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		dummyFile   string
		testZip     string
		wantZipHash string
	}{
		{
			name: "test1",
			args: args{
				volumeHandler: &VolumeHandler{},
				resultWriter:  ZipResultWriter{},
				listOfSearchKeywords: listOfSearchTerms{
					0: searchTerms{
						fullPathString: `c:\$mft`,
						fullPathRegex:  nil,
						fileNameString: "$mft",
						fileNameRegex:  nil,
					},
					1: searchTerms{
						fullPathString: `c:\\$mftmirr`,
						fullPathRegex:  nil,
						fileNameString: "$mftmirr",
						fileNameRegex:  nil,
					},
				},
			},
			dummyFile:   `test\testdata\dummyntfs`,
			testZip:     `test\testdata\getFilesTest1.zip`,
			wantErr:     false,
			wantZipHash: "e37081c5c97884bd419cfadaa281f77a",
		},
		{
			name: "test2",
			args: args{
				volumeHandler: &VolumeHandler{},
				resultWriter:  ZipResultWriter{},
				listOfSearchKeywords: listOfSearchTerms{
					0: searchTerms{
						fullPathString: `c:\\$mftmirr`,
						fullPathRegex:  nil,
						fileNameString: "$mftmirr",
						fileNameRegex:  nil,
					},
				},
			},
			dummyFile:   `test\testdata\dummyntfs`,
			testZip:     `test\testdata\getFilesTest2.zip`,
			wantErr:     false,
			wantZipHash: "04c3f56fb7388624c0119eee3c97cae2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileHandle, _ := os.Create(tt.testZip)
			zipWriter := zip.NewWriter(fileHandle)
			tt.args.resultWriter = ZipResultWriter{
				ZipWriter:  zipWriter,
				FileHandle: fileHandle,
			}
			dummyHandle := &dummyHandler{
				Handle:               nil,
				VolumeLetter:         "c",
				Vbr:                  vbr.VolumeBootRecord{},
				mftReader:            nil,
				lastReadVolumeOffset: 0,
				filePath:             tt.dummyFile,
			}
			var err error
			*tt.args.volumeHandler, err = GetVolumeHandler("c", dummyHandle)
			if err != nil {
				log.Panic(err)
			}
			defer tt.args.volumeHandler.Handle.Close()

			_ = getFiles(tt.args.volumeHandler, &tt.args.resultWriter, tt.args.listOfSearchKeywords)

			// Get file hash
			file, _ := os.Open(tt.testZip)
			defer file.Close()
			hash := md5.New()
			_, _ = io.Copy(hash, file)
			hashInBytes := hash.Sum(nil)[:]
			gotZipHash := hex.EncodeToString(hashInBytes)
			if gotZipHash != tt.wantZipHash {
				t.Errorf("getFiles() gotZipHash = %v, want %v", gotZipHash, tt.wantZipHash)
			}
		})
	}
}
