/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package GoFor

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestParseVolumeBootRecord(t *testing.T) {
	vbrBytes, _ := hex.DecodeString("EB52904E5446532020202000020800000000000000F800003F00FF0000A80F000000000080008000666C1D770000000000000C00000000000200000000000000F600000001000000B8DF5ECA0E5FCA6000000000FA33C08ED0BC007CFB68C0071F1E686600CB88160E0066813E03004E5446537515B441BBAA55CD13720C81FB55AA7506F7C101007503E9DD001E83EC18681A00B4488A160E008BF4161FCD139F83C4189E581F72E13B060B0075DBA30F00C12E0F00041E5A33DBB900202BC866FF06110003160F008EC2FF061600E84B002BC877EFB800BBCD1A6623C0752D6681FB54435041752481F90201721E166807BB166852111668090066536653665516161668B80166610E07CD1A33C0BF0A13B9F60CFCF3AAE9FE01909066601E0666A111006603061C001E66680000000066500653680100681000B4428A160E00161F8BF4CD1366595B5A665966591F0F82160066FF06110003160F008EC2FF0E160075BC071F6661C3A1F601E80900A1FA01E80300F4EBFD8BF0AC3C007409B40EBB0700CD10EBF2C30D0A41206469736B2072656164206572726F72206F63637572726564000D0A424F4F544D475220697320636F6D70726573736564000D0A5072657373204374726C2B416C742B44656C20746F20726573746172740D0A000000000000000000000000000000000000000000008A01A701BF01000055AA")

	type args struct {
		volumeBootRecordBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantVbr VolumeBootRecord
		wantErr bool
	}{
		{
			name: "Testing VBR",
			args: args{
				volumeBootRecordBytes: vbrBytes,
			},
			wantVbr: VolumeBootRecord{
				BytesPerSector:         512,
				BytesPerCluster:        4096,
				SectorsPerCluster:      8,
				MftRecordSize:          1024,
				MftByteOffset:          3221225472,
				ClustersPerIndexRecord: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVbr, err := ParseVolumeBootRecord(tt.args.volumeBootRecordBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVolumeBootRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVbr, tt.wantVbr) {
				t.Errorf("ParseVolumeBootRecord() = %v, want %v", gotVbr, tt.wantVbr)
			}
		})
	}
}
