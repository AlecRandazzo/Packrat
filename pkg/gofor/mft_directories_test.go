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
	"sync"
	"testing"
)

func TestMasterFileTableRecord_quickDirectoryCheck(t *testing.T) {
	mftBytes, _ := hex.DecodeString("46494C45300003007D7A5185210000000500010038000300800300000004000000000000000000000A00000005000000E207650000000000100000006000000000001800000000004800000018000000E4B9B331DC12D101F6771A76D414D501F6771A76D414D501F6771A76D414D501060000000000000000000000000000000000000009010000000000000000000088C17A71030000003000000060000000000018000000010044000000180001000500000000000500EA24CD4A74D4D101EA24CD4A74D4D101EA24CD4A74D4D101EA24CD4A74D4D10100000000000000000000000000000000060000100000000001032E000000000090000000A00100000004180000000600800100002000000024004900330030003000000001000000001000000100000010000000700100007001000001000000435C01000000030070005200010000000500000000000500F29B8C6971D4D101F29B8C6971D4D101F29B8C6971D4D101F29B8C6971D4D1010000000000000000000000000000000006240010030000A0080244004F00430055004D0045007E0031002E007300790000000000000000008F7C01000000020070005200010000000500000000000500A8DDEB4306B8D301A8DDEB4306B8D301A8DDEB4306B8D301A8DDEB4306B8D301000000000000000000000000000000000000001000000000080369004400E207660065006E00730065000000000001000100000000000000EF0200000000020068004E00010000000500000000000500805D5D3A03D4D101805D5D3A03D4D101805D5D3A03D4D101805D5D3A03D4D10100000000000000000000000000000000000000100000000006034E005600490044004900410076000300000000000000000000000000000018000000030000000200000000000000A0000000500000000104400000000800000000000000000003000000000000004800000000000000004000000000000000400000000000000040000000000000240049003300300031048DB200000000B0000000280000000004180000000700080000002000000024004900330030000F0000000000000000010000680000000009180000000900380000003000000024005400580046005F004400410054004100000000000000050000000000050001000000010000000000000000000000035253004D000000014453004D000000001472006200000002000000371E0000FFFFFFFF0000000002000000371E0000FFFFFFFF00000000FFFFFFFF00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000E207")

	tests := []struct {
		name              string
		mftRecord         *MasterFileTableRecord
		wantFlagDirectory bool
	}{
		{
			name: "Testing MFT record 5",
			mftRecord: &MasterFileTableRecord{
				MftRecordBytes: mftBytes,
			},
			wantFlagDirectory: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mftRecord.quickDirectoryCheck()
			assert.Equal(t, tt.wantFlagDirectory, tt.mftRecord.RecordHeader.FlagDirectory)
		})
	}
}

func Test_createDirectoryList(t *testing.T) {
	type args struct {
		inboundBuffer        *chan []byte
		directoryListChannel *chan map[uint64]Directory
		waitGroup            *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createDirectoryList(tt.args.inboundBuffer, tt.args.directoryListChannel, tt.args.waitGroup)
		})
	}
}

func TestMftFile_combineDirectoryInformation(t *testing.T) {
	type args struct {
		directoryListChannel        *chan map[uint64]Directory
		waitForDirectoryCombination *sync.WaitGroup
	}
	tests := []struct {
		name string
		file *MftFile
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.file.combineDirectoryInformation(tt.args.directoryListChannel, tt.args.waitForDirectoryCombination)
		})
	}
}

func TestVolumeHandle_combineDirectoryInformation(t *testing.T) {
	type args struct {
		directoryListChannel        *chan map[uint64]Directory
		waitForDirectoryCombination *sync.WaitGroup
	}
	tests := []struct {
		name   string
		volume *VolumeHandle
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.volume.combineDirectoryInformation(tt.args.directoryListChannel, tt.args.waitForDirectoryCombination)
		})
	}
}

func TestMftFile_BuildDirectoryTree(t *testing.T) {
	tests := []struct {
		name    string
		file    *MftFile
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.file.buildDirectoryTree(); (err != nil) != tt.wantErr {
				t.Errorf("MftFile.buildDirectoryTree() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//func TestVolumeHandle_BuildDirectoryTree(t *testing.T) {
//	tests := []struct {
//		name    string
//		volume  *VolumeHandle
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := tt.volume.buildDirectoryTree(); (err != nil) != tt.wantErr {
//				t.Errorf("VolumeHandle.buildDirectoryTree() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
