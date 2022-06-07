package collect

import (
	"archive/zip"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

type Zip struct {
	file   *os.File
	writer *zip.Writer
}

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
