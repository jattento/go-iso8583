package iso8583_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jattento/go-iso8583/pkg/bitmap"
	"github.com/jattento/go-iso8583/pkg/iso8583"

	"github.com/stretchr/testify/assert"
)

// TODO Nil field test case
// TODO second bitmap only if follow field are present
func TestMarshal(t *testing.T) {
	exampleString := iso8583.VAR("1234")
	testList := []struct {
		Name        string
		Run         bool
		Input       interface{}
		OutputBytes []byte
		OutputError string
	}{
		{
			Name: "simple_one_field_ebcdic",
			Run:  true,
			Input: struct {
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1 iso8583.VAR    `iso8583:"mti,encoding:ebcdic"`
				Field2 iso8583.VAR    `iso8583:"1"`
			}{
				Field1: "1234",
				Field2: "1234",
			},
			OutputError: "",
			OutputBytes: []byte{0xf1, 0xf2, 0xf3, 0xf4, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x31, 0x32, 0x33, 0x34},
		},
		{
			Name: "unused_bitmap",
			Run:  true,
			Input: struct {
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1 iso8583.VAR    `iso8583:"mti"`
				Field2 iso8583.BITMAP `iso8583:"1,length:64"`
				Field3 iso8583.VAR    `iso8583:"2"`
			}{
				Field1: "1000",
				Field3: "1234",
			},
			OutputError: "",
			OutputBytes: []byte{0x31, 0x30, 0x30, 0x30, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x31, 0x32, 0x33, 0x34},
		},
		{
			Name: "simple_one_field_string_one_bytes",
			Run:  true,
			Input: struct {
				MTI    string         `iso8583:"mti"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field2 []byte         `iso8583:"2"`
			}{
				MTI:    "field1",
				Field2: []byte("field2"),
			},
			OutputError: "",
			OutputBytes: appendBytes([]byte("field1"), []byte{0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, []byte("field2")),
		},
		{
			Name: "simple_one_field_nil",
			Run:  true,
			Input: struct {
				MTI    string         `iso8583:"mti"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1 *iso8583.VAR   `iso8583:"1"`
				Field2 iso8583.VAR    `iso8583:"2"`
			}{
				MTI:    "1000",
				Field1: nil,
				Field2: "1000",
			},
			OutputError: "",
			OutputBytes: []byte{0x31, 0x30, 0x30, 0x30, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x31, 0x30, 0x30, 0x30},
		},
		{
			Name: "simple_one_field_one_denied",
			Run:  true,
			Input: struct {
				MTI    iso8583.VAR    `iso8583:"mti"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field2 iso8583.VAR    `iso8583:"-"`
				Field3 iso8583.VAR    `iso8583:"2"`
			}{
				MTI:    "1234",
				Field2: "1234",
				Field3: "1234",
			},
			OutputError: "",
			OutputBytes: []byte{0x31, 0x32, 0x33, 0x34, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x31, 0x32, 0x33, 0x34},
		},
		{
			Name: "simple_one_field_one_private_one_anonymous_one_without_tag_one_omitempty",
			Run:  true,
			Input: struct {
				MTI         iso8583.VAR    `iso8583:"mti"`
				Bitmap      iso8583.BITMAP `iso8583:"bitmap,length:64"`
				iso8583.VAR `iso8583:"2"`
				field3      iso8583.VAR `iso8583:"3"`
				Field4      iso8583.VAR
				Field5      iso8583.VAR `iso8583:"5,omitempty"`
				Field6      iso8583.VAR `iso8583:"6,omitempty"`
			}{
				MTI:    "1234",
				VAR:    "1234",
				field3: "1234",
				Field4: "1234",
				Field5: "",
				Field6: "1234",
			},
			OutputError: "",
			OutputBytes: []byte{0x31, 0x32, 0x33, 0x34, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x31, 0x32, 0x33, 0x34},
		},
		{
			Name: "simple_three_field",
			Run:  true,
			Input: struct {
				MTI    iso8583.VAR    `iso8583:"mti"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1 iso8583.VAR    `iso8583:"1"`
				Field2 iso8583.VAR    `iso8583:"2"`
			}{
				Field1: "1234",
				Field2: "1234",
				MTI:    "1234",
			},
			OutputError: "",
			OutputBytes: appendBytes([]byte("1234"), []byte{0xc0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, []byte("12341234")),
		},
		{
			Name: "simple_one_field_pointer",
			Run:  true,
			Input: &struct {
				MTI    *iso8583.VAR   `iso8583:"mti"`
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1 iso8583.VAR    `iso8583:"1"`
			}{
				MTI:    &exampleString,
				Field1: "1234",
			},
			OutputError: "",
			OutputBytes: []byte{0x31, 0x32, 0x33, 0x34, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x31, 0x32, 0x33, 0x34},
		},
		{
			Name: "simple_one_field_with_mti_bitmap",
			Run:  true,
			Input: struct {
				Bitmap BMAPWithoutMarshalerBitmap `iso8583:"bitmap"`
				MTI    iso8583.VAR                `iso8583:"mti"`
				Field1 iso8583.VAR                `iso8583:"1"`
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
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				MTI    iso8583.VAR    `iso8583:"mti"`
				Field1 iso8583.VAR    `iso8583:"1"`
				Field2 iso8583.VAR    `iso8583:"2"`
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
				Bitmap  iso8583.BITMAP `iso8583:"bitmap,length:64"`
				MTI     iso8583.VAR    `iso8583:"mti"`
				Field1  iso8583.BITMAP `iso8583:"1,length:64"`
				Field2  iso8583.VAR    `iso8583:"2"`
				Field66 iso8583.VAR    `iso8583:"66"`
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
				Bitmap   iso8583.BITMAP `iso8583:"bitmap,length:64"`
				MTI      iso8583.VAR    `iso8583:"mti"`
				Field1   iso8583.BITMAP `iso8583:"1,length:64"`
				Field2   iso8583.VAR    `iso8583:"2"`
				Field32  iso8583.VAR    `iso8583:"32"`
				Field64  iso8583.VAR    `iso8583:"64"`
				Field65  iso8583.BITMAP `iso8583:"65,length:64"`
				Field66  iso8583.VAR    `iso8583:"66"`
				Field130 iso8583.VAR    `iso8583:"130"`
				Field192 iso8583.VAR    `iso8583:"192"`
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
				Bitmap   iso8583.BITMAP `iso8583:"bitmap,length:64"`
				MTI      iso8583.VAR    `iso8583:"mti"`
				Field1   iso8583.BITMAP `iso8583:"1,length:64"`
				Field2   iso8583.VAR    `iso8583:"2"`
				Field32  iso8583.VAR    `iso8583:"32"`
				Field64  iso8583.VAR    `iso8583:"64"`
				Field65  iso8583.BITMAP `iso8583:"65,length:32"`
				Field66  iso8583.VAR    `iso8583:"66"`
				Field96  iso8583.VAR    `iso8583:"96"`
				Field129 iso8583.BITMAP `iso8583:"129,length:64"`
				Field160 iso8583.VAR    `iso8583:"160"`
				Field162 iso8583.VAR    `iso8583:"162"`
				Field192 iso8583.VAR    `iso8583:"192"`
				Field224 iso8583.VAR    `iso8583:"224"`
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
			Name: "example_1",
			Run:  true,
			Input: struct {
				FirstBitmap           iso8583.BITMAP `iso8583:"bitmap,length:64"`
				SecondBitmap          iso8583.BITMAP `iso8583:"1,length:64"`
				MTI                   iso8583.LLVAR  `iso8583:"mti"`
				ProcessingCode        iso8583.VAR    `iso8583:"3"`
				Amount                iso8583.VAR    `iso8583:"4"`
				ICC                   iso8583.LLLVAR `iso8583:"55"`
				SettlementCode        iso8583.VAR    `iso8583:"66"`
				MessageNumber         iso8583.VAR    `iso8583:"71"`
				TransactionDescriptor iso8583.VAR    `iso8583:"104"`
			}{
				MTI:                   iso8583.LLVAR("1234567891234567"),
				ProcessingCode:        iso8583.VAR("1000"),
				Amount:                iso8583.VAR("0001000"),
				ICC:                   iso8583.LLLVAR("ABCDEFGH123456789"),
				SettlementCode:        iso8583.VAR("8"),
				MessageNumber:         iso8583.VAR("1"),
				TransactionDescriptor: iso8583.VAR("JUST A PURCHASE"),
			},
			OutputError: "",
			OutputBytes: appendBytes(
				[]byte("16"), []byte("1234567891234567"), // MTI
				[]byte{0xb0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0}, // First bitmap
				[]byte{0x42, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0}, // Second bitmap
				[]byte("1000"),                                  // Processing code
				[]byte("0001000"),                               // Amount
				[]byte("017"), []byte("ABCDEFGH123456789"),      // ICC
				[]byte("8"),               // Settlement code
				[]byte("1"),               // Message number
				[]byte("JUST A PURCHASE"), // Transaction Descriptor

			),
		},
		{
			Name: "field_marshal_bitmap_len_0",
			Run:  true,
			Input: struct {
				Field1  iso8583.VAR    `iso8583:"mti,encoding:ascii"`
				Bitmap  iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Bitmap2 BmapMock       `iso8583:"2,length:64"`
				Field3  iso8583.VAR    `iso8583:"3,encoding:ascii"`
			}{
				Field1:  "1",
				Bitmap2: BmapMock{},
				Field3:  "999",
			},
			OutputError: "",
			OutputBytes: []byte{0x31, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x39, 0x39, 0x39},
		},
		{
			Name: "error_0_field",
			Run:  true,
			Input: struct {
				Field0 iso8583.VAR `iso8583:"0"`
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
				Field1  iso8583.VAR `iso8583:"1"`
				Field01 iso8583.VAR `iso8583:"1"`
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
				Field1 iso8583.VAR `iso8583:"asd"`
			}{
				Field1: "1234",
			},
			OutputError: "iso8583.marshal: field asd does not have a valid field name",
			OutputBytes: nil,
		},
		{
			Name:        "nil_input",
			Run:         true,
			Input:       nil,
			OutputError: "iso8583.marshal: nil input",
			OutputBytes: nil,
		},
		{
			Name: "length_not_int_error",
			Run:  true,
			Input: struct {
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1 iso8583.VAR    `iso8583:"mti,encoding:ebcdic,length:a"`
			}{
				Field1: "1234",
			},
			OutputError: "iso8583.marshal: field mti does not have a valid length",
			OutputBytes: nil,
		},
		{
			Name: "field_marshal_error",
			Run:  true,
			Input: struct {
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1 iso8583.VAR    `iso8583:"mti,encoding:whale_song"`
			}{
				Field1: "1234",
			},
			OutputError: "iso8583.marshal: field mti cant be marshaled: encoder 'whale_song' does not exist",
			OutputBytes: nil,
		},
		{
			Name: "not_implemented_type_error",
			Run:  true,
			Input: struct {
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1 rune           `iso8583:"mti,encoding:ascii"`
			}{
				Field1: 'q',
			},
			OutputError: "iso8583.marshal: field mti does not implement Marshaler interface, is a string or slice of " +
				"bytes but does have iso8583 tags",
			OutputBytes: nil,
		},
		{
			Name: "mti_repeated_error",
			Run:  true,
			Input: struct {
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1 iso8583.VAR    `iso8583:"mti,encoding:ascii"`
				Field2 iso8583.VAR    `iso8583:"mti,encoding:ascii"`
			}{
				Field1: "1",
				Field2: "2",
			},
			OutputError: "iso8583.marshal: field mti is repeated",
			OutputBytes: nil,
		},
		{
			Name: "bitmap_repeated_not_bitmap_marshal_error",
			Run:  true,
			Input: struct {
				Bitmap  []byte      `iso8583:"bitmap,length:64"`
				Bitmap2 []byte      `iso8583:"bitmap,length:64"`
				Field1  iso8583.VAR `iso8583:"mti,encoding:ascii"`
			}{
				Field1:  "1",
				Bitmap:  []byte{1},
				Bitmap2: []byte{1},
			},
			OutputError: "iso8583.marshal: field bitmap is repeated",
			OutputBytes: nil,
		},
		{
			Name: "bitmap_repeated_error",
			Run:  true,
			Input: struct {
				Bitmap  iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Bitmap2 iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1  iso8583.VAR    `iso8583:"mti,encoding:ascii"`
			}{
				Field1: "1",
			},
			OutputError: "iso8583.marshal: field bitmap is repeated",
			OutputBytes: nil,
		},
		{
			Name: "bitmap_invalid_length_error",
			Run:  true,
			Input: struct {
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:a"`
				Field1 iso8583.VAR    `iso8583:"mti,encoding:ascii"`
			}{
				Field1: "1",
			},
			OutputError: "iso8583.marshal: field bitmap is implements bitmap interface and does not have a valid length",
			OutputBytes: nil,
		},
		{
			Name: "bitmap_no_content_error",
			Run:  true,
			Input: struct {
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Field1 iso8583.VAR    `iso8583:"mti,encoding:ascii"`
			}{
				Field1: "1",
			},
			OutputError: "iso8583.marshal: first bitmap present but without content",
			OutputBytes: nil,
		},
		{
			Name: "bitmap_invalid_name_error",
			Run:  true,
			Input: struct {
				Bitmap       iso8583.BITMAP `iso8583:"bitmap,length:64"`
				BitmapSecond iso8583.BITMAP `iso8583:"bmap,length:64"`
				Field1       iso8583.VAR    `iso8583:"mti,encoding:ascii"`
			}{
				Field1: "1",
			},
			OutputError: "iso8583: unrecognized field: bmap",
			OutputBytes: nil,
		},
		{
			Name: "bitmap_invalid_name_0_error",
			Run:  true,
			Input: struct {
				Bitmap iso8583.BITMAP `iso8583:"0,length:64"`
				Field1 iso8583.VAR    `iso8583:"mti,encoding:ascii"`
			}{
				Field1: "1",
			},
			OutputError: "iso8583.marshal: field 0 not allowed",
			OutputBytes: nil,
		},
		{
			Name: "bitmap_repeated_error",
			Run:  true,
			Input: struct {
				Bitmap  iso8583.BITMAP `iso8583:"bitmap,length:64"`
				Bitmap2 iso8583.BITMAP `iso8583:"2,length:64"`
				Bitmap3 iso8583.BITMAP `iso8583:"2,length:64"`
				Field1  iso8583.VAR    `iso8583:"mti,encoding:ascii"`
			}{
				Field1: "1",
			},
			OutputError: "iso8583.marshal: field 2 is repeated",
			OutputBytes: nil,
		},
		{
			Name: "mti_not_present_error",
			Run:  true,
			Input: struct {
				Bitmap iso8583.BITMAP `iso8583:"bitmap,length:64"`
			}{},
			OutputError: "iso8583.marshal: MTI not present",
			OutputBytes: nil,
		},
		{
			Name: "bitmap_not_present",
			Run:  true,
			Input: struct {
				Field1 iso8583.VAR `iso8583:"mti,encoding:ascii"`
			}{
				Field1: "1",
			},
			OutputError: "iso8583.marshal: first bitmap no present",
			OutputBytes: nil,
		},
		{
			Name: "field_marshal_bitmap_failed",
			Run:  true,
			Input: struct {
				Field1 iso8583.VAR `iso8583:"mti,encoding:ascii"`
				Bitmap BmapMock    `iso8583:"bitmap,length:64"`
			}{
				Field1: "1",
				Bitmap: BmapMock{returnError: errors.New("forced_error")},
			},
			OutputError: "iso8583.marshal: field bitmap cant be marshaled: forced_error",
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

type BmapMock struct {
	returnError              error
	returnValueMarshalBmap   []byte
	returnValueUnmarshalBmap int
}

func (b *BmapMock) UnmarshalISO8583(byt []byte, length int, encoding string) (int, error) {
	if b.returnValueUnmarshalBmap != 0 {
		return b.returnValueUnmarshalBmap, nil
	}

	return 0, b.returnError
}

func (b BmapMock) MarshalISO8583(length int, encoding string) ([]byte, error) {
	return nil, b.returnError
}

func (b BmapMock) Bits() (map[int]bool, error) {
	return nil, b.returnError
}

func (b BmapMock) MarshalISO8583Bitmap(m map[int]bool, encoding string) ([]byte, error) {
	return b.returnValueMarshalBmap, b.returnError
}
