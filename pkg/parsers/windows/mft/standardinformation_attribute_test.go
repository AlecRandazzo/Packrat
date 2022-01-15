// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"reflect"
	"testing"
	"time"
)

func Test_getStandardInformationAttribute(t *testing.T) {
	tests := []struct {
		name    string
		want    StandardInformationAttribute
		input   []byte
		wantErr bool
	}{
		{
			name:    "test 1",
			wantErr: false,
			input:   []byte{16, 0, 0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 49, 147, 66, 169, 237, 209, 211, 1, 49, 147, 66, 169, 237, 209, 211, 1, 44, 238, 221, 229, 226, 245, 211, 1, 49, 147, 66, 169, 237, 209, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 253, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 168, 220, 169, 88, 0, 0, 0, 0},
			want: StandardInformationAttribute{
				Created:      time.Date(2018, 4, 11, 23, 34, 40, 104324900, time.UTC),
				Modified:     time.Date(2018, 4, 11, 23, 34, 40, 104324900, time.UTC),
				Accessed:     time.Date(2018, 4, 11, 23, 34, 40, 104324900, time.UTC),
				Changed:      time.Date(2018, 5, 27, 17, 48, 19, 181726000, time.UTC),
				FlagResident: true,
			},
		},
		{
			name:    "nil input",
			wantErr: true,
			input:   nil,
		},
		{
			name:    "non-resident",
			wantErr: true,
			input:   []byte{16, 0, 0, 0, 96, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 72, 0, 0, 0, 24, 0, 0, 0, 49, 147, 66, 169, 237, 209, 211, 1, 49, 147, 66, 169, 237, 209, 211, 1, 44, 238, 221, 229, 226, 245, 211, 1, 49, 147, 66, 169, 237, 209, 211, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 253, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 168, 220, 169, 88, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si, err := getStandardInformationAttribute(tt.input)
			if !reflect.DeepEqual(si, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("Test %v failed \ngot = %v, \nwant = %v", tt.name, si, tt.want)
			}
		})
	}
}
