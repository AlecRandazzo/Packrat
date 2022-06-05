// Copyright (c) 2022 Alec Randazzo

package binxml

import (
	"encoding/hex"
)

// These are effectively constants
// Various binary xml tokens
var (
	magicNumber, _ = hex.DecodeString("dfff") // 0xdfff
)

// Various binary xml tokens. Source: https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-even6/c73573ae-1c90-43a2-a65f-ad7501155956
const (
	eofToken                  = 0x00
	openStartElementToken     = 0x01
	closeStartElementToken    = 0x02
	closeEmptyElementToken    = 0x03
	endElementToken           = 0x04
	valueTextToken            = 0x05
	attributeToken            = 0x06
	charReftoken              = 0x08
	entityRefToken            = 0x09
	pITargetToken             = 0x0A
	pIDataToken               = 0x0B
	templateInstanceToken     = 0x0C
	normalSubstitutionToken   = 0x0D
	optionalSubstitutionToken = 0x0E
	fragmentHeaderToken       = 0x0F
	nullType                  = 0x00
	stringType                = 0x01
	ansiStringType            = 0x02
	int8Type                  = 0x03
	uInt8Type                 = 0x04
	int16Type                 = 0x05
	uInt16Type                = 0x06
	int32Type                 = 0x07
	uInt32Type                = 0x08
	int64Type                 = 0x09
	uInt64Type                = 0x0A
	real32Type                = 0x0B
	real64Type                = 0x0C
	boolType                  = 0x0D
	binaryType                = 0x0E
	guidType                  = 0x0F
	sizeTType                 = 0x10
	fileTimeType              = 0x11
	sysTimeType               = 0x12
	sidType                   = 0x13
	hexInt32Type              = 0x14
	hexInt64Type              = 0x15
	binXmlType                = 0x21
	stringArrayType           = 0x81
	ansiStringArrayType       = 0x82
	int8ArrayType             = 0x83
	uInt8ArrayType            = 0x84
	int16ArrayType            = 0x85
	uInt16ArrayType           = 0x86
	int32ArrayType            = 0x87
	uInt32ArrayType           = 0x88
	int64ArrayType            = 0x89
	uInt64ArrayType           = 0x8A
	real32ArrayType           = 0x8B
	real64ArrayType           = 0x8C
	boolArrayType             = 0x8D
	guidArrayType             = 0x8F
	sizeTArrayType            = 0x90
	fileTimeArrayType         = 0x91
	sysTimeArrayType          = 0x92
	sidArrayType              = 0x93
	hexInt32ArrayType         = 0x0094
	hexInt64ArrayType         = 0x0095
)

// ParseBinXmlOrdered will convert a byte slice of binary XML data and return the data in an ordered key value pairs.
func Parse(inBytes []byte) (map[int]map[string]string, error) {

	return nil, nil
}
