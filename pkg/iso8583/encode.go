package iso8583

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Marshaler interface for iso8583 fields.
// Tags must be contained in `iso8583:"xxx"` with each argument separated
// by ',', for example: `iso8583:"19,length:5,encoding:ascii,omitempty"`
// - name: The name of the field must be numeric of exception of the "bitmap"
// (first bitmap) and "mti". `iso8583:"bitmap"` or `iso8583:"2"`
// - length: arrives to the MarshalISO8583 method through parameter.
// In case of bitmap it indicates the amount of representative bits contained
// by the bitmap. For example in a classic 8 byte bitmap it would be 64,
// the same bitmap in hex string can occupy 16 bytes but would anyway contain
// 64 representative bits. For example: `iso8583:"length:10"`
// - encoding: arrives to the MarshalISO8583 method through parameter.
// For example: `iso8583:"encoding:ascii"`
// - omitempty: if present field will be marshaled only if its not in zero value.
// This tag does not affect MarshalerBitmap.
// For example: `iso8583:"omitempty"`
// - disesteem: if present will be ignores by Marshal().
// For example: `iso8583:"-"`
type Marshaler interface {
	MarshalISO8583(length int, encoding string) ([]byte, error)
}

// MarshalerBitmap allows the bitmap to self charge, this means that a bitmap without this implementation should be
// charged while constructing the marshal objective struct, but if this interface is implemented by the bitmaps the field
// need only to be declared in the struct, later its loaded with the LoadBits method and marshaled.
// Note that if one of the bitmap of the main message implements this interface, all should do it, otherwise behaviour
// is unexpected
// Length tag must represent the amount of representative bits.
// Omitempty tag does not affect MarshalerBitmap.
// NOTE: This implementation is not a must in bitmaps.
type MarshalerBitmap interface {
	// b should always represent which bits are on what not necessarily is equal to what fields are present.
	// For example: In a traditional 8 byte bitmap. Te presence of field 65 would be represented in with bit 1 on.
	// len(b) == length tag.
	// the output is the bitmap in bytes ready to put in the message.
	MarshalISO8583Bitmap(b map[int]bool, encoding string) ([]byte, error)
}

// Marshal returns the JSON encoding of v.
//
// If field is a string with valid tags and does not implement Marshaler ascii encoding is assumed.
// If field is a []byte with valid tags and does not implement Marshaler, its content is used as value.
func Marshal(v interface{}) ([]byte, error) {
	var mti []byte
	var firstBitmap []byte
	var fields = make(map[int][]byte)
	var bitmaps = make(map[tags]MarshalerBitmap)

	// Obtain value and type of input.
	inputValue, inputType := reflectContext(v)

	// Input must be a struct or a pointer to one.
	if inputType.Kind() != reflect.Struct {
		return nil, errors.New("iso8583.marshal: input is not a struct or is pointing to one")
	}

	// Iterate over all fields of input struct.
	for index := 0; index < inputType.NumField(); index++ {
		fieldValue := inputValue.Field(index)
		fieldType := inputType.Field(index)

		// If field is nil its not considerated for marshaling.
		if isNil(fieldValue) {
			continue
		}

		// If field is a pointer Elem() its applied.
		for fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		// Only exported and non-anonymous field are considerated.
		if fieldType.PkgPath != "" || fieldType.Anonymous {
			continue
		}

		// Read tags from current field, if there aren't field is ignored.
		tags := readTags(fieldType.Tag)
		if tags == nil || tags.Disesteem {
			continue
		}

		// Bitmap fields are saved in a map, they must be marshaled at latest when all fields are known
		if bmapInterface, isBitmapInterface := fieldValue.Interface().(MarshalerBitmap); isBitmapInterface {
			bitmaps[*tags] = bmapInterface
			continue
		}

		// If omitempty tag is present and field value is zero value the field is ignored.
		if tags.OmitEmpty && fieldValue.IsZero() {
			continue
		}

		// Resolve field value that not necessary implements iso8583.Marshaler interface.
		b, err := resolveMarshalFieldValue(fieldValue, *tags)
		if err != nil {
			return nil, err
		}

		// Check if bytes are mti, firstBitmap or a field value
		if strings.ToLower(tags.Field) == _tagMTI {
			// Same field can't be two time in same message.
			if len(mti) > 0 {
				return nil, fmt.Errorf("iso8583.marshal: field %s is repeated", _tagMTI)
			}

			mti = b
			continue
		}

		// This block exist for cases where first bitmap does not implement iso8583.BitmapMarshaler.
		if strings.ToLower(tags.Field) == _tagBITMAP {
			// Same field can't be two time in same message.
			if len(firstBitmap) > 0 {
				return nil, fmt.Errorf("iso8583.marshal: field %s is repeated", _tagBITMAP)
			}

			firstBitmap = b
			continue
		}

		// If field isn't MTI or BITMAP its name must be numeric.
		fk, err := tags.fieldINT()
		if err != nil {
			return nil, err
		}

		// 0 is not a valid field key for iso fields.
		if fk == 0 {
			return nil, errors.New("iso8583.marshal: field 0 not allowed")
		}

		// Length of marshal result is check here because maybe under strange condition the implementation don't
		// want to send mti or first bitmap but in case of a field it must be discarded here to avoid break
		// bitmaps.
		if len(b) > 0 {
			if _, exist := fields[fk]; exist {
				// Same field can't be two time in same message.
				return nil, fmt.Errorf("iso8583.marshal: field %v is repeated", fk)
			}
			fields[fk] = b
		}
	}

	// Resolve all bitmaps that implement iso8583.MarshalerBitmap once all fields are known.
	if err := resolveBitmaps(fields, &firstBitmap, bitmaps); err != nil {
		return nil, err
	}

	// Order of bytes must be: mti -> first bitmap -> fields (include all other bitmaps).
	return append(append(append([]byte{}, mti...), firstBitmap...), fieldsToBytes(fields)...), nil
}

// reflectContext returns the underlying type and value.
// Elem() its applied up to to ground type.
func reflectContext(i interface{}) (reflect.Value, reflect.Type) {
	value := reflect.ValueOf(i)
	Type := reflect.TypeOf(i)

	// Check that input is a struct or is pointing to one
	for Type.Kind() == reflect.Ptr {
		Type = Type.Elem()
		value = value.Elem()
	}

	return value, Type
}

// resolveMarshalFieldValue resolves Marshal return value of a field that must not necessary be a marshaler.
func resolveMarshalFieldValue(v reflect.Value, tag tags) ([]byte, error) {
	marshaler, isMarshaler := v.Interface().(Marshaler)
	isBytes := v.Kind() == reflect.Slice && v.Type() == reflect.TypeOf([]byte(nil))
	isString := v.Kind() == reflect.String

	var length int
	if tag.Length != "" {
		var err error
		length, err = tag.lenINT()
		if err != nil {
			return nil, err
		}
	}

	// Priority of marshaling order is marshaler -> bytes -> string
	switch {
	case isMarshaler:
		b, err := marshaler.MarshalISO8583(length, tag.Encoding)
		if err != nil {
			return nil, fmt.Errorf("iso8583.marshal: field %s cant be marshaled: %w", tag.Field, err)
		}
		return b, nil
	case isBytes:
		return v.Bytes(), nil
	case isString:
		// ASCII assumed.
		return []byte(v.String()), nil
	}

	// value does not implement marshal interface nether is a byte slice
	return nil, fmt.Errorf("iso8583.marshal: field %s does not implement Marshaler interface, "+
		"is a string or slice of bytes but does have iso8583 tags", tag.Field)

}

// transform the fields map to []byte using the correct order
func fieldsToBytes(m map[int][]byte) []byte {
	fieldsByte := make([]byte, 0)

	// Every loop fiends the lowest key
	for n, lenFields := 0, len(m); n < lenFields; n++ {
		var lowestField int
		firstIteration := true

		for fieldKey := range m {
			if firstIteration {
				lowestField = fieldKey
				firstIteration = false
				continue
			}
			if fieldKey < lowestField {
				lowestField = fieldKey
			}
		}

		// Value from lowest key is added at the end of the slice and key is deleted from map
		// to avoid duplicity.
		fieldsByte = append(fieldsByte, m[lowestField]...)
		delete(m, lowestField)
	}

	return fieldsByte
}

// Resolve bitmaps marshal values.
// Reads bitmaps from "bitmaps" parameter and save them in "fields" and "firstBitmap" variables.
func resolveBitmaps(fields map[int][]byte, firstBitmap *[]byte, bitmaps map[tags]MarshalerBitmap) error {
	type bitmap struct {
		tag      tags
		value    MarshalerBitmap
		capacity int
	}

	// Create slice of bitmaps and
	var orderedBitmaps []bitmap
	for t, m := range bitmaps {
		orderedBitmaps = append(orderedBitmaps, bitmap{tag: t, value: m})
	}

	sort.SliceStable(orderedBitmaps, func(i, j int) bool {
		// first bitmap must be always element 0.
		if orderedBitmaps[i].tag.Field == _tagBITMAP {
			return true
		}
		if orderedBitmaps[j].tag.Field == _tagBITMAP {
			return false
		}

		// These errors are not accessible by the public api of this package
		ii, _ := strconv.Atoi(orderedBitmaps[i].tag.Field)
		ji, _ := strconv.Atoi(orderedBitmaps[j].tag.Field)

		return ji > ii
	})

	// Add capacity to every bitmap struct
	for n := 0; n < len(orderedBitmaps); n++ {
		l, err := strconv.Atoi(orderedBitmaps[n].tag.Length)
		if err != nil {
			return fmt.Errorf("iso8583.marshal: field %s is implements bitmap interface and "+
				"does not have a valid length", orderedBitmaps[n].tag.Field)
		}
		orderedBitmaps[n].capacity = l
	}

	// Resolve all bitmaps starting from last one to first,
	// this is because every bitmaps indicates the presence of the next one.
	for n := len(orderedBitmaps) - 1; n >= 0; n-- {
		startPosition := 1
		for m := n - 1; m >= 0; m-- {
			startPosition += orderedBitmaps[m].capacity
		}

		// Create new map and only add elements that apply to the current bitmap.
		present := make(map[int]bool)
		for f := range fields {
			if f >= startPosition && f < startPosition+orderedBitmaps[n].capacity {
				present[f-startPosition+1] = true
			}
		}

		// All none selected bits under the capacity are setted false.
		for m := 1; m <= orderedBitmaps[n].capacity; m++ {
			if _, exist := present[m]; !exist {
				present[m] = false
			}
		}

		// Marshal bitmap...
		b, err := orderedBitmaps[n].value.MarshalISO8583Bitmap(present, orderedBitmaps[n].tag.Encoding)
		if err != nil {
			return fmt.Errorf("iso8583.marshal: field %s cant be marshaled: %w", orderedBitmaps[n].tag.Field, err)
		}

		if len(b) == 0 {
			// If length is 0m the field is ommited.
			continue
		}

		if orderedBitmaps[n].tag.Field == _tagBITMAP {
			// If field is first bitmap, its safed in indicated parameter pointer.
			*firstBitmap = b
			continue
		}

		// For all other cases, its added in the map along with all other fields.
		fk, err := strconv.Atoi(orderedBitmaps[n].tag.Field)
		if err != nil {
			return fmt.Errorf("iso8583: unrecognized field: %s", orderedBitmaps[n].tag.Field)
		}
		if fk == 0 {
			return errors.New("iso8583.marshal: field 0 not allowed")
		}

		if _, exist := fields[fk]; exist {
			return errors.New("iso8583.marshal: field can't be repeated")
		}

		fields[fk] = b
	}

	return nil
}

func isNil(v reflect.Value) bool {
	if k := v.Kind(); k == reflect.Ptr || k == reflect.Map || k == reflect.Chan || k == reflect.Func ||
		k == reflect.UnsafePointer || k == reflect.Interface || k == reflect.Slice {
		return v.IsNil()
	}

	return false
}
