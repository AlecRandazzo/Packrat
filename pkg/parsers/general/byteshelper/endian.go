// Copyright (c) 2022 Alec Randazzo

package byteshelper

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// LittleEndianBinaryToInt64 converts a little endian byte slice to int64.
func LittleEndianBinaryToInt64(inBytes []byte) (outInt64 int64, err error) {
	inBytesLength := len(inBytes)
	if inBytesLength > 8 || inBytesLength == 0 {
		return 0, fmt.Errorf("LittleEndianBinaryToInt64() received %d bytes but expected 1-8 bytes", inBytesLength)

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
	return int64(binary.LittleEndian.Uint64(inBytes)), nil
}

// LittleEndianBinaryToInt32 converts a little endian byte slice to int32.
func LittleEndianBinaryToInt32(inBytes []byte) (int32, error) {
	inBytesLength := len(inBytes)
	if inBytesLength > 4 || inBytesLength == 0 {
		return 0, fmt.Errorf("LittleEndianBinaryToInt32() received %d bytes but expected 1-4 bytes", inBytesLength)
	}

	// Check if the number is negative
	if inBytes[inBytesLength-1] >= 0x80 {
		// Pad it to get to 4 bytes
		if inBytesLength < 4 {
			bytesToPad := 4 - inBytesLength
			for i := 0; i < bytesToPad; i++ {
				inBytes = append(inBytes, 0xff)
			}
		}
	} else {
		// Pad it to get to 4 bytes
		if inBytesLength < 4 {
			bytesToPad := 4 - inBytesLength
			for i := 0; i < bytesToPad; i++ {
				inBytes = append(inBytes, 0x00)
			}
		}
	}
	return int32(binary.LittleEndian.Uint32(inBytes)), nil
}

// LittleEndianBinaryToUInt64 converts a little endian byte slice to uint64.
func LittleEndianBinaryToUInt64(inBytes []byte) (outUInt64 uint64, err error) {
	inBytesLength := len(inBytes)
	if inBytesLength > 8 || inBytesLength == 0 {
		return 0, fmt.Errorf("LittleEndianBinaryToUInt64() received %d bytes but expected 1-8 bytes", inBytesLength)
	}

	// Pad it to get to 8 bytes
	if inBytesLength < 8 {
		bytesToPad := 8 - inBytesLength
		for i := 0; i < bytesToPad; i++ {
			inBytes = append(inBytes, 0x00)
		}
	}
	return binary.LittleEndian.Uint64(inBytes), nil
}

// LittleEndianBinaryToInt16 converts a little endian byte slice to int16.
func LittleEndianBinaryToInt16(inBytes []byte) (int16, error) {
	inBytesLength := len(inBytes)
	if inBytesLength > 2 || inBytesLength == 0 {
		return 0, fmt.Errorf("LittleEndianBinaryToInt16() received %d bytes but expected 1-2 bytes", inBytesLength)
	}

	// Check if the number is negative
	if inBytes[inBytesLength-1] >= 0x80 {
		// Pad it to get to 2 bytes
		if inBytesLength < 2 {
			bytesToPad := 2 - inBytesLength
			for i := 0; i < bytesToPad; i++ {
				inBytes = append(inBytes, 0xff)
			}
		}
	} else {
		// Pad it to get to 2 bytes
		if inBytesLength < 2 {
			bytesToPad := 2 - inBytesLength
			for i := 0; i < bytesToPad; i++ {
				inBytes = append(inBytes, 0x00)
			}
		}
	}
	return int16(binary.LittleEndian.Uint16(inBytes)), nil
}

// UnicodeBytesToASCII converts a byte slice of unicode characters to ASCII
func UnicodeBytesToASCII(unicodeBytes []byte) (asciiString string, err error) {
	inBytesLength := len(unicodeBytes)
	if inBytesLength == 0 {
		return "", errors.New("UnicodeBytesToASCII() received no bytes")
	}
	return strings.Replace(string(unicodeBytes), "\x00", "", -1), nil
}
