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
	"encoding/binary"
)

// Convert a byte slice to a little endian int64.
func convertLittleEndianByteSliceToInt64(inBytes []byte) (outInt64 int64) {
	inBytesLength := len(inBytes)
	if inBytesLength == 0 {
		outInt64 = 0
		return
	}

	// Pad it to get to 8 bytes
	if inBytes[inBytesLength-1] >= 0x80 {
		if inBytesLength < 8 {
			bytesToPad := 8 - inBytesLength
			for i := 0; i < bytesToPad; i++ {
				inBytes = append(inBytes, 0xff)
			}
		}
	} else {
		if inBytesLength < 8 {
			bytesToPad := 8 - inBytesLength
			for i := 0; i < bytesToPad; i++ {
				inBytes = append(inBytes, 0x00)
			}
		}
	}
	outInt64 = int64(binary.LittleEndian.Uint64(inBytes))
	return
}

// Convert a byte slice to a little endian uint64.
func ConvertLittleEndianByteSliceToUInt64(inBytes []byte) (outUint64 uint64) {
	inBytesLength := len(inBytes)
	if inBytesLength == 0 {
		outUint64 = 0
		return
	}

	// Pad it to get to 8 bytes
	if inBytesLength < 8 {
		bytesToPad := 8 - inBytesLength
		for i := 0; i < bytesToPad; i++ {
			inBytes = append(inBytes, 0x00)
		}
	}
	outUint64 = binary.LittleEndian.Uint64(inBytes)
	return
}
