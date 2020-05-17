package iso8583

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// LLVAR: For use of different encoding for 'LL' and 'VAR' separate both encodings with a slash,
// where first element is the ll encoding and the second the var encoding.
// For Unmarshal length indicate the amount of byte that contain the LL value
// For example:
// 	`iso8583:"2,length:3,encoding:ascii/ebcdic"`
type LLVAR string

// MarshalISO8583 allows to use this type in structs and be able tu iso8583.Marshal it.
func (v LLVAR) MarshalISO8583(length int, enc string) ([]byte, error) {
	return lengthMarshal(2, string(v), enc)
}

// UnmarshalISO8583 allows to use this type in structs and be able tu iso8583.Unmarshal it.
func (v *LLVAR) UnmarshalISO8583(b []byte, length int, enc string) (int, error) {
	if b == nil {
		return 0, errors.New("bytes input is nil")
	}

	str, n, err := lengthUnmarshal(2, b, length, enc)
	*v = LLVAR(str)
	return n, err
}

func lengthMarshal(l int, v string, enc string) ([]byte, error) {
	varContent := []byte(v)

	llEncoding := enc
	varEncoding := enc

	if encS := strings.Split(enc, "/"); len(encS) > 1 {
		llEncoding = encS[0]
		varEncoding = encS[1]
	}

	llValue := strconv.Itoa(len(varContent))
	if len(llValue) > l {
		return nil, fmt.Errorf("content length exceeded the %s limit for %s var elements",
			strings.Repeat("9", l), strings.Repeat("L", l))
	}

	for len(llValue) < l {
		llValue = "0" + llValue
	}

	llContent, err := applyEncoding([]byte(llValue), llEncoding, MarshalEncodings)
	if err != nil {
		return nil, err
	}

	varContent, err = applyEncoding(varContent, varEncoding, MarshalEncodings)
	if err != nil {
		return nil, err
	}

	return append(llContent, varContent...), nil
}

func lengthUnmarshal(l int, b []byte, length int, enc string) (string, int, error) {
	if len(b) < length {
		return "", 0, fmt.Errorf("message remain (%v bytes) is shorter than %s byte length (%v)",
			len(b), strings.Repeat("L", l), length)
	}

	llContent := make([]byte, length)
	copy(llContent, b[:length])

	llEncoding := enc
	varEncoding := enc

	if encS := strings.Split(enc, "/"); len(encS) > 1 {
		llEncoding = encS[0]
		varEncoding = encS[1]
	}

	llContent, err := applyEncoding(llContent, llEncoding, UnmarshalDecodings)
	if err != nil {
		return "", 0, err
	}

	llValue, err := strconv.Atoi(string(llContent))
	if err != nil {
		return "", 0, fmt.Errorf("obtained %s after decoding is not a valid integer: %v",
			strings.Repeat("L", l), string(llContent))
	}

	if len(b)-length < llValue {
		return "", 0, fmt.Errorf("message remain (%v bytes) is shorter than %s indicated length (%v)",
			len(b)-length, strings.Repeat("L", l), llValue)
	}

	varContent := make([]byte, llValue)
	copy(varContent, b[length:length+llValue])

	varContent, err = applyEncoding(varContent, varEncoding, UnmarshalDecodings)
	if err != nil {
		return "", 0, err
	}

	return strings.TrimFunc(string(varContent), func(r rune) bool {
		return !unicode.IsGraphic(r)
	}), length + llValue, nil
}
