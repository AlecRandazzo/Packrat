/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package gofor

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMasterFileTableRecord_getFileNameAttributes(t *testing.T) {
	attributeBytes, _ := hex.DecodeString("300000006800000000001800000003004A000000180001000500000000000500EA24CD4A74D4D101EA24CD4A74D4D101EA24CD4A74D4D101EA24CD4A74D4D101004000000000000000400000000000000600000000000000040324004D0046005400000000000000")

	tests := []struct {
		name          string
		mftRecord     *MasterFileTableRecord
		wantMftRecord *MasterFileTableRecord
	}{
		{
			name: "Testing MFT Record 0",
			mftRecord: &MasterFileTableRecord{
				AttributeInfo: []AttributeInfo{
					{
						AttributeBytes: attributeBytes,
						AttributeType:  0x30,
					},
				},
			},
			wantMftRecord: &MasterFileTableRecord{
				FileNameAttributes: []FileNameAttributes{
					{
						FnCreated:               "2016-07-02T15:13:30Z",
						FnModified:              "2016-07-02T15:13:30Z",
						FnAccessed:              "2016-07-02T15:13:30Z",
						FnChanged:               "2016-07-02T15:13:30Z",
						FlagResident:            true,
						FlagNamed:               false,
						NamedSize:               0,
						AttributeSize:           104,
						ParentDirRecordNumber:   5,
						ParentDirSequenceNumber: 0,
						LogicalFileSize:         16384,
						PhysicalFileSize:        16384,
						FileName:                "$MFT",
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
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mftRecord.getFileNameAttributes()
			assert.Equal(t, tt.wantMftRecord.FileNameAttributes, tt.mftRecord.FileNameAttributes, "Filename attributes should be equal.")
		})
	}
}

func Test_resolveFileFlags(t *testing.T) {
	fnFlagBytes, _ := hex.DecodeString("0600000000")

	type args struct {
		flagBytes []byte
	}
	tests := []struct {
		name            string
		args            args
		wantParsedFlags FileNameFlags
	}{
		{
			name: "Testing MFT record 0",
			args: args{
				flagBytes: fnFlagBytes,
			},
			wantParsedFlags: FileNameFlags{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotParsedFlags := resolveFileFlags(tt.args.flagBytes); !reflect.DeepEqual(gotParsedFlags, tt.wantParsedFlags) {
				t.Errorf("resolveFileFlags() = %v, want %v", gotParsedFlags, tt.wantParsedFlags)
			}
		})
	}
}
