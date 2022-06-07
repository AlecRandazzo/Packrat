// Copyright (c) 2022 Alec Randazzo

package mft

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_getDataAttribute(t *testing.T) {
	type args struct {
		input           []byte
		bytesPerCluster uint
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "nonresident data attribute test 1",
			args: args{
				input:           []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
				bytesPerCluster: 4096,
			},
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
			wantErr: false,
		},
		{
			name: "null input data attribute",
			args: args{
				input:           nil,
				bytesPerCluster: 4096,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "resident data attribute test 1",
			wantErr: false,
			args: args{
				input:           []byte{128, 0, 0, 0, 136, 0, 0, 0, 0, 0, 24, 0, 0, 0, 1, 0, 106, 0, 0, 0, 24, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 172, 3, 0, 0, 0, 0, 0, 0, 48, 238, 136, 38, 104, 47, 213, 1, 39, 0, 0, 0, 67, 0, 58, 0, 92, 0, 85, 0, 115, 0, 101, 0, 114, 0, 115, 0, 92, 0, 80, 0, 117, 0, 98, 0, 108, 0, 105, 0, 99, 0, 92, 0, 68, 0, 101, 0, 115, 0, 107, 0, 116, 0, 111, 0, 112, 0, 92, 0, 66, 0, 97, 0, 116, 0, 116, 0, 108, 0, 101, 0, 46, 0, 110, 0, 101, 0, 116, 0, 46, 0, 108, 0, 110, 0, 107, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				bytesPerCluster: 4096,
			},
			want: ResidentDataAttribute([]byte{2, 0, 0, 0, 0, 0, 0, 0, 172, 3, 0, 0, 0, 0, 0, 0, 48, 238, 136, 38, 104, 47, 213, 1, 39, 0, 0, 0, 67, 0, 58, 0, 92, 0, 85, 0, 115, 0, 101, 0, 114, 0, 115, 0, 92, 0, 80, 0, 117, 0, 98, 0, 108, 0, 105, 0, 99, 0, 92, 0, 68, 0, 101, 0, 115, 0, 107, 0, 116, 0, 111, 0, 112, 0, 92, 0, 66, 0, 97, 0, 116, 0, 116, 0, 108, 0, 101, 0, 46, 0, 110, 0, 101, 0, 116, 0, 46, 0, 108, 0, 110, 0, 107, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
		},
		{
			name: "zero input per cluster",
			args: args{
				input:           []byte{128, 0, 0, 0, 136, 0, 0, 0, 0, 0, 24, 0, 0, 0, 1, 0, 106, 0, 0, 0, 24, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 172, 3, 0, 0, 0, 0, 0, 0, 48, 238, 136, 38, 104, 47, 213, 1, 39, 0, 0, 0, 67, 0, 58, 0, 92, 0, 85, 0, 115, 0, 101, 0, 114, 0, 115, 0, 92, 0, 80, 0, 117, 0, 98, 0, 108, 0, 105, 0, 99, 0, 92, 0, 68, 0, 101, 0, 115, 0, 107, 0, 116, 0, 111, 0, 112, 0, 92, 0, 66, 0, 97, 0, 116, 0, 116, 0, 108, 0, 101, 0, 46, 0, 110, 0, 101, 0, 116, 0, 46, 0, 108, 0, 110, 0, 107, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				bytesPerCluster: 0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not enough resident data input",
			args: args{
				input:           []byte{128, 0, 0, 0, 136, 0, 0, 0, 0},
				bytesPerCluster: 4096,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not enough non resident data input",
			args: args{
				input:           []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0},
				bytesPerCluster: 4096,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDataAttribute(tt.args.input, tt.args.bytesPerCluster)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngoterr = %v, \nwanterr = %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_getDataRuns(t *testing.T) {
	type args struct {
		input           []byte
		bytesPerCluster uint
	}
	tests := []struct {
		name    string
		args    args
		want    DataRuns
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				input:           []byte{51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				bytesPerCluster: 4096,
			},
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
			name:    "null input",
			wantErr: true,
			want:    DataRuns{},
			args: args{
				input:           nil,
				bytesPerCluster: 4096,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDataRuns(tt.args.input, tt.args.bytesPerCluster)
			if !reflect.DeepEqual(got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf(cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_getNonResidentDataAttribute(t *testing.T) {
	type args struct {
		input           []byte
		bytesPerCluster uint
	}
	tests := []struct {
		name    string
		args    args
		want    NonResidentDataAttribute
		wantErr bool
	}{
		{
			name: "TestNonResidentDataAttribute_Parse test 1",
			args: args{
				input:           []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
				bytesPerCluster: 4096,
			},
			wantErr: false,
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
			name:    "null input in",
			wantErr: true,
			args: args{
				input:           nil,
				bytesPerCluster: 4096,
			},
		},
		{
			name:    "attribute offset longer than length",
			wantErr: true,
			args: args{
				input:           []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
				bytesPerCluster: 4096,
			},
		},
		{
			name:    "less than 18 input",
			wantErr: true,
			args: args{
				input:           []byte{128, 0, 0, 0},
				bytesPerCluster: 4096,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNonResidentDataAttribute(tt.args.input, tt.args.bytesPerCluster)
			if !reflect.DeepEqual(got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf(cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_getResidentDataAttribute(t *testing.T) {
	tests := []struct {
		name    string
		want    ResidentDataAttribute
		input   []byte
		wantErr bool
	}{
		{
			name:    "TestResidentDataAttribute_Parse test 1",
			input:   []byte{128, 0, 0, 0, 120, 0, 0, 0, 1, 0, 64, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0},
			want:    ResidentDataAttribute([]byte{63, 55, 5, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 0, 0, 116, 83, 0, 0, 0, 0, 51, 32, 200, 0, 0, 0, 12, 67, 109, 148, 1, 212, 133, 226, 1, 67, 54, 210, 0, 106, 250, 123, 9, 66, 253, 12, 241, 48, 8, 245, 66, 69, 99, 201, 78, 228, 8, 67, 97, 209, 0, 235, 81, 198, 1, 67, 218, 198, 0, 17, 228, 150, 1, 0, 0, 0}),
			wantErr: false,
		},
		{
			name:    "null input in",
			wantErr: true,
			want:    ResidentDataAttribute{},
			input:   nil,
		},
		{
			name:    "less than 18 input",
			wantErr: true,
			want:    ResidentDataAttribute{},
			input:   []byte{128, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getResidentDataAttribute(tt.input)
			if !reflect.DeepEqual(got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_splitDataRunByte(t *testing.T) {
	tests := []struct {
		name  string
		want  dataRunSplit
		input byte
	}{
		{
			name:  "Split 0x33",
			input: byte(0x33),
			want: dataRunSplit{
				offsetByteCount: 3,
				lengthByteCount: 3,
			},
		},
		{
			name:  "Split 0x04",
			input: byte(0x04),
			want: dataRunSplit{
				offsetByteCount: 0,
				lengthByteCount: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitDataRunByte(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, got, tt.want)
			}
		})
	}
}
