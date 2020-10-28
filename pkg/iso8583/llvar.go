package iso8583

import (
	"errors"
)

// LLVAR field type.
// For use of different encoding for 'LL' and 'VAR' separate both encodings with a slash,
// where first element is the ll encoding and the second the var encoding.
// For Unmarshal length indicate the amount of byte that contain the LL value
// For example:
// 	`iso8583:"2,length:3,encoding:ascii/ebcdic"`
type LLVAR string

// MarshalISO8583 allows to use this type in structs and be able tu iso8583.Marshal it.
func (v LLVAR) MarshalISO8583(length int, enc string) ([]byte, error) {
	content := []byte(v)

	llEncoding, varEncoding := ReadSplitEncodings(enc)

	content, err := applyEncoding(content, varEncoding, MarshalEncodings)
	if err != nil {
		return nil, err
	}

	return LengthMarshal(2, content, llEncoding)
}

// UnmarshalISO8583 allows to use this type in structs and be able tu iso8583.Unmarshal it.
func (v *LLVAR) UnmarshalISO8583(b []byte, length int, enc string) (int, error) {
	if b == nil {
		return 0, errors.New("bytes input is nil")
	}

	llEncoding, varEncoding := ReadSplitEncodings(enc)

	n, b, err := LengthUnmarshal(2, b, length, llEncoding)
	if err != nil {
		return 0, err
	}

	b, err = applyEncoding(b, varEncoding, UnmarshalDecodings)
	if err != nil {
		return 0, err
	}

	*v = LLVAR(b)
	return n, err
}
