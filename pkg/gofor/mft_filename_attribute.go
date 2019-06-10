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
	"encoding/hex"
	"github.com/pkg/errors"
	"strings"
)

type fileNameAttributes struct {
	FnCreated               string
	FnModified              string
	FnAccessed              string
	FnChanged               string
	FlagResident            bool
	FlagNamed               bool
	NamedSize               byte
	AttributeSize           uint32
	ParentDirRecordNumber   uint64
	ParentDirSequenceNumber uint16
	LogicalFileSize         uint64
	PhysicalFileSize        uint64
	FileNameFlags           fileNameFlags
	FileNameLength          byte
	FileNamespace           string
	FileName                string
}

type fileNameFlags struct {
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

func (mftRecord *masterFileTableRecord) getFileNameAttributes() (err error) {
	const codeFileName = 0x30

	const offsetAttributeSize = 0x04
	const lengthAttributeSize = 0x04

	const offsetResidentFlag = 0x08
	const offsetNameLength = 0x09

	const offsetParentRecordNumber = 0x18
	const lengthParentRecordNumber = 0x06

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

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("Failed to parse filename attribute")
		}
	}()

	for _, attribute := range mftRecord.AttributeInfo {
		if attribute.AttributeType == codeFileName {
			// The filename attribute has a minimum length of 0x44
			if len(attribute.AttributeBytes) < 0x44 {
				return
			}
			var fileNameAttributes fileNameAttributes
			fileNameAttributes.AttributeSize = binary.LittleEndian.Uint32(attribute.AttributeBytes[offsetAttributeSize : offsetAttributeSize+lengthAttributeSize])

			switch attribute.AttributeBytes[offsetResidentFlag] {
			case 0x00:
				fileNameAttributes.FlagResident = true
			default:
				fileNameAttributes.FlagResident = false
				err = errors.Errorf("\nparseFileNameAttribute(): non-resident filename attribute encountered, hex dump: %s", hex.EncodeToString(attribute.AttributeBytes))
				return
			}

			var nameLengthOffsetModifier byte
			switch attribute.AttributeBytes[offsetNameLength] {
			case 0x00:
				fileNameAttributes.FlagNamed = false
				nameLengthOffsetModifier = 0x00
			default:
				fileNameAttributes.FlagNamed = true
				fileNameAttributes.NamedSize = attribute.AttributeBytes[offsetNameLength]
				nameLengthOffsetModifier = nameLengthOffsetModifier * 2 // x2 to account for unicode
			}

			fileNameAttributes.ParentDirRecordNumber = ConvertLittleEndianByteSliceToUInt64(attribute.AttributeBytes[offsetParentRecordNumber+nameLengthOffsetModifier : offsetParentRecordNumber+lengthParentRecordNumber+nameLengthOffsetModifier])
			if fileNameAttributes.ParentDirRecordNumber == 0 {
				err = errors.Wrap(err, "failed to convert filename attribute's parent dir record number")
				return
			}

			fileNameAttributes.ParentDirSequenceNumber = binary.LittleEndian.Uint16(attribute.AttributeBytes[offsetParentDirSequenceNumber+nameLengthOffsetModifier : offsetParentDirSequenceNumber+lengthParentDirSequenceNumber+nameLengthOffsetModifier])

			fileNameAttributes.FnCreated = parseTimestamp(attribute.AttributeBytes[offsetFnCreated+nameLengthOffsetModifier : offsetFnCreated+lengthFnCreated+nameLengthOffsetModifier])
			if fileNameAttributes.FnCreated == "" {
				err = errors.Wrap(err, "could not parse fn created timestamp")
				return
			}

			fileNameAttributes.FnModified = parseTimestamp(attribute.AttributeBytes[offsetFnModified+nameLengthOffsetModifier : offsetFnModified+lengthFnModified+nameLengthOffsetModifier])
			if fileNameAttributes.FnModified == "" {
				err = errors.Wrap(err, "could not parse fn modified timestamp")
				return
			}

			fileNameAttributes.FnChanged = parseTimestamp(attribute.AttributeBytes[offsetFnChanged+nameLengthOffsetModifier : offsetFnChanged+lengthFnChanged+nameLengthOffsetModifier])
			if fileNameAttributes.FnChanged == "" {
				err = errors.Wrap(err, "could not parse fn changed timestamp")
				return
			}

			fileNameAttributes.FnAccessed = parseTimestamp(attribute.AttributeBytes[offsetFnAccessed+nameLengthOffsetModifier : offsetFnAccessed+lengthFnAccessed+nameLengthOffsetModifier])
			if fileNameAttributes.FnAccessed == "" {
				err = errors.Wrap(err, "could not parse fn accessed timestamp")
				return
			}

			fileNameAttributes.LogicalFileSize = binary.LittleEndian.Uint64(attribute.AttributeBytes[offsetLogicalFileSize+nameLengthOffsetModifier : offsetLogicalFileSize+lengthLogicalFileSize+nameLengthOffsetModifier])

			fileNameAttributes.PhysicalFileSize = binary.LittleEndian.Uint64(attribute.AttributeBytes[offSetPhysicalFileSize+nameLengthOffsetModifier : offSetPhysicalFileSize+lengthPhysicalFileSize+nameLengthOffsetModifier])

			fnFlags := attribute.AttributeBytes[offsetFnFlags+nameLengthOffsetModifier : offsetFnFlags+lengthFnFlags+nameLengthOffsetModifier]
			fileNameAttributes.FileNameFlags = resolveFileFlags(fnFlags)

			fileNameAttributes.FileNameLength = attribute.AttributeBytes[offsetFileNameLength+nameLengthOffsetModifier] * 2 // times two to account for unicode characters

			fileNamespaceFlag := attribute.AttributeBytes[offsetFileNameSpace+nameLengthOffsetModifier]
			switch fileNamespaceFlag {
			case 0x00:
				fileNameAttributes.FileNamespace = "POSIX"
			case 0x01:
				fileNameAttributes.FileNamespace = "WIN32"
			case 0x02:
				fileNameAttributes.FileNamespace = "DOS"
			case 0x03:
				fileNameAttributes.FileNamespace = "WIN32 & DOS"
			default:
				fileNameAttributes.FileNamespace = ""
			}

			unicodeFileName := string(attribute.AttributeBytes[offsetFileName+nameLengthOffsetModifier : offsetFileName+fileNameAttributes.FileNameLength+nameLengthOffsetModifier])

			fileNameAttributes.FileName = strings.Replace(unicodeFileName, "\x00", "", -1)
			mftRecord.FileNameAttributes = append(mftRecord.FileNameAttributes, fileNameAttributes)
		}
	}
	return
}

func resolveFileFlags(flagBytes []byte) (parsedFlags fileNameFlags) {
	unparsedFlags := binary.LittleEndian.Uint32(flagBytes)

	//init values
	parsedFlags.ReadOnly = false
	parsedFlags.Hidden = false
	parsedFlags.System = false
	parsedFlags.Archive = false
	parsedFlags.Device = false
	parsedFlags.Normal = false
	parsedFlags.Temporary = false
	parsedFlags.Sparse = false
	parsedFlags.Reparse = false
	parsedFlags.Compressed = false
	parsedFlags.Offline = false
	parsedFlags.NotContentIndexed = false
	parsedFlags.Encrypted = false
	parsedFlags.Directory = false
	parsedFlags.IndexView = false

	if unparsedFlags&0x0001 != 0 {
		parsedFlags.ReadOnly = true
	}
	if unparsedFlags&0x0002 != 0 {
		parsedFlags.Hidden = true
	}
	if unparsedFlags&0x0004 != 0 {
		parsedFlags.System = true
	}
	if unparsedFlags&0x0010 != 0 {
		parsedFlags.Directory = true
	}
	if unparsedFlags&0x0020 != 0 {
		parsedFlags.Archive = true
	}
	if unparsedFlags&0x0040 != 0 {
		parsedFlags.Device = true
	}
	if unparsedFlags&0x0080 != 0 {
		parsedFlags.Normal = true
	}
	if unparsedFlags&0x0100 != 0 {
		parsedFlags.Temporary = true
	}
	if unparsedFlags&0x0200 != 0 {
		parsedFlags.Sparse = true
	}
	if unparsedFlags&0x0400 != 0 {
		parsedFlags.Reparse = true
	}
	if unparsedFlags&0x0800 != 0 {
		parsedFlags.Compressed = true
	}
	if unparsedFlags&0x1000 != 0 {
		parsedFlags.Offline = true
	}
	if unparsedFlags&0x2000 != 0 {
		parsedFlags.NotContentIndexed = true
	}
	if unparsedFlags&0x4000 != 0 {
		parsedFlags.Encrypted = true
	}
	if unparsedFlags&0x10000000 != 0 {
		parsedFlags.Directory = true
	}
	if unparsedFlags&0x20000000 != 0 {
		parsedFlags.IndexView = true
	}
	return
}
