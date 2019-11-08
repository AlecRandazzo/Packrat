package windowscollector

import (
	"golang.org/x/sys/windows"
	"reflect"
	"testing"
)

func TestGetVolumeHandler(t *testing.T) {
	type args struct {
		volumeLetter string
	}
	tests := []struct {
		name       string
		args       args
		wantVolume VolumeHandler
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVolume, err := GetVolumeHandler(tt.args.volumeLetter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVolumeHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVolume, tt.wantVolume) {
				t.Errorf("GetVolumeHandler() gotVolume = %v, want %v", gotVolume, tt.wantVolume)
			}
		})
	}
}

func Test_getHandle(t *testing.T) {
	type args struct {
		volumeLetter string
	}
	tests := []struct {
		name       string
		args       args
		wantHandle windows.Handle
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHandle, err := getHandle(tt.args.volumeLetter)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHandle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHandle != tt.wantHandle {
				t.Errorf("getHandle() gotHandle = %v, want %v", gotHandle, tt.wantHandle)
			}
		})
	}
}

func Test_identifyVolumesOfInterest(t *testing.T) {
	type args struct {
		exportList *ListOfFilesToExport
	}
	tests := []struct {
		name                  string
		args                  args
		wantVolumesOfInterest []string
		wantErr               bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVolumesOfInterest, err := identifyVolumesOfInterest(tt.args.exportList)
			if (err != nil) != tt.wantErr {
				t.Errorf("identifyVolumesOfInterest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVolumesOfInterest, tt.wantVolumesOfInterest) {
				t.Errorf("identifyVolumesOfInterest() gotVolumesOfInterest = %v, want %v", gotVolumesOfInterest, tt.wantVolumesOfInterest)
			}
		})
	}
}

func Test_isLetter(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := isLetter(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("isLetter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("isLetter() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
