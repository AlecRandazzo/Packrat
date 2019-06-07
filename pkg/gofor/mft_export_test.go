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
	"sync"
	"testing"
)

func TestMftFile_mftToCSV(t *testing.T) {
	type args struct {
		outFileName string
		waitgroup   *sync.WaitGroup
	}
	tests := []struct {
		name    string
		file    MftFile
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.file.mftToCSV(tt.args.outFileName, tt.args.waitgroup); (err != nil) != tt.wantErr {
				t.Errorf("MftFile.mftToCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
