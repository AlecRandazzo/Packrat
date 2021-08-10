// Copyright (c) 2020 Alec Randazzo

package writers

import "io"

type TsvWriter struct {
	writer        *io.Writer
	headerWritten bool
}

func NewTsvWriter(writer *io.Writer) TsvWriter {
	return TsvWriter{writer: writer}
}

func (tsv TsvWriter) Write(data interface{}) error {

	return nil
}
