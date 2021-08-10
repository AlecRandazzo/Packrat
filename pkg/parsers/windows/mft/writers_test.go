// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"sync"
	"testing"
	"time"
)

type DummyResultWriter struct {
	AggregatedData []byte
}

func (dummyWriter *DummyResultWriter) Write(inData []byte) (n int, err error) {
	for _, value := range inData {
		dummyWriter.AggregatedData = append(dummyWriter.AggregatedData, value)
	}
	return
}

func TestCsvWriter_Write(t *testing.T) {
	type args struct {
		outputChannel chan UsefulMftFields
		waitGroup     sync.WaitGroup
		writer        CsvResultWriter
		streamer      DummyResultWriter
	}
	tests := []struct {
		name            string
		writer          CsvResultWriter
		args            args
		usefulMftFields []UsefulMftFields
		want            []byte
	}{
		{
			name: "test1",
			args: args{
				outputChannel: nil,
				waitGroup:     sync.WaitGroup{},
				writer:        CsvResultWriter{},
				streamer:      DummyResultWriter{},
			},
			usefulMftFields: []UsefulMftFields{
				0: {
					RecordNumber:     1,
					FilePath:         "\\",
					FullPath:         "\\$MFTMirr",
					FileName:         "$MFTMirr",
					SystemFlag:       true,
					HiddenFlag:       true,
					ReadOnlyFlag:     false,
					DirectoryFlag:    false,
					DeletedFlag:      false,
					FnCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					PhysicalFileSize: 4096,
				},
				1: {
					RecordNumber:     2,
					FilePath:         "\\",
					FullPath:         "\\$MFTMirr2",
					FileName:         "$MFTMirr2",
					SystemFlag:       true,
					HiddenFlag:       true,
					ReadOnlyFlag:     false,
					DirectoryFlag:    false,
					DeletedFlag:      false,
					FnCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					FnChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiCreated:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiModified:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiAccessed:       time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					SiChanged:        time.Date(2018, 2, 25, 0, 10, 45, 642455000, time.UTC),
					PhysicalFileSize: 4096,
				},
			},
			want: []byte{82, 101, 99, 111, 114, 100, 32, 78, 117, 109, 98, 101, 114, 124, 68, 105, 114, 101, 99, 116, 111, 114, 121, 124, 83, 121, 115, 116, 101, 109, 32, 70, 105, 108, 101, 124, 72, 105, 100, 100, 101, 110, 124, 82, 101, 97, 100, 45, 111, 110, 108, 121, 124, 68, 101, 108, 101, 116, 101, 100, 124, 70, 105, 108, 101, 32, 80, 97, 116, 104, 124, 70, 105, 108, 101, 32, 78, 97, 109, 101, 124, 70, 105, 108, 101, 32, 83, 105, 122, 101, 124, 70, 105, 108, 101, 32, 67, 114, 101, 97, 116, 101, 100, 124, 70, 105, 108, 101, 32, 77, 111, 100, 105, 102, 105, 101, 100, 124, 70, 105, 108, 101, 32, 65, 99, 99, 101, 115, 115, 101, 100, 124, 70, 105, 108, 101, 32, 69, 110, 116, 114, 121, 32, 77, 111, 100, 105, 102, 105, 101, 100, 124, 70, 105, 108, 101, 78, 97, 109, 101, 32, 67, 114, 101, 97, 116, 101, 100, 124, 70, 105, 108, 101, 78, 97, 109, 101, 32, 77, 111, 100, 105, 102, 105, 101, 100, 124, 70, 105, 108, 101, 110, 97, 109, 101, 32, 65, 99, 99, 101, 115, 115, 101, 100, 124, 70, 105, 108, 101, 110, 97, 109, 101, 32, 69, 110, 116, 114, 121, 32, 77, 111, 100, 105, 102, 105, 101, 100, 10, 49, 124, 102, 97, 108, 115, 101, 124, 116, 114, 117, 101, 124, 116, 114, 117, 101, 124, 102, 97, 108, 115, 101, 124, 102, 97, 108, 115, 101, 124, 92, 124, 36, 77, 70, 84, 77, 105, 114, 114, 124, 52, 48, 57, 54, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 10, 50, 124, 102, 97, 108, 115, 101, 124, 116, 114, 117, 101, 124, 116, 114, 117, 101, 124, 102, 97, 108, 115, 101, 124, 102, 97, 108, 115, 101, 124, 92, 124, 36, 77, 70, 84, 77, 105, 114, 114, 50, 124, 52, 48, 57, 54, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 124, 50, 48, 49, 56, 45, 48, 50, 45, 50, 53, 84, 48, 48, 58, 49, 48, 58, 52, 53, 90, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.outputChannel = make(chan UsefulMftFields, 2)
			tt.args.waitGroup.Add(1)
			tt.args.streamer.AggregatedData = make([]byte, 0)
			go tt.writer.ResultWriter(&tt.args.streamer, &tt.args.outputChannel, &tt.args.waitGroup)
			for _, value := range tt.usefulMftFields {
				tt.args.outputChannel <- value
			}
			close(tt.args.outputChannel)
			tt.args.waitGroup.Wait()

			if !reflect.DeepEqual(tt.args.streamer.AggregatedData, tt.want) {
				t.Errorf(cmp.Diff(tt.args.streamer.AggregatedData, tt.want))
			}

		})
	}
}
