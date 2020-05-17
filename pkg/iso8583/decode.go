package iso8583

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Unmarshaler interface for iso8583 fields.
// Tags must be contained in `iso8583:"xxx"` with each argument separated
// by ',', for example: `iso8583:"19,length:5,encoding:ascii"`
// - name: The name of the field must be numeric of exception of the "bitmap"
// (first bitmap) and "mti". `iso8583:"bitmap"` or `iso8583:"2"`
// - length: arrives to the UnmarshalISO8583 method through parameter.
// In case of bitmap it indicates the amount of representative bits contained
// by the bitmap. For example in a classic 8 byte bitmap it would be 64,
// the same bitmap in hex string can occupy 16 bytes but would anyway contain
// 64 representative bits. For example: `iso8583:"length:10"`
// - encoding: arrives to the UnmarshalISO8583 method through parameter.
// For example: `iso8583:"encoding:ascii"`
type Unmarshaler interface {
	// UnmarshalISO8583 should not modify the input slice.
	// The input is the message starting from first yet not consumed byte to final.
	// Returns the amount of bytes consumed
	UnmarshalISO8583(b []byte, length int, encoding string) (n int, err error)
}

// Unmarshaler interface for iso8583 bitmaps.
//
// Length tag: It indicates the amount of representative bits contained
// by the bitmap. For example in a classic 8 byte bitmap it would be 64,
// the same bitmap in hex string can occupy 16 bytes but would anyway contain
// 64 representative bits.
type UnmarshalerBitmap interface {
	Unmarshaler

	// Bits should always return which bits are on what not necessarily is equal to what fields are present.
	// For example: In a traditional 8 byte bitmap. Te presence of field 65 would be represented in with bit 1 on.
	// Not present keys are assumed false.
	Bits() (m map[int]bool, err error)
}

// Unmarshal parses the ISO-8583 data and stores the result
// in the struct pointed to by v. If v is nil or not a pointer to a struct,
// Unmarshal returns an error.
//
// For all fields in the message must exist a field in the destination struct,
// otherwise an error is returned because there is no way to know the length
// of the field and therefore the start/end index of the following one.
//
// Unmarshal do not modify data input.
// returns the amount of bytes consumed from original message. If unused bytes remain from input
// its not considerate an error.
// If an error is encountered a counter with consumed bytes up to the moment is returned.
func Unmarshal(data []byte, v interface{}) (int, error) {
	input := reflect.ValueOf(v)
	// bitnapN works like an index that allows to know which fields are already
	// mapped to a bitmap. For example: If bitmap = 10 means that all fields from 1 to 10
	// are contemplated by already unmarshaled bitmaps and the next bitmap should consider
	// from 11 onwards.
	var bitmapN int

	// Because only each implementation of iso8583.Unmarshaler know how many bytes it needs to consume
	// and can potentially only know it while reading the field raw data; all currently un-read bytes
	// are given to the fields which return how many bytes from the message it consumed, value that is
	// subtracted from unconsumedBytes.
	unconsumedBytes := make([]byte, len(data))
	copy(unconsumedBytes, data)

	// consumed returns the amount of consumed bytes
	consumed := func() int {
		return len(data) - len(unconsumedBytes)
	}

	// consume reduces unconsumedBytes by the indicated amount with pointer safety.
	consume := func(m int, f string) error {
		if m > len(unconsumedBytes) {
			return fmt.Errorf("iso8583.unmarshal: Unmarshaler from field %s returned a n higher than uncosumed "+
				"bytes", f)
		}
		unconsumedBytes = unconsumedBytes[m:]
		return nil
	}

	// mustLen returns the capacity of a field ignoring all possible errors.
	mustLen := func(f string) int {
		_, tag, _ := getUnmarshaler(input, f)
		l, _ := tag.lenINT()
		return l
	}

	// Only pointer to structs can be unmarshaled.
	if !isPointerToStruct(input) {
		return consumed(), errors.New("iso8583.unmarshal: interface input is not a pointer to a structure")
	}

	// Unmarshal MTI.
	m, err := unmarshalField(input, _tagMTI, unconsumedBytes)
	if err != nil {
		return consumed(), err
	}

	if err := consume(m, _tagMTI); err != nil {
		return consumed(), err
	}

	// Unmarshal first bitmap.
	m, err = unmarshalField(input, _tagBITMAP, unconsumedBytes)
	if err != nil {
		return consumed(), err
	}

	if err := consume(m, _tagBITMAP); err != nil {
		return consumed(), err
	}

	// Get Bits method from first bitmap type to know which fields
	// to expect.
	firstFieldsMethod, err := getBitsMethod(input, _tagBITMAP)
	if err != nil {
		return consumed(), err
	}

	// Execute Bits method.
	firstFieldsList, err := firstFieldsMethod()
	if err != nil {
		return consumed(), fmt.Errorf("iso8583.unmarshal: failed reading first bitmap: %w", err)
	}

	// Obtain the length tag from first bitmap to know how many
	// fields presences are indicated by him.
	// These operation are safe because there were executed during unmarshaling.
	bitmapN += mustLen(_tagBITMAP)

	// This block is needed because if we use the original "firstFieldList" variable, we would copy only the reference
	// to the map and modify the value of the "bitmap" field.
	fieldList := make(map[int]bool)
	for k, v := range firstFieldsList {
		fieldList[k] = v
	}

	// Iterate over all possible fields up to the currently highest indicated by bitmaps
	// highestValue return might vary over iterations.
	for n := 1; n <= highestValue(fieldList); n++ {
		// Field is considerate only if its present and ON in the bitmaps.
		fieldExist, ok := fieldList[n]
		if !ok || !fieldExist {
			continue
		}

		// Unmarshal current field.
		m, err := unmarshalField(input, strconv.Itoa(n), unconsumedBytes)
		if err != nil {
			return consumed(), err
		}

		if err := consume(m, strconv.Itoa(n)); err != nil {
			return consumed(), err
		}

		// Check if current field is a bitmap.
		if bitsMethod, err := getBitsMethod(input, strconv.Itoa(n)); bitsMethod != nil && err == nil {
			// If bits method is present current field is a bitmap and the method is executed.
			fieldsListExpansion, err := bitsMethod()
			if err != nil {
				return consumed(), fmt.Errorf("iso8583.unmarshal: failed reading field %v bitmap: %w", n, err)
			}

			// Expand current field list with obtained information from current field.
			expandFieldList(fieldList, fieldsListExpansion, bitmapN)

			// Obtain the length tag from first bitmap to know how many
			// fields presences are indicated by him.
			// These operation are safe because there were executed during unmarshaling.
			// The value is added to index of mapped fields.
			bitmapN += mustLen(strconv.Itoa(n))
		}
	}

	// Return the amount of consumed fields.
	return consumed(), nil
}

// unmarshalField and save the value. Returns the amount of consumed bytes.
func unmarshalField(strct reflect.Value, fieldName string, bytes []byte) (int, error) {
	fieldInterface, tag, err := getUnmarshaler(strct, fieldName)
	if err != nil {
		return 0, err
	}

	if fieldInterface == nil {
		return 0, fmt.Errorf("iso8583.unmarshal: unknown field in message '%v', cant resolve upcomming fields",
			fieldName)
	}

	return executeUnmarshal(fieldInterface, bytes, *tag)
}

// executeUnmarshal calls unmarshal method of objective, obtaining parameters from tags.
// Returns consumed bytes from implementation.
func executeUnmarshal(field Unmarshaler, b []byte, tag tags) (int, error) {
	// Transform length from string to int type.
	length, err := tag.lenINT()
	if err != nil {
		return 0, err
	}

	// Execute unmarshal.
	n, err := field.UnmarshalISO8583(b, length, tag.Encoding)
	if err != nil {
		return 0, fmt.Errorf("iso8583.unmarshal: cant unmarshal field %v: %w", tag.Field, err)
	}

	// return consumed bytes.
	return n, nil
}

// getUnmarshaler iterates over target struct until it finds a field that matches with the field name.
// If field is not found all values are returned in nil.
func getUnmarshaler(v reflect.Value, n string) (Unmarshaler, *tags, error) {
	var strct reflect.Value
	var tag *tags

	// Obtain underlying struct.
	for strct = v; strct.Kind() == reflect.Ptr; {
		strct = strct.Elem()
	}

	// Iterate over struct until a field match with the tag name.
	var vField reflect.Value
	for index, strctType := 0, strct.Type(); index < strctType.NumField(); index++ {
		if tags := readTags(strctType.Field(index).Tag); tags != nil && tags.Field == n {
			vField = strct.Field(index)
			if vField.Kind() != reflect.Ptr {
				vField = vField.Addr()
			}
			tag = tags
		}
	}

	// If field is not found all values are returned nil.
	if reflect.ValueOf(vField).IsZero() {
		// Field not in struct
		return nil, nil, nil
	}

	// founded field must implement Unmarshaler, otherwise an error is returned.
	field, isValidUnmarshaler := vField.Interface().(Unmarshaler)
	if !isValidUnmarshaler {
		return nil, nil, fmt.Errorf(
			"iso8583.unmarshal: %s is present is struct but does not implement Unmarshaler interface", n)
	}

	// Return field data.
	return field, tag, nil
}

// getBitsMethod searches a iso8583.UnmarshalerBitmap implementation and returns it Bits method.
func getBitsMethod(strct reflect.Value, fieldName string) (func() (map[int]bool, error), error) {
	// Reuse getUnmarshaler because iso8583.Unmarshaler contains iso8583.UnmarshalerBitmap
	bitmapUnmarshaler, _, err := getUnmarshaler(strct, fieldName)
	if err != nil {
		return nil, err
	}

	// Field must be present.
	if bitmapUnmarshaler == nil {
		return nil, fmt.Errorf("iso8583.unmarshal: %s field is not present", fieldName)
	}

	// Field must implement iso8583.UnmarshalerBitmap.
	bitmap, ok := bitmapUnmarshaler.(UnmarshalerBitmap)
	if !ok {
		return nil, fmt.Errorf("iso8583.unmarshal: %s field is present but does not implement UnmarshalerBitmap",
			fieldName)
	}

	// return method.
	return bitmap.Bits, nil
}

func isPointerToStruct(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct
}

func expandFieldList(list map[int]bool, expansion map[int]bool, bitmapN int) map[int]bool {
	for fieldNumber, v := range expansion {
		list[fieldNumber+bitmapN] = v
	}
	return list
}

func highestValue(i map[int]bool) (h int) {
	for v := range i {
		if v > h {
			h = v
		}
	}
	return h
}
