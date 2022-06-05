// Copyright (c) 2022 Alec Randazzo

package byteshelper

import (
	"reflect"
	"testing"
)

func TestGetValue(t *testing.T) {
	type args struct {
		bytes        []byte
		dataLocation DataLocation
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				bytes: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
				dataLocation: DataLocation{
					Offset: 0,
					Length: 2,
				},
			},
			want:    []byte{0x00, 0x01},
			wantErr: false,
		},
		{
			name: "nil bytes",
			args: args{
				bytes: nil,
				dataLocation: DataLocation{
					Offset: 0x00,
					Length: 0x02,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil DataLocation",
			args: args{
				bytes:        []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
				dataLocation: DataLocation{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "out of bounds",
			args: args{
				bytes: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
				dataLocation: DataLocation{
					Offset: 0x10,
					Length: 0x04,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetValue(tt.args.bytes, tt.args.dataLocation)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}
