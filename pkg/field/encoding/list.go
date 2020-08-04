package encoding

import (
	"github.com/jattento/go-iso8583/pkg/field/encoding/ebcdic"
)

var Unmarshal = map[string]func([]byte) ([]byte, error){
	"ebcdic": errWrapper(ebcdic.Decode),
	"ascii":  nop,
}

var Marshal = map[string]func([]byte) ([]byte, error){
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
