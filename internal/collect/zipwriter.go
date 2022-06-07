package collect

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// Zip is used for writing files to a zip.
type Zip struct {
	file   *os.File
	writer *zip.Writer
}

// NewZipWriter creates a new zip writer.
func NewZipWriter(file string) (Zip, error) {
	var err error
	z := Zip{}
	z.file, err = os.Create(file)
	if err != nil {
		return Zip{}, fmt.Errorf("failed to create zip file %s: %w", file, err)
	}

	z.writer = zip.NewWriter(z.file)

	return z, nil
}

// Write writes a file to a zip.
func (z Zip) Write(reader io.Reader, fullPath string) error {
	writer, err := z.writer.Create(fullPath)
	if err != nil {
		log.Errorf("failed to create file in zip: %v", err)
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		log.Errorf("failed to write %s", fullPath)
	}

	return nil
}

// Close closes the handles to the zip writer and to the file.
func (z Zip) Close() error {
	var errRollup error
	err := z.writer.Close()
	if err != nil {
		errRollup = fmt.Errorf("%s, %w", errRollup.Error(), err)
	}
	err = z.file.Close()
	if err != nil {
		errRollup = err
	}

	return errRollup
}
