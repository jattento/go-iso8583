package bitmap

import (
	"errors"
	"fmt"
)

type Bitmap = map[int]bool

const (
	_bitmapLength = 8
	_bitsInByte   = 8
)

var (
	ErrBitmapISOWrongLength       = errors.New("wrong bitmap length input")
	ErrBitmapISOBadBitmapPosition = errors.New("bad bitmap position input")
)

// ISO8583FromBytes indicates which elements of a ISO8583 message are present.
// It receives a 8 byte long ISO8583 bitmap with a position (to indicate if its the first or second).
// Returns a map[int]bool to allow searching by element.
func ISO8583FromBytes(b []byte, bitmapPosition int) (presentElements Bitmap, nextBitmapPresent bool, returnErr error) {
	const nextBitmapIndicator = 1

	// Validate input
	if len(b) != _bitmapLength {
		return nil, false,
			fmt.Errorf("%w: should be %v, but its %v", ErrBitmapISOWrongLength, _bitmapLength, len(b))
	}

	if bitmapPosition < 1 {
		return nil, false,
			fmt.Errorf("%w: should not be lower than 1, but its %v", ErrBitmapISOBadBitmapPosition,bitmapPosition)
	}

	rawBitmap := FromBytes(b)
	isoBitmap := make(Bitmap)

	for k, v := range rawBitmap {
		// Next bitmap indicator is returned separately in the second return value
		if k != nextBitmapIndicator {
			isoBitmap[_bitsInByte*_bitmapLength*(bitmapPosition-1)+k] = v
		}
	}

	return isoBitmap, rawBitmap[nextBitmapIndicator], nil
}

// FromBytes given bytes it returns a map[int]bool indicating which biy is on or off. Most left is 1.
func FromBytes(b []byte) Bitmap {
	const (
		firstByteIndex = 0
		lastByteIndex  = 7
		firstBitOffset = 7
		lastBitOffset  = 0
	)

	availableElements := make(Bitmap)

	// Iterate over each bit of the input from most left to most right
	for bytePosition := firstByteIndex; bytePosition <= lastByteIndex; bytePosition++ {
		for bitOffset := firstBitOffset; bitOffset >= lastBitOffset; bitOffset-- {
			// Calculate element position and save in map
			previousBytesSummary := _bitsInByte * bytePosition
			bitPosition := _bitsInByte - bitOffset

			availableElements[bitPosition+previousBytesSummary] = hasBitSet(b[bytePosition], uint(bitOffset))
		}
	}

	return availableElements
}

func hasBitSet(n byte, pos uint) bool {
	val := n & (1 << pos)
	return val > 0
}
