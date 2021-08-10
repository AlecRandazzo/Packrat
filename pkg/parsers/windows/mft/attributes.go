// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// []byte alias containing bytes of a raw MFT record attribute.
type rawAttribute []byte

// RawAttributes contains a slice of rawAttribute []byte aliases. Used primarily as a method receiver for the parse() method.
// See here for a handy list of attributes: https://flatcap.org/linux-ntfs/ntfs/attributes/index.html
type RawAttributes []rawAttribute

// Parse parses a slice of raw attributes and returns its filename, standard information, and dat attributes. It takes an argument for bytes per cluster (typically 4096) which is used for computing data run information in a data attributes.
func (rawAttributes RawAttributes) Parse(bytesPerCluster int64) (fileNameAttributes FileNameAttributes, standardInformationAttribute StandardInformationAttribute, dataAttribute DataAttribute, attributeListAttributes AttributeListAttributes, err error) {
	// Sanity check to make sure that the method received valid data
	sizeOfRawAttributesSlice := len(rawAttributes)
	if sizeOfRawAttributesSlice == 0 {
		err = errors.New("nil sized rawAttributes slice")
		return
	} else if bytesPerCluster == 0 {
		err = errors.New("received bytesPerCluster value of 0")
		return
	}

	// These constants are the "magic number" aka first byte for each type of attribute.
	const codeStandardInformation = 0x10
	const codeattributeList = 0x20
	const codeFileName = 0x30
	const codeData = 0x80

	// Determine what each raw attribute is and parse it accordingly.
	for _, rawAttribute := range rawAttributes {

		// Sanity check to make sure the attribute actually has bytes in it.
		sizeOfRawAttribute := len(rawAttribute)
		if sizeOfRawAttribute == 0 {
			err = errors.New("came across a rawAttribute with a nil size")
			fileNameAttributes = nil
			standardInformationAttribute = StandardInformationAttribute{}
			dataAttribute = DataAttribute{}
			return
		}

		// Check the first byte to see if it is one of the "magic number" bytes we care about. If it is, we parse those raw attributes accordingly.
		switch rawAttribute[0x00] {
		case codeFileName:
			rawFileNameAttribute := RawFileNameAttribute(make([]byte, len(rawAttribute)))
			copy(rawFileNameAttribute, rawAttribute)
			var fileNameAttribute FileNameAttribute
			fileNameAttribute, err = rawFileNameAttribute.Parse()
			if err != nil {
				err = fmt.Errorf("failed to get filename Attribute %w", err)
				fileNameAttributes = nil
				standardInformationAttribute = StandardInformationAttribute{}
				dataAttribute = DataAttribute{}
				return
			}
			fileNameAttributes = append(fileNameAttributes, fileNameAttribute)
		case codeStandardInformation:
			rawStandardInformationAttribute := RawStandardInformationAttribute(make([]byte, len(rawAttribute)))
			copy(rawStandardInformationAttribute, rawAttribute)
			standardInformationAttribute, err = rawStandardInformationAttribute.Parse()
			if err != nil {
				err = fmt.Errorf("failed to get standard info Attribute %w", err)
				fileNameAttributes = nil
				standardInformationAttribute = StandardInformationAttribute{}
				dataAttribute = DataAttribute{}
				return
			}
		case codeData:
			rawDataAttribute := RawDataAttribute(make([]byte, len(rawAttribute)))
			copy(rawDataAttribute, rawAttribute)
			dataAttribute.NonResidentDataAttribute, dataAttribute.ResidentDataAttribute, err = rawDataAttribute.Parse(bytesPerCluster)
			if err != nil {
				err = fmt.Errorf("failed to get data Attribute %w", err)
				fileNameAttributes = nil
				standardInformationAttribute = StandardInformationAttribute{}
				dataAttribute = DataAttribute{}
				return
			}
		case codeattributeList:
			rawAttributeListAttribute := RawAttributeListAttribute(make([]byte, len(rawAttribute)))
			copy(rawAttributeListAttribute, rawAttribute)
			attributeListAttributes, err = rawAttributeListAttribute.Parse()
			if err != nil {
				err = fmt.Errorf("failed to get attribute list Attribute %w", err)
				attributeListAttributes = AttributeListAttributes{}
				return
			}
		}
	}
	return
}

// GetRawAttributes returns the attribute bytes from an unparsed mft record which is the method receiver. It takes recordHeader as an argument since the record header contains the offset for the start of the attributes.
func (rawMftRecord RawMasterFileTableRecord) GetRawAttributes(recordHeader RecordHeader) (rawAttributes RawAttributes, err error) {
	// Doing some sanity checks
	if len(rawMftRecord) == 0 {
		err = errors.New("received nil bytes")
		return
	}
	if recordHeader.AttributesOffset == 0 {
		err = errors.New("record header argument has an attribute offset value of 0")
		return
	}

	const offsetAttributeSize = 0x04
	const lengthAttributeSize = 0x04

	// Init variable that tracks how far to the next Attribute
	var distanceToNextAttribute uint16
	offset := recordHeader.AttributesOffset
	sizeOfRawMftRecord := len(rawMftRecord)

	for {
		// Calculate offset to next Attribute
		offset = offset + distanceToNextAttribute

		// Break if the offset is beyond the byte slice
		if offset > uint16(sizeOfRawMftRecord) || offset+0x04 > uint16(sizeOfRawMftRecord) {
			break
		}

		// Verify if the byte slice is actually an MFT Attribute
		shouldWeContinue := isThisAnAttribute(rawMftRecord[offset])
		if shouldWeContinue == false {
			break
		}

		attributeSize := binary.LittleEndian.Uint16(rawMftRecord[offset+offsetAttributeSize : offset+offsetAttributeSize+lengthAttributeSize])
		end := offset + attributeSize

		rawAttribute := rawAttribute(make([]byte, attributeSize))
		copy(rawAttribute, rawMftRecord[offset:end])

		// Append the rawAttributes to the RawAttributes struct
		rawAttributes = append(rawAttributes, rawAttribute)

		// Track the distance to the next Attribute based on the size of the current Attribute
		distanceToNextAttribute = binary.LittleEndian.Uint16(rawMftRecord[offset+offsetAttributeSize : offset+offsetAttributeSize+lengthAttributeSize])
	}

	return
}

// Checks if the byte value equals a valid attribute type. We only do things with a few of these.
func isThisAnAttribute(attributeHeaderToCheck byte) (result bool) {
	// Init a byte slice that tracks all possible valid MFT Attribute types.
	// We'll be used this to verify if what we are looking at is actually an MFT Attribute.
	const codeStandardInformation = 0x10
	const codeAttributeList = 0x20
	const codeFileName = 0x30
	const codeVolumeVersion = 0x40
	const codeSecurityDescriptor = 0x50
	const codeVolumeName = 0x60
	const codeVolumeInformation = 0x70
	const codeData = 0x80
	const codeIndexRoot = 0x90
	const codeIndexAllocation = 0xA0
	const codeBitmap = 0xB0
	const codeSymbolicLink = 0xC0
	const codeReparsePoint = 0xD0
	const codeEaInformation = 0xE0
	const codePropertySet = 0xF0

	validAttributeTypes := []byte{
		codeStandardInformation,
		codeAttributeList,
		codeFileName,
		codeVolumeVersion,
		codeSecurityDescriptor,
		codeVolumeName,
		codeVolumeInformation,
		codeData,
		codeIndexRoot,
		codeIndexAllocation,
		codeBitmap,
		codeSymbolicLink,
		codeReparsePoint,
		codeEaInformation,
		codePropertySet,
	}

	// Verify if the byte slice is actually an MFT Attribute
	for _, validType := range validAttributeTypes {
		if attributeHeaderToCheck == validType {
			result = true
			break
		} else {
			result = false
		}
	}

	return
}
