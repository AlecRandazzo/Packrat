// Copyright (c) 2022 Alec Randazzo

package mft

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/AlecRandazzo/Packrat/pkg/general/byteshelper"
	"github.com/AlecRandazzo/Packrat/pkg/general/sanitycheck"
	"github.com/AlecRandazzo/Packrat/pkg/general/timestamp"
)

// FileNameAttributes is a slice that contains a list of filename attributes.
type FileNameAttributes []FileNameAttribute

// FileNamespace is the namespace for a filename attribute.
type FileNamespace string

// FileNameAttribute contains information about a filename attribute.
type FileNameAttribute struct {
	Created                 time.Time
	Modified                time.Time
	Accessed                time.Time
	Changed                 time.Time
	FlagResident            bool
	NameLength              NameLength
	AttributeSize           uint32
	ParentDirRecordNumber   uint32
	ParentDirSequenceNumber uint16
	LogicalFileSize         uint64
	PhysicalFileSize        uint64
	FileNameFlags           FileNameFlags
	FileNameLength          byte
	FileNamespace           FileNamespace
	FileName                string
}

// NameLength contains information whether a filename attribute is named and its size.
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

var (
	fnAttributeSizeLocation           = byteshelper.NewDataLocation(0x04, 0x04)
	fnResidentFlagLocation            = byteshelper.NewDataLocation(0x08, 0x01)
	fnParentRecordNumberLocation      = byteshelper.NewDataLocation(0x18, 0x04)
	fnParentDirSequenceNumberLocation = byteshelper.NewDataLocation(0x1E, 0x02)
	fnCreatedLocation                 = byteshelper.NewDataLocation(0x20, 0x08)
	fnModifiedLocation                = byteshelper.NewDataLocation(0x28, 0x08)
	fnChangedLocation                 = byteshelper.NewDataLocation(0x30, 0x08)
	fnAccessedLocation                = byteshelper.NewDataLocation(0x38, 0x08)
	fnLogicalFileSizeLocation         = byteshelper.NewDataLocation(0x40, 0x08)
	fnPhysicalFileSizeLocation        = byteshelper.NewDataLocation(0x48, 0x08)
	fnFlagsLocation                   = byteshelper.NewDataLocation(0x50, 0x04)
	fnFileNameLengthLocation          = byteshelper.NewDataLocation(0x58, 0x01)
	fnFileNameSpaceLocation           = byteshelper.NewDataLocation(0x59, 0x01)
)

const fnFileNameOffset = 0x5A

// getFileNameAttribute parses the raw filename attribute.
func getFileNameAttribute(input []byte) (FileNameAttribute, error) {
	// Sanity checks
	err := sanitycheck.Bytes(input, 0x44)
	if err != nil {
		return FileNameAttribute{}, fmt.Errorf("invalid input: %w", err)
	}

	// Get filename values
	var fn FileNameAttribute
	var buffer []byte
	buffer, _ = byteshelper.GetValue(input, fnResidentFlagLocation)
	fn.FlagResident = checkResidency(buffer[0])
	if !fn.FlagResident {
		return FileNameAttribute{}, errors.New("non-resident filename attribute")
	}

	buffer, _ = byteshelper.GetValue(input, fnAttributeSizeLocation)
	fn.AttributeSize = binary.LittleEndian.Uint32(buffer)

	buffer, _ = byteshelper.GetValue(input, fnParentRecordNumberLocation)
	fn.ParentDirRecordNumber = binary.LittleEndian.Uint32(buffer)

	buffer, _ = byteshelper.GetValue(input, fnParentDirSequenceNumberLocation)
	fn.ParentDirSequenceNumber = binary.LittleEndian.Uint16(buffer)

	buffer, _ = byteshelper.GetValue(input, fnCreatedLocation)
	fn.Created, _ = timestamp.Parse(buffer)

	buffer, _ = byteshelper.GetValue(input, fnModifiedLocation)
	fn.Modified, _ = timestamp.Parse(buffer)

	buffer, _ = byteshelper.GetValue(input, fnChangedLocation)
	fn.Changed, _ = timestamp.Parse(buffer)

	buffer, _ = byteshelper.GetValue(input, fnAccessedLocation)
	fn.Accessed, _ = timestamp.Parse(buffer)

	buffer, _ = byteshelper.GetValue(input, fnLogicalFileSizeLocation)
	fn.LogicalFileSize, _ = byteshelper.LittleEndianBinaryToUInt64(buffer)

	buffer, _ = byteshelper.GetValue(input, fnPhysicalFileSizeLocation)
	fn.PhysicalFileSize, _ = byteshelper.LittleEndianBinaryToUInt64(buffer)

	buffer, _ = byteshelper.GetValue(input, fnFlagsLocation)
	fn.FileNameFlags = getFileNameFlags(buffer)

	buffer, _ = byteshelper.GetValue(input, fnFileNameLengthLocation)
	fn.FileNameLength = buffer[0] * 2 // times two to account for unicode characters

	buffer, _ = byteshelper.GetValue(input, fnFileNameSpaceLocation)
	fn.FileNamespace = getFileNameSpace(buffer[0])

	fileNameLocation := byteshelper.NewDataLocation(fnFileNameOffset, fn.FileNameLength)
	buffer, _ = byteshelper.GetValue(input, fileNameLocation)
	fn.FileName, _ = byteshelper.UnicodeBytesToASCII(buffer)

	return fn, nil
}

// getFileNameSpace parses the raw file namespace.
func getFileNameSpace(input byte) FileNamespace {
	switch input {
	case 0x00:
		return "POSIX"
	case 0x01:
		return "WIN32"
	case 0x02:
		return "DOS"
	case 0x03:
		return "WIN32 & DOS"
	default:
		return ""
	}
}

// getFileNameFlags parses the raw filename flags.
func getFileNameFlags(input []byte) FileNameFlags {
	unparsedFlags := binary.LittleEndian.Uint32(input)

	var flags FileNameFlags

	if unparsedFlags&0x0001 != 0 {
		flags.ReadOnly = true

	}
	if unparsedFlags&0x0002 != 0 {
		flags.Hidden = true
	}
	if unparsedFlags&0x0004 != 0 {
		flags.System = true
	}
	if unparsedFlags&0x0020 != 0 {
		flags.Archive = true
	}
	if unparsedFlags&0x0040 != 0 {
		flags.Device = true
	}
	if unparsedFlags&0x0080 != 0 {
		flags.Normal = true
	}
	if unparsedFlags&0x0100 != 0 {
		flags.Temporary = true
	}
	if unparsedFlags&0x0200 != 0 {
		flags.Sparse = true
	}
	if unparsedFlags&0x0400 != 0 {
		flags.Reparse = true
	}
	if unparsedFlags&0x0800 != 0 {
		flags.Compressed = true
	}
	if unparsedFlags&0x1000 != 0 {
		flags.Offline = true
	}
	if unparsedFlags&0x2000 != 0 {
		flags.NotContentIndexed = true
	}
	if unparsedFlags&0x4000 != 0 {
		flags.Encrypted = true
	}
	if unparsedFlags&0x10000000 != 0 {
		flags.Directory = true
	}
	if unparsedFlags&0x20000000 != 0 {
		flags.IndexView = true
	}
	return flags
}
