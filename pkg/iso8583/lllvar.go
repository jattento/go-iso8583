package iso8583

import (
	"errors"
)

// LLLVAR field type.
// For use of different encoding for 'LLL' and 'VAR' separate both encodings with a slash,
// where first element is the lll encoding and the second the var encoding.
// For Unmarshal length indicate the amount of byte that contain the LLL value
// For example:
// 	`iso8583:"2,length:3,encoding:ascii/ebcdic"`
type LLLVAR string

// MarshalISO8583 allows to use this type in structs and be able tu iso8583.Marshal it.
func (v LLLVAR) MarshalISO8583(length int, enc string) ([]byte, error) {
	content := []byte(v)

	lllEncoding, varContent := ReadSplitEncodings(enc)

	content, err := applyEncoding(content, varContent, MarshalEncodings)
	if err != nil {
		return nil, err
	}

	return LengthMarshal(3, content, lllEncoding)
}

// UnmarshalISO8583 allows to use this type in structs and be able tu iso8583.Unmarshal it.
func (v *LLLVAR) UnmarshalISO8583(b []byte, length int, enc string) (int, error) {
	if b == nil {
		return 0, errors.New("bytes input is nil")
	}

	lllEncoding, varEncoding := ReadSplitEncodings(enc)

	n, b, err := LengthUnmarshal(3, b, length, lllEncoding)
	if err != nil {
		return 0, err
	}

	b, err = applyEncoding(b, varEncoding, UnmarshalDecodings)
	if err != nil {
		return 0, err
	}

	*v = LLLVAR(b)
	return n, err
}
