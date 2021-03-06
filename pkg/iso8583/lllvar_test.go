package iso8583_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jattento/go-iso8583/pkg/encoding/ebcdic"
	"github.com/jattento/go-iso8583/pkg/iso8583"

	"github.com/stretchr/testify/assert"
)

func TestLLLVAR_MarshalISO8583(t *testing.T) {
	testList := []struct {
		Name        string
		V           iso8583.LLLVAR
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
			OutputBytes: append(ebcdic.V1047.FromGoString("014"), []byte("ascii_standard")...),
		},
		{
			Name:        "ebcdic_standard",
			V:           "ebcdic",
			Encoding:    "ebcdic",
			OutputError: "",
			OutputBytes: append(ebcdic.V1047.FromGoString("006"), ebcdic.V1047.FromGoString("ebcdic")...),
		},
		{
			Name:        "encoding_fail",
			V:           "text",
			Encoding:    "force_error",
			OutputError: "force_error",
			OutputBytes: nil,
		},
	}

	iso8583.MarshalEncodings["force_error"] = func(bytes []byte) ([]byte, error) {
		if string(bytes) != "04" {
			return nil, errors.New("forced_error")
		}
		return iso8583.MarshalEncodings["ascii"](bytes)
	}
	defer delete(iso8583.MarshalEncodings, "force_error")

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("lllvar_to_bytes_%s", testCase.Name), func(t *testing.T) {
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
			InputBytes:    append(ebcdic.V1047.FromGoString("6"), []byte("ebcdic")...),
			ExpectedRead:  7,
		},
		{
			Name:          "nil_bytes_error",
			InputEncoding: "ascii",
			InputLength:   1,
			OutputContent: "ebcdic",
			OutputError:   "bytes input is nil",
			InputBytes:    nil,
			ExpectedRead:  7,
		},
		{
			Name:          "ll_encoding_error",
			InputEncoding: "force_error/ascii",
			InputLength:   1,
			OutputContent: "text",
			OutputError:   "encoder 'force_error' returned error: forced_error",
			InputBytes:    []byte{1, 2, 3},
			ExpectedRead:  7,
		},
		{
			Name:          "var_encoding_error",
			InputEncoding: "ascii/force_error",
			InputLength:   1,
			OutputContent: "text",
			OutputError:   "encoder 'force_error' returned error: forced_error",
			InputBytes:    []byte("123123"),
			ExpectedRead:  7,
		},
	}

	iso8583.UnmarshalDecodings["force_error"] = func(bytes []byte) ([]byte, error) {
		if string(bytes) != "1" {
			return nil, errors.New("forced_error")
		}
		return iso8583.UnmarshalDecodings["ascii"](bytes)
	}
	defer delete(iso8583.UnmarshalDecodings, "force_error")

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("lllvar_to_bytes_%s", testCase.Name), func(t *testing.T) {
			var v iso8583.LLLVAR

			n, err := v.UnmarshalISO8583(testCase.InputBytes, testCase.InputLength, testCase.InputEncoding)
			if testCase.OutputError != "" {
				if assert.NotNil(t, err) {
					assert.Equal(t, err.Error(), testCase.OutputError)
				}
				return
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
