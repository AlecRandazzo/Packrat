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
	"testing"
)

func TestParseTimestamp(t *testing.T) {
	timestampBytes, _ := hex.DecodeString("EA24CD4A74D4D101")

	type args struct {
		timestampBytes []byte
	}
	tests := []struct {
		name          string
		args          args
		wantTimestamp string
	}{
		{
			name: "Timestamp test",
			args: args{
				timestampBytes: timestampBytes,
			},
			wantTimestamp: "2016-07-02T15:13:30Z",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTimestamp := parseTimestamp(tt.args.timestampBytes); gotTimestamp != tt.wantTimestamp {
				t.Errorf("parseTimestamp() = %v, want %v", gotTimestamp, tt.wantTimestamp)
			}
		})
	}
}
