// Copyright (c) 2020 Alec Randazzo

package packrat

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
			wantZipHash:       "d333bc8a8de2682d40e8db32ffb090d8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.zipResultWriter.FileHandle, _ = os.Create(tt.zipToCreate)
			tt.zipResultWriter.ZipWriter = zip.NewWriter(tt.zipResultWriter.FileHandle)
			reader := bytes.NewReader(tt.dummyData)
			tt.listOfFileReaders = make([]fileReader, 1)
			tt.listOfFileReaders[0] = fileReader{
				fullPath: "test",
				reader:   reader,
			}
			tt.args.waitForFileCopying.Add(1)
			tt.args.fileReaders = make(chan fileReader, 0)
			go tt.zipResultWriter.ResultWriter(tt.args.fileReaders, tt.args.waitForFileCopying)
			for _, each := range tt.listOfFileReaders {
				tt.args.fileReaders <- each
			}
			close(tt.args.fileReaders)
			tt.args.waitForFileCopying.Wait()

			tt.zipResultWriter.ZipWriter.Close()
			tt.zipResultWriter.FileHandle.Close()

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
