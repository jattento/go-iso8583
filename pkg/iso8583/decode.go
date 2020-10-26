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
//
// If you want to add a new encoding to use in the inbuilt types, just add them to the iso8583.UnmarshalDecodings
// variable.
//
// If you want to add a completely new field, just copy the most similar from the existing one
// and modify what ever you want.
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

var errStructFieldNonExistent = errors.New("non existent struct field")

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
	strctInput := reflect.ValueOf(v)
	// bitnapN works like an index that allows to know which fields are already
	// mapped to a bitmap. For example: If bitmap = 10 means that all fields from 1 to 10
	// are contemplated by already unmarshaled bitmaps and the next bitmap should consider
	// from 11 onwards.
	var bitmapN int

	// Because only each implementation of iso8583.Unmarshaler know how many bytes it needs to consume
	// and can potentially only know it while reading the field raw data; all currently un-read bytes
	// are given to the fields which return how many bytes from the message it consumed.
	buffer := newUnmarshalBuffer(data)

	// Only pointer to structs can be unmarshaled.
	if !isPointerToStruct(strctInput) {
		return buffer.UntilNowConsumed(), errors.New("iso8583.unmarshal: interface input is not a pointer to a structure")
	}

	representativeBits, err := readMessageHeader(buffer, strctInput)
	if err != nil {
		return buffer.UntilNowConsumed(), err
	}

	// Obtain the length tag from first bitmap to know how many
	// fields presences are indicated by him.
	bitmapN += representativeBits

	// Get Bits method from first bitmap type to know which fields
	// to expect.
	firstFieldsMethod, err := getBitsMethod(strctInput, _tagBITMAP)
	if err != nil {
		return buffer.UntilNowConsumed(), err
	}

	// Execute Bits method.
	firstFieldsList, err := firstFieldsMethod()
	if err != nil {
		return buffer.UntilNowConsumed(), fmt.Errorf("iso8583.unmarshal: failed reading first bitmap: %w", err)
	}

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
		m, tagValues, err := unmarshalField(strctInput, strconv.Itoa(n), buffer.Bytes())
		if err != nil {
			return buffer.UntilNowConsumed(), err
		}

		if err := buffer.IncrementConsumedCounter(m, "struct "+strconv.Itoa(n)); err != nil {
			return buffer.UntilNowConsumed(), err
		}

		// Check if current field is a bitmap.
		if bitsMethod, err := getBitsMethod(strctInput, strconv.Itoa(n)); bitsMethod != nil && err == nil {
			// If bits method is present current field is a bitmap and the method is executed.
			fieldsListExpansion, err := bitsMethod()
			if err != nil {
				return buffer.UntilNowConsumed(), fmt.Errorf("iso8583.unmarshal: failed reading field %v bitmap: %w", n, err)
			}

			// Expand current field list with obtained information from current field.
			expandFieldList(fieldList, fieldsListExpansion, bitmapN)

			// Obtain the length tag from first bitmap to know how many
			// fields presences are indicated by him.
			// The value is added to index of mapped fields.
			bitmapN += tagValues.Length
		}
	}

	// Return the amount of consumed fields.
	return buffer.UntilNowConsumed(), nil
}

type unmarshalBuffer struct {
	data        []byte
	placeholder int
}

func newUnmarshalBuffer(data []byte) *unmarshalBuffer {
	buffer := unmarshalBuffer{data: make([]byte, len(data))}
	copy(buffer.data, data)
	return &buffer
}

// incrementConsumed increments the placeholder considering the data length.
func (buffer *unmarshalBuffer) IncrementConsumedCounter(n int, fieldName string) error {
	if len(buffer.data[buffer.placeholder:]) < n {
		// This error can only be caused by bad unmarshal implementations
		return fmt.Errorf("iso8583.unmarshal: Unmarshaler from field %s returned a n higher than unconsumed "+
			"bytes", fieldName)
	}

	buffer.placeholder += n
	return nil
}

// UntilNowConsumed returns the amount of until now consumed bytes
func (buffer *unmarshalBuffer) UntilNowConsumed() int { return buffer.placeholder }

// Bytes returns a copy from the input data bytes remainder.
func (buffer *unmarshalBuffer) Bytes() []byte {
	data := make([]byte, len(buffer.data[buffer.placeholder:]))
	copy(data, buffer.data[buffer.placeholder:])
	return data
}

// readMessageHeader reads the message MTI and the first bitmap.
// Returns the amount of representative bits of the bitmap.
func readMessageHeader(buffer *unmarshalBuffer, strct reflect.Value) (int, error) {
	// Unmarshal MTI.
	consumed, _, err := unmarshalField(strct, _tagMTI, buffer.Bytes())
	if err != nil {
		return 0, err
	}

	if err := buffer.IncrementConsumedCounter(consumed, _tagMTI); err != nil {
		return 0, err
	}

	return readFirstBitmap(buffer, strct)
}

// readFirstBitmap reads the first bitmap and returns the amount of representative bits.
func readFirstBitmap(buffer *unmarshalBuffer, strct reflect.Value) (int, error) {
	// Unmarshal first bitmap.
	consumed, tagValues, err := unmarshalField(strct, _tagBITMAP, buffer.Bytes())
	if err != nil {
		return 0, err
	}

	if err := buffer.IncrementConsumedCounter(consumed, _tagBITMAP); err != nil {
		return 0, err
	}

	return tagValues.Length, nil
}

// unmarshalField and save the value. Returns the amount of consumed bytes.
func unmarshalField(strct reflect.Value, fieldName string, bytes []byte) (int, tags, error) {
	fieldValue, tag, err := searchStructField(strct, fieldName)
	if err != nil {
		if errors.Is(err, errStructFieldNonExistent) {
			err = fmt.Errorf("unknown field in message '%v', cant resolve upcomming fields",
				fieldName)
		}

		return 0, tags{}, fmt.Errorf("iso8583.unmarshal: %w", err)
	}

	// founded field must implement Unmarshaler, otherwise an error is returned.
	fieldInterface, isValidUnmarshaler := fieldValue.Interface().(Unmarshaler)
	if !isValidUnmarshaler {
		return 0, tags{}, fmt.Errorf(
			"iso8583.unmarshal: field %s is present but does not implement Unmarshaler interface", fieldName)
	}

	consumed, err := executeUnmarshal(fieldInterface, bytes, tag)

	return consumed, tag, err
}

// executeUnmarshal calls unmarshal method of objective, obtaining parameters from tags.
// Returns consumed bytes from implementation.
func executeUnmarshal(field Unmarshaler, b []byte, tag tags) (int, error) {
	// Execute unmarshal.
	n, err := field.UnmarshalISO8583(b, tag.Length, tag.Encoding)
	if err != nil {
		return 0, fmt.Errorf("iso8583.unmarshal: cant unmarshal field %v: %w", tag.Field, err)
	}

	// return consumed bytes.
	return n, nil
}

// Returns the reflect value of the indicated field with its tags.
// If v isn't a struct this function panics.
// Returns errStructFieldNonExistent if the field is non existing.
func searchStructField(v reflect.Value, n string) (reflect.Value, tags, error) {
	var strct reflect.Value
	var tag tags

	// Obtain underlying struct.
	for strct = v; strct.Kind() == reflect.Ptr; {
		strct = strct.Elem()
	}

	// Iterate over struct until a field match with the tag name.
	var vField reflect.Value
	for index, strctType := 0, strct.Type(); index < strctType.NumField(); index++ {
		if structFieldValue, structFieldTags, err := getStructFieldData(strct, index); structFieldTags.Field == n {
			if err != nil {
				return reflect.Value{}, tags{}, err
			}

			if !reflect.ValueOf(vField).IsZero() {
				return reflect.Value{}, tags{}, fmt.Errorf("field %v is repeteated in struct", n)
			}

			if structFieldValue.Kind() != reflect.Ptr {
				structFieldValue = structFieldValue.Addr()
			}

			vField = structFieldValue
			tag = structFieldTags
		}
	}

	// If field is not found all values are returned nil.
	if reflect.ValueOf(vField).IsZero() {
		// Field not in struct
		return reflect.Value{}, tags{}, errStructFieldNonExistent
	}

	return vField, tag, nil
}

// getStructFieldData returns the specified struct field data using the struct index.
func getStructFieldData(parentStructValue reflect.Value, index int) (reflect.Value, tags, error) {
	// search tags from current field, if there aren't field is ignored.
	tags, err := searchTags(parentStructValue.Type().Field(index))

	return parentStructValue.Field(index), tags, err
}

// getBitsMethod searches a iso8583.UnmarshalerBitmap implementation and returns it Bits method.
func getBitsMethod(strct reflect.Value, fieldName string) (func() (map[int]bool, error), error) {
	// Reuse getUnmarshaler because iso8583.Unmarshaler contains iso8583.UnmarshalerBitmap
	bitmapUnmarshaler, _, _ := searchStructField(strct, fieldName)

	// Field must implement iso8583.UnmarshalerBitmap.
	bitmap, ok := bitmapUnmarshaler.Interface().(UnmarshalerBitmap)
	if !ok {
		return nil, fmt.Errorf("iso8583.unmarshal: %s field is present but does not implement UnmarshalerBitmap",
			fieldName)
	}

	// return method.
	return bitmap.Bits, nil
}

func isPointerToStruct(v reflect.Value) bool {
	underlyingObj := v
	for underlyingObj.Kind() == reflect.Ptr {
		underlyingObj = underlyingObj.Elem()
	}

	return v.Kind() == reflect.Ptr && underlyingObj.Kind() == reflect.Struct
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
