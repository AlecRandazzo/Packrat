package windowscollector

import (
	mft "github.com/AlecRandazzo/GoFor-MFT-Parser"
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
