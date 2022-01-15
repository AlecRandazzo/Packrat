// Copyright (c) 2020 Alec Randazzo

package sanitycheck

import "testing"

func Test_Bytes(t *testing.T) {
	type args struct {
		input        []byte
		expectedSize int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid slice",
			args: args{
				input:        []byte{0x00, 0x00, 0x00, 0x00},
				expectedSize: 4,
			},
			wantErr: false,
		},
		{
			name:    "nil bytes",
			args:    args{},
			wantErr: true,
		},
		{
			name: "not enough bytes",
			args: args{
				input:        []byte{0x00, 0x00},
				expectedSize: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Bytes(tt.args.input, tt.args.expectedSize); (err != nil) != tt.wantErr {
				t.Errorf("Bytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
