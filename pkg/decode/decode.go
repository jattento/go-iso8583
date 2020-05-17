package decode

import (
	"errors"
)

const (
	_mtiKey    = "MTI"
	_mtiLength = 4
)

func Decode(b []byte) (map[string][]byte, error) {
	const minLength = 4

	content := make(map[string][]byte)

	if len(b) < minLength {
		return nil, errors.New("message to short")
	}

	content[_mtiKey] = b[:_mtiLength]

	return content, nil
}
