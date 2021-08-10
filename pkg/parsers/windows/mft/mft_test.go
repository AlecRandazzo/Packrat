// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestRawMasterFileTableRecord_TrimSlackSpace(t *testing.T) {
	tests := []struct {
		name string
		want RawMasterFileTableRecord
		got  RawMasterFileTableRecord
	}{
		{
			name: "test1",
			got:  RawMasterFileTableRecord([]byte{0xba, 0xdb, 0xff, 0xff, 0xff, 0xff, 0x00}),
			want: RawMasterFileTableRecord([]byte{0xba, 0xdb}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.got.trimSlackSpace()
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestRawMasterFileTableRecord_Parse(t *testing.T) {
	type args struct {
		bytesPerCluster int64
	}
	tests := []struct {
		name          string
		rawMftRecord  RawMasterFileTableRecord
		args          args
		wantMftRecord MasterFileTableRecord
		wantErr       bool
	}{
		{
			name:         "test1",
			args:         args{bytesPerCluster: 4096},
			rawMftRecord: []byte{70, 73, 76, 69, 48, 0, 3, 0, 113, 250, 76, 78, 8, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 216, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 199, 5, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 81, 3, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 50, 96, 5, 194, 0, 56, 67, 16, 219, 0, 78, 89, 133, 0, 66, 176, 108, 91, 31, 119, 255, 66, 192, 69, 205, 200, 190, 0, 66, 0, 56, 8, 170, 148, 0, 66, 128, 80, 188, 200, 136, 1, 66, 64, 25, 2, 118, 2, 253, 66, 64, 85, 48, 135, 101, 2, 0, 176, 0, 0, 0, 80, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 27, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 192, 1, 0, 0, 0, 0, 0, 8, 176, 1, 0, 0, 0, 0, 0, 8, 176, 1, 0, 0, 0, 0, 0, 49, 25, 115, 210, 0, 65, 3, 176, 243, 197, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 49, 1, 255, 255, 11, 49, 1, 38, 0, 244, 0, 0, 0, 0, 199, 5, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 199, 5},
			wantErr:      false,
			wantMftRecord: MasterFileTableRecord{
				RecordHeader: RecordHeader{
					AttributesOffset: 56,
					RecordNumber:     0,
					Flags: RecordHeaderFlags{
						FlagDeleted:   false,
						FlagDirectory: false,
					},
				},
				StandardInformationAttributes: StandardInformationAttribute{
					SiCreated:    time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
					SiModified:   time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
					SiAccessed:   time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
					SiChanged:    time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
					FlagResident: true,
				},
				FileNameAttributes: FileNameAttributes{
					0: FileNameAttribute{
						FnCreated:               time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						FnModified:              time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						FnAccessed:              time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						FnChanged:               time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						FlagResident:            true,
						NameLength:              NameLength{},
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
				DataAttribute: DataAttribute{
					TotalSize:             0,
					FlagResident:          false,
					ResidentDataAttribute: nil,
					NonResidentDataAttribute: NonResidentDataAttribute{
						DataRuns: DataRuns{
							0: DataRun{
								AbsoluteOffset: 3221225472,
								Length:         209846272,
							},
							1: DataRun{
								AbsoluteOffset: 18254405632,
								Length:         5636096,
							},
							2: DataRun{
								AbsoluteOffset: 54049964032,
								Length:         229703680,
							},
							3: DataRun{
								AbsoluteOffset: 17307185152,
								Length:         113967104,
							},
							4: DataRun{
								AbsoluteOffset: 68520476672,
								Length:         73138176,
							},
							5: DataRun{
								AbsoluteOffset: 108427214848,
								Length:         58720256,
							},
							6: DataRun{
								AbsoluteOffset: 213864398848,
								Length:         84410368,
							},
							7: DataRun{
								AbsoluteOffset: 8366579712,
								Length:         26476544,
							},
							8: DataRun{
								AbsoluteOffset: 173059268608,
								Length:         89391104,
							},
						},
					},
				},
			},
		},
		{
			name:         "nil bytes",
			rawMftRecord: RawMasterFileTableRecord(nil),
			wantErr:      true,
			args:         args{bytesPerCluster: 4096},
		},
		{
			name:         "bytes per cluster arg of 0",
			rawMftRecord: []byte{70, 73, 76, 69, 48, 0, 3, 0, 113, 250, 76, 78, 8, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 216, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 199, 5, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 81, 3, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 50, 96, 5, 194, 0, 56, 67, 16, 219, 0, 78, 89, 133, 0, 66, 176, 108, 91, 31, 119, 255, 66, 192, 69, 205, 200, 190, 0, 66, 0, 56, 8, 170, 148, 0, 66, 128, 80, 188, 200, 136, 1, 66, 64, 25, 2, 118, 2, 253, 66, 64, 85, 48, 135, 101, 2, 0, 176, 0, 0, 0, 80, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 27, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 192, 1, 0, 0, 0, 0, 0, 8, 176, 1, 0, 0, 0, 0, 0, 8, 176, 1, 0, 0, 0, 0, 0, 49, 25, 115, 210, 0, 65, 3, 176, 243, 197, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 49, 1, 255, 255, 11, 49, 1, 38, 0, 244, 0, 0, 0, 0, 199, 5, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 199, 5},
			wantErr:      true,
			args:         args{bytesPerCluster: 0},
		},
		{
			name:         "raw record not large enough 1",
			rawMftRecord: []byte{70, 73, 76, 69, 48},
			wantErr:      true,
			args:         args{bytesPerCluster: 4096},
		},
		{
			name:         "raw record not large enough 2",
			rawMftRecord: []byte{70, 73, 76},
			wantErr:      true,
			args:         args{bytesPerCluster: 4096},
		},
		{
			name:         "invalid mft record",
			rawMftRecord: []byte{0, 73, 76, 69, 48, 0, 3, 0, 113, 250, 76, 78, 8, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 216, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 199, 5, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 81, 3, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 50, 96, 5, 194, 0, 56, 67, 16, 219, 0, 78, 89, 133, 0, 66, 176, 108, 91, 31, 119, 255, 66, 192, 69, 205, 200, 190, 0, 66, 0, 56, 8, 170, 148, 0, 66, 128, 80, 188, 200, 136, 1, 66, 64, 25, 2, 118, 2, 253, 66, 64, 85, 48, 135, 101, 2, 0, 176, 0, 0, 0, 80, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 27, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 192, 1, 0, 0, 0, 0, 0, 8, 176, 1, 0, 0, 0, 0, 0, 8, 176, 1, 0, 0, 0, 0, 0, 49, 25, 115, 210, 0, 65, 3, 176, 243, 197, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 49, 1, 255, 255, 11, 49, 1, 38, 0, 244, 0, 0, 0, 0, 199, 5, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 199, 5},
			wantErr:      true,
			args:         args{bytesPerCluster: 4096},
		},
		{
			name:         "attribute offset in record header is 0",
			rawMftRecord: []byte{70, 73, 76, 69, 48, 0, 3, 0, 113, 250, 76, 78, 8, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 216, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 199, 5, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 102, 248, 4, 21, 205, 173, 211, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 81, 3, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 0, 0, 32, 53, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 50, 96, 5, 194, 0, 56, 67, 16, 219, 0, 78, 89, 133, 0, 66, 176, 108, 91, 31, 119, 255, 66, 192, 69, 205, 200, 190, 0, 66, 0, 56, 8, 170, 148, 0, 66, 128, 80, 188, 200, 136, 1, 66, 64, 25, 2, 118, 2, 253, 66, 64, 85, 48, 135, 101, 2, 0, 176, 0, 0, 0, 80, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 27, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 192, 1, 0, 0, 0, 0, 0, 8, 176, 1, 0, 0, 0, 0, 0, 8, 176, 1, 0, 0, 0, 0, 0, 49, 25, 115, 210, 0, 65, 3, 176, 243, 197, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 49, 1, 255, 255, 11, 49, 1, 38, 0, 244, 0, 0, 0, 0, 199, 5, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 199, 5},
			wantErr:      true,
			args:         args{bytesPerCluster: 4096},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMftRecord, err := tt.rawMftRecord.Parse(tt.args.bytesPerCluster)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMftRecord, tt.wantMftRecord) {
				t.Errorf("parse() \ngotMftRecord = %v\nwant %v", gotMftRecord, tt.wantMftRecord)
			}
		})
	}
}

func TestGetUsefulMftFields(t *testing.T) {
	type args struct {
		mftRecord     MasterFileTableRecord
		directoryTree DirectoryTree
	}
	tests := []struct {
		name                string
		args                args
		wantUseFulMftFields UsefulMftFields
	}{
		{
			name: "test1",
			args: args{
				mftRecord: MasterFileTableRecord{
					RecordHeader: RecordHeader{
						AttributesOffset: 56,
						RecordNumber:     0,
						Flags: RecordHeaderFlags{
							FlagDeleted:   false,
							FlagDirectory: false,
						},
					},
					StandardInformationAttributes: StandardInformationAttribute{
						SiCreated:    time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						SiModified:   time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						SiAccessed:   time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						SiChanged:    time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						FlagResident: true,
					},
					FileNameAttributes: FileNameAttributes{
						0: FileNameAttribute{
							FnCreated:               time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
							FnModified:              time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
							FnAccessed:              time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
							FnChanged:               time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
							FlagResident:            true,
							NameLength:              NameLength{},
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
					DataAttribute: DataAttribute{
						TotalSize:             0,
						FlagResident:          false,
						ResidentDataAttribute: nil,
						NonResidentDataAttribute: NonResidentDataAttribute{
							DataRuns: DataRuns{
								0: DataRun{
									AbsoluteOffset: 3221225472,
									Length:         209846272,
								},
								1: DataRun{
									AbsoluteOffset: 18254405632,
									Length:         5636096,
								},
								2: DataRun{
									AbsoluteOffset: 54049964032,
									Length:         229703680,
								},
								3: DataRun{
									AbsoluteOffset: 17307185152,
									Length:         113967104,
								},
								4: DataRun{
									AbsoluteOffset: 68520476672,
									Length:         73138176,
								},
								5: DataRun{
									AbsoluteOffset: 108427214848,
									Length:         58720256,
								},
								6: DataRun{
									AbsoluteOffset: 213864398848,
									Length:         84410368,
								},
								7: DataRun{
									AbsoluteOffset: 8366579712,
									Length:         26476544,
								},
								8: DataRun{
									AbsoluteOffset: 173059268608,
									Length:         89391104,
								},
							},
						},
					},
				},
				directoryTree: DirectoryTree{},
			},
			wantUseFulMftFields: UsefulMftFields{
				RecordNumber:     0,
				FilePath:         "$ORPHANFILE\\",
				FullPath:         "$ORPHANFILE\\$MFT",
				FileName:         "$MFT",
				SystemFlag:       true,
				HiddenFlag:       true,
				ReadOnlyFlag:     false,
				DirectoryFlag:    false,
				DeletedFlag:      false,
				FnCreated:        time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				FnModified:       time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				FnAccessed:       time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				FnChanged:        time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				SiCreated:        time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				SiModified:       time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				SiAccessed:       time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				SiChanged:        time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				PhysicalFileSize: 16384,
			},
		},
		{
			name: "test2",
			args: args{
				mftRecord: MasterFileTableRecord{
					RecordHeader: RecordHeader{
						AttributesOffset: 56,
						RecordNumber:     0,
						Flags: RecordHeaderFlags{
							FlagDeleted:   false,
							FlagDirectory: false,
						},
					},
					StandardInformationAttributes: StandardInformationAttribute{
						SiCreated:    time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						SiModified:   time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						SiAccessed:   time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						SiChanged:    time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
						FlagResident: true,
					},
					FileNameAttributes: FileNameAttributes{
						0: FileNameAttribute{
							FnCreated:               time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
							FnModified:              time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
							FnAccessed:              time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
							FnChanged:               time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
							FlagResident:            true,
							NameLength:              NameLength{},
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
					DataAttribute: DataAttribute{
						TotalSize:             0,
						FlagResident:          false,
						ResidentDataAttribute: nil,
						NonResidentDataAttribute: NonResidentDataAttribute{
							DataRuns: DataRuns{
								0: DataRun{
									AbsoluteOffset: 3221225472,
									Length:         209846272,
								},
								1: DataRun{
									AbsoluteOffset: 18254405632,
									Length:         5636096,
								},
								2: DataRun{
									AbsoluteOffset: 54049964032,
									Length:         229703680,
								},
								3: DataRun{
									AbsoluteOffset: 17307185152,
									Length:         113967104,
								},
								4: DataRun{
									AbsoluteOffset: 68520476672,
									Length:         73138176,
								},
								5: DataRun{
									AbsoluteOffset: 108427214848,
									Length:         58720256,
								},
								6: DataRun{
									AbsoluteOffset: 213864398848,
									Length:         84410368,
								},
								7: DataRun{
									AbsoluteOffset: 8366579712,
									Length:         26476544,
								},
								8: DataRun{
									AbsoluteOffset: 173059268608,
									Length:         89391104,
								},
							},
						},
					},
				},
				directoryTree: DirectoryTree{
					5: "\\",
				},
			},
			wantUseFulMftFields: UsefulMftFields{
				RecordNumber:     0,
				FilePath:         "\\",
				FullPath:         "\\$MFT",
				FileName:         "$MFT",
				SystemFlag:       true,
				HiddenFlag:       true,
				ReadOnlyFlag:     false,
				DirectoryFlag:    false,
				DeletedFlag:      false,
				FnCreated:        time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				FnModified:       time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				FnAccessed:       time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				FnChanged:        time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				SiCreated:        time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				SiModified:       time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				SiAccessed:       time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				SiChanged:        time.Date(2018, 2, 25, 00, 10, 45, 642455000, time.UTC),
				PhysicalFileSize: 16384,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotUseFulMftFields := GetUsefulMftFields(tt.args.mftRecord, tt.args.directoryTree); !reflect.DeepEqual(gotUseFulMftFields, tt.wantUseFulMftFields) {
				t.Errorf("GetUsefulMftFields() = %v, want %v", gotUseFulMftFields, tt.wantUseFulMftFields)
			}
		})
	}
}

type WriteToSlice []UsefulMftFields

func (writer *WriteToSlice) ResultWriter(streamer io.Writer, outputChannel *chan UsefulMftFields, waitGroup *sync.WaitGroup) {
	openChannel := true
	for openChannel != false {
		var record UsefulMftFields
		record, openChannel = <-*outputChannel
		*writer = append(*writer, record)
	}
	waitGroup.Done()
	return
}

func TestParseMftRecords(t *testing.T) {
	type args struct {
		reader          io.Reader
		bytesPerCluster int64
		directoryTree   DirectoryTree
		outputChannel   chan UsefulMftFields
	}
	tests := []struct {
		name string
		args args
		want []UsefulMftFields
	}{
		{
			name: "test1",
			args: args{
				reader:          bytes.NewReader([]byte{46, 0x49, 0x4C, 0x45, 0x30, 0x00, 0x03, 0x00, 0x71, 0xFA, 0x4C, 0x4E, 0x08, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x38, 0x00, 0x01, 0x00, 0xD8, 0x01, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xC7, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x48, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x68, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x03, 0x00, 0x4A, 0x00, 0x00, 0x00, 0x18, 0x00, 0x01, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x03, 0x24, 0x00, 0x4D, 0x00, 0x46, 0x00, 0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x01, 0x00, 0x40, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0x51, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x35, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x35, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x35, 0x00, 0x00, 0x00, 0x00, 0x33, 0x20, 0xC8, 0x00, 0x00, 0x00, 0x0C, 0x32, 0x60, 0x05, 0xC2, 0x00, 0x38, 0x43, 0x10, 0xDB, 0x00, 0x4E, 0x59, 0x85, 0x00, 0x42, 0xB0, 0x6C, 0x5B, 0x1F, 0x77, 0xFF, 0x42, 0xC0, 0x45, 0xCD, 0xC8, 0xBE, 0x00, 0x42, 0x00, 0x38, 0x08, 0xAA, 0x94, 0x00, 0x42, 0x80, 0x50, 0xBC, 0xC8, 0x88, 0x01, 0x42, 0x40, 0x19, 0x02, 0x76, 0x02, 0xFD, 0x42, 0x40, 0x55, 0x30, 0x87, 0x65, 0x02, 0x00, 0xB0, 0x00, 0x00, 0x00, 0x50, 0x00, 0x00, 0x00, 0x01, 0x00, 0x40, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xC0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0xB0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0xB0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x31, 0x19, 0x73, 0xD2, 0x00, 0x41, 0x03, 0xB0, 0xF3, 0xC5, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x31, 0x01, 0xFF, 0xFF, 0x0B, 0x31, 0x01, 0x26, 0x00, 0xF4, 0x00, 0x00, 0x00, 0x00, 0xC7, 0x05, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xC7, 0x05, 0x46, 0x49, 0x4C, 0x45, 0x30, 0x00, 0x03, 0x00, 0x40, 0x14, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x38, 0x00, 0x01, 0x00, 0x58, 0x01, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0xC7, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x48, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x70, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x02, 0x00, 0x52, 0x00, 0x00, 0x00, 0x18, 0x00, 0x01, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x66, 0xF8, 0x04, 0x15, 0xCD, 0xAD, 0xD3, 0x01, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x03, 0x24, 0x00, 0x4D, 0x00, 0x46, 0x00, 0x54, 0x00, 0x4D, 0x00, 0x69, 0x00, 0x72, 0x00, 0x72, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x48, 0x00, 0x00, 0x00, 0x01, 0x00, 0x40, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x20, 0x00, 0x00, 0x00, 0x20, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x48, 0x00, 0x00, 0x00, 0x01, 0x00, 0x40, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xC7, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xC7, 0x05, 0x46}),
				bytesPerCluster: 4096,
				directoryTree: DirectoryTree{
					5: ".\\",
				},
				outputChannel: nil,
			},
			want: []UsefulMftFields{
				0: {
					RecordNumber:     1,
					FilePath:         ".\\",
					FullPath:         ".\\$MFTMirr",
					FileName:         "$MFTMirr",
					SystemFlag:       true,
					HiddenFlag:       true,
					ReadOnlyFlag:     false,
					DirectoryFlag:    false,
					DeletedFlag:      false,
					FnCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					PhysicalFileSize: 4096,
				},
				1: {
					RecordNumber:     0,
					FilePath:         "",
					FullPath:         "",
					FileName:         "",
					SystemFlag:       false,
					HiddenFlag:       false,
					ReadOnlyFlag:     false,
					DirectoryFlag:    false,
					DeletedFlag:      false,
					PhysicalFileSize: 0,
					FnCreated:        time.Time{},
					FnModified:       time.Time{},
					FnAccessed:       time.Time{},
					FnChanged:        time.Time{},
					SiCreated:        time.Time{},
					SiModified:       time.Time{},
					SiAccessed:       time.Time{},
					SiChanged:        time.Time{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.outputChannel = make(chan UsefulMftFields, 100)
			ParseMftRecords(tt.args.reader, tt.args.bytesPerCluster, tt.args.directoryTree, &tt.args.outputChannel)
			got := make([]UsefulMftFields, 0)
			openChannel := true
			for openChannel == true {
				var result UsefulMftFields
				result, openChannel = <-tt.args.outputChannel
				got = append(got, result)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestParseMFT(t *testing.T) {
	type args struct {
		fileHandle      *os.File
		writer          WriteToSlice
		streamer        io.Writer
		bytesPerCluster int64
		volumeLetter    string
	}
	tests := []struct {
		name     string
		args     args
		testFile string
		want     WriteToSlice
	}{
		{
			name: "test1",
			args: args{
				fileHandle:      nil,
				writer:          nil,
				bytesPerCluster: 4096,
				volumeLetter:    "C",
			},
			testFile: filepath.FromSlash("../../../../test/testdata/mft-lite"),
			want: WriteToSlice{
				0: {
					RecordNumber:     0,
					FilePath:         "C:\\",
					FullPath:         "C:\\$MFT",
					FileName:         "$MFT",
					SystemFlag:       true,
					HiddenFlag:       true,
					ReadOnlyFlag:     false,
					DirectoryFlag:    false,
					DeletedFlag:      false,
					FnCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					PhysicalFileSize: 16384,
				},
				1: {
					RecordNumber:     1,
					FilePath:         "C:\\",
					FullPath:         "C:\\$MFTMirr",
					FileName:         "$MFTMirr",
					SystemFlag:       true,
					HiddenFlag:       true,
					ReadOnlyFlag:     false,
					DirectoryFlag:    false,
					DeletedFlag:      false,
					FnCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					PhysicalFileSize: 4096,
				},
				2: {
					RecordNumber:     2,
					FilePath:         "C:\\",
					FullPath:         "C:\\$LogFile",
					FileName:         "$LogFile",
					SystemFlag:       true,
					HiddenFlag:       true,
					ReadOnlyFlag:     false,
					DirectoryFlag:    false,
					DeletedFlag:      false,
					FnCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					PhysicalFileSize: 67108864,
				},
				3: {
					RecordNumber:     3,
					FilePath:         "C:\\",
					FullPath:         "C:\\$Volume",
					FileName:         "$Volume",
					SystemFlag:       true,
					HiddenFlag:       true,
					ReadOnlyFlag:     false,
					DirectoryFlag:    false,
					DeletedFlag:      false,
					FnCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					PhysicalFileSize: 0,
				},
				4: {
					RecordNumber:     4,
					FilePath:         "C:\\",
					FullPath:         "C:\\$AttrDef",
					FileName:         "$AttrDef",
					SystemFlag:       true,
					HiddenFlag:       true,
					ReadOnlyFlag:     false,
					DirectoryFlag:    false,
					DeletedFlag:      false,
					FnCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					PhysicalFileSize: 2400,
				},
				5: {
					RecordNumber:     5,
					FilePath:         "C:\\",
					FullPath:         "C:\\.",
					FileName:         ".",
					SystemFlag:       true,
					HiddenFlag:       true,
					ReadOnlyFlag:     false,
					DirectoryFlag:    true,
					DeletedFlag:      false,
					FnCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiCreated:        time.Date(2017, 9, 29, 8, 45, 11, 680123300, time.UTC),
					SiModified:       time.Date(2019, 9, 8, 14, 53, 21, 936932600, time.UTC),
					SiAccessed:       time.Date(2019, 9, 8, 14, 53, 21, 936932600, time.UTC),
					SiChanged:        time.Date(2019, 9, 8, 14, 53, 21, 936932600, time.UTC),
					PhysicalFileSize: 0,
				},
				6: UsefulMftFields{
					RecordNumber:     0,
					FilePath:         "",
					FullPath:         "",
					FileName:         "",
					SystemFlag:       false,
					HiddenFlag:       false,
					ReadOnlyFlag:     false,
					DirectoryFlag:    false,
					DeletedFlag:      false,
					FnCreated:        time.Time{},
					FnModified:       time.Time{},
					FnAccessed:       time.Time{},
					FnChanged:        time.Time{},
					SiCreated:        time.Time{},
					SiModified:       time.Time{},
					SiAccessed:       time.Time{},
					SiChanged:        time.Time{},
					PhysicalFileSize: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.args.fileHandle, err = os.Open(tt.testFile)
			if err != nil {
				t.Errorf("failed to open handle to test file: %s", err.Error())
				return
			}
			ParseMFT(tt.args.volumeLetter, tt.args.fileHandle, &tt.args.writer, tt.args.streamer, tt.args.bytesPerCluster)
			if !reflect.DeepEqual(tt.args.writer, tt.want) {
				t.Errorf(cmp.Diff(tt.args.writer, tt.want))
			}
		})
	}
}
