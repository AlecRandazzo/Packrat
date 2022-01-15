// Copyright (c) 2020 Alec Randazzo

package timestamp

import (
	"testing"
	"time"
)

func Test_Parse(t *testing.T) {
	tests := []struct {
		name           string
		timestamp      time.Time
		want           string
		wantErr        bool
		timestampBytes []byte
	}{
		{
			name:           "Timestamp test - valid timestamp",
			timestampBytes: []byte{234, 36, 205, 74, 116, 212, 209, 1},
			timestamp:      time.Time{},
			want:           "2016-07-02T15:13:30Z",
			wantErr:        false,
		},
		{
			name:           "Timestamp test - null timestamp",
			timestampBytes: []byte{},
			timestamp:      time.Time{},
			want:           "0001-01-01T00:00:00Z",
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got time.Time
			got, err := Parse(tt.timestampBytes)
			ts := got.Format("2006-01-02T15:04:05Z")
			if ts != tt.want || (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want = %v", ts, tt.want)
			}
		})
	}
}
