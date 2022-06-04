// Copyright (c) 2020 Alec Randazzo

//go:build windows

package windows

import (
	"testing"
)

func Test_GetHandle(t *testing.T) {
	type args struct {
		volumeLetter string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "no error",
			args:    args{volumeLetter: "C"},
			wantErr: false,
		},
		{
			name:    "nil string input",
			args:    args{volumeLetter: ""},
			wantErr: true,
		},
		{
			name:    "bad input",
			args:    args{volumeLetter: "CD"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handle := NewVolumeHandler(tt.args.volumeLetter)
			err := handle.GetHandle()
			if (err != nil) != tt.wantErr {
				t.Errorf("getHandle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			handle.Handle().Close()
		})
	}
}
