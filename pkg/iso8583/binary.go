package iso8583

import (
	"errors"
	"fmt"
)

// BINARY is a []byte implementation of a field,
// it does not contain any special behaviour more than unload all bytes on marshaling and
// reading the specified length on unmarshaling.
type BINARY []byte

// MarshalISO8583 returns a copy of binary content. Encoding and length input are ignored.
func (binary BINARY) MarshalISO8583(length int, enc string) ([]byte, error) {
	binaryCopy := make([]byte, len(binary))
	copy(binaryCopy, binary)

	return binaryCopy, nil
}

// UnmarshalISO8583 reads the length indicated amount of bytes from b and load the BINARY field with it.
// Encoding is ignored.
func (binary *BINARY) UnmarshalISO8583(b []byte, length int, enc string) (int, error) {
	if b == nil {
		return 0, errors.New("bytes input is nil")
	}

	if len(b) < length {
		return 0, fmt.Errorf("message remain (%v bytes) is shorter than indicated length: %v",
			len(b), length)
	}

	*binary = make([]byte, length)
	copy(*binary, b[:length])

	return length, nil
}
