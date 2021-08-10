// Copyright (c) 2020 Alec Randazzo

package collector

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/windows/vbr"
)

func TestCollect(t *testing.T) {
	type args struct {
		exportList   ListOfFilesToExport
		resultWriter ZipResultWriter
		handler      handler
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		zipTestOutput string
		wantZipHash   string
	}{
		{
			name: "test1",
			args: args{
				exportList: ListOfFilesToExport{
					0: {
						FullPath:        `%SYSTEMDRIVE%:\$MFT`,
						IsFullPathRegex: false,
						FileName:        `$MFT`,
						IsFileNameRegex: false,
					},
				},
				resultWriter: ZipResultWriter{},
				handler: dummyHandler{
					Handle:               nil,
					VolumeLetter:         "",
					Vbr:                  vbr.VolumeBootRecord{},
					mftReader:            nil,
					lastReadVolumeOffset: 0,
					filePath:             filepath.FromSlash("../../test/testdata/dummyntfs"),
				},
			},
			wantErr:       false,
			zipTestOutput: filepath.FromSlash("../../test/testdata/collecttestzip.zip"),
			wantZipHash:   "29f689d96a790b68df7e84c9e04ef741",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileHandle, _ := os.Create(tt.zipTestOutput)
			zipWriter := zip.NewWriter(fileHandle)
			tt.args.resultWriter = ZipResultWriter{
				ZipWriter:  zipWriter,
				FileHandle: fileHandle,
			}
			_ = Collect(tt.args.handler, tt.args.exportList, &tt.args.resultWriter)
			zipWriter.Close()
			fileHandle.Close()

			// Get file hash
			file, _ := os.Open(tt.zipTestOutput)
			hash := md5.New()
			_, _ = io.Copy(hash, file)
			hashInBytes := hash.Sum(nil)[:]
			gotZipHash := hex.EncodeToString(hashInBytes)
			file.Close()
			if gotZipHash != tt.wantZipHash {
				t.Errorf("collect() gotZipHash = %v, want %v", gotZipHash, tt.wantZipHash)
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
			dummyFile:   filepath.FromSlash("../../test/testdata/dummyntfs"),
			testZip:     filepath.FromSlash("../../test/testdata/getFilesTest1.zip"),
			wantErr:     false,
			wantZipHash: "a50b885249c709ae97eeba0e2d6ec78d",
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
			dummyFile:   filepath.FromSlash("../../test/testdata/dummyntfs"),
			testZip:     filepath.FromSlash("../../test/testdata/getFilesTest2.zip"),
			wantErr:     false,
			wantZipHash: "75c57f05d2879cb723dbec6e2e1e8f83",
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
				t.Error("failed to get volume handle")
				return
			}
			defer tt.args.volumeHandler.Handle.Close()

			_ = getFiles(tt.args.volumeHandler, &tt.args.resultWriter, tt.args.listOfSearchKeywords)
			zipWriter.Close()
			fileHandle.Close()

			// Get file hash
			file, _ := os.Open(tt.testZip)
			hash := md5.New()
			_, _ = io.Copy(hash, file)
			hashInBytes := hash.Sum(nil)[:]
			gotZipHash := hex.EncodeToString(hashInBytes)
			file.Close()
			if gotZipHash != tt.wantZipHash {
				t.Errorf("getFiles() gotZipHash = %v, want %v", gotZipHash, tt.wantZipHash)
			}
		})
	}
}
