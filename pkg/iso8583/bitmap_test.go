package iso8583_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jattento/go-iso8583/pkg/bitmap"
	"github.com/jattento/go-iso8583/pkg/iso8583"
)

func TestBITMAP_MarshalISO8583(t *testing.T) {
	bmap := iso8583.BITMAP{
		Bitmap: bitmap.Bitmap{
			1: true,
			2: true,
			3: true,
		},
	}

	bmapMarshaled, bmapErr := bmap.MarshalISO8583(0, "")

	assert.Equal(t, bitmap.ToBytes(bmap.Bitmap), bmapMarshaled)
	assert.Nil(t, bmapErr)
}

func TestBITMAP_UnmarshalISO8583_too_short_input(t *testing.T) {
	var bmap iso8583.BITMAP

	n, bmapErr := bmap.UnmarshalISO8583([]byte{1, 1, 1}, 64, "ascii")

	assert.Equal(t, 0, n)
	if assert.NotNil(t, bmapErr) {
		assert.Equal(t, bmapErr.Error(), "bitmap should be 8 bytes long but only 3 bytes are avaiable")
	}
}

func TestBITMAP_UnmarshalISO8583_nil_input(t *testing.T) {
	var bmap iso8583.BITMAP

	n, bmapErr := bmap.UnmarshalISO8583(nil, 64, "ascii")

	assert.Equal(t, 0, n)
	if assert.NotNil(t, bmapErr) {
		assert.Equal(t, bmapErr.Error(), "bytes input is nil")
	}
}
