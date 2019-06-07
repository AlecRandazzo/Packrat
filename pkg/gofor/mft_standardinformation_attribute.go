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
	"encoding/hex"
	"github.com/pkg/errors"
)

type StandardInformationAttributes struct {
	SiCreated    string
	SiModified   string
	SiAccessed   string
	SiChanged    string
	FlagResident bool
}

func (mftRecord *MasterFileTableRecord) getStandardInformationAttribute() (err error) {
	const codeStandardInformation = 0x10

	const offsetResidentFlag = 0x08

	const offsetSiCreated = 0x18
	const lengthSiCreated = 0x08

	const offsetSiModified = 0x20
	const lengthSiModified = 0x08

	const offsetSiChanged = 0x28
	const lengthSiChanged = 0x08

	const offsetSiAccessed = 0x30
	const lengthSiAccessed = 0x08

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("failed to parse standard info attribute")
		}
	}()

	for _, attribute := range mftRecord.AttributeInfo {
		if attribute.AttributeType == codeStandardInformation {
			// The standard information attribute has a minimum length of 0x30
			if len(attribute.AttributeBytes) < 0x30 {
				return
			}

			// Check to see if the standard information attribute is resident to the MFT or not
			if attribute.AttributeBytes[offsetResidentFlag] == 0x00 {
				mftRecord.StandardInformationAttributes.FlagResident = true
			} else {
				mftRecord.StandardInformationAttributes.FlagResident = false
				err = errors.Errorf("non resident standard information flag found, hex dump: %s", hex.EncodeToString(attribute.AttributeBytes))
				return
			}

			// Convert timestamps from bytes to time.Time

			mftRecord.StandardInformationAttributes.SiCreated = ParseTimestamp(attribute.AttributeBytes[offsetSiCreated : offsetSiCreated+lengthSiCreated])
			if mftRecord.StandardInformationAttributes.SiCreated == "" {
				err = errors.Wrap(err, "could not parse si created timestamp")
				return
			}

			mftRecord.StandardInformationAttributes.SiModified = ParseTimestamp(attribute.AttributeBytes[offsetSiModified : offsetSiModified+lengthSiModified])
			if mftRecord.StandardInformationAttributes.SiModified == "" {
				err = errors.Wrap(err, "could not parse si modified timestamp")
				return
			}

			mftRecord.StandardInformationAttributes.SiChanged = ParseTimestamp(attribute.AttributeBytes[offsetSiChanged : offsetSiChanged+lengthSiChanged])
			if mftRecord.StandardInformationAttributes.SiChanged == "" {
				err = errors.Wrap(err, "could not parse si changed timestamp")
				return
			}

			mftRecord.StandardInformationAttributes.SiAccessed = ParseTimestamp(attribute.AttributeBytes[offsetSiAccessed : offsetSiAccessed+lengthSiAccessed])
			if mftRecord.StandardInformationAttributes.SiAccessed == "" {
				err = errors.Wrap(err, "could not parse si accessed timestamp")
				return
			}
			return
		}
	}
	return
}
