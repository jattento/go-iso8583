package iso8583

import (
	"errors"
	"math"

	"github.com/jattento/go-iso8583/pkg/bitmap"
)

// BITMAP wrapps the bitmap.bitmap type to match the iso8583.MarshalerBitmap,
// and iso8583.UnmarshalerBitmap interfaces.
type BITMAP struct {
	bitmap.Bitmap
}

// UnmarshalISO8583 wrapps bitmap.FromBytes to match iso8583.Unmarshal interface.
func (b *BITMAP) UnmarshalISO8583(byt []byte, length int, encoding string) (int, error) {
	const bitsInByte = 8

	if byt == nil {
		return 0, errors.New("bytes input is nil")
	}

	bcap := int(math.Ceil(float64(length) / float64(bitsInByte)))
	b.Bitmap = bitmap.FromBytes(byt[:bcap])
	return bcap, nil
}

// MarshalISO8583 wrapps bitmap.ToBytes to match iso8583.Marshal interface.
func (b BITMAP) MarshalISO8583(length int, encoding string) ([]byte, error) {
	return bitmap.ToBytes(b.Bitmap), nil
}

// Bits returns which bits are on, key values are between 1 and 64, both included.
// First value is bit 1.
func (b BITMAP) Bits() (map[int]bool, error) {
	return b.Bitmap, nil
}

// MarshalISO8583Bitmap returns a empty slice if all bytes are 0x0
func (b BITMAP) MarshalISO8583Bitmap(m map[int]bool, encoding string) ([]byte, error) {
	bytes := bitmap.ToBytes(m)
	for _, b := range bytes {
		if b != 0x0 {
			// Only if some byte has information, they are returned
			return bytes, nil
		}
	}

	return []byte{}, nil
}
