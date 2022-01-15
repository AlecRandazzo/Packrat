// Copyright (c) 2020 Alec Randazzo

package byteshelper

import (
	"errors"
	"reflect"
)

type DataLocation struct {
	Offset byte
	Length byte
}

func NewDataLocation(offset, length byte) DataLocation {
	return DataLocation{
		Offset: offset,
		Length: length,
	}
}

func GetValue(bytes []byte, dataLocation DataLocation) ([]byte, error) {
	// Sanity Checks
	dataSize := len(bytes)
	if dataSize < int(dataLocation.Length)+int(dataLocation.Offset) {
		return nil, errors.New("GetValue() received a []byte that is not large enough to contain the dataLocation")
	}
	if reflect.DeepEqual(dataLocation, DataLocation{}) {
		return nil, errors.New("GetValue() received a nil DataLocation")
	}

	return bytes[dataLocation.Offset : int(dataLocation.Offset)+int(dataLocation.Length)], nil
}
