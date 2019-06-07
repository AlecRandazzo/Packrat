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
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMasterFileTableRecord_getDataAttribute(t *testing.T) {
	tests := []struct {
		name          string
		mftRecord     *MasterFileTableRecord
		wantMftRecord *MasterFileTableRecord
	}{
		{
			name: "Testing mft record 0.",
			mftRecord: &MasterFileTableRecord{
				AttributeInfo: []AttributeInfo{
					{
						AttributeType:  128,
						AttributeBytes: []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
					},
				},
			},
			wantMftRecord: &MasterFileTableRecord{
				DataAttributes: DataAttributes{
					FlagResident: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mftRecord.getDataAttribute()
			assert.Equal(t, tt.wantMftRecord.DataAttributes.FlagResident, tt.mftRecord.DataAttributes.FlagResident, "Flag resident should be equal.")
		})
	}
}

func Test_getResidentDataAttribute(t *testing.T) {
	type args struct {
		attributeBytes []byte
	}
	tests := []struct {
		name                       string
		args                       args
		wantResidentDataAttributes ResidentDataAttributes
		wantErr                    bool
	}{
		{
			name: "Test mft record 0.",
			args: args{
				attributeBytes: []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
			},
			wantResidentDataAttributes: ResidentDataAttributes{
				ResidentData: []byte{63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResidentDataAttributes, err := getResidentDataAttribute(tt.args.attributeBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("getResidentDataAttribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResidentDataAttributes, tt.wantResidentDataAttributes) {
				t.Errorf("getResidentDataAttribute() = %v, want %v", gotResidentDataAttributes, tt.wantResidentDataAttributes)
			}
		})
	}
}

func Test_getNonResidentDataAttribute(t *testing.T) {
	type args struct {
		attributeBytes  []byte
		bytesPerCluster int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test mft record 0.",
			args: args{
				attributeBytes: []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getNonResidentDataAttribute(tt.args.attributeBytes, tt.args.bytesPerCluster)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNonResidentDataAttribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_getDataRuns(t *testing.T) {
	type args struct {
		dataRunBytes    []byte
		bytesPerCluster int64
	}
	tests := []struct {
		name         string
		args         args
		wantDataRuns map[int]DataRun
	}{
		{
			name: "Test with MFT record 0.",
			args: args{
				dataRunBytes:    []byte{51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
				bytesPerCluster: 4096,
			},
			wantDataRuns: map[int]DataRun{
				0: {
					AbsoluteOffset: 3221225472,
					Length:         209846272,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDataRuns := getDataRuns(tt.args.dataRunBytes, tt.args.bytesPerCluster)
			assert.Equal(t, tt.wantDataRuns[0], gotDataRuns[0], "This datarun should be equal.")
		})
	}
}

func Test_getDataRunSplit(t *testing.T) {
	type args struct {
		dataRunByte byte
	}
	tests := []struct {
		name                string
		args                args
		wantOffsetByteCount int
		wantLengthByteCount int
	}{
		{
			name: "Testing with MFT record 0",
			args: args{
				dataRunByte: 0x43,
			},
			wantOffsetByteCount: 4,
			wantLengthByteCount: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOffsetByteCount, gotLengthByteCount := getDataRunSplit(tt.args.dataRunByte)
			if gotOffsetByteCount != tt.wantOffsetByteCount {
				t.Errorf("getDataRunSplit() gotOffsetByteCount = %v, want %v", gotOffsetByteCount, tt.wantOffsetByteCount)
			}
			if gotLengthByteCount != tt.wantLengthByteCount {
				t.Errorf("getDataRunSplit() gotLengthByteCount = %v, want %v", gotLengthByteCount, tt.wantLengthByteCount)
			}
		})
	}
}
