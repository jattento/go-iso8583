package bitmap_test

import (
	"errors"
	"fmt"
	"github.com/iso-lib/pkg/bitmap"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitmap_Read(t *testing.T) {
	const (
		binaryBase = 2
		intBitSize = 0
	)

	testList := []struct {
		name               string
		binaryInput        string
		bitmapPosition     int
		expectedError      error
		expectedOutput     map[int]bool
		expectedNextBinary bool
	}{
		{
			name:           "first_bitmap_all_set",
			binaryInput:    "11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111",
			bitmapPosition: 1,
			expectedError:  nil,
			expectedOutput: map[int]bool{2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true,
				10: true, 11: true, 12: true, 13: true, 14: true, 15: true, 16: true, 17: true, 18: true, 19: true,
				20: true, 21: true, 22: true, 23: true, 24: true, 25: true, 26: true, 27: true, 28: true, 29: true,
				30: true, 31: true, 32: true, 33: true, 34: true, 35: true, 36: true, 37: true, 38: true, 39: true,
				40: true, 41: true, 42: true, 43: true, 44: true, 45: true, 46: true, 47: true, 48: true, 49: true,
				50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 58: true, 59: true,
				60: true, 61: true, 62: true, 63: true, 64: true},
			expectedNextBinary: true,
		},
		{
			name:           "second_bitmap_all_set",
			binaryInput:    "11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111",
			bitmapPosition: 2,
			expectedError:  nil,
			expectedOutput: map[int]bool{66: true, 67: true, 68: true, 69: true, 70: true, 71: true, 72: true, 73: true,
				74: true, 75: true, 76: true, 77: true, 78: true, 79: true, 80: true, 81: true, 82: true, 83: true,
				84: true, 85: true, 86: true, 87: true, 88: true, 89: true, 90: true, 91: true, 92: true, 93: true,
				94: true, 95: true, 96: true, 97: true, 98: true, 99: true, 100: true, 101: true, 102: true, 103: true,
				104: true, 105: true, 106: true, 107: true, 108: true, 109: true, 110: true, 111: true, 112: true,
				113: true, 114: true, 115: true, 116: true, 117: true, 118: true, 119: true, 120: true, 121: true,
				122: true, 123: true, 124: true, 125: true, 126: true, 127: true, 128: true},
			expectedNextBinary: true,
		},
		{
			name:           "first_bitmap_all_set_no_next_bitmap",
			binaryInput:    "01111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111",
			bitmapPosition: 1,
			expectedError:  nil,
			expectedOutput: map[int]bool{2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true,
				10: true, 11: true, 12: true, 13: true, 14: true, 15: true, 16: true, 17: true, 18: true, 19: true,
				20: true, 21: true, 22: true, 23: true, 24: true, 25: true, 26: true, 27: true, 28: true, 29: true,
				30: true, 31: true, 32: true, 33: true, 34: true, 35: true, 36: true, 37: true, 38: true, 39: true,
				40: true, 41: true, 42: true, 43: true, 44: true, 45: true, 46: true, 47: true, 48: true, 49: true,
				50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 58: true, 59: true,
				60: true, 61: true, 62: true, 63: true, 64: true},
			expectedNextBinary: false,
		},
		{
			name:           "first_bitmap_last_bit_of_byte_off",
			binaryInput:    "11111110 11111110 11111110 11111110 11111110 11111110 11111110 11111110",
			bitmapPosition: 1,
			expectedError:  nil,
			expectedOutput: map[int]bool{2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: true,
				10: true, 11: true, 12: true, 13: true, 14: true, 15: true, 16: false, 17: true, 18: true, 19: true,
				20: true, 21: true, 22: true, 23: true, 24: false, 25: true, 26: true, 27: true, 28: true, 29: true,
				30: true, 31: true, 32: false, 33: true, 34: true, 35: true, 36: true, 37: true, 38: true, 39: true,
				40: false, 41: true, 42: true, 43: true, 44: true, 45: true, 46: true, 47: true, 48: false, 49: true,
				50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: false, 57: true, 58: true, 59: true,
				60: true, 61: true, 62: true, 63: true, 64: false},
			expectedNextBinary: true,
		},
		{
			name:           "error_too_short_input",
			binaryInput:    "11111110",
			bitmapPosition: 1,
			expectedError:  bitmap.ErrBitmapISOWrongLength,
		},
		{
			name:           "error_bad_bitmap-position",
			binaryInput:    "11111110 11111110 11111110 11111110 11111110 11111110 11111110 11111110",
			bitmapPosition: 0,
			expectedError:  bitmap.ErrBitmapISOBadBitmapPosition,
		},
	}

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("bitmap_read_%s", testCase.name), func(t *testing.T) {
			binaries := strings.Split(testCase.binaryInput, " ")

			var readInput []byte
			for _, b := range binaries {
				i, err := strconv.ParseInt(b, binaryBase, intBitSize)
				if err != nil {
					t.Fatal("bad test input, cant parse binary string: ", err)
				}

				readInput = append(readInput, byte(i))
			}

			availableElements, nextElement, err := bitmap.ISO8583FromBytes(readInput, testCase.bitmapPosition)
			if testCase.expectedError != nil {
				assert.True(t, errors.Is(err, testCase.expectedError))
			} else {
				if !assert.Nil(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, testCase.expectedNextBinary, nextElement)
			assert.Equal(t, testCase.expectedOutput, availableElements)
		})
	}
}
