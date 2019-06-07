/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package main

import (
	"github.com/AlecRandazzo/GoFor/pkg/gofor"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Log configuration
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)
}

func main() {
	err := gofor.ParseMFT("MFT", "out.csv")
	if err != nil {
		log.Fatal(err)
	}

}
