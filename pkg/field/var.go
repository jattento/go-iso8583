package field

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/jattento/go-iso8583/pkg/field/encoding"
)

// VAR type should be used for fixed length fields.
type VAR string

// MarshalISO8583 allows to use this type in structs and be able tu iso8583.Marshal it.
func (v VAR) MarshalISO8583(length int, enc string) ([]byte, error) {
	content := []byte(v)

	content, err := applyEncoding(content, enc, encoding.Marshal)
	if err != nil{
		return nil, err
	}

	return content, nil
}

// UnmarshalISO8583 allows to use this type in structs and be able tu iso8583.Unmarshal it. //TODO CHECK LENGTH
func (v *VAR) UnmarshalISO8583(b []byte, length int, enc string) (int, error) {
	byt := make([]byte, length)
	copy(byt, b[:length])

	byt, err := applyEncoding(byt, enc, encoding.Unmarshal)
	if err != nil{
		return 0, err
	}

	*v = VAR(strings.TrimFunc(string(byt), func(r rune) bool {
		return !unicode.IsGraphic(r)
	}))

	return length, nil
}

func applyEncoding(bytes []byte, enc string, encodings map[string]func([]byte) ([]byte, error)) ([]byte, error) {
	b := make([]byte, len(bytes))
	copy(b, bytes)

	if enc != "" {
		encoder, exist := encodings[enc]
		if !exist {
			return nil, fmt.Errorf("encoder '%s' does not exist", enc)
		}

		var err error
		b, err = encoder(b)
		if err != nil {
			return nil, fmt.Errorf("encoder '%s' returned error: %w", enc, err)
		}
	}

	return b, nil
}
