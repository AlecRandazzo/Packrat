// Copyright (c) 2022 Alec Randazzo

package main

import (
	"os"

	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
)

var CLI struct {
	Collect CollectCmd `cmd:"" help:"Collect forensic data."`
	Parse   ParseCmd   `cmd:"" help:"Parse forensic data."`
	Debug   bool       `help:"Enable debug mode."`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)

	if CLI.Debug {
		debugLog, err := os.Create("debug.log")
		if err != nil {
			panic(err)
		}
		log.SetOutput(debugLog)
		log.SetLevel(log.DebugLevel)
	}

	ctx := kong.Parse(&CLI)
	err := ctx.Run()
	if err != nil {
		ctx.FatalIfErrorf(err)
	}
}
