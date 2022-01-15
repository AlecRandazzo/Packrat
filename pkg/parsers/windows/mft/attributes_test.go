// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"reflect"
	"testing"
	"time"
)

func Test_validateAttribute(t *testing.T) {
	type args struct {
		attributeHeaderToCheck byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Header 0x10",
			args: args{
				attributeHeaderToCheck: 0x10,
			},
			wantErr: false,
		},
		{
			name: "Header 0x11",
			args: args{
				attributeHeaderToCheck: 0x11,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAttribute(tt.args.attributeHeaderToCheck); (err != nil) != tt.wantErr {
				t.Errorf("GetRawAttributes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getAttributes(t *testing.T) {
	type args struct {
		bytesPerCluster int64
	}
	tests := []struct {
		name                             string
		rawAttributes                    [][]byte
		args                             args
		wantFileNameAttributes           FileNameAttributes
		wantStandardInformationAttribute StandardInformationAttribute
		wantDataAttribute                DataAttribute
		wantAttributeListAttribute       AttributeListAttributes
		wantErr                          bool
	}{
		{
			name: "parse all attribute",
			rawAttributes: [][]byte{
				0: {16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				1: {48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0},
				2: {128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
				3: {176, 0, 0, 0, 72, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 42, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 176, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 49, 43, 103, 244, 2, 0, 0, 0},
			},
			args: args{bytesPerCluster: 4096},
			wantFileNameAttributes: FileNameAttributes{
				FileNameAttribute{
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
			wantStandardInformationAttribute: StandardInformationAttribute{
				Created:      time.Date(2016, 7, 2, 15, 13, 30, 670820200, time.UTC),
				Modified:     time.Date(2016, 7, 2, 15, 13, 30, 670820200, time.UTC),
				Accessed:     time.Date(2016, 7, 2, 15, 13, 30, 670820200, time.UTC),
				Changed:      time.Date(2016, 7, 2, 15, 13, 30, 670820200, time.UTC),
				FlagResident: true,
			},
			wantDataAttribute: DataAttribute{
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
							AbsoluteOffset: 132747444224,
							Length:         424071168,
						},
						2: DataRun{
							AbsoluteOffset: 784502874112,
							Length:         220422144,
						},
						3: DataRun{
							AbsoluteOffset: 30787432448,
							Length:         13619200,
						},
						4: DataRun{
							AbsoluteOffset: 641829142528,
							Length:         104091648,
						},
						5: DataRun{
							AbsoluteOffset: 763784736768,
							Length:         219549696,
						},
						6: DataRun{
							AbsoluteOffset: 873008676864,
							Length:         208510976,
						},
					},
				},
			},
			wantAttributeListAttribute: AttributeListAttributes{},
			wantErr:                    false,
		},
		{
			name: "nil input for an attribute",
			rawAttributes: [][]byte{
				0: nil,
				1: {48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0},
				2: {128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
				3: {176, 0, 0, 0, 72, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 42, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 176, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 49, 43, 103, 244, 2, 0, 0, 0},
			},
			args:                             args{bytesPerCluster: 4096},
			wantFileNameAttributes:           FileNameAttributes{},
			wantStandardInformationAttribute: StandardInformationAttribute{},
			wantDataAttribute:                DataAttribute{},
			wantAttributeListAttribute:       AttributeListAttributes{},
			wantErr:                          true,
		},
		{
			name:                             "nil rawAttributes slice",
			rawAttributes:                    nil,
			args:                             args{bytesPerCluster: 4096},
			wantFileNameAttributes:           FileNameAttributes{},
			wantStandardInformationAttribute: StandardInformationAttribute{},
			wantDataAttribute:                DataAttribute{},
			wantAttributeListAttribute:       AttributeListAttributes{},
			wantErr:                          true,
		},
		{
			name: "zero input per cluster input",
			rawAttributes: [][]byte{
				0: {16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				1: {48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0},
				2: {128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
				3: {176, 0, 0, 0, 72, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 42, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 176, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 49, 43, 103, 244, 2, 0, 0, 0},
			},
			args:                             args{bytesPerCluster: 0},
			wantFileNameAttributes:           FileNameAttributes{},
			wantStandardInformationAttribute: StandardInformationAttribute{},
			wantDataAttribute:                DataAttribute{},
			wantAttributeListAttribute:       AttributeListAttributes{},
			wantErr:                          true,
		},
		{
			name: "invalid filename attribute",
			rawAttributes: [][]byte{
				0: {48},
			},
			args:                             args{bytesPerCluster: 4096},
			wantFileNameAttributes:           FileNameAttributes{},
			wantStandardInformationAttribute: StandardInformationAttribute{},
			wantDataAttribute:                DataAttribute{},
			wantAttributeListAttribute:       AttributeListAttributes{},
			wantErr:                          false,
		},
		{
			name: "invalid standard info attribute",
			rawAttributes: [][]byte{
				0: {16},
			},
			args:                             args{bytesPerCluster: 4096},
			wantFileNameAttributes:           FileNameAttributes{},
			wantStandardInformationAttribute: StandardInformationAttribute{},
			wantDataAttribute:                DataAttribute{},
			wantAttributeListAttribute:       AttributeListAttributes{},
			wantErr:                          false,
		},
		{
			name: "invalid data attribute",
			rawAttributes: [][]byte{
				0: {128},
			},
			args:                             args{bytesPerCluster: 4096},
			wantFileNameAttributes:           FileNameAttributes{},
			wantStandardInformationAttribute: StandardInformationAttribute{},
			wantDataAttribute:                DataAttribute{},
			wantAttributeListAttribute:       AttributeListAttributes{},
			wantErr:                          true,
		},
		{
			name: "mft record with data list attribute",
			rawAttributes: [][]byte{
				0: {0x20, 0x00, 0x00, 0x00, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x80, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x87, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x44, 0x43, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x87, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0xB7, 0x15, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x45, 0xB7, 0x15, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
			args:                             args{bytesPerCluster: 4096},
			wantFileNameAttributes:           FileNameAttributes{},
			wantStandardInformationAttribute: StandardInformationAttribute{},
			wantDataAttribute:                DataAttribute{},
			wantAttributeListAttribute: AttributeListAttributes{
				AttributeListAttribute{
					Type:                     0x10,
					MFTReferenceRecordNumber: 493401,
				},
				AttributeListAttribute{
					Type:                     0x30,
					MFTReferenceRecordNumber: 493401,
				},
				AttributeListAttribute{
					Type:                     0x80,
					MFTReferenceRecordNumber: 1423172,
				},
				AttributeListAttribute{
					Type:                     0x80,
					MFTReferenceRecordNumber: 1423173,
				},
			},
			wantErr: false,
		},
		{
			name: "mft record with invalid data list attribute",
			rawAttributes: [][]byte{
				0: {0x20, 0x00, 0x00, 0x00, 0x98, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x80, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x87, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x44, 0x43, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x87, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0xB7, 0x15, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x1A, 0x00, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x45, 0xB7, 0x15, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
			args:                             args{bytesPerCluster: 4096},
			wantFileNameAttributes:           FileNameAttributes{},
			wantStandardInformationAttribute: StandardInformationAttribute{},
			wantDataAttribute:                DataAttribute{},
			wantAttributeListAttribute:       AttributeListAttributes{},
			wantErr:                          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileNameAttributes, gotStandardInformationAttribute, gotDataAttribute, gotAttributeListAttribute, err := GetAttributes(tt.rawAttributes, tt.args.bytesPerCluster)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFileNameAttributes, tt.wantFileNameAttributes) {
				t.Errorf("parse() \ngotFileNameAttributes = %v, \nwant %v", gotFileNameAttributes, tt.wantFileNameAttributes)
			}
			if !reflect.DeepEqual(gotStandardInformationAttribute, tt.wantStandardInformationAttribute) {
				t.Errorf("parse() \ngotStandardInformationAttribute = %v, \nwant %v", gotStandardInformationAttribute, tt.wantStandardInformationAttribute)
			}
			if !reflect.DeepEqual(gotDataAttribute, tt.wantDataAttribute) {
				t.Errorf("parse() \ngotDataAttribute = %v, \nwant %v", gotDataAttribute, tt.wantDataAttribute)
			}
			if !reflect.DeepEqual(gotAttributeListAttribute, tt.wantAttributeListAttribute) {
				t.Errorf("parse() \ngotDataAttribute = %v, \nwant %v", gotAttributeListAttribute, tt.wantAttributeListAttribute)
			}

		})
	}
}

func Test_getRawAttributes(t *testing.T) {
	type args struct {
		recordHeader RecordHeader
	}
	tests := []struct {
		name              string
		rawMftRecord      []byte
		args              args
		wantRawAttributes [][]byte
		wantErr           bool
	}{
		{
			name:         "test1",
			rawMftRecord: []byte{70, 73, 76, 69, 48, 0, 3, 0, 155, 21, 101, 188, 33, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 200, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 29, 7, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0, 176, 0, 0, 0, 72, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 42, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 176, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 49, 43, 103, 244, 2, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 49, 1, 255, 255, 11, 17, 1, 255, 0, 0, 0, 0, 0, 0, 29, 7, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 29, 7},
			args: args{recordHeader: RecordHeader{
				AttributesOffset: 0x38,
				RecordNumber:     0,
				Flags:            RecordHeaderFlags{},
			}},
			wantErr: false,
			wantRawAttributes: [][]byte{
				0: {16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				1: {48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0},
				2: {128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
				3: {176, 0, 0, 0, 72, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 42, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 176, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 49, 43, 103, 244, 2, 0, 0, 0},
			},
		},
		{
			name:         "nil input",
			rawMftRecord: nil,
			wantErr:      true,
		},
		{
			name:         "no record header input",
			rawMftRecord: []byte{70, 73, 76, 69, 48, 0, 3, 0, 155, 21, 101, 188, 33, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 200, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 29, 7, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 24, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 48, 0, 0, 0, 104, 0, 0, 0, 0, 0, 24, 0, 0, 0, 3, 0, 74, 0, 0, 0, 24, 0, 1, 0, 5, 0, 0, 0, 0, 0, 5, 0, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 234, 36, 205, 74, 116, 212, 209, 1, 0, 64, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 4, 3, 36, 0, 77, 0, 70, 0, 84, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0, 176, 0, 0, 0, 72, 0, 0, 0, 1, 0, 64, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 42, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 176, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 8, 160, 2, 0, 0, 0, 0, 0, 49, 43, 103, 244, 2, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 8, 16, 0, 0, 0, 0, 0, 0, 49, 1, 255, 255, 11, 17, 1, 255, 0, 0, 0, 0, 0, 0, 29, 7, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 29, 7},
			args: args{recordHeader: RecordHeader{
				AttributesOffset: 0,
				RecordNumber:     0,
				Flags:            RecordHeaderFlags{},
			}},
			wantErr: true,
		},
		{
			name:         "Break if the attribute offset is beyond the byte slice",
			rawMftRecord: []byte{70, 73, 76, 69, 48, 0, 3, 0, 155, 21, 101, 188, 33, 0, 0, 0, 1, 0, 1, 0, 56, 0, 1, 0, 200, 1, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 29, 7, 0, 0, 0, 0, 0, 0, 255, 0, 0, 0, 96},
			args: args{recordHeader: RecordHeader{
				AttributesOffset: 255,
				RecordNumber:     0,
				Flags:            RecordHeaderFlags{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRawAttributes, err := GetRawAttributes(tt.rawMftRecord, tt.args.recordHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRawAttributes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRawAttributes, tt.wantRawAttributes) {
				t.Errorf("GetRawAttributes() gotRawAttributes = %v, want %v", gotRawAttributes, tt.wantRawAttributes)
			}
		})
	}
}
