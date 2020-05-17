package iso8583_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jattento/go-iso8583/pkg/bitmap"
	"github.com/jattento/go-iso8583/pkg/field"
	"github.com/jattento/go-iso8583/pkg/iso8583"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	testList := []struct {
		Name string
		Run  bool

		InputByte            []byte
		InputStruct          interface{}
		ExpectedOutputStruct interface{}
		ExpectedOutputError  string
		ExpectedRemaining    int
	}{
		{
			Name:              "simple_one_field",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte:         append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}), []byte("asd")...)...),
			InputStruct: &struct {
				Mti    field.VAR    `iso8583:"mti,length:4"`
				Bitmap field.BITMAP `iso8583:"bitmap,length:64"`
				Field2 field.VAR    `iso8583:"2,length:3"`
			}{},
			ExpectedOutputError: "",
			ExpectedOutputStruct: struct {
				Mti    field.VAR    `iso8583:"mti,length:4"`
				Bitmap field.BITMAP `iso8583:"bitmap,length:64"`
				Field2 field.VAR    `iso8583:"2,length:3"`
			}{
				Field2: "asd",
				Bitmap: field.BITMAP{map[int]bool{1: false, 2: true, 3: false, 4: false, 5: false, 6: false, 7: false, 8: false, 9: false, 10: false, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false, 29: false, 30: false, 31: false, 32: false, 33: false, 34: false, 35: false, 36: false, 37: false, 38: false, 39: false, 40: false, 41: false, 42: false, 43: false, 44: false, 45: false, 46: false, 47: false, 48: false, 49: false, 50: false, 51: false, 52: false, 53: false, 54: false, 55: false, 56: false, 57: false, 58: false, 59: false, 60: false, 61: false, 62: false, 63: false, 64: false}},
				Mti:    "1000",
			},
		},
		{
			Name:              "simple_two_bitmaps",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte:         append(append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: true, 2: true, 64: false}), bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false})...)...), []byte("asdfgh")...),
			InputStruct: &struct {
				Mti     field.VAR    `iso8583:"mti,length:4"`
				Bitmap  field.BITMAP `iso8583:"bitmap,length:64"`
				Field1  field.BITMAP `iso8583:"1,length:64"`
				Field2  field.VAR    `iso8583:"2,length:3"`
				Field66 field.VAR    `iso8583:"66,length:3"`
			}{},
			ExpectedOutputError: "",
			ExpectedOutputStruct: struct {
				Mti     field.VAR    `iso8583:"mti,length:4"`
				Bitmap  field.BITMAP `iso8583:"bitmap,length:64"`
				Field1  field.BITMAP `iso8583:"1,length:64"`
				Field2  field.VAR    `iso8583:"2,length:3"`
				Field66 field.VAR    `iso8583:"66,length:3"`
			}{
				Field2:  "asd",
				Field66: "fgh",
				Field1:  field.BITMAP{map[int]bool{1: false, 2: true, 3: false, 4: false, 5: false, 6: false, 7: false, 8: false, 9: false, 10: false, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false, 29: false, 30: false, 31: false, 32: false, 33: false, 34: false, 35: false, 36: false, 37: false, 38: false, 39: false, 40: false, 41: false, 42: false, 43: false, 44: false, 45: false, 46: false, 47: false, 48: false, 49: false, 50: false, 51: false, 52: false, 53: false, 54: false, 55: false, 56: false, 57: false, 58: false, 59: false, 60: false, 61: false, 62: false, 63: false, 64: false}},
				Bitmap:  field.BITMAP{map[int]bool{1: true, 2: true, 3: false, 4: false, 5: false, 6: false, 7: false, 8: false, 9: false, 10: false, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false, 29: false, 30: false, 31: false, 32: false, 33: false, 34: false, 35: false, 36: false, 37: false, 38: false, 39: false, 40: false, 41: false, 42: false, 43: false, 44: false, 45: false, 46: false, 47: false, 48: false, 49: false, 50: false, 51: false, 52: false, 53: false, 54: false, 55: false, 56: false, 57: false, 58: false, 59: false, 60: false, 61: false, 62: false, 63: false, 64: false}},
				Mti:     "1000",
			},
		},
		{
			Name: "example_1", //TODO BITMAP AUTO LENGTH 64
			Run:  true,
			InputStruct: &struct {
				MTI                   field.MTI    `iso8583:"mti,length:4"`
				FirstBitmap           field.BITMAP `iso8583:"bitmap,length:64"`
				SecondBitmap          field.BITMAP `iso8583:"1,length:64"`
				PAN                   field.LLVAR  `iso8583:"2,length:2"`
				ProcessingCode        field.VAR    `iso8583:"3,length:4"`
				Amount                field.VAR    `iso8583:"4,length:7"`
				ICC                   field.LLLVAR `iso8583:"55,length:3"`
				SettlementCode        field.VAR    `iso8583:"66,length:1"`
				MessageNumber         field.VAR    `iso8583:"71,length:1"`
				TransactionDescriptor field.VAR    `iso8583:"104,length:15"`
			}{},
			ExpectedOutputStruct: struct {
				MTI                   field.MTI    `iso8583:"mti,length:4"`
				FirstBitmap           field.BITMAP `iso8583:"bitmap,length:64"`
				SecondBitmap          field.BITMAP `iso8583:"1,length:64"`
				PAN                   field.LLVAR  `iso8583:"2,length:2"`
				ProcessingCode        field.VAR    `iso8583:"3,length:4"`
				Amount                field.VAR    `iso8583:"4,length:7"`
				ICC                   field.LLLVAR `iso8583:"55,length:3"`
				SettlementCode        field.VAR    `iso8583:"66,length:1"`
				MessageNumber         field.VAR    `iso8583:"71,length:1"`
				TransactionDescriptor field.VAR    `iso8583:"104,length:15"`
			}{
				MTI:                   field.MTI("1000"),
				FirstBitmap:           field.BITMAP{Bitmap: map[int]bool{1: true, 2: true, 3: true, 4: true, 5: false, 6: false, 7: false, 8: false, 9: false, 10: false, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false, 29: false, 30: false, 31: false, 32: false, 33: false, 34: false, 35: false, 36: false, 37: false, 38: false, 39: false, 40: false, 41: false, 42: false, 43: false, 44: false, 45: false, 46: false, 47: false, 48: false, 49: false, 50: false, 51: false, 52: false, 53: false, 54: false, 55: true, 56: false, 57: false, 58: false, 59: false, 60: false, 61: false, 62: false, 63: false, 64: false}},
				SecondBitmap:          field.BITMAP{Bitmap: map[int]bool{1: false, 2: true, 3: false, 4: false, 5: false, 6: false, 7: true, 8: false, 9: false, 10: false, 11: false, 12: false, 13: false, 14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false, 25: false, 26: false, 27: false, 28: false, 29: false, 30: false, 31: false, 32: false, 33: false, 34: false, 35: false, 36: false, 37: false, 38: false, 39: false, 40: true, 41: false, 42: false, 43: false, 44: false, 45: false, 46: false, 47: false, 48: false, 49: false, 50: false, 51: false, 52: false, 53: false, 54: false, 55: false, 56: false, 57: false, 58: false, 59: false, 60: false, 61: false, 62: false, 63: false, 64: false}},
				PAN:                   field.LLVAR("1234567891234567"),
				ProcessingCode:        field.VAR("1000"),
				Amount:                field.VAR("0001000"),
				ICC:                   field.LLLVAR("ABCDEFGH123456789"),
				SettlementCode:        field.VAR("8"),
				MessageNumber:         field.VAR("1"),
				TransactionDescriptor: field.VAR("JUST A PURCHASE"),
			},
			ExpectedOutputError: "",
			InputByte: appendBytes(
				[]byte("1000"),                                  // MTI
				[]byte{0xf0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0}, // First bitmap
				[]byte{0x42, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0}, // Second bitmap
				[]byte("16"), []byte("1234567891234567"),        // PAN
				[]byte("1000"),                             // Processing code
				[]byte("0001000"),                          // Amount
				[]byte("017"), []byte("ABCDEFGH123456789"), // ICC
				[]byte("8"),               // Settlement code
				[]byte("1"),               // Message number
				[]byte("JUST A PURCHASE"), // Transaction Descriptor
			),
			ExpectedRemaining: 0,
		},
	}

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("unmarshal_%s", testCase.Name), func(t *testing.T) {
			if !testCase.Run {
				t.Skip()
				return
			}
			n, err := iso8583.Unmarshal(testCase.InputByte, testCase.InputStruct)
			if testCase.ExpectedOutputError != "" {
				assert.EqualError(t, err, testCase.ExpectedOutputError)
				return
			} else {
				if !assert.Nil(t, err) {
					t.FailNow()
				}
			}

			assert.Equal(t, testCase.ExpectedOutputStruct, reflect.ValueOf(testCase.InputStruct).Elem().Interface())
			assert.Equal(t, len(testCase.InputByte)-n, testCase.ExpectedRemaining)
		})
	}
}
