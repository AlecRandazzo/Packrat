// Copyright (c) 2020 Alec Randazzo

package byteshelper

import (
	"fmt"
	"reflect"
	"testing"
)

func TestLittleEndianBinaryToInt64(t *testing.T) {
	type args struct {
		inBytes []byte
	}
	tests := []struct {
		name         string
		args         args
		wantOutInt64 int64
		wantErr      bool
	}{
		{
			name:         "Testing with a random byte slice.",
			args:         args{inBytes: []byte{20, 21, 255}},
			wantOutInt64: -60140,
			wantErr:      false,
		},
		{
			name:         "Testing with null byte slice.",
			args:         args{inBytes: []byte{}},
			wantOutInt64: 0,
			wantErr:      true,
		},
		{
			name:         "Testing with byte slice larger than 8 bytes.",
			args:         args{inBytes: []byte{20, 21, 22, 23, 24, 25, 26, 27, 28, 29}},
			wantOutInt64: 0,
			wantErr:      true,
		},
		{
			name:         "Testing with byte slice with 8 bytes.",
			args:         args{inBytes: []byte{1}},
			wantOutInt64: 1,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutInt64, err := LittleEndianBinaryToInt64(tt.args.inBytes)
			fmt.Println(gotOutInt64, tt.wantOutInt64)
			if gotOutInt64 != tt.wantOutInt64 || (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want %v", gotOutInt64, tt.wantOutInt64)
			}

		})
	}
}

func TestLittleEndianBinaryToUInt64(t *testing.T) {
	type args struct {
		inBytes []byte
	}
	tests := []struct {
		name          string
		args          args
		wantOutUInt64 uint64
		wantErr       bool
	}{
		{
			name:          "Testing with a random byte slice.",
			args:          args{inBytes: []byte{20, 21, 255}},
			wantOutUInt64: 16717076,
			wantErr:       false,
		},
		{
			name:          "Testing with null byte slice.",
			args:          args{inBytes: []byte{}},
			wantOutUInt64: 0,
			wantErr:       true,
		},
		{
			name:          "Testing with byte slice larger than 8 bytes.",
			args:          args{inBytes: []byte{20, 20, 20, 20, 20, 20, 20, 20, 20, 20}},
			wantOutUInt64: 0,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOutUInt64, err := LittleEndianBinaryToUInt64(tt.args.inBytes); !reflect.DeepEqual(gotOutUInt64, tt.wantOutUInt64) || (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want %v", gotOutUInt64, tt.wantOutUInt64)
			}
		})
	}
}

func TestUnicodeBytesToASCII(t *testing.T) {
	type args struct {
		unicodeBytes []byte
	}
	tests := []struct {
		name            string
		args            args
		wantASCIIString string
		wantErr         bool
	}{
		{
			name:            "Testing with a random byte slice.",
			args:            args{unicodeBytes: []byte{116, 0, 101, 0, 115, 0, 116}},
			wantASCIIString: "test",
			wantErr:         false,
		},
		{
			name:            "Testing with null byte slice.",
			args:            args{unicodeBytes: []byte{}},
			wantASCIIString: "",
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotASCIIString, err := UnicodeBytesToASCII(tt.args.unicodeBytes); !reflect.DeepEqual(gotASCIIString, tt.wantASCIIString) || (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want %v", gotASCIIString, tt.wantASCIIString)
			}
		})
	}
}

func TestLittleEndianBinaryToInt16(t *testing.T) {
	type args struct {
		inBytes []byte
	}
	tests := []struct {
		name         string
		args         args
		wantOutInt16 int16
		wantErr      bool
	}{
		{
			name:         "Testing with a random byte slice.",
			args:         args{inBytes: []byte{20}},
			wantOutInt16: 20,
			wantErr:      false,
		},
		{
			name:         "Testing with null byte slice.",
			args:         args{inBytes: []byte{}},
			wantOutInt16: 0,
			wantErr:      true,
		},
		{
			name:         "Testing with byte slice larger than 2 bytes.",
			args:         args{inBytes: []byte{20, 21, 22, 23, 24, 25, 26, 27, 28, 29}},
			wantOutInt16: 0,
			wantErr:      true,
		},
		{
			name:         "Testing with byte slice with 2 bytes.",
			args:         args{inBytes: []byte{49, 255}},
			wantOutInt16: -207,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutInt16, err := LittleEndianBinaryToInt16(tt.args.inBytes)
			if gotOutInt16 != tt.wantOutInt16 || (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want %v", gotOutInt16, tt.wantOutInt16)
			}

		})
	}
}

func TestLittleEndianBinaryToInt32(t *testing.T) {
	type args struct {
		inBytes []byte
	}
	tests := []struct {
		name         string
		args         args
		wantOutInt32 int32
		wantErr      bool
	}{
		{
			name:         "Testing with a random byte slice.",
			args:         args{inBytes: []byte{20, 21, 255}},
			wantOutInt32: -60140,
			wantErr:      false,
		},
		{
			name:         "Testing with null byte slice.",
			args:         args{inBytes: []byte{}},
			wantOutInt32: 0,
			wantErr:      true,
		},
		{
			name:         "Testing with byte slice larger than 4 bytes.",
			args:         args{inBytes: []byte{20, 21, 22, 23, 24, 25, 26, 27, 28, 29}},
			wantOutInt32: 0,
			wantErr:      true,
		},
		{
			name:         "Testing with byte slice with 4 bytes.",
			args:         args{inBytes: []byte{1, 2, 3, 4}},
			wantOutInt32: 67305985,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutInt32, err := LittleEndianBinaryToInt32(tt.args.inBytes)
			if gotOutInt32 != tt.wantOutInt32 || (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want %v", gotOutInt32, tt.wantOutInt32)
			}

		})
	}
}
