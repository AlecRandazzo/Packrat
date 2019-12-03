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
	"os"
	"sync"
	"testing"
)

type dummyResultWriter struct{}

func (dummy dummyResultWriter) ResultWriter(*chan fileReader, *sync.WaitGroup, *sync.WaitGroup) (err error) {

	return
}

func TestZipResultWriter_ResultWriter(t *testing.T) {
	type fields struct {
		ZipWriter  *zip.Writer
		FileHandle *os.File
	}
	type args struct {
		fileReaders           *chan fileReader
		waitForInitialization *sync.WaitGroup
		waitForFileCopying    *sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zipResultWriter := &ZipResultWriter{
				ZipWriter:  tt.fields.ZipWriter,
				FileHandle: tt.fields.FileHandle,
			}
			if err := zipResultWriter.ResultWriter(tt.args.fileReaders, tt.args.waitForFileCopying); (err != nil) != tt.wantErr {
				t.Errorf("ResultWriter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
