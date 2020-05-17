package ebcdic_test

import (
	"github.com/jattento/go-iso8583/pkg/encoding/ebcdic"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncoding_Translate(t *testing.T) {
	e := "qwertyuiop1234567890asdfghjklñ´çzxcvbnm,.-`+"
	if n := ebcdic.V1047.ToGoString(ebcdic.V1047.FromGoString(e)); n != e {
		t.Fatal("example test isnt the same encoded and decoded: ", n)
	}
}

func TestEncoding_Translate_unknown_rune(t *testing.T) {
	e := "鲸鱼歌"
	assert.Equal(t, []byte{ebcdic.NULL, ebcdic.NULL, ebcdic.NULL}, ebcdic.V1047.FromGoString(e))
}
