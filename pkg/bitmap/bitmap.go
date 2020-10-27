package bitmap

import (
	"errors"
	"fmt"
	"math"
)

// Bitmap is a alias for map[int]bool used for better code reading.
type Bitmap = map[int]bool

const (
	_bitmapLength = 8
	_bitsInByte   = 8

	_firstByteIndex = 0
	_firstBitOffset = 7
	_lastBitOffset  = 0
)

var (
	// ErrBitmapISOWrongLength exported error for asserting.
	ErrBitmapISOWrongLength = errors.New("wrong bitmap length input")
	// ErrBitmapISOBadBitmapPosition exported error for asserting.
	ErrBitmapISOBadBitmapPosition = errors.New("bad bitmap position input")
	// ErrBitmapISOImpossibleBitmap exported error for asserting.
	ErrBitmapISOImpossibleBitmap = errors.New("impossible generate bitmap, lowest and highest limits too far")
	// ErrBitmapISOFirstBitProhibited exported error for asserting.
	ErrBitmapISOFirstBitProhibited = errors.New("first bit can be setted manually in input")
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
			fmt.Errorf("%w: should not be lower than 1, but its %v", ErrBitmapISOBadBitmapPosition, bitmapPosition)
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

// ISO8583ToBytes creates a bitmap in byte format.
// Map key 1 must not be present.
func ISO8583ToBytes(b Bitmap, nextBitmapPresent bool) ([]byte, error) {
	// Find the highest and lowest element in map
	lowestElement, highestElement := Extremities(b)

	inferiorLimit := 1
	for checkedLimit := 1; checkedLimit < lowestElement; checkedLimit += 64 {
		inferiorLimit = checkedLimit
	}

	superiorLimit := inferiorLimit + 63

	if superiorLimit < highestElement {
		return nil, fmt.Errorf("%w: lowest limit %v (element %v), highest limit %v (element %v)",
			ErrBitmapISOImpossibleBitmap, inferiorLimit, lowestElement, superiorLimit, highestElement)
	}

	if _, exist := b[inferiorLimit]; exist {
		return nil, fmt.Errorf("%w: position %v", ErrBitmapISOFirstBitProhibited, inferiorLimit)
	}

	bmap := b
	bmap[inferiorLimit] = nextBitmapPresent

	if _, exist := bmap[superiorLimit]; !exist {
		bmap[superiorLimit] = false
	}

	byt := ToBytes(bmap)

	return byt[len(byt)-_bitmapLength:], nil
}

// FromBytes given bytes it returns a map[int]bool indicating which biy is on or off. Most left is 1.
func FromBytes(b []byte) Bitmap {
	availableElements := make(Bitmap)

	// Iterate over each bit of the input from most left to most right
	for bytePosition := _firstByteIndex; bytePosition < len(b); bytePosition++ {
		for bitOffset := _firstBitOffset; bitOffset >= _lastBitOffset; bitOffset-- {
			// Calculate element position and save in map
			previousBytesSummary := _bitsInByte * bytePosition
			bitPosition := _bitsInByte - bitOffset

			availableElements[bitPosition+previousBytesSummary] = hasBitSet(b[bytePosition], uint(bitOffset))
		}
	}

	return availableElements
}

// ToBytes creates a bitmap in []byte format. Most left is 1.
func ToBytes(b Bitmap) []byte {
	getPosition := func(byt, bit int) int { return bit + byt*_bitsInByte }

	_, highestElement := Extremities(b)

	bmap := make([]byte, int(math.Ceil(float64(highestElement)/_bitsInByte)))
	for bytePosition := _firstByteIndex; bytePosition < len(bmap); bytePosition++ {
		for bitOffset := _firstBitOffset; bitOffset >= _lastBitOffset; bitOffset-- {
			if isOn, exist := b[getPosition(bytePosition, _bitsInByte-bitOffset)]; exist && isOn {
				bmap[bytePosition] = setBit(bmap[bytePosition], uint(bitOffset))
			}
		}
	}

	return bmap
}

// Returns if the indicated bit is on.
func hasBitSet(n byte, pos uint) bool {
	val := n & (1 << pos)
	return val > 0
}

// Sets the bit at pos in the integer n.
func setBit(n byte, pos uint) byte {
	n |= 1 << pos
	return n
}

// Extremities returns the lowest and highest elements of a bitmap.
func Extremities(b Bitmap) (low, high int) {
	firstIteration := true
	for k := range b {
		if k > high {
			high = k
		}
		if k < low || firstIteration {
			low = k
		}
		firstIteration = false
	}

	return low, high
}
