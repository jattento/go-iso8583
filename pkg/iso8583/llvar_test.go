package iso8583_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jattento/go-iso8583/pkg/encoding/ebcdic"
	"github.com/jattento/go-iso8583/pkg/iso8583"

	"github.com/stretchr/testify/assert"
)

func TestLLVAR_MarshalISO8583(t *testing.T) {
	testList := []struct {
		Name        string
		V           iso8583.LLVAR
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
			OutputBytes: append(ebcdic.V1047.FromGoString("14"), []byte("ascii_standard")...),
		},
		{
			Name:        "ebcdic_standard",
			V:           "ebcdic",
			Encoding:    "ebcdic",
			OutputError: "",
			OutputBytes: append(ebcdic.V1047.FromGoString("06"), ebcdic.V1047.FromGoString("ebcdic")...),
		},
		{
			Name: "too_long",
			V: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890" +
				"1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890" +
				"1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			Encoding:    "ascii",
			OutputError: "content length exceeded the 99 limit for LL var elements",
			OutputBytes: nil,
		},
		{
			Name:        "ll_encoding_error",
			V:           "11111",
			Encoding:    "force_error",
			OutputError: "forced_error",
			OutputBytes: nil,
		},
		{
			Name:        "var_encoding_error",
			V:           "0000",
			Encoding:    "force_error",
			OutputError: "forced_error",
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

func TestLLVAR_UnmarshalISO8583(t *testing.T) {
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
			Name:          "too_short",
			InputEncoding: "ascii",
			InputLength:   10,
			OutputContent: "",
			OutputError:   "message remain (3 bytes) is shorter than LL byte length (10)",
			InputBytes:    []byte("asd"),
			ExpectedRead:  0,
		},
		{
			Name:          "var_decoding_error",
			InputEncoding: "force_error",
			InputLength:   1,
			OutputContent: "",
			OutputError:   "encoder 'force_error' returned error: forced_error",
			InputBytes:    []byte("10234"),
			ExpectedRead:  0,
		},
		{
			Name:          "ll_decoding_error",
			InputEncoding: "force_error",
			InputLength:   2,
			OutputContent: "",
			OutputError:   "encoder 'force_error' returned error: forced_error",
			InputBytes:    []byte("10234"),
			ExpectedRead:  0,
		},
		{
			Name:          "length_not_numeric_error",
			InputEncoding: "ascii",
			InputLength:   1,
			OutputContent: "",
			OutputError:   "obtained LL after decoding is not a valid integer: a",
			InputBytes:    []byte("a0234"),
			ExpectedRead:  0,
		},
		{
			Name:          "too_short_value_content",
			InputEncoding: "ascii",
			InputLength:   2,
			OutputContent: "",
			OutputError:   "message remain (4 bytes) is shorter than LL indicated length (99)",
			InputBytes:    []byte("990234"),
			ExpectedRead:  0,
		},
		{
			Name:          "nil_bytes_input",
			InputEncoding: "ascii",
			InputLength:   2,
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
		return iso8583.MarshalEncodings["ascii"](bytes)
	}
	defer delete(iso8583.UnmarshalDecodings, "force_error")

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("var_to_bytes_%s", testCase.Name), func(t *testing.T) {
			var v iso8583.LLVAR

			n, err := v.UnmarshalISO8583(testCase.InputBytes, testCase.InputLength, testCase.InputEncoding)
			if testCase.OutputError != "" {
				if assert.NotNil(t, err) {
					assert.Equal(t, err.Error(), testCase.OutputError)
				}
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
