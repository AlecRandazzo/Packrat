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
	"testing"
)

func TestConvertLittleEndianByteSliceToInt64(t *testing.T) {
	type args struct {
		inBytes []byte
	}
	tests := []struct {
		name         string
		args         args
		wantOutInt64 int64
	}{
		{
			name:         "Testing with a random byte slice.",
			args:         args{inBytes: []byte{20, 21, 255}},
			wantOutInt64: -60140,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOutInt64 := convertLittleEndianByteSliceToInt64(tt.args.inBytes); gotOutInt64 != tt.wantOutInt64 {
				t.Errorf("convertLittleEndianByteSliceToInt64() = %v, want %v", gotOutInt64, tt.wantOutInt64)
			}
		})
	}
}

func TestConvertLittleEndianByteSliceToUInt64(t *testing.T) {
	type args struct {
		inBytes []byte
	}
	tests := []struct {
		name          string
		args          args
		wantOutUint64 uint64
	}{
		{
			name:          "Testing with a random byte slice.",
			args:          args{inBytes: []byte{20, 21, 255}},
			wantOutUint64: 16717076,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOutUint64 := convertLittleEndianByteSliceToUInt64(tt.args.inBytes); gotOutUint64 != tt.wantOutUint64 {
				t.Errorf("convertLittleEndianByteSliceToUInt64() = %v, want %v", gotOutUint64, tt.wantOutUint64)
			}
		})
	}
}
