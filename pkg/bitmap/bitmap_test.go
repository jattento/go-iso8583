package bitmap_test

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/jattento/go-iso8583/pkg/bitmap"

	"github.com/stretchr/testify/assert"
)

func TestISO8583FromBytes(t *testing.T) {
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
		t.Run(fmt.Sprintf("bitmap_ISO8583FromBytes_%s", testCase.name), func(t *testing.T) {
			byt := binaryStringToBytes(t, testCase.binaryInput)

			availableElements, nextElement, err := bitmap.ISO8583FromBytes(byt, testCase.bitmapPosition)
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

func TestISO8583ToBytes(t *testing.T) {
	testList := []struct {
		name            string
		mapInput        bitmap.Bitmap
		nextBitmapInput bool
		expectedError   error
		expectedOutput  string
	}{
		{
			name:           "all_set",
			expectedOutput: "11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111",
			expectedError:  nil,
			mapInput: map[int]bool{2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true,
				10: true, 11: true, 12: true, 13: true, 14: true, 15: true, 16: true, 17: true, 18: true, 19: true,
				20: true, 21: true, 22: true, 23: true, 24: true, 25: true, 26: true, 27: true, 28: true, 29: true,
				30: true, 31: true, 32: true, 33: true, 34: true, 35: true, 36: true, 37: true, 38: true, 39: true,
				40: true, 41: true, 42: true, 43: true, 44: true, 45: true, 46: true, 47: true, 48: true, 49: true,
				50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 58: true, 59: true,
				60: true, 61: true, 62: true, 63: true, 64: true},
			nextBitmapInput: true,
		},
		{
			name:           "all_set_no_next_bitmap",
			expectedOutput: "01111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111",
			expectedError:  nil,
			mapInput: map[int]bool{2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true,
				10: true, 11: true, 12: true, 13: true, 14: true, 15: true, 16: true, 17: true, 18: true, 19: true,
				20: true, 21: true, 22: true, 23: true, 24: true, 25: true, 26: true, 27: true, 28: true, 29: true,
				30: true, 31: true, 32: true, 33: true, 34: true, 35: true, 36: true, 37: true, 38: true, 39: true,
				40: true, 41: true, 42: true, 43: true, 44: true, 45: true, 46: true, 47: true, 48: true, 49: true,
				50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 58: true, 59: true,
				60: true, 61: true, 62: true, 63: true, 64: true},
			nextBitmapInput: false,
		},
		{
			name:           "last_bit_of_byte_off",
			expectedOutput: "11111110 11111110 11111110 11111110 11111110 11111110 11111110 11111110",
			expectedError:  nil,
			mapInput: map[int]bool{2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: true,
				10: true, 11: true, 12: true, 13: true, 14: true, 15: true, 16: false, 17: true, 18: true, 19: true,
				20: true, 21: true, 22: true, 23: true, 24: false, 25: true, 26: true, 27: true, 28: true, 29: true,
				30: true, 31: true, 32: false, 33: true, 34: true, 35: true, 36: true, 37: true, 38: true, 39: true,
				40: false, 41: true, 42: true, 43: true, 44: true, 45: true, 46: true, 47: true, 48: false, 49: true,
				50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: false, 57: true, 58: true, 59: true,
				60: true, 61: true, 62: true, 63: true, 64: false},
			nextBitmapInput: true,
		},
		{
			name:            "only_130_element",
			expectedOutput:  "01000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
			expectedError:   nil,
			mapInput:        map[int]bool{130: true},
			nextBitmapInput: false,
		},
		{
			name:            "only_192_element",
			expectedOutput:  "00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001",
			expectedError:   nil,
			mapInput:        map[int]bool{192: true},
			nextBitmapInput: false,
		},
		{
			name:          "error_impossible_bitmap",
			mapInput:      map[int]bool{1: true, 100: true},
			expectedError: bitmap.ErrBitmapISOImpossibleBitmap,
		},
		{
			name:          "error_first_bit",
			mapInput:      map[int]bool{1: true, 2: true, 3: true},
			expectedError: bitmap.ErrBitmapISOFirstBitProhibited,
		},
		{
			name:            "extremities",
			expectedOutput:  "01000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001",
			expectedError:   nil,
			mapInput:        map[int]bool{130: true, 192: true},
			nextBitmapInput: false,
		},
		{
			name:            "extremities_1_bit_out",
			expectedError:   bitmap.ErrBitmapISOImpossibleBitmap,
			mapInput:        map[int]bool{130: true, 193: true},
			nextBitmapInput: false,
		},
		{
			name:            "empty_map_no_next_bitmap",
			expectedOutput:  "00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
			expectedError:   nil,
			mapInput:        map[int]bool{},
			nextBitmapInput: false,
		},
		{
			name:            "empty_map_with_next_bitmap",
			expectedOutput:  "10000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
			expectedError:   nil,
			mapInput:        map[int]bool{},
			nextBitmapInput: true,
		},
	}

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("bitmap_ISO8583ToBytes_%s", testCase.name), func(t *testing.T) {
			b, err := bitmap.ISO8583ToBytes(testCase.mapInput, testCase.nextBitmapInput)
			if testCase.expectedError != nil {
				assert.True(t, errors.Is(err, testCase.expectedError))
			} else {
				if !assert.Nil(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, testCase.expectedOutput,
				strings.TrimSuffix(strings.TrimPrefix(fmt.Sprintf("%08b ", b), "["), "] "))
		})
	}
}

func TestToBytes_FromBytes(t *testing.T) {
	testList := []struct {
		name         string
		bmap         bitmap.Bitmap
		binaryString string
	}{
		{
			name:         "8_byte_all_set",
			binaryString: "11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111",
			bmap: map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true,
				10: true, 11: true, 12: true, 13: true, 14: true, 15: true, 16: true, 17: true, 18: true, 19: true,
				20: true, 21: true, 22: true, 23: true, 24: true, 25: true, 26: true, 27: true, 28: true, 29: true,
				30: true, 31: true, 32: true, 33: true, 34: true, 35: true, 36: true, 37: true, 38: true, 39: true,
				40: true, 41: true, 42: true, 43: true, 44: true, 45: true, 46: true, 47: true, 48: true, 49: true,
				50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: true, 57: true, 58: true, 59: true,
				60: true, 61: true, 62: true, 63: true, 64: true},
		},
		{
			name:         "empty",
			binaryString: "",
			bmap:         map[int]bool{},
		},
		{
			name:         "last_bit_of_byte_off",
			binaryString: "11111110 11111110 11111110 11111110 11111110 11111110 11111110 11111110",
			bmap: map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: false, 9: true,
				10: true, 11: true, 12: true, 13: true, 14: true, 15: true, 16: false, 17: true, 18: true, 19: true,
				20: true, 21: true, 22: true, 23: true, 24: false, 25: true, 26: true, 27: true, 28: true, 29: true,
				30: true, 31: true, 32: false, 33: true, 34: true, 35: true, 36: true, 37: true, 38: true, 39: true,
				40: false, 41: true, 42: true, 43: true, 44: true, 45: true, 46: true, 47: true, 48: false, 49: true,
				50: true, 51: true, 52: true, 53: true, 54: true, 55: true, 56: false, 57: true, 58: true, 59: true,
				60: true, 61: true, 62: true, 63: true, 64: false},
		},
		{
			name:         "64_bit_off",
			binaryString: "00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
			bmap:         map[int]bool{64: false},
		},
		{
			name:         "63_bit_on",
			binaryString: "00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000010",
			bmap:         map[int]bool{63: true},
		},
		{
			name:         "66_bit_off",
			binaryString: "00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000",
			bmap:         map[int]bool{66: false},
		},
		{
			name:         "9_bit_all_on",
			binaryString: "11111111 10000000",
			bmap:         map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true},
		},
		{
			name: "200_bit_on",
			binaryString: "00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 " +
				"00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 " +
				"00000000 00000000 00000000 00000000 00000001",
			bmap: map[int]bool{200: true},
		},
	}

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("bitmap_ISO8583ToBytes_%s", testCase.name), func(t *testing.T) {

			// Check ToByes...
			bitmapToBytesResult := bitmap.ToBytes(testCase.bmap)
			assert.Equal(t, testCase.binaryString,
				strings.TrimSuffix(strings.TrimPrefix(
					fmt.Sprintf("%08b ", bitmapToBytesResult), "["), "] "),
				"failed doing ToBytes transformation")

			// Check FromBytes...
			bitmapFromBytesResult := bitmap.FromBytes(bitmapToBytesResult)
			assertBitmap(t, testCase.bmap, bitmapFromBytesResult, "failed doing FromBytes transformation")
		})
	}
}

func assertBitmap(t *testing.T, expected, actual bitmap.Bitmap, msgAndArgs ...interface{}) bool {
	defer func() {
		if t.Failed() {
			t.Logf("%v", msgAndArgs)
		}
	}()

	if assert.Equal(&testing.T{}, expected, actual) {
		return true
	}

	for k, v := range actual {
		if expected[k] != v {
			t.Logf("actual bitmap key %v is on but not in expected", k)
		}
	}

	for k, v := range expected {
		if actual[k] != v {
			t.Logf("expected bitmap key %v is on but not in expected", k)
		}
	}

	return true
}

func binaryStringToBytes(t *testing.T, s string) []byte {
	const (
		binaryBase = 2
		intBitSize = 0
	)

	if s == "" {
		return []byte{}
	}

	var byt []byte
	for _, b := range strings.Split(s, " ") {
		i, err := strconv.ParseInt(b, binaryBase, intBitSize)
		if err != nil {
			t.Fatal("bad test input, cant parse binary string: ", err)
		}

		byt = append(byt, byte(i))
	}

	return byt
}
