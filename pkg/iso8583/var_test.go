package iso8583_test

import (
	"errors"
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
			OutputBytes: ebcdic.V1047.FromGoString("ebcdic"),
		},
		{
			Name:        "encoding_error",
			V:           "ebcdic",
			Encoding:    "force_error",
			OutputError: "encoder 'force_error' returned error: forced_error",
			OutputBytes: nil,
		},
	}

	iso8583.MarshalEncodings["force_error"] = func(bytes []byte) ([]byte, error) { return nil, errors.New("forced_error") }
	defer delete(iso8583.MarshalEncodings, "force_error")

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("var_to_bytes_%s", testCase.Name), func(t *testing.T) {
			o, err := testCase.V.MarshalISO8583(testCase.Length, testCase.Encoding)
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
			InputBytes:    ebcdic.V1047.FromGoString("ebcdic"),
		},
		{
			Name:          "bytes_is_nil_error",
			InputEncoding: "ascii",
			InputLength:   6,
			OutputContent: "",
			OutputError:   "bytes input is nil",
			InputBytes:    nil,
		},
		{
			Name:          "encoding_error",
			InputEncoding: "force_error",
			InputLength:   6,
			OutputContent: "",
			OutputError:   "encoder 'force_error' returned error: forced_error",
			InputBytes:    []byte("123123"),
		},
		{
			Name:          "unexisting_encoding",
			InputEncoding: "whale_song",
			InputLength:   6,
			OutputContent: "",
			OutputError:   "encoder 'whale_song' does not exist",
			InputBytes:    []byte("123123"),
		},
	}

	iso8583.UnmarshalDecodings["force_error"] = func(bytes []byte) ([]byte, error) { return nil, errors.New("forced_error") }
	defer delete(iso8583.UnmarshalDecodings, "force_error")

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("var_to_bytes_%s", testCase.Name), func(t *testing.T) {
			var v iso8583.VAR

			_, err := v.UnmarshalISO8583(testCase.InputBytes, testCase.InputLength, testCase.InputEncoding)
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
		})
	}
}
