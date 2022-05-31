// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// These constants are the "magic number" aka first byte for each type of attribute.
const (
	magicNumberStandardInformation = 0x10
	magicNumberAttributeList       = 0x20
	magicNumberFileName            = 0x30
	magicNumberVolumeVersion       = 0x40
	magicNumberSecurityDescriptor  = 0x50
	magicNumberVolumeName          = 0x60
	magicNumberVolumeInformation   = 0x70
	magicNumberData                = 0x80
	magicNumberIndexRoot           = 0x90
	magicNumberIndexAllocation     = 0xA0
	magicNumberBitmap              = 0xB0
	magicNumberSymbolicLink        = 0xC0
	magicNumberReparsePoint        = 0xD0
	magicNumberEaInformation       = 0xE0
	magicNumberPropertySet         = 0xF0
)

var validAttributeTypes = []byte{
	magicNumberStandardInformation,
	magicNumberAttributeList,
	magicNumberFileName,
	magicNumberVolumeVersion,
	magicNumberSecurityDescriptor,
	magicNumberVolumeName,
	magicNumberVolumeInformation,
	magicNumberData,
	magicNumberIndexRoot,
	magicNumberIndexAllocation,
	magicNumberBitmap,
	magicNumberSymbolicLink,
	magicNumberReparsePoint,
	magicNumberEaInformation,
	magicNumberPropertySet,
}

// GetAttributes parses a slice of raw attributes and returns its filename, standard information, and dat attributes. It takes an argument for input per cluster (typically 4096) which is used for computing data run information in a data attributes.
func GetAttributes(input [][]byte, bytesPerCluster uint) (FileNameAttributes, StandardInformationAttribute, DataAttribute, AttributeListAttributes, error) {
	// Sanity checks
	size := len(input)
	if size == 0 {
		return FileNameAttributes{},
			StandardInformationAttribute{},
			DataAttribute{},
			AttributeListAttributes{},
			errors.New("nil sized rawAttributes slice")
	} else if bytesPerCluster == 0 {
		return FileNameAttributes{},
			StandardInformationAttribute{},
			DataAttribute{},
			AttributeListAttributes{},
			errors.New("received bytesPerCluster value of 0")
	}

	// init return variables
	fnAttributes := make(FileNameAttributes, 0)
	var siAttribute StandardInformationAttribute
	var dataAttribute DataAttribute
	attributesList := make(AttributeListAttributes, 0)

	// Determine what each raw attribute is and parse it accordingly.
	for _, rawAttribute := range input {

		// Sanity check to make sure the attribute actually has input in it.
		sizeOfRawAttribute := len(rawAttribute)
		if sizeOfRawAttribute == 0 {
			return FileNameAttributes{},
				StandardInformationAttribute{},
				DataAttribute{},
				AttributeListAttributes{},
				errors.New("came across a rawAttribute with a nil size")
		}

		// Check the first byte to see if it is one of the "magic number" input we care about. If it is, we parse those raw attributes accordingly.
		switch rawAttribute[0x00] {
		case magicNumberFileName:
			fnAttribute, err := getFileNameAttribute(rawAttribute)
			if err == nil {
				fnAttributes = append(fnAttributes, fnAttribute)
			}
		case magicNumberStandardInformation:
			siAttribute, _ = getStandardInformationAttribute(rawAttribute)
		case magicNumberData:
			var err error
			var result interface{}
			result, err = getDataAttribute(rawAttribute, bytesPerCluster)
			if err != nil {
				return FileNameAttributes{},
					StandardInformationAttribute{},
					DataAttribute{},
					AttributeListAttributes{},
					fmt.Errorf("failed to get data Attribute %w", err)
			}

			switch v := result.(type) {
			case ResidentDataAttribute:
				dataAttribute.ResidentDataAttribute = v
			case NonResidentDataAttribute:
				dataAttribute.NonResidentDataAttribute = v
			}

		case magicNumberAttributeList:
			var err error
			attributesList, err = getAttributeListAttribute(rawAttribute)
			if err != nil {
				err = fmt.Errorf("failed to get attribute list Attribute %w", err)
				return FileNameAttributes{},
					StandardInformationAttribute{},
					DataAttribute{},
					AttributeListAttributes{},
					fmt.Errorf("failed to get attribute list Attribute %w", err)
			}
		}
	}
	return fnAttributes,
		siAttribute,
		dataAttribute,
		attributesList,
		nil
}

const (
	offsetAttributeSize = 0x04
	lengthAttributeSize = 0x04
)

// GetRawAttributes returns the attribute from an unparsed mft record. It takes recordHeader as an argument since the record header contains the offset for the start of the attributes.
func GetRawAttributes(input []byte, recordHeader RecordHeader) (rawAttributes [][]byte, err error) {
	// sanity checks
	size := len(input)
	if size == 0 {
		return nil, errors.New("received nil input")
	}
	if recordHeader.AttributesOffset == 0 {
		return nil, errors.New("record header argument has an attribute offset value of 0")
	}

	// Init variable that tracks how far to the next Attribute
	var distanceToNextAttribute uint16
	offset := uint16(recordHeader.AttributesOffset)

	for {
		// Calculate offset to next Attribute
		offset += distanceToNextAttribute

		// Break if the offset is beyond the byte slice
		if offset > uint16(size) || offset+0x04 > uint16(size) {
			break
		}

		// Verify if the byte slice is actually an MFT Attribute
		err = validateAttribute(input[offset])
		if err != nil {
			break
		}

		attributeSize := binary.LittleEndian.Uint16(input[offset+offsetAttributeSize : offset+offsetAttributeSize+lengthAttributeSize])
		end := offset + attributeSize

		rawAttribute := make([]byte, attributeSize)
		copy(rawAttribute, input[offset:end])

		// Append the rawAttributes to the rawAttributes struct
		rawAttributes = append(rawAttributes, rawAttribute)

		// Track the distance to the next Attribute based on the size of the current Attribute
		distanceToNextAttribute = binary.LittleEndian.Uint16(input[offset+offsetAttributeSize : offset+offsetAttributeSize+lengthAttributeSize])
	}

	return rawAttributes, nil
}

// Checks if the byte value equals a valid attribute type.
func validateAttribute(attributeHeaderToCheck byte) error {
	for _, validType := range validAttributeTypes {
		if attributeHeaderToCheck == validType {
			return nil
		}
	}
	return errors.New("invalid attribute magic number")
}

func checkResidency(input byte) bool {
	if input == 0x00 {
		return true
	}
	return false
}
