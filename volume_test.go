// Copyright (c) 2020 Alec Randazzo

package packrat

import (
	"errors"
	vbr "github.com/AlecRandazzo/VBR-Parser"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestGetVolumeHandler(t *testing.T) {
	type args struct {
		volumeLetter  string
		volumeHandler dummyHandler
	}
	tests := []struct {
		name     string
		args     args
		wantVBR  vbr.VolumeBootRecord
		wantErr  bool
		filePath string
	}{
		{
			name: "success",
			args: args{
				volumeLetter:  "C",
				volumeHandler: dummyHandler{},
			},
			wantVBR: vbr.VolumeBootRecord{
				VolumeLetter:           "",
				BytesPerSector:         512,
				SectorsPerCluster:      8,
				BytesPerCluster:        4096,
				MftByteOffset:          4096,
				MftRecordSize:          1024,
				ClustersPerIndexRecord: 1,
			},
			wantErr:  false,
			filePath: `test\testdata\dummyntfs`,
		},
		{
			name: "bad volume letter",
			args: args{
				volumeLetter:  "error",
				volumeHandler: dummyHandler{},
			},
			wantErr:  true,
			filePath: `test\testdata\dummyntfs`,
		},
		{
			name: "bad vbr1",
			args: args{
				volumeLetter:  "C",
				volumeHandler: dummyHandler{},
			},
			wantErr:  true,
			filePath: `test\testdata\dummyntfs-badvbr1`,
		},
		{
			name: "bad vbr2",
			args: args{
				volumeLetter:  "C",
				volumeHandler: dummyHandler{},
			},
			wantErr:  true,
			filePath: `test\testdata\dummyntfs-badvbr2`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.volumeHandler.filePath = tt.filePath
			gotVolume, err := GetVolumeHandler(tt.args.volumeLetter, tt.args.volumeHandler)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVolumeHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVolume.Vbr, tt.wantVBR) {
				t.Errorf("GetVolumeHandler() gotVBR = %+v, want %+v", gotVolume.Vbr, tt.wantVBR)
			}
		})
	}
}

type dummyHandler struct {
	Handle               *os.File
	VolumeLetter         string
	Vbr                  vbr.VolumeBootRecord
	mftReader            io.Reader
	lastReadVolumeOffset int64
	filePath             string
}

func (dummy dummyHandler) GetHandle(volumeLetter string) (handle *os.File, err error) {
	if volumeLetter == "error" {
		err = errors.New("faux error")
		return
	}
	handle, _ = os.Open(dummy.filePath)
	return
}

func Test_GetHandle(t *testing.T) {
	type args struct {
		volumeLetter string
	}
	tests := []struct {
		name    string
		args    args
		volume  VolumeHandler
		wantErr bool
	}{
		{
			name:    "no error",
			args:    args{volumeLetter: "C"},
			wantErr: false,
		},
		{
			name:    "nil string input",
			args:    args{volumeLetter: ""},
			wantErr: true,
		},
		{
			name:    "bad input",
			args:    args{volumeLetter: "CD"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.volume.GetHandle(tt.args.volumeLetter)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHandle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_identifyVolumesOfInterest(t *testing.T) {
	type args struct {
		exportList *ListOfFilesToExport
	}
	tests := []struct {
		name                  string
		args                  args
		wantVolumesOfInterest []string
		wantErr               bool
	}{
		{
			name: "systemdrive and d",
			args: args{exportList: &ListOfFilesToExport{
				0: FileToExport{
					FullPath:        `%SYSTEMDRIVE%:\$MFT`,
					IsFullPathRegex: false,
					FileName:        "$MFT",
					IsFileNameRegex: false,
				},
				1: FileToExport{
					FullPath:        `D:\$MFT`,
					IsFullPathRegex: false,
					FileName:        "$MFT",
					IsFileNameRegex: false,
				},
				2: FileToExport{
					FullPath:        `D:\blah`,
					IsFullPathRegex: false,
					FileName:        "blah",
					IsFileNameRegex: false,
				},
			}},
			wantVolumesOfInterest: []string{"C", "d"},
			wantErr:               false,
		},
		{
			name: "not a real volume",
			args: args{exportList: &ListOfFilesToExport{
				0: FileToExport{
					FullPath:        `1:\$MFT`,
					IsFullPathRegex: false,
					FileName:        "$MFT",
					IsFileNameRegex: false,
				},
			}},
			wantVolumesOfInterest: nil,
			wantErr:               true,
		},
		{
			name: "bad input",
			args: args{exportList: &ListOfFilesToExport{
				0: FileToExport{
					FullPath:        `CD:\$MFT`,
					IsFullPathRegex: false,
					FileName:        "$MFT",
					IsFileNameRegex: false,
				},
			}},
			wantVolumesOfInterest: nil,
			wantErr:               true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVolumesOfInterest, err := identifyVolumesOfInterest(tt.args.exportList)
			if (err != nil) != tt.wantErr {
				t.Errorf("identifyVolumesOfInterest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVolumesOfInterest, tt.wantVolumesOfInterest) {
				t.Errorf("identifyVolumesOfInterest() gotVolumesOfInterest = %v, want %v", gotVolumesOfInterest, tt.wantVolumesOfInterest)
			}
		})
	}
}

func Test_isLetter(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
		wantErr    bool
	}{
		{
			name:       "letter c",
			args:       args{s: "C"},
			wantResult: true,
			wantErr:    false,
		},
		{
			name:       "nil input",
			args:       args{s: ""},
			wantResult: false,
			wantErr:    true,
		},
		{
			name:       "string length of 2",
			args:       args{s: "CC"},
			wantResult: false,
			wantErr:    true,
		},
		{
			name:       "number input",
			args:       args{s: "1"},
			wantResult: false,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := isLetter(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("isLetter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("isLetter() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
