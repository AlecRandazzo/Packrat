// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"reflect"
	"testing"
)

func TestRawDataAttribute_Parse(t *testing.T) {
	type args struct {
		bytesPerCluster int64
	}
	tests := []struct {
		name             string
		gotResident      ResidentDataAttribute
		gotNonResident   NonResidentDataAttribute
		args             args
		wantResident     ResidentDataAttribute
		wantNonResident  NonResidentDataAttribute
		wantErr          bool
		rawDataAttribute RawDataAttribute
	}{
		{
			name: "nonresident data attribute test 1",
			args: args{
				bytesPerCluster: 4096,
			},
			rawDataAttribute: []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
			wantNonResident: NonResidentDataAttribute{
				DataRuns: DataRuns{
					0: {
						AbsoluteOffset: 3221225472,
						Length:         209846272,
					},
					1: {
						AbsoluteOffset: 132747444224,
						Length:         424071168,
					},
					2: {
						AbsoluteOffset: 784502874112,
						Length:         220422144,
					},
					3: {
						AbsoluteOffset: 30787432448,
						Length:         13619200,
					},
					4: {
						AbsoluteOffset: 641829142528,
						Length:         104091648,
					},
					5: {
						AbsoluteOffset: 763784736768,
						Length:         219549696,
					},
					6: {
						AbsoluteOffset: 873008676864,
						Length:         208510976,
					},
				},
			},
			wantResident: nil,
			wantErr:      false,
		},
		{
			name: "null bytes data attribute",
			args: args{
				bytesPerCluster: 4096,
			},
			rawDataAttribute: nil,
			wantErr:          true,
		},
		{
			name:    "resident data attribute test 1",
			wantErr: false,
			args: args{
				bytesPerCluster: 4096,
			},
			rawDataAttribute: []byte{128, 0, 0, 0, 136, 0, 0, 0, 0, 0, 24, 0, 0, 0, 1, 0, 106, 0, 0, 0, 24, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 172, 3, 0, 0, 0, 0, 0, 0, 48, 238, 136, 38, 104, 47, 213, 1, 39, 0, 0, 0, 67, 0, 58, 0, 92, 0, 85, 0, 115, 0, 101, 0, 114, 0, 115, 0, 92, 0, 80, 0, 117, 0, 98, 0, 108, 0, 105, 0, 99, 0, 92, 0, 68, 0, 101, 0, 115, 0, 107, 0, 116, 0, 111, 0, 112, 0, 92, 0, 66, 0, 97, 0, 116, 0, 116, 0, 108, 0, 101, 0, 46, 0, 110, 0, 101, 0, 116, 0, 46, 0, 108, 0, 110, 0, 107, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantResident:     ResidentDataAttribute([]byte{2, 0, 0, 0, 0, 0, 0, 0, 172, 3, 0, 0, 0, 0, 0, 0, 48, 238, 136, 38, 104, 47, 213, 1, 39, 0, 0, 0, 67, 0, 58, 0, 92, 0, 85, 0, 115, 0, 101, 0, 114, 0, 115, 0, 92, 0, 80, 0, 117, 0, 98, 0, 108, 0, 105, 0, 99, 0, 92, 0, 68, 0, 101, 0, 115, 0, 107, 0, 116, 0, 111, 0, 112, 0, 92, 0, 66, 0, 97, 0, 116, 0, 116, 0, 108, 0, 101, 0, 46, 0, 110, 0, 101, 0, 116, 0, 46, 0, 108, 0, 110, 0, 107, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			wantNonResident:  NonResidentDataAttribute{},
		},
		{
			name: "zero bytes per cluster",
			args: args{
				bytesPerCluster: 0,
			},
			rawDataAttribute: []byte{128, 0, 0, 0, 136, 0, 0, 0, 0, 0, 24, 0, 0, 0, 1, 0, 106, 0, 0, 0, 24, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 172, 3, 0, 0, 0, 0, 0, 0, 48, 238, 136, 38, 104, 47, 213, 1, 39, 0, 0, 0, 67, 0, 58, 0, 92, 0, 85, 0, 115, 0, 101, 0, 114, 0, 115, 0, 92, 0, 80, 0, 117, 0, 98, 0, 108, 0, 105, 0, 99, 0, 92, 0, 68, 0, 101, 0, 115, 0, 107, 0, 116, 0, 111, 0, 112, 0, 92, 0, 66, 0, 97, 0, 116, 0, 116, 0, 108, 0, 101, 0, 46, 0, 110, 0, 101, 0, 116, 0, 46, 0, 108, 0, 110, 0, 107, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:          true,
		},
		{
			name: "not enough resident data bytes",
			args: args{
				bytesPerCluster: 4096,
			},
			rawDataAttribute: []byte{128, 0, 0, 0, 136, 0, 0, 0, 0},
			wantErr:          true,
		},
		{
			name: "not enough non resident data bytes",
			args: args{
				bytesPerCluster: 4096,
			},
			rawDataAttribute: []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0},
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.gotNonResident, tt.gotResident, err = tt.rawDataAttribute.Parse(tt.args.bytesPerCluster)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngoterr = %v, \nwanterr = %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.gotNonResident, tt.wantNonResident) {
				t.Errorf("Test %v failed \ngotNonResident = %v, \nwantNonResident = %v", tt.name, tt.gotNonResident, tt.wantNonResident)
			}
			if !reflect.DeepEqual(tt.gotResident, tt.wantResident) {
				t.Errorf("Test %v failed \ngotResident = %v, \nwantResident = %v", tt.name, tt.gotResident, tt.wantResident)
			}
		})
	}
}

func TestRawDataRuns_Parse(t *testing.T) {
	type args struct {
		bytesPerCluster int64
	}
	tests := []struct {
		name        string
		got         DataRuns
		args        args
		want        DataRuns
		wantErr     bool
		rawDataRuns RawDataRuns
	}{
		{
			name: "TestDataRuns_Parse test 1",
			args: args{
				bytesPerCluster: 4096,
			},
			rawDataRuns: []byte{51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want: DataRuns{
				0: {
					AbsoluteOffset: 3221225472,
					Length:         209846272,
				},
				1: {
					AbsoluteOffset: 132747444224,
					Length:         424071168,
				},
				2: {
					AbsoluteOffset: 784502874112,
					Length:         220422144,
				},
				3: {
					AbsoluteOffset: 30787432448,
					Length:         13619200,
				},
				4: {
					AbsoluteOffset: 641829142528,
					Length:         104091648,
				},
				5: {
					AbsoluteOffset: 763784736768,
					Length:         219549696,
				},
				6: {
					AbsoluteOffset: 873008676864,
					Length:         208510976,
				},
			},
			wantErr: false,
		},
		{
			name:        "null bytes",
			wantErr:     true,
			rawDataRuns: nil,
			args: args{
				bytesPerCluster: 4096,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.got = make(map[int]DataRun)
			tt.got, err = tt.rawDataRuns.Parse(tt.args.bytesPerCluster)
			if !reflect.DeepEqual(tt.got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestRawNonResidentDataAttribute_Parse(t *testing.T) {
	type args struct {
		bytesPerCluster int64
	}
	tests := []struct {
		name                        string
		want                        NonResidentDataAttribute
		args                        args
		got                         NonResidentDataAttribute
		wantErr                     bool
		rawNonResidentDataAttribute RawNonResidentDataAttribute
	}{
		{
			name: "TestNonResidentDataAttribute_Parse test 1",
			args: args{
				bytesPerCluster: 4096,
			},
			rawNonResidentDataAttribute: RawNonResidentDataAttribute([]byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0}),
			wantErr:                     false,
			want: NonResidentDataAttribute{
				DataRuns: DataRuns{
					0: {
						AbsoluteOffset: 3221225472,
						Length:         209846272,
					},
					1: {
						AbsoluteOffset: 132747444224,
						Length:         424071168,
					},
					2: {
						AbsoluteOffset: 784502874112,
						Length:         220422144,
					},
					3: {
						AbsoluteOffset: 30787432448,
						Length:         13619200,
					},
					4: {
						AbsoluteOffset: 641829142528,
						Length:         104091648,
					},
					5: {
						AbsoluteOffset: 763784736768,
						Length:         219549696,
					},
					6: {
						AbsoluteOffset: 873008676864,
						Length:         208510976,
					},
				},
			},
		},
		{
			name:    "null bytes in",
			wantErr: true,
			args: args{
				bytesPerCluster: 4096,
			},
			rawNonResidentDataAttribute: nil,
		},
		{
			name:                        "attribute offset longer than length",
			wantErr:                     true,
			rawNonResidentDataAttribute: []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
		},
		{
			name:                        "less than 18 bytes",
			wantErr:                     true,
			rawNonResidentDataAttribute: []byte{128, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.got, err = tt.rawNonResidentDataAttribute.Parse(tt.args.bytesPerCluster)
			if !reflect.DeepEqual(tt.got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestRawResidentDataAttribute_Parse(t *testing.T) {
	tests := []struct {
		name                     string
		rawResidentDataAttribute RawResidentDataAttribute
		want                     ResidentDataAttribute
		got                      ResidentDataAttribute
		wantErr                  bool
	}{
		{
			name:                     "TestResidentDataAttribute_Parse test 1",
			rawResidentDataAttribute: RawResidentDataAttribute([]byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0}),
			want:                     ResidentDataAttribute([]byte{63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0}),
			wantErr:                  false,
		},
		{
			name:                     "null bytes in",
			wantErr:                  true,
			rawResidentDataAttribute: nil,
		},
		{
			name:                     "less than 18 bytes",
			wantErr:                  true,
			rawResidentDataAttribute: RawResidentDataAttribute([]byte{128, 0, 0, 0}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.got, err = tt.rawResidentDataAttribute.Parse()
			if !reflect.DeepEqual(tt.got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func Test_RawDataRunSplitByte_Parse(t *testing.T) {
	tests := []struct {
		name                string
		got                 dataRunSplit
		want                dataRunSplit
		rawDataRunSplitByte rawDataRunSplitByte
	}{
		{
			name:                "Split 0x33",
			rawDataRunSplitByte: rawDataRunSplitByte(byte(0x33)),
			want: dataRunSplit{
				offsetByteCount: 3,
				lengthByteCount: 3,
			},
		},
		{
			name:                "Split 0x04",
			rawDataRunSplitByte: rawDataRunSplitByte(byte(0x04)),
			want: dataRunSplit{
				offsetByteCount: 0,
				lengthByteCount: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.got = tt.rawDataRunSplitByte.parse()
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, tt.got, tt.want)
			}
		})
	}
}
