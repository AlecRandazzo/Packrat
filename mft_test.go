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
	mft "github.com/Go-Forensics/MFT-Parser"
	vbr "github.com/Go-Forensics/VBR-Parser"
	"log"
	"reflect"
	"testing"
	"time"
)

func Test_parseMFTRecord0(t *testing.T) {
	type args struct {
		volume *VolumeHandler
	}
	tests := []struct {
		name           string
		args           args
		wantMftRecord0 mft.MasterFileTableRecord
		wantErr        bool
		dummyFile      string
	}{
		{
			name:      "test1",
			dummyFile: `test\testdata\dummyntfs`,
			wantErr:   false,
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
					SiCreated:    time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiModified:   time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiAccessed:   time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiChanged:    time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FlagResident: true,
				},
				FileNameAttributes: mft.FileNameAttributes{
					0: mft.FileNameAttribute{
						FnCreated:    time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
						FnModified:   time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
						FnAccessed:   time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
						FnChanged:    time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
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
				AttributeList: nil,
			},
			args: args{volume: &VolumeHandler{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handle := dummyHandler{
				Handle:               nil,
				VolumeLetter:         "",
				Vbr:                  vbr.VolumeBootRecord{},
				mftReader:            nil,
				lastReadVolumeOffset: 0,
				filePath:             tt.dummyFile,
			}
			var err error
			*tt.args.volume, err = GetVolumeHandler("c", handle)
			if err != nil {
				log.Panic(err)
			}
			defer tt.args.volume.Handle.Close()
			gotMftRecord0, err := parseMFTRecord0(tt.args.volume)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMFTRecord0() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMftRecord0, tt.wantMftRecord0) {
				t.Errorf("parseMFTRecord0() gotMftRecord0 = %v, want %v", gotMftRecord0, tt.wantMftRecord0)
			}
		})
	}
}
