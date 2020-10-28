package iso8583_test

import (
	"github.com/jattento/go-iso8583/pkg/iso8583"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBINARY_UnmarshalISO8583_too_short_input(t *testing.T) {
	var binary iso8583.BINARY

	n, bmapErr := binary.UnmarshalISO8583([]byte{1, 1, 1}, 64, "ascii")

	assert.Equal(t, 0, n)
	if assert.NotNil(t, bmapErr) {
		assert.Equal(t, bmapErr.Error(), "message remain (3 bytes) is shorter than indicated length: 64")
	}
}

func TestBINARY_UnmarshalISO8583_nil_input(t *testing.T) {
	var binary iso8583.BINARY

	n, bmapErr := binary.UnmarshalISO8583(nil, 64, "ascii")

	assert.Equal(t, 0, n)
	if assert.NotNil(t, bmapErr) {
		assert.Equal(t, bmapErr.Error(), "bytes input is nil")
	}
}
