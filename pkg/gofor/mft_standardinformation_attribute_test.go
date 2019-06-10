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
	"testing"
)

func TestMasterFileTableRecord_getStandardInformationAttribute(t *testing.T) {
	attributeBytes, _ := hex.DecodeString("100000006000000000001800000000004800000018000000EA24CD4A74D4D101EA24CD4A74D4D101EA24CD4A74D4D101EA24CD4A74D4D10106000000000000000000000000000000000000000001000000000000000000000000000000000000")
	tests := []struct {
		name          string
		mftRecord     *masterFileTableRecord
		wantMftRecord *masterFileTableRecord
	}{
		{
			name: "Test MFT record 0.",
			mftRecord: &masterFileTableRecord{
				AttributeInfo: []attributeInfo{
					{
						AttributeType:  0x10,
						AttributeBytes: attributeBytes,
					},
				},
			},
			wantMftRecord: &masterFileTableRecord{
				StandardInformationAttributes: standardInformationAttributes{
					SiCreated:    "2016-07-02T15:13:30Z",
					SiModified:   "2016-07-02T15:13:30Z",
					SiAccessed:   "2016-07-02T15:13:30Z",
					SiChanged:    "2016-07-02T15:13:30Z",
					FlagResident: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mftRecord.getStandardInformationAttribute()
			assert.Equal(t, tt.wantMftRecord.StandardInformationAttributes, tt.mftRecord.StandardInformationAttributes, "Standard information should be equal.")
		})
	}
}
