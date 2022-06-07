// Copyright (c) 2022 Alec Randazzo

//go:build windows

package volume

import (
	"testing"
)

func Test_GetHandle(t *testing.T) {
	type args struct {
		letter string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "no error",
			args:    args{letter: "C"},
			wantErr: false,
		},
		{
			name:    "nil string input",
			args:    args{letter: ""},
			wantErr: true,
		},
		{
			name:    "bad input",
			args:    args{letter: "CD"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handle := NewHandler(tt.args.letter)
			err := handle.GetHandle()
			if (err != nil) != tt.wantErr {
				t.Errorf("getHandle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			handle.Handle().Close()
		})
	}
}
