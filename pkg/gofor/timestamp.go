/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package gofor

import (
	"time"
)

// Parse a byte slice containing a unix timestamp and convert it to a timestamp string.
func parseTimestamp(timestampBytes []byte) (timestamp string) {

	var delta = time.Date(1970-369, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()

	// Convert the byte slice to little endian int64 and then convert it to a string
	timestampInt64 := convertLittleEndianByteSliceToInt64(timestampBytes)
	if timestampInt64 == 0 {
		timestamp = ""
		return
	}

	timestamp = time.Unix(0, int64(timestampInt64)*100+delta).UTC().Format("2006-01-02T15:04:05Z")

	return
}
