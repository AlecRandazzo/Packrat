/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package GoFor_Collector

import (
	mft "github.com/AlecRandazzo/GoFor-MFT-Parser"
	"sync"
	"testing"
)

func TestVolumeHandle_CombineDirectoryInformation(t *testing.T) {
	type args struct {
		directoryListChannel        *chan map[uint64]mft.Directory
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
			tt.volume.CombineDirectoryInformation(tt.args.directoryListChannel, tt.args.waitForDirectoryCombination)
		})
	}
}
