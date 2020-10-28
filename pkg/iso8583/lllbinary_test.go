package iso8583_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jattento/go-iso8583/pkg/iso8583"
)

func TestLLLBINARY_MarshalISO8583(t *testing.T) {
	testList := []struct {
		Name        string
		V           iso8583.LLLBINARY
		Encoding    string
		Length      int
		OutputBytes []byte
		OutputError string
	}{
		{
			Name:        "ascii",
			V:           []byte("text"),
			Encoding:    "ascii",
			OutputError: "",
			OutputBytes: []byte("004text"),
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
		t.Run(fmt.Sprintf("lllbinary_to_bytes_%s", testCase.Name), func(t *testing.T) {
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

func TestLLLBINARY_UnmarshalISO8583(t *testing.T) {
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
			InputLength:   3,
			OutputContent: "ascii_standard",
			OutputError:   "",
			InputBytes:    []byte("014ascii_standard"),
			ExpectedRead:  17,
		},
		{
			Name:          "lll_encoding_error",
			InputEncoding: "force_error/ascii",
			InputLength:   1,
			OutputContent: "text",
			OutputError:   "encoder 'force_error' returned error: forced_error",
			InputBytes:    []byte{1, 2, 3},
			ExpectedRead:  7,
		},
		{
			Name:          "nil_bytes_input",
			InputEncoding: "ascii",
			InputLength:   3,
			OutputContent: "",
			OutputError:   "bytes input is nil",
			InputBytes:    nil,
			ExpectedRead:  0,
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
		t.Run(fmt.Sprintf("lllbinary_to_bytes_%s", testCase.Name), func(t *testing.T) {
			v := iso8583.LLLBINARY("")

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
