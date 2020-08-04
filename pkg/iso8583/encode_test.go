package iso8583_test

import (
	"fmt"
	"testing"

	"github.com/jattento/go-iso8583/pkg/bitmap"
	"github.com/jattento/go-iso8583/pkg/field"
	"github.com/jattento/go-iso8583/pkg/iso8583"

	"github.com/stretchr/testify/assert"
)

// TODO Nil field test case
func TestMarshal(t *testing.T) {
	exampleString := field.VAR("1234")
	testList := []struct {
		Name        string
		Run         bool
		Input       interface{}
		OutputBytes []byte
		OutputError string
	}{
		{
			Name: "simple_one_field",
			Run:  true,
			Input: struct {
				Field1 field.VAR `iso8583:"1"`
			}{
				Field1: "1234",
			},
			OutputError: "",
			OutputBytes: []byte("1234"),
		}, {
			Name: "simple_one_field_string_one_bytes",
			Run:  true,
			Input: struct {
				Field1 string `iso8583:"1"`
				Field2 []byte `iso8583:"2"`
			}{
				Field1: "field1",
				Field2: []byte("field2"),
			},
			OutputError: "",
			OutputBytes: []byte("field1field2"),
		},
		{
			Name: "simple_one_field_nil",
			Run:  true,
			Input: struct {
				Field1 *field.VAR `iso8583:"1"`
			}{
				Field1: nil,
			},
			OutputError: "",
			OutputBytes: []byte(""),
		},
		{
			Name: "simple_one_field_one_denied",
			Run:  true,
			Input: struct {
				Field1 field.VAR `iso8583:"1"`
				Field2 field.VAR `iso8583:"-"`
			}{
				Field1: "1234",
				Field2: "1234",
			},
			OutputError: "",
			OutputBytes: []byte("1234"),
		},
		{
			Name: "simple_one_field_one_private_one_anonymous_one_without_tag_one_omitempty",
			Run:  true,
			Input: struct {
				Field1    field.VAR `iso8583:"1"`
				field.VAR `iso8583:"2"`
				field3    field.VAR `iso8583:"3"`
				Field4    field.VAR
				Field5    field.VAR `iso8583:"5,omitempty"`
			}{
				Field1: "1234",
				VAR:    "1234",
				field3: "1234",
				Field4: "1234",
				Field5: "",
			},
			OutputError: "",
			OutputBytes: []byte("1234"),
		},
		{
			Name: "simple_three_field",
			Run:  true,
			Input: struct {
				Field3 field.VAR `iso8583:"3"`
				Field1 field.VAR `iso8583:"1"`
				Field2 field.VAR `iso8583:"2"`
			}{
				Field1: "1234",
				Field2: "1234",
				Field3: "1234",
			},
			OutputError: "",
			OutputBytes: []byte("123412341234"),
		},
		{
			Name: "simple_one_field_pointer",
			Run:  true,
			Input: &struct {
				Field1 *field.VAR `iso8583:"1"`
			}{
				Field1: &exampleString,
			},
			OutputError: "",
			OutputBytes: []byte("1234"),
		},
		{
			Name: "simple_one_field_with_mti_bitmap",
			Run:  true,
			Input: struct {
				Bitmap BMAPWithoutMarshalerBitmap `iso8583:"bitmap"`
				MTI    field.VAR                  `iso8583:"mti"`
				Field1 field.VAR                  `iso8583:"1"`
			}{
				MTI:    "1000",
				Bitmap: BMAPWithoutMarshalerBitmap{Bitmap: bitmap.FromBytes([]byte{126})},
				Field1: "12345",
			},
			OutputError: "",
			OutputBytes: append(append([]byte("1000"), 126), []byte("12345")...),
		},
		{
			Name: "simple_marhsaler_bitmap",
			Run:  true,
			Input: struct {
				Bitmap field.BITMAP `iso8583:"bitmap,length:64"`
				MTI    field.VAR    `iso8583:"mti"`
				Field1 field.VAR    `iso8583:"1"`
				Field2 field.VAR    `iso8583:"2"`
			}{
				MTI:    "1000",
				Field1: "12345",
				Field2: "678",
			},
			OutputError: "",
			OutputBytes: append(append([]byte("1000"), []byte{0xc0, 0, 0, 0, 0, 0, 0, 0}...), []byte("12345678")...),
		},
		{
			Name: "two_marhsaler_bitmap",
			Run:  true,
			Input: struct {
				Bitmap  field.BITMAP `iso8583:"bitmap,length:64"`
				MTI     field.VAR    `iso8583:"mti"`
				Field1  field.BITMAP `iso8583:"1,length:64"`
				Field2  field.VAR    `iso8583:"2"`
				Field66 field.VAR    `iso8583:"66"`
			}{
				MTI:     "1000",
				Field2:  "123",
				Field66: "456",
			},
			OutputError: "",
			OutputBytes: append(append(append([]byte("1000"), // MTI.
				[]byte{0xc0, 0, 0, 0, 0, 0, 0, 0}...), []byte{0x40, 0, 0, 0, 0, 0, 0, 0}...), // First and second bmap.
				[]byte("123456")...), // Field 2 and 66.
		},
		{
			Name: "three_marhsaler_bitmap",
			Run:  true,
			Input: struct {
				Bitmap   field.BITMAP `iso8583:"bitmap,length:64"`
				MTI      field.VAR    `iso8583:"mti"`
				Field1   field.BITMAP `iso8583:"1,length:64"`
				Field2   field.VAR    `iso8583:"2"`
				Field32  field.VAR    `iso8583:"32"`
				Field64  field.VAR    `iso8583:"64"`
				Field65  field.BITMAP `iso8583:"65,length:64"`
				Field66  field.VAR    `iso8583:"66"`
				Field130 field.VAR    `iso8583:"130"`
				Field192 field.VAR    `iso8583:"192"`
			}{
				MTI:      "1000",
				Field2:   "11",
				Field32:  "22",
				Field64:  "33",
				Field66:  "44",
				Field130: "55",
				Field192: "66",
			},
			OutputError: "",
			OutputBytes: append(append(append(append([]byte("1000"), // MTI.
				[]byte{0xc0, 0, 0, 0x1, 0, 0, 0, 0x1, 0xc0, 0, 0, 0, 0, 0, 0, 0}...), // First and second bmap.
				[]byte("112233")...), // Fields 2, 32 and 64.
				[]byte{0x40, 0, 0, 0, 0, 0, 0, 0x1}...), // Third bitmap.
				[]byte("445566")...), // Fields 66, 130,192.
		},
		{
			Name: "four_marhsaler_bitmap_third_with_half_length",
			Run:  true,
			Input: struct {
				Bitmap   field.BITMAP `iso8583:"bitmap,length:64"`
				MTI      field.VAR    `iso8583:"mti"`
				Field1   field.BITMAP `iso8583:"1,length:64"`
				Field2   field.VAR    `iso8583:"2"`
				Field32  field.VAR    `iso8583:"32"`
				Field64  field.VAR    `iso8583:"64"`
				Field65  field.BITMAP `iso8583:"65,length:32"`
				Field66  field.VAR    `iso8583:"66"`
				Field96  field.VAR    `iso8583:"96"`
				Field129 field.BITMAP `iso8583:"129,length:64"`
				Field160 field.VAR    `iso8583:"160"`
				Field162 field.VAR    `iso8583:"162"`
				Field192 field.VAR    `iso8583:"192"`
				Field224 field.VAR    `iso8583:"224"`
			}{
				MTI:      "1000",
				Field2:   "11",
				Field32:  "22",
				Field64:  "33",
				Field66:  "44",
				Field96:  "55",
				Field160: "66",
				Field162: "77",
				Field192: "88",
				Field224: "99",
			},
			OutputError: "",
			OutputBytes: append(append(append(append(append(append([]byte("1000"), // MTI.
				[]byte{0xc0, 0, 0, 0x1, 0, 0, 0, 0x1, 0xc0, 0, 0, 0x1, 0, 0, 0, 0x0}...), // First and second bmap.
				[]byte("112233")...), // Fields 2, 32 and 64.
				[]byte{0x80, 0, 0, 0x1}...), // Third bitmap.
				[]byte("4455")...), // Fields 66 and 96.
				[]byte{0x40, 0, 0, 0x1, 0, 0, 0, 0x1}...), // Fourth bitmap.
				[]byte("66778899")...), // Fields 162, 192,224.
		},
		{
			Name: "example_1", //TODO BITMAP AUTO LENGTH 64
			Run:  true,
			Input: struct {
				FirstBitmap           field.BITMAP `iso8583:"bitmap,length:64"`
				SecondBitmap          field.BITMAP `iso8583:"1,length:64"`
				PAN                   field.LLVAR  `iso8583:"2"`
				ProcessingCode        field.VAR    `iso8583:"3"`
				Amount                field.VAR    `iso8583:"4"`
				ICC                   field.LLLVAR `iso8583:"55"`
				SettlementCode        field.VAR    `iso8583:"66"`
				MessageNumber         field.VAR    `iso8583:"71"`
				TransactionDescriptor field.VAR    `iso8583:"104"`
			}{
				PAN:                   field.LLVAR("1234567891234567"),
				ProcessingCode:        field.VAR("1000"),
				Amount:                field.VAR("0001000"),
				ICC:                   field.LLLVAR("ABCDEFGH123456789"),
				SettlementCode:        field.VAR("8"),
				MessageNumber:         field.VAR("1"),
				TransactionDescriptor: field.VAR("JUST A PURCHASE"),
			},
			OutputError: "",
			OutputBytes: appendBytes(
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
		},
		{
			Name: "error_0_field",
			Run:  true,
			Input: struct {
				Field0 field.VAR `iso8583:"0"`
			}{
				Field0: "12345",
			},
			OutputError: "iso8583.marshal: field 0 not allowed",
			OutputBytes: nil,
		},
		{
			Name: "error_duplicated_field",
			Run:  true,
			Input: struct {
				Field1  field.VAR `iso8583:"1"`
				Field01 field.VAR `iso8583:"1"`
			}{
				Field1:  "12345",
				Field01: "12345",
			},
			OutputError: "iso8583.marshal: field 1 is repeated",
			OutputBytes: nil,
		},
		{
			Name:        "error_not_struct",
			Run:         true,
			Input:       "string",
			OutputError: "iso8583.marshal: input is not a struct or is pointing to one",
			OutputBytes: nil,
		},
		{
			Name: "error_unrecognized_field",
			Run:  true,
			Input: struct {
				Field1 field.VAR `iso8583:"asd"`
			}{
				Field1: "1234",
			},
			OutputError: "iso8583.marshal: field asd does not have a valid field name",
			OutputBytes: nil,
		},
	}

	for _, testCase := range testList {
		t.Run(fmt.Sprintf("marshal_%s", testCase.Name), func(t *testing.T) {
			if !testCase.Run {
				t.Skip()
				return
			}
			o, err := iso8583.Marshal(testCase.Input)
			if testCase.OutputError != "" {
				assert.EqualError(t, err, testCase.OutputError)
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

// BITMAP clone without MarshalerBitmap for testing purpose.
type BMAPWithoutMarshalerBitmap struct {
	bitmap.Bitmap
}

func (b *BMAPWithoutMarshalerBitmap) Bits() (map[int]bool, error) {
	return b.Bitmap, nil
}

func (b *BMAPWithoutMarshalerBitmap) UnmarshalISO8583(byt []byte, length int, encoding string) (int, error) {
	const bitmapLength = 8
	b.Bitmap = bitmap.FromBytes(byt[:bitmapLength])
	return bitmapLength, nil
}

func (b BMAPWithoutMarshalerBitmap) MarshalISO8583(length int, encoding string) ([]byte, error) {
	return bitmap.ToBytes(b.Bitmap), nil
}
func appendBytes(b ...[]byte) (bb []byte) {
	for _, byt := range b {
		bb = append(bb, byt...)
	}
	return bb
}
