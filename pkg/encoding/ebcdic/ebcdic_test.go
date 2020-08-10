package ebcdic_test

import (
	"github.com/jattento/go-iso8583/pkg/encoding/ebcdic"
	"testing"
)

func TestEncoding_Translate(t *testing.T) {
	e := "qwertyuiop1234567890asdfghjklñ´çzxcvbnm,.-`+"
	if e != string(ebcdic.V1047.ToASCII(ebcdic.V1047.FromASCII([]byte(e)))){
		t.Fatal("example test isnt the same encoded and decoded")
	}
}