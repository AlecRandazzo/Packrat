// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_getFileNameAttribute(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    FileNameAttribute
		wantErr bool
	}{
		{
			name:    "TestFileNameAttribute_Parse test 1",
			wantErr: false,
			input:   []byte{48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0},
			want: FileNameAttribute{
				Created:      time.Date(2016, 7, 2, 15, 13, 30, 670820200, time.UTC),
				Modified:     time.Date(2016, 7, 2, 15, 13, 30, 670820200, time.UTC),
				Accessed:     time.Date(2016, 7, 2, 15, 13, 30, 670820200, time.UTC),
				Changed:      time.Date(2016, 7, 2, 15, 13, 30, 670820200, time.UTC),
				FlagResident: true,
				NameLength: NameLength{
					FlagNamed: false,
					NamedSize: 0,
				},
				AttributeSize:           104,
				ParentDirRecordNumber:   5,
				ParentDirSequenceNumber: 5,
				LogicalFileSize:         16384,
				PhysicalFileSize:        16384,
				FileNameFlags: FileNameFlags{
					ReadOnly:          false,
					Hidden:            true,
					System:            true,
					Archive:           false,
					Device:            false,
					Normal:            false,
					Temporary:         false,
					Sparse:            false,
					Reparse:           false,
					Compressed:        false,
					Offline:           false,
					NotContentIndexed: false,
					Encrypted:         false,
					Directory:         false,
					IndexView:         false,
				},
				FileNameLength: 8,
				FileNamespace:  "WIN32 & DOS",
				FileName:       "$MFT",
			},
		},
		{
			name:    "null input in",
			wantErr: true,
			input:   nil,
		},
		{
			name:    "non-resident",
			wantErr: true,
			input:   []byte{48, 0, 0, 0, 104, 0, 0, 1, 1, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0},
			want:    FileNameAttribute{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := getFileNameAttribute(tt.input)
			if !reflect.DeepEqual(fn, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf(cmp.Diff(fn, tt.want))
			}
		})
	}
}

func Test_getFileNameFlags(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  FileNameFlags
	}{
		{
			name:  "flag test 1",
			input: []byte{6, 0, 0, 0},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            true,
				System:            true,
				Archive:           false,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         false,
				IndexView:         false,
			},
		},
		{
			name:  "flag test 2",
			input: []byte{1, 0, 0, 16},
			want: FileNameFlags{
				ReadOnly:          true,
				Hidden:            false,
				System:            false,
				Archive:           false,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         true,
				IndexView:         false,
			},
		},
		{
			name:  "flag test 3",
			input: []byte{32, 0, 0, 0},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           true,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         false,
				IndexView:         false,
			},
		},
		{
			name:  "flag test 4",
			input: []byte{32, 33, 0, 0},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           true,
				Device:            false,
				Normal:            false,
				Temporary:         true,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: true,
				Encrypted:         false,
				Directory:         false,
				IndexView:         false,
			},
		},
		{
			name:  "flag test 5",
			input: []byte{32, 6, 0, 0},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           true,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            true,
				Reparse:           true,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         false,
				IndexView:         false,
			},
		},
		{
			name:  "flag test 6",
			input: []byte{6, 36, 0, 16},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            true,
				System:            true,
				Archive:           false,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           true,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: true,
				Encrypted:         false,
				Directory:         true,
				IndexView:         false,
			},
		},
		{
			name:  "flag test 7",
			input: []byte{0, 8, 0, 16},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           false,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        true,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         true,
				IndexView:         false,
			},
		},
		{
			name:  "flag test 8",
			input: []byte{32, 18, 64, 0},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           true,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            true,
				Reparse:           false,
				Compressed:        false,
				Offline:           true,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         false,
				IndexView:         false,
			},
		},
		{
			name:  "flag test 9",
			input: []byte{0, 32, 0, 16},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           false,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: true,
				Encrypted:         false,
				Directory:         true,
				IndexView:         false,
			},
		},
		{
			name:  "flag test 10",
			input: []byte{0, 0, 0, 16},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           false,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         true,
				IndexView:         false,
			},
		},
		{
			name:  "flag test 10",
			input: []byte{38, 0, 0, 32},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            true,
				System:            true,
				Archive:           true,
				Device:            false,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         false,
				IndexView:         true,
			},
		},
		{
			name:  "flag test 11",
			input: []byte{0x40, 0x40, 0x10, 0x10},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           false,
				Device:            true,
				Normal:            false,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         true,
				Directory:         true,
				IndexView:         false,
			},
		},
		{
			name:  "flag test 12",
			input: []byte{0x80, 0x00, 0x00, 0x00},
			want: FileNameFlags{
				ReadOnly:          false,
				Hidden:            false,
				System:            false,
				Archive:           false,
				Device:            false,
				Normal:            true,
				Temporary:         false,
				Sparse:            false,
				Reparse:           false,
				Compressed:        false,
				Offline:           false,
				NotContentIndexed: false,
				Encrypted:         false,
				Directory:         false,
				IndexView:         false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := getFileNameFlags(tt.input)
			if !reflect.DeepEqual(flags, tt.want) {
				t.Errorf(cmp.Diff(flags, tt.want))
			}
		})
	}
}

func Test_getFileNameSpace(t *testing.T) {
	tests := []struct {
		name  string
		input byte
		want  FileNamespace
	}{
		{
			name:  "POSIX",
			input: 0x00,
			want:  "POSIX",
		},
		{
			name:  "WIN32",
			input: 0x01,
			want:  "WIN32",
		},
		{
			name:  "DOS",
			input: 0x02,
			want:  "DOS",
		},
		{
			name:  "WIN32 & DOS",
			input: 0x03,
			want:  "WIN32 & DOS",
		},
		{
			name:  "null",
			input: 0x04,
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namespace := getFileNameSpace(tt.input)
			if namespace != tt.want {
				t.Errorf("Parse() = %v, want %v", namespace, tt.want)
			}
		})
	}
}
