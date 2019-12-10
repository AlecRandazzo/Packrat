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
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"sync"
	"testing"
)

func TestZipResultWriter_ResultWriter(t *testing.T) {
	type args struct {
		fileReaders        chan fileReader
		waitForFileCopying *sync.WaitGroup
	}
	tests := []struct {
		name              string
		args              args
		wantErr           bool
		dummyData         []byte
		listOfFileReaders []fileReader
		zipToCreate       string
		wantZipHash       string
		zipResultWriter   ZipResultWriter
	}{
		{
			name: "test1",
			zipResultWriter: ZipResultWriter{
				ZipWriter:  nil,
				FileHandle: nil,
			},
			wantErr: false,
			args: args{
				fileReaders:        nil,
				waitForFileCopying: &sync.WaitGroup{},
			},
			dummyData:         []byte{0x00, 0x00, 0x00},
			listOfFileReaders: []fileReader{},
			zipToCreate:       `test\testdata\test.zip`,
			wantZipHash:       "84a75bc35ad74c12cf7225a0fe802f07",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileHandle, _ := os.Create(tt.zipToCreate)
			zipWriter := zip.NewWriter(fileHandle)
			tt.zipResultWriter.FileHandle = fileHandle
			tt.zipResultWriter.ZipWriter = zipWriter
			reader := bytes.NewReader(tt.dummyData)
			tt.listOfFileReaders = make([]fileReader, 1)
			tt.listOfFileReaders[0] = fileReader{
				fullPath: "test",
				reader:   reader,
			}
			tt.args.waitForFileCopying.Add(1)
			channel := make(chan fileReader, 0)
			tt.args.fileReaders = channel
			go tt.zipResultWriter.ResultWriter(tt.args.fileReaders, tt.args.waitForFileCopying)
			for _, each := range tt.listOfFileReaders {
				tt.args.fileReaders <- each
			}
			close(tt.args.fileReaders)
			tt.args.waitForFileCopying.Wait()

			// Get file hash
			file, _ := os.Open(tt.zipToCreate)
			defer file.Close()
			hash := md5.New()
			_, _ = io.Copy(hash, file)
			hashInBytes := hash.Sum(nil)[:]
			gotZipHash := hex.EncodeToString(hashInBytes)
			if gotZipHash != tt.wantZipHash {
				t.Errorf("ZipResultWriter.resultWriter() gotZipHash = %v, want %v", gotZipHash, tt.wantZipHash)
			}
		})
	}
}
