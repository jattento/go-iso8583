package iso8583_test

import (
	"fmt"
	"testing"

	"github.com/jattento/go-iso8583/pkg/encoding/ebcdic"
	"github.com/jattento/go-iso8583/pkg/iso8583"

	"github.com/stretchr/testify/assert"
)

func TestVAR_MarshalISO8583(t *testing.T) {
	testList := []struct {
		Name        string
		V           iso8583.VAR
		Encoding    string
		Length      int
		OutputBytes []byte
		OutputError string
	}{
		{
			Name:        "ascii_standard",
			V:           "ascii_standard",
			Encoding:    "ascii",
			OutputError: "",
			OutputBytes: []byte("ascii_standard"),
		},
		{
			Name:        "ebcdic_standard",
			V:           "ebcdic",
			Encoding:    "ebcdic",
			OutputError: "",
			OutputBytes: ebcdic.V1047.FromASCII([]byte("ebcdic")),
		},
	}

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("var_to_bytes_%s", testCase.Name), func(t *testing.T) {
			o, err := testCase.V.MarshalISO8583(testCase.Length, testCase.Encoding)
			if testCase.OutputError != "" {
				assert.Errorf(t, err, testCase.OutputError)
			} else {
				if !assert.Nil(t, err) {
					t.FailNow()
				}
			}
			assert.Equal(t, testCase.OutputBytes, o)
		})
	}
}

func TestVAR_UnmarshalISO8583(t *testing.T) {
	testList := []struct {
		Name          string
		InputBytes    []byte
		InputEncoding string
		InputLength   int
		OutputContent string
		OutputError   string
	}{
		{
			Name:          "ascii_standard",
			InputEncoding: "ascii",
			InputLength:   14,
			OutputContent: "ascii_standard",
			OutputError:   "",
			InputBytes:    []byte("ascii_standard"),
		},
		{
			Name:          "ebcdic_standard",
			InputEncoding: "ebcdic",
			InputLength:   6,
			OutputContent: "ebcdic",
			OutputError:   "",
			InputBytes:    ebcdic.V1047.FromASCII([]byte("ebcdic")),
		},
	}

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("var_to_bytes_%s", testCase.Name), func(t *testing.T) {
			var v iso8583.VAR

			_, err := v.UnmarshalISO8583(testCase.InputBytes, testCase.InputLength, testCase.InputEncoding)
			if testCase.OutputError != "" {
				assert.Errorf(t, err, testCase.OutputError)
			} else {
				if !assert.Nil(t, err) {
					t.FailNow()
				}
			}
			assert.Equal(t, testCase.OutputContent, string(v))
		})
	}
}
