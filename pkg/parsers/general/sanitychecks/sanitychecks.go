// Copyright (c) 2020 Alec Randazzo

package sanitychecks

import (
	"errors"
	"fmt"
)

// SanityCheckByteSlice will verify that the byte slice is not empty and check for the expected size. Disable size checking by setting expectedSize to 0.
func SanityCheckByteSlice(input []byte, expectedSize int) error {
	size := len(input)
	if size == 0 {
		return errors.New("nil byte slice")
	}

	if expectedSize == 0 {
		return nil
	} else if size < expectedSize {
		return fmt.Errorf("expected size of byte slice was %d, actual size was %d", expectedSize, size)
	}

	return nil
}
