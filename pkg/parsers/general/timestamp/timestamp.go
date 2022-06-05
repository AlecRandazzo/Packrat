// Copyright (c) 2022 Alec Randazzo

package timestamp

import (
	"fmt"
	"time"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/byteshelper"
)

// Parse a byte slice containing a unix timestamp and convert it to a timestamp string.
func Parse(raw []byte) (timestamp time.Time, err error) {
	// verify that we are getting valid data
	if len(raw) != 8 {
		return time.Time{}, fmt.Errorf("received %d bytes, not 8 bytes", len(raw))
	}

	var delta = time.Date(1970-369, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()
	// Convert the byte slice to little endian int64 and then convert it to a string
	timestampInt64, _ := byteshelper.LittleEndianBinaryToInt64(raw)

	return time.Unix(0, timestampInt64*100+delta).UTC(), nil
}
