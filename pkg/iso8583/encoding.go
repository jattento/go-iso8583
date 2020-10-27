package iso8583

import (
	"github.com/jattento/go-iso8583/pkg/encoding/ebcdic"
)

// UnmarshalDecodings is the unmarshal encodings map used by inbuilt ISO fields, you can append more encoding for extended functionality.
var UnmarshalDecodings = map[string]func([]byte) ([]byte, error){
	"ebcdic": errWrapper(func(bytes []byte) []byte { return []byte(ebcdic.V1047.ToGoString(bytes)) }),
	"ascii":  nop,
}

// MarshalEncodings is the marshal encodings map used by inbuilt ISO fields, you can append more encoding for extended functionality.
var MarshalEncodings = map[string]func([]byte) ([]byte, error){
	"ebcdic": errWrapper(func(bytes []byte) []byte {
		return ebcdic.V1047.FromGoString(string(bytes))
	}),
	"ascii": nop,
}

func errWrapper(Func func([]byte) []byte) func([]byte) ([]byte, error) {
	return func(bytes []byte) ([]byte, error) {
		return Func(bytes), nil
	}
}

func nop(bytes []byte) ([]byte, error) {
	return bytes, nil
}
