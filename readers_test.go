package windowscollector

import (
	mft "github.com/AlecRandazzo/GoFor-MFT-Parser"
	"io"
	"reflect"
	"testing"
)

func TestDataRunsReader_Read(t *testing.T) {
	type fields struct {
		VolumeHandler                 *VolumeHandler
		DataRuns                      mft.DataRuns
		fileName                      string
		dataRunTracker                int
		dataRunBytesLeftToReadTracker int64
		totalFileSize                 int64
		totalByesRead                 int64
		initialized                   bool
	}
	type args struct {
		byteSliceToPopulate []byte
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantNumberOfBytesRead int
		wantErr               bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataRunReader := &DataRunsReader{
				VolumeHandler:                 tt.fields.VolumeHandler,
				DataRuns:                      tt.fields.DataRuns,
				fileName:                      tt.fields.fileName,
				dataRunTracker:                tt.fields.dataRunTracker,
				dataRunBytesLeftToReadTracker: tt.fields.dataRunBytesLeftToReadTracker,
				totalFileSize:                 tt.fields.totalFileSize,
				totalByesRead:                 tt.fields.totalByesRead,
				initialized:                   tt.fields.initialized,
			}
			gotNumberOfBytesRead, err := dataRunReader.Read(tt.args.byteSliceToPopulate)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNumberOfBytesRead != tt.wantNumberOfBytesRead {
				t.Errorf("Read() gotNumberOfBytesRead = %v, want %v", gotNumberOfBytesRead, tt.wantNumberOfBytesRead)
			}
		})
	}
}

func Test_apiFileReader(t *testing.T) {
	type args struct {
		file foundFile
	}
	tests := []struct {
		name       string
		args       args
		wantReader io.Reader
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReader, err := apiFileReader(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiFileReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotReader, tt.wantReader) {
				t.Errorf("apiFileReader() gotReader = %v, want %v", gotReader, tt.wantReader)
			}
		})
	}
}

func Test_rawFileReader(t *testing.T) {
	type args struct {
		handler *VolumeHandler
		file    foundFile
	}
	tests := []struct {
		name       string
		args       args
		wantReader io.Reader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotReader := rawFileReader(tt.args.handler, tt.args.file); !reflect.DeepEqual(gotReader, tt.wantReader) {
				t.Errorf("rawFileReader() = %v, want %v", gotReader, tt.wantReader)
			}
		})
	}
}
