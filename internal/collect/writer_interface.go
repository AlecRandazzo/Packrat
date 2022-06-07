package collect

import (
	"io"
)

type Writer interface {
	Write(io.Reader, string) error
	Close() error
}
