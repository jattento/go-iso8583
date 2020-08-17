package iso8583util_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jattento/go-iso8583/pkg/iso8583util"

	"github.com/stretchr/testify/assert"
)

func TestStringBitsToBytes(t *testing.T) {
	testList := []struct {
		input string
	}{
		{
			input: "00000001 00000001 00000001",
		},
		{
			input: "00000001 00000001 00000001 1",
		},
		{
			input: "00000001 00000001 00000001 1000",
		},
		{
			input: "00010001 10000001 00010001 1000",
		},
		{
			input: "00010101100000010010001 1000",
		},
		{
			input: "00010101100000010010001 1000",
		},
		{
			input: "00010101100000010010001 1000",
		},
		{
			input: "00010101100000010010001 1000",
		},
		{
			input: "1111111111111111111111101111111111111111011111111111111111111111011111111111111111111111111111111",
		},
	}

	for _, testCase := range testList {
		t.Run("string_bits_to_bytes", func(t *testing.T) {
			expected := strings.ReplaceAll(testCase.input, " ", "")
			for len(expected)%8 != 0 {
				expected += "0"
			}
			assert.Equal(t, expected, strings.TrimSuffix(strings.TrimPrefix(strings.ReplaceAll(
				fmt.Sprintf("%08b", iso8583util.StringBitsToBytes(testCase.input)),
				" ", ""), "["), "]"))
		})
	}
}

func TestBytesToStringBits(t *testing.T) {
	testList := []struct {
		input []byte
	}{
		{
			input: []byte{12, 3, 4, 5, 6, 76, 87, 89, 9, 6, 5, 4, 1, 1},
		},
		{
			input: []byte{200},
		},
		{
			input: []byte{128, 128, 128},
		},
		{
			input: []byte{1, 1, 1, 1, 1, 1},
		},
		{
			input: []byte{221},
		},
		{
			input: []byte{12, 101, 51, 41, 57, 71, 211, 200, 219},
		},
	}

	for _, testCase := range testList {
		t.Run("string_bits_to_bytes", func(t *testing.T) {
			assert.Equal(t, strings.TrimSuffix(strings.TrimPrefix(
				fmt.Sprintf("%08b", testCase.input), "["), "]"),
				iso8583util.BytesToStringBits(testCase.input, true))
			assert.Equal(t, strings.ReplaceAll(strings.TrimSuffix(strings.TrimPrefix(
				fmt.Sprintf("%08b", testCase.input), "["), "]"), " ", ""),
				iso8583util.BytesToStringBits(testCase.input, false))
		})
	}
}

func TestStringBitsToBytes_Panic(t *testing.T) {
	assert.Panics(t, func() { iso8583util.StringBitsToBytes("whale_song") })
}
