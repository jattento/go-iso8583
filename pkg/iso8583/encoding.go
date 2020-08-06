package iso8583

import (
	"github.com/jattento/go-iso8583/pkg/encoding/ebcdic"
)

var UnmarshalDecodings = map[string]func([]byte) ([]byte, error){
	"ebcdic": errWrapper(ebcdic.Decode),
	"ascii":  nop,
}

var MarshalEncodings = map[string]func([]byte) ([]byte, error){
	"ebcdic": errWrapper(ebcdic.Encode),
	"ascii":  nop,
}

func errWrapper(Func func([]byte) []byte) func([]byte) ([]byte, error) {
	return func(bytes []byte) ([]byte, error) {
		return Func(bytes), nil
	}
}

func nop(bytes []byte) ([]byte, error) {
	return bytes, nil
}
