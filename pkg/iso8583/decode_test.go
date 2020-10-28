package iso8583_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/jattento/go-iso8583/pkg/bitmap"
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
			Name:              "simple_one_field_and_private_one_and_disesteem",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}),
				[]byte("asd")...)...),
			InputStruct: &struct {
				Mti    iso8583.VAR    `iso8583:"mti,length:4"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field2 iso8583.VAR    `iso8583:"2,length:3"`
				field3 iso8583.VAR    `iso8583:"3,length:3"`
			}{},
			ExpectedOutputError: "",
			ExpectedOutputStruct: struct {
				Mti    iso8583.VAR    `iso8583:"mti,length:4"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field2 iso8583.VAR    `iso8583:"2,length:3"`
				field3 iso8583.VAR    `iso8583:"3,length:3"`
			}{
				Field2: "asd",
				Bitmap: iso8583.BITMAP{Bitmap: map[int]bool{1: false, 2: true, 3: false, 4: false, 5: false, 6: false,
					7: false, 8: false, 9: false, 10: false, 11: false, 12: false, 13: false, 14: false, 15: false,
					16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false,
					25: false, 26: false, 27: false, 28: false, 29: false, 30: false, 31: false, 32: false, 33: false,
					34: false, 35: false, 36: false, 37: false, 38: false, 39: false, 40: false, 41: false, 42: false,
					43: false, 44: false, 45: false, 46: false, 47: false, 48: false, 49: false, 50: false, 51: false,
					52: false, 53: false, 54: false, 55: false, 56: false, 57: false, 58: false, 59: false, 60: false,
					61: false, 62: false, 63: false, 64: false}},
				Mti: "1000",
			},
		},
		{
			Name:              "simple_two_bitmaps",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append(append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: true, 2: true,
				64: false}), bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false})...)...), []byte("asdfgh")...),
			InputStruct: &struct {
				Mti     iso8583.VAR    `iso8583:"mti,length:4"`
				Bitmap  iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1  iso8583.BITMAP `iso8583:"1,length:64"`
				Field2  iso8583.VAR    `iso8583:"2,length:3"`
				Field66 iso8583.VAR    `iso8583:"66,length:3"`
			}{},
			ExpectedOutputError: "",
			ExpectedOutputStruct: struct {
				Mti     iso8583.VAR    `iso8583:"mti,length:4"`
				Bitmap  iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1  iso8583.BITMAP `iso8583:"1,length:64"`
				Field2  iso8583.VAR    `iso8583:"2,length:3"`
				Field66 iso8583.VAR    `iso8583:"66,length:3"`
			}{
				Field2:  "asd",
				Field66: "fgh",
				Field1: iso8583.BITMAP{Bitmap: map[int]bool{1: false, 2: true, 3: false, 4: false, 5: false, 6: false,
					7: false, 8: false, 9: false, 10: false, 11: false, 12: false, 13: false, 14: false, 15: false,
					16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false,
					25: false, 26: false, 27: false, 28: false, 29: false, 30: false, 31: false, 32: false, 33: false,
					34: false, 35: false, 36: false, 37: false, 38: false, 39: false, 40: false, 41: false, 42: false,
					43: false, 44: false, 45: false, 46: false, 47: false, 48: false, 49: false, 50: false, 51: false,
					52: false, 53: false, 54: false, 55: false, 56: false, 57: false, 58: false, 59: false, 60: false,
					61: false, 62: false, 63: false, 64: false}},
				Bitmap: iso8583.BITMAP{Bitmap: map[int]bool{1: true, 2: true, 3: false, 4: false, 5: false, 6: false,
					7: false, 8: false, 9: false, 10: false, 11: false, 12: false, 13: false, 14: false, 15: false,
					16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false, 24: false,
					25: false, 26: false, 27: false, 28: false, 29: false, 30: false, 31: false, 32: false, 33: false,
					34: false, 35: false, 36: false, 37: false, 38: false, 39: false, 40: false, 41: false, 42: false,
					43: false, 44: false, 45: false, 46: false, 47: false, 48: false, 49: false, 50: false, 51: false,
					52: false, 53: false, 54: false, 55: false, 56: false, 57: false, 58: false, 59: false, 60: false,
					61: false, 62: false, 63: false, 64: false}},
				Mti: "1000",
			},
		},
		{
			Name: "example_1",
			Run:  true,
			InputStruct: &struct {
				MTI                   iso8583.MTI    `iso8583:"mti,length:4"`
				FirstBitmap           iso8583.BITMAP `iso8583:"bitmap,length:64"`
				SecondBitmap          iso8583.BITMAP `iso8583:"1,length:64"`
				PAN                   iso8583.LLVAR  `iso8583:"2,length:2"`
				ProcessingCode        iso8583.BINARY `iso8583:"3,length:4"`
				Amount                iso8583.VAR    `iso8583:"4,length:7"`
				ICC                   iso8583.LLLVAR `iso8583:"55,length:3"`
				SettlementCode        iso8583.VAR    `iso8583:"66,length:1"`
				MessageNumber         iso8583.VAR    `iso8583:"71,length:1"`
				TransactionDescriptor iso8583.VAR    `iso8583:"104,length:15"`
			}{},
			ExpectedOutputStruct: struct {
				MTI                   iso8583.MTI    `iso8583:"mti,length:4"`
				FirstBitmap           iso8583.BITMAP `iso8583:"bitmap,length:64"`
				SecondBitmap          iso8583.BITMAP `iso8583:"1,length:64"`
				PAN                   iso8583.LLVAR  `iso8583:"2,length:2"`
				ProcessingCode        iso8583.BINARY `iso8583:"3,length:4"`
				Amount                iso8583.VAR    `iso8583:"4,length:7"`
				ICC                   iso8583.LLLVAR `iso8583:"55,length:3"`
				SettlementCode        iso8583.VAR    `iso8583:"66,length:1"`
				MessageNumber         iso8583.VAR    `iso8583:"71,length:1"`
				TransactionDescriptor iso8583.VAR    `iso8583:"104,length:15"`
			}{
				MTI: iso8583.MTI{MTI: "1000"},
				FirstBitmap: iso8583.BITMAP{Bitmap: map[int]bool{1: true, 2: true, 3: true, 4: true, 5: false,
					6: false, 7: false, 8: false, 9: false, 10: false, 11: false, 12: false, 13: false, 14: false,
					15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false, 23: false,
					24: false, 25: false, 26: false, 27: false, 28: false, 29: false, 30: false, 31: false, 32: false,
					33: false, 34: false, 35: false, 36: false, 37: false, 38: false, 39: false, 40: false, 41: false,
					42: false, 43: false, 44: false, 45: false, 46: false, 47: false, 48: false, 49: false, 50: false,
					51: false, 52: false, 53: false, 54: false, 55: true, 56: false, 57: false, 58: false, 59: false,
					60: false, 61: false, 62: false, 63: false, 64: false}},
				SecondBitmap: iso8583.BITMAP{Bitmap: map[int]bool{1: false, 2: true, 3: false, 4: false,
					5: false, 6: false, 7: true, 8: false, 9: false, 10: false, 11: false, 12: false, 13: false,
					14: false, 15: false, 16: false, 17: false, 18: false, 19: false, 20: false, 21: false, 22: false,
					23: false, 24: false, 25: false, 26: false, 27: false, 28: false, 29: false, 30: false, 31: false,
					32: false, 33: false, 34: false, 35: false, 36: false, 37: false, 38: false, 39: false, 40: true,
					41: false, 42: false, 43: false, 44: false, 45: false, 46: false, 47: false, 48: false, 49: false,
					50: false, 51: false, 52: false, 53: false, 54: false, 55: false, 56: false, 57: false, 58: false,
					59: false, 60: false, 61: false, 62: false, 63: false, 64: false}},
				PAN:                   iso8583.LLVAR("1234567891234567"),
				ProcessingCode:        iso8583.BINARY("1000"),
				Amount:                iso8583.VAR("0001000"),
				ICC:                   iso8583.LLLVAR("ABCDEFGH123456789"),
				SettlementCode:        iso8583.VAR("8"),
				MessageNumber:         iso8583.VAR("1"),
				TransactionDescriptor: iso8583.VAR("JUST A PURCHASE"),
			},
			ExpectedOutputError: "",
			InputByte: appendBytes(
				[]byte("1000"), // MTI
				[]byte{0xf0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0}, // First bitmap
				[]byte{0x42, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0}, // Second bitmap
				[]byte("16"), []byte("1234567891234567"), // PAN
				[]byte("1000"),                             // Processing code
				[]byte("0001000"),                          // Amount
				[]byte("017"), []byte("ABCDEFGH123456789"), // ICC
				[]byte("8"),               // Settlement code
				[]byte("1"),               // Message number
				[]byte("JUST A PURCHASE"), // Transaction Descriptor
			),
			ExpectedRemaining: 0,
		},
		{
			Name:              "length_not_int_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}),
				[]byte("asd")...)...),
			InputStruct: &struct {
				Mti    iso8583.VAR    `iso8583:"mti,length:a"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field2 iso8583.VAR    `iso8583:"2,length:3"`
			}{},
			ExpectedOutputError:  "iso8583.unmarshal: field mti: invalid length: strconv.Atoi: parsing \"a\": invalid syntax",
			ExpectedOutputStruct: nil,
		},
		{
			Name:              "repeated_field_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}),
				[]byte("asd")...)...),
			InputStruct: &struct {
				Mti    iso8583.VAR    `iso8583:"mti,length:4"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field2 iso8583.VAR    `iso8583:"2,length:3"`
				Field3 iso8583.VAR    `iso8583:"2,length:3"`
			}{},
			ExpectedOutputError:  "iso8583.unmarshal: field 2 is repeteated in struct",
			ExpectedOutputStruct: nil,
		},
		{
			Name:              "insufficient_bytes_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}),
				[]byte("asd")...)...),
			InputStruct: &struct {
				Mti    iso8583.VAR    `iso8583:"mti,length:4"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field2 VarMock        `iso8583:"2,length:3"`
			}{
				Field2: VarMock{returnUnmarshal: 999},
			},
			ExpectedOutputError:  "iso8583.unmarshal: Unmarshaler from field struct 2 returned a n higher than unconsumed bytes",
			ExpectedOutputStruct: nil,
		},
		{
			Name:              "input_nil_struct_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}),
				[]byte("asd")...)...),
			InputStruct:          nil,
			ExpectedOutputError:  "iso8583.unmarshal: interface input is not a pointer to a structure",
			ExpectedOutputStruct: nil,
		},
		{
			Name:              "insufficient_bytes_mti_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}),
				[]byte("asd")...)...),
			InputStruct: &struct {
				MTI    VarMock        `iso8583:"mti,length:3"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
			}{
				MTI: VarMock{returnUnmarshal: 999},
			},
			ExpectedOutputError:  "iso8583.unmarshal: Unmarshaler from field mti returned a n higher than unconsumed bytes",
			ExpectedOutputStruct: nil,
		},
		{
			Name:              "insufficient_bytes_bitmap_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}),
				[]byte("asd")...)...),
			InputStruct: &struct {
				MTI    iso8583.VAR `iso8583:"mti,length:3"`
				Bitmap VarMock     `iso8583:"bitmap,length:64"`
			}{
				Bitmap: VarMock{returnUnmarshal: 999},
			},
			ExpectedOutputError:  "iso8583.unmarshal: Unmarshaler from field bitmap returned a n higher than unconsumed bytes",
			ExpectedOutputStruct: nil,
		},
		{
			Name:              "field_unmarshal_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}),
				[]byte("asd")...)...),
			InputStruct: &struct {
				MTI    iso8583.VAR    `iso8583:"mti,length:3,encoding:whale_song"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
			}{},
			ExpectedOutputError:  "iso8583.unmarshal: cant unmarshal field mti: encoder 'whale_song' does not exist",
			ExpectedOutputStruct: nil,
		},
		{
			Name:              "field_unmarshal_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}),
				[]byte("asd")...)...),
			InputStruct: &struct {
				MTI    iso8583.VAR `iso8583:"mti,length:3"`
				Bitmap BmapMock    `iso8583:"bitmap,length:64"`
			}{
				Bitmap: BmapMock{returnError: errors.New("forced_error")},
			},
			ExpectedOutputError:  "iso8583.unmarshal: cant unmarshal field bitmap: forced_error",
			ExpectedOutputStruct: nil,
		},
		{
			Name:              "bitmap_bits_method_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}),
				[]byte("asd")...)...),
			InputStruct: &struct {
				MTI    iso8583.VAR `iso8583:"mti,length:3"`
				Bitmap BmapMock    `iso8583:"bitmap,length:64"`
			}{
				Bitmap: BmapMock{returnError: errors.New("forced_error"), returnValueUnmarshalBmap: 8},
			},
			ExpectedOutputError:  "iso8583.unmarshal: failed reading first bitmap: forced_error",
			ExpectedOutputStruct: nil,
		},
		{
			Name:              "bitmap_get_bits_method_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}),
				[]byte("asd")...)...),
			InputStruct: &struct {
				MTI    iso8583.VAR                `iso8583:"mti,length:3"`
				Bitmap BMAPWithoutMarshalerBitmap `iso8583:"bitmap,length:64"`
			}{},
			ExpectedOutputError:  "iso8583.unmarshal: bitmap field is present but does not implement UnmarshalerBitmap",
			ExpectedOutputStruct: nil,
		},
		{
			Name:              "getting_unmarshaler_field_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: false, 2: true, 64: false}),
				[]byte("asd")...)...),
			InputStruct: &struct {
				MTI    iso8583.VAR    `iso8583:"mti,length:4"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1 uint8          `iso8583:"2,length:3"`
			}{},
			ExpectedOutputError:  "iso8583.unmarshal: field 2 is present but does not implement Unmarshaler interface",
			ExpectedOutputStruct: nil,
		},
		{
			Name:              "getting_unmarshaler_field_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append(append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: true, 2: true, 64: false}),
				bitmap.ToBytes(map[int]bool{1: false, 64: false})...)...), []byte("asd")...),
			InputStruct: &struct {
				MTI          iso8583.VAR    `iso8583:"mti,length:4"`
				Bitmap       iso8583.BITMAP `iso8583:"bitmap,length:64"`
				SecondBitmap BmapMock       `iso8583:"1,length:64"`
				Field1       iso8583.VAR    `iso8583:"2,length:3"`
			}{
				SecondBitmap: BmapMock{returnValueUnmarshalBmap: 8, returnError: errors.New("forced_error")},
			},
			ExpectedOutputError:  "iso8583.unmarshal: failed reading field 1 bitmap: forced_error",
			ExpectedOutputStruct: nil,
		},
		{
			Name:              "unexpected_field_income_error",
			Run:               true,
			ExpectedRemaining: 0,
			InputByte: append(append([]byte("1000"), append(bitmap.ToBytes(map[int]bool{1: true, 64: false}),
				bitmap.ToBytes(map[int]bool{1: false, 64: false})...)...), []byte("asd")...),
			InputStruct: &struct {
				MTI    iso8583.VAR    `iso8583:"mti,length:4"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
			}{},
			ExpectedOutputError:  "iso8583.unmarshal: unknown field in message '1', cant resolve upcomming fields",
			ExpectedOutputStruct: nil,
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

// VAR type should be used for fixed length fields.
type VarMock struct {
	returnError     error
	returnMarshal   []byte
	returnUnmarshal int
}

// MarshalISO8583 allows to use this type in structs and be able tu iso8583.Marshal it.
func (v VarMock) MarshalISO8583(length int, enc string) ([]byte, error) {
	return v.returnMarshal, v.returnError
}

// UnmarshalISO8583 allows to use this type in structs and be able tu iso8583.Unmarshal it. //TODO CHECK LENGTH
func (v *VarMock) UnmarshalISO8583(b []byte, length int, enc string) (int, error) {
	return v.returnUnmarshal, v.returnError
}
