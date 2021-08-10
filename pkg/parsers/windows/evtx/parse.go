// Copyright (c) 2020 Alec Randazzo

package evtx

import (
	"io"

	"github.com/AlecRandazzo/Packrat/internal/writers"
)

type Evt struct {
	reader       io.Reader
	resultWriter writers.ResultWriter
	fileHeader   fileHeader
}

func (evt *Evt) Parse() {

}

func NewEvtxReader(reader io.Reader, resultWriter writers.ResultWriter) Evt {
	return Evt{
		reader:       reader,
		resultWriter: resultWriter,
	}
}
