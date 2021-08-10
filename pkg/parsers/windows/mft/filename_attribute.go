// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/byteshelper"
	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/timestamp"
)

// FileNameAttributes is a slice that contains a list of filename attributes.
type FileNameAttributes []FileNameAttribute

type flagResidency bool

// FileNameAttribute contains information about a filename attribute.
type FileNameAttribute struct {
	FnCreated               time.Time
	FnModified              time.Time
	FnAccessed              time.Time
	FnChanged               time.Time
	FlagResident            bool
	NameLength              NameLength
	AttributeSize           uint32
	ParentDirRecordNumber   uint32
	ParentDirSequenceNumber uint16
	LogicalFileSize         uint64
	PhysicalFileSize        uint64
	FileNameFlags           FileNameFlags
	FileNameLength          byte
	FileNamespace           string
	FileName                string
}

type NameLength struct {
	FlagNamed bool
	NamedSize byte
}

// FileNameFlags contains possible filename flags a file may have.
type FileNameFlags struct {
	ReadOnly          bool
	Hidden            bool
	System            bool
	Archive           bool
	Device            bool
	Normal            bool
	Temporary         bool
	Sparse            bool
	Reparse           bool
	Compressed        bool
	Offline           bool
	NotContentIndexed bool
	Encrypted         bool
	Directory         bool
	IndexView         bool
}

// RawFilenameFlags is a []byte alias for raw filename flags. Used with the Parse() method.
type RawFilenameFlags []byte

// RawResidencyFlag is a byte alias for raw residency flag. Used with the Parse() method.
type RawResidencyFlag byte

// RawFileNameAttribute is a []byte alias for raw filename attribute. Used with the Parse() method.
type RawFileNameAttribute []byte

// RawFilenameNameSpaceFlag is a byte alias for raw filename namespace flag. Used with the Parse() method.
type RawFilenameNameSpaceFlag byte

// Parse parses the raw filename attribute receiver and returns a parsed filename attribute.
func (rawFileNameAttribute RawFileNameAttribute) Parse() (filenameAttribute FileNameAttribute, err error) {
	const offsetAttributeSize = 0x04
	const lengthAttributeSize = 0x04

	const offsetResidentFlag = 0x08

	const offsetParentRecordNumber = 0x18
	const lengthParentRecordNumber = 0x04

	const offsetParentDirSequenceNumber = 0x1e
	const lengthParentDirSequenceNumber = 0x02

	const offsetFnCreated = 0x20
	const lengthFnCreated = 0x08

	const offsetFnModified = 0x28
	const lengthFnModified = 0x08

	const offsetFnChanged = 0x30
	const lengthFnChanged = 0x08

	const offsetFnAccessed = 0x38
	const lengthFnAccessed = 0x08

	const offsetLogicalFileSize = 0x40
	const lengthLogicalFileSize = 0x08

	const offSetPhysicalFileSize = 0x48
	const lengthPhysicalFileSize = 0x08

	const offsetFnFlags = 0x50
	const lengthFnFlags = 0x04

	const offsetFileNameLength = 0x58
	const offsetFileNameSpace = 0x59
	const offsetFileName = 0x5a

	// Sanity check that we have data to work with
	attributeLength := len(rawFileNameAttribute)
	if attributeLength < 0x44 {
		err = errors.New("FileNameAttribute.parse() did not receive valid bytes")
		return
	}

	rawResidencyFlag := RawResidencyFlag(rawFileNameAttribute[offsetResidentFlag])
	filenameAttribute.FlagResident = rawResidencyFlag.Parse()
	if filenameAttribute.FlagResident == false {
		err = errors.New("parseFileNameAttribute(): non-resident filename Attribute encountered")
		return
	}
	filenameAttribute.AttributeSize, _ = byteshelper.LittleEndianBinaryToUInt32(rawFileNameAttribute[offsetAttributeSize : offsetAttributeSize+lengthAttributeSize])
	filenameAttribute.ParentDirRecordNumber, _ = byteshelper.LittleEndianBinaryToUInt32(rawFileNameAttribute[offsetParentRecordNumber : offsetParentRecordNumber+lengthParentRecordNumber])
	filenameAttribute.ParentDirSequenceNumber, _ = byteshelper.LittleEndianBinaryToUInt16(rawFileNameAttribute[offsetParentDirSequenceNumber : offsetParentDirSequenceNumber+lengthParentDirSequenceNumber])
	filenameAttribute.FnCreated, _ = timestamp.Parse(rawFileNameAttribute[offsetFnCreated : offsetFnCreated+lengthFnCreated])
	filenameAttribute.FnModified, _ = timestamp.Parse(rawFileNameAttribute[offsetFnModified : offsetFnModified+lengthFnModified])
	filenameAttribute.FnChanged, _ = timestamp.Parse(rawFileNameAttribute[offsetFnChanged : offsetFnChanged+lengthFnChanged])
	filenameAttribute.FnAccessed, _ = timestamp.Parse(rawFileNameAttribute[offsetFnAccessed : offsetFnAccessed+lengthFnAccessed])
	filenameAttribute.LogicalFileSize, _ = byteshelper.LittleEndianBinaryToUInt64(rawFileNameAttribute[offsetLogicalFileSize : offsetLogicalFileSize+lengthLogicalFileSize])
	filenameAttribute.PhysicalFileSize, _ = byteshelper.LittleEndianBinaryToUInt64(rawFileNameAttribute[offSetPhysicalFileSize : offSetPhysicalFileSize+lengthPhysicalFileSize])
	flagBytes := RawFilenameFlags(rawFileNameAttribute[offsetFnFlags : offsetFnFlags+lengthFnFlags])
	filenameAttribute.FileNameFlags = flagBytes.Parse()
	filenameAttribute.FileNameLength = rawFileNameAttribute[offsetFileNameLength] * 2 // times two to account for unicode characters
	rawFilenameNameSpaceFlag := RawFilenameNameSpaceFlag(rawFileNameAttribute[offsetFileNameSpace])
	filenameAttribute.FileNamespace = rawFilenameNameSpaceFlag.Parse()
	filenameAttribute.FileName, _ = byteshelper.UnicodeBytesToASCII(rawFileNameAttribute[offsetFileName : offsetFileName+int(filenameAttribute.FileNameLength)])
	return
}

// Parse parses the raw residency flag receiver and returns a flag residency value.
func (byteToCheck RawResidencyFlag) Parse() (flagResidency bool) {
	switch byteToCheck {
	case 0x00:
		flagResidency = true
	default:
		flagResidency = false
	}
	return
}

// Parse parses the raw file namespace flag receiver and returns a file namespace value.
func (fileNamespaceFlag RawFilenameNameSpaceFlag) Parse() (fileNameSpace string) {
	switch fileNamespaceFlag {
	case 0x00:
		fileNameSpace = "POSIX"
	case 0x01:
		fileNameSpace = "WIN32"
	case 0x02:
		fileNameSpace = "DOS"
	case 0x03:
		fileNameSpace = "WIN32 & DOS"
	default:
		fileNameSpace = ""
	}

	return
}

// Parse parses the raw filename flags receiver and returns filename flags.
func (flagBytes RawFilenameFlags) Parse() (fileNameFlags FileNameFlags) {
	unparsedFlags := binary.LittleEndian.Uint32(flagBytes)
	//init values
	fileNameFlags.ReadOnly = false
	fileNameFlags.Hidden = false
	fileNameFlags.System = false
	fileNameFlags.Archive = false
	fileNameFlags.Device = false
	fileNameFlags.Normal = false
	fileNameFlags.Temporary = false
	fileNameFlags.Sparse = false
	fileNameFlags.Reparse = false
	fileNameFlags.Compressed = false
	fileNameFlags.Offline = false
	fileNameFlags.NotContentIndexed = false
	fileNameFlags.Encrypted = false
	fileNameFlags.Directory = false
	fileNameFlags.IndexView = false

	if unparsedFlags&0x0001 != 0 {
		fileNameFlags.ReadOnly = true

	}
	if unparsedFlags&0x0002 != 0 {
		fileNameFlags.Hidden = true
	}
	if unparsedFlags&0x0004 != 0 {
		fileNameFlags.System = true
	}
	if unparsedFlags&0x0020 != 0 {
		fileNameFlags.Archive = true
	}
	if unparsedFlags&0x0040 != 0 {
		fileNameFlags.Device = true
	}
	if unparsedFlags&0x0080 != 0 {
		fileNameFlags.Normal = true
	}
	if unparsedFlags&0x0100 != 0 {
		fileNameFlags.Temporary = true
	}
	if unparsedFlags&0x0200 != 0 {
		fileNameFlags.Sparse = true
	}
	if unparsedFlags&0x0400 != 0 {
		fileNameFlags.Reparse = true
	}
	if unparsedFlags&0x0800 != 0 {
		fileNameFlags.Compressed = true
	}
	if unparsedFlags&0x1000 != 0 {
		fileNameFlags.Offline = true
	}
	if unparsedFlags&0x2000 != 0 {
		fileNameFlags.NotContentIndexed = true
	}
	if unparsedFlags&0x4000 != 0 {
		fileNameFlags.Encrypted = true
	}
	if unparsedFlags&0x10000000 != 0 {
		fileNameFlags.Directory = true
	}
	if unparsedFlags&0x20000000 != 0 {
		fileNameFlags.IndexView = true
	}
	return
}
