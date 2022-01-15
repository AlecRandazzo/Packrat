// Copyright (c) 2020 Alec Randazzo

package collector

import (
	"github.com/google/go-cmp/cmp"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/windows/mft"
	"github.com/AlecRandazzo/Packrat/pkg/parsers/windows/vbr"
)

func Test_parseMFTRecord0(t *testing.T) {
	type args struct {
		dummyHandler *dummyHandler
	}
	tests := []struct {
		name           string
		args           args
		wantMftRecord0 mft.MasterFileTableRecord
		wantErr        bool
	}{
		{
			name:    "test1",
			wantErr: false,
			wantMftRecord0: mft.MasterFileTableRecord{
				RecordHeader: mft.RecordHeader{
					AttributesOffset: 56,
					RecordNumber:     0,
					Flags: mft.RecordHeaderFlags{
						FlagDeleted:   false,
						FlagDirectory: false,
					},
				},
				StandardInformationAttributes: mft.StandardInformationAttribute{
					Created:      time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					Modified:     time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					Accessed:     time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					Changed:      time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FlagResident: true,
				},
				FileNameAttributes: mft.FileNameAttributes{
					0: mft.FileNameAttribute{
						Created:      time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
						Modified:     time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
						Accessed:     time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
						Changed:      time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
						FlagResident: true,
						NameLength: mft.NameLength{
							FlagNamed: false,
							NamedSize: 0,
						},
						AttributeSize:           104,
						ParentDirRecordNumber:   5,
						ParentDirSequenceNumber: 5,
						LogicalFileSize:         16384,
						PhysicalFileSize:        16384,
						FileNameFlags: mft.FileNameFlags{
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
				DataAttribute: mft.DataAttribute{
					TotalSize:             0,
					FlagResident:          false,
					ResidentDataAttribute: nil,
					NonResidentDataAttribute: mft.NonResidentDataAttribute{
						DataRuns: mft.DataRuns{
							0: mft.DataRun{
								AbsoluteOffset: 4096,
								Length:         32768,
							},
						},
					},
				},
				AttributeList: mft.AttributeListAttributes{},
			},
			args: args{dummyHandler: &dummyHandler{
				handle:       nil,
				volumeLetter: "",
				vbr:          vbr.VolumeBootRecord{},
				reader:       nil,
				lastOffset:   0,
				filePath:     filepath.FromSlash("../../test/testdata/dummyntfs"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.dummyHandler.GetHandle()
			if err != nil {
				t.Errorf("failed to open dummyHandler file %s: %v", tt.args.dummyHandler.filePath, err)
				return
			}
			defer tt.args.dummyHandler.handle.Close()
			var gotMftRecord0 mft.MasterFileTableRecord
			gotMftRecord0, err = parseMFTRecord0(tt.args.dummyHandler)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMFTRecord0() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMftRecord0, tt.wantMftRecord0) {
				t.Errorf(cmp.Diff(gotMftRecord0, tt.wantMftRecord0))
			}
		})
	}
}
