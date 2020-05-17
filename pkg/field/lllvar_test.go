package field_test

import (
	"fmt"
	"testing"

	"github.com/jattento/go-iso8583/pkg/field"
	"github.com/jattento/go-iso8583/pkg/field/encoding/ebcdic"

	"github.com/stretchr/testify/assert"
)

func TestLLLVAR_MarshalISO8583(t *testing.T) {
	testList := []struct {
		Name        string
		V           field.LLLVAR
		Encoding    string
		Length      int
		OutputBytes []byte
		OutputError string
	}{
		{
			Name:        "ebcdic_ascii",
			V:           "ascii_standard",
			Encoding:    "ebcdic/ascii",
			OutputError: "",
			OutputBytes: append(ebcdic.Encode([]byte("014")), []byte("ascii_standard")...),
		},
		{
			Name:        "ebcdic_standard",
			V:           "ebcdic",
			Encoding:    "ebcdic",
			OutputError: "",
			OutputBytes: append(ebcdic.Encode([]byte("006")), ebcdic.Encode([]byte("ebcdic"))...),
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

func TestLLLVAR_UnmarshalISO8583(t *testing.T) {
	testList := []struct {
		Name          string
		InputBytes    []byte
		InputEncoding string
		InputLength   int
		OutputContent string
		OutputError   string
		ExpectedRead  int
	}{
		{
			Name:          "ascii_standard",
			InputEncoding: "ascii",
			InputLength:   2,
			OutputContent: "ascii_standard",
			OutputError:   "",
			InputBytes:    []byte("14ascii_standard"),
			ExpectedRead:  16,
		},
		{
			Name:          "ebcdic_ascii",
			InputEncoding: "ebcdic/ascii",
			InputLength:   1,
			OutputContent: "ebcdic",
			OutputError:   "",
			InputBytes:    append(ebcdic.Encode([]byte("6")), []byte("ebcdic")...),
			ExpectedRead:  7,
		},
	}

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("var_to_bytes_%s", testCase.Name), func(t *testing.T) {
			var v field.LLLVAR

			n, err := v.UnmarshalISO8583(testCase.InputBytes, testCase.InputLength, testCase.InputEncoding)
			if testCase.OutputError != "" {
				assert.Errorf(t, err, testCase.OutputError)
			} else {
				if !assert.Nil(t, err) {
					t.FailNow()
				}
			}
			assert.Equal(t, testCase.OutputContent, string(v))
			assert.Equal(t, testCase.ExpectedRead, n)
		})
	}
}
