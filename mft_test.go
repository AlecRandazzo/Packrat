/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package windowscollector

import (
	mft "github.com/Go-Forensics/MFT-Parser"
	"reflect"
	"testing"
)

func Test_parseMFTRecord0(t *testing.T) {
	type args struct {
		volume *VolumeHandler
	}
	tests := []struct {
		name           string
		args           args
		wantMftRecord0 mft.MasterFileTableRecord
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMftRecord0, err := parseMFTRecord0(tt.args.volume)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMFTRecord0() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMftRecord0, tt.wantMftRecord0) {
				t.Errorf("parseMFTRecord0() gotMftRecord0 = %v, want %v", gotMftRecord0, tt.wantMftRecord0)
			}
		})
	}
}
