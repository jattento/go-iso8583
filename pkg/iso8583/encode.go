package iso8583

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
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
//
// If you want to add a new encoding to use in the inbuilt types, just add them to the iso8583.MarshalEncodings
// variable.
//
// If you want to add a completely new field, just copy the most similar from the existing one
// and modify what ever you want.
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
	// MarshalISO8583Bitmap must consume b map information and return the corresponding byte slice.
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
	processedFields := make(map[string]struct{})

	if v == nil {
		return nil, errors.New("iso8583.marshal: nil input")
	}

	// Obtain value and type of input.
	inputValue := reflect.ValueOf(v)
	for inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}

	// Input must be a struct or a pointer to one.
	if inputValue.Kind() != reflect.Struct {
		return nil, errors.New("iso8583.marshal: input is not a struct or is pointing to one")
	}

	msg := newMarshalerMessage()

	// Iterate over all fields of input struct.
	for index := 0; index < inputValue.Type().NumField(); index++ {
		structFieldValue, tag, err := getStructFieldData(inputValue, index)
		if errors.Is(err, errUnexportedField) || errors.Is(err, errAnonymousField) || errors.Is(err, errTagsNotFound) ||
			tag.Disesteem || isNil(structFieldValue) {
			continue
		}

		if err != nil {
			return nil, fmt.Errorf("iso8583.marshal: %w", err)
		}

		if tag.Field != _tagMTI && tag.Field != _tagBITMAP {
			// Validate field names
			n, err := strconv.Atoi(tag.Field)
			if err != nil || n < 1 {
				return nil, fmt.Errorf("iso8583.marshal: invalid field name: %s", tag.Field)
			}
		}

		// Check if field is repeated
		if _, existsAlready := processedFields[tag.Field]; existsAlready {
			return nil, fmt.Errorf("iso8583.marshal: field %s is repeated", tag.Field)
		}

		processedFields[tag.Field] = struct{}{}

		// Bitmap fields are saved in a map, they must be marshaled at latest when all fields are known
		bmapInterface, isBitmapInterface := structFieldValue.Interface().(MarshalerBitmap)
		if isBitmapInterface {
			msg.addBitmap(bmapInterface, tag)
			continue
		}

		// If omitempty tag is present and field value is zero value the field is ignored.
		if tag.OmitEmpty && structFieldValue.IsZero() {
			continue
		}

		msg.addField(structFieldValue, tag)
	}

	return msg.Bytes()
}

// resolveMarshalFieldValue resolves Marshal return value of a field that must not necessary be a marshaler.
func resolveMarshalFieldValue(v reflect.Value, tag tags) ([]byte, error) {
	marshaler, isMarshaler := v.Interface().(Marshaler)

	// Priority of marshaling order is marshaler -> bytes -> string
	if !isMarshaler {
		return nil, fmt.Errorf("iso8583.marshal: field %s does not implement Marshaler interface "+
			"but does have iso8583 tags", tag.Field)
	}

	b, err := marshaler.MarshalISO8583(tag.Length, tag.Encoding)
	if err != nil {
		return nil, fmt.Errorf("iso8583.marshal: field %s cant be marshaled: %w", tag.Field, err)
	}

	return b, nil
}

// isNil Checks the kind of the value since reflect.IsNil method could panic at some.
func isNil(v reflect.Value) (isNil bool) {
	reflectKind := v.Kind()

	for _, kind := range []reflect.Kind{
		reflect.Ptr, reflect.Map, reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Interface, reflect.Slice,
	} {
		if reflectKind == kind {
			isNil = v.IsNil()
		}
	}

	return isNil
}

// field represents a iso8583 field in bytes format
type field struct {
	name  string
	bytes []byte
}

// marshalerMessage represents a iso8583 message before its marshaled
type marshalerMessage struct {
	Bitmaps []isoMarshalerBitmap
	Fields  []isoMarshalerField
}

// isoMarshalerField represents a iso8583 field before its marshaled
type isoMarshalerField struct {
	Marshaler reflect.Value
	tags
}

// isoMarshalerBitmap represents a iso8583 field (that implements the marshaler bitmap interface) before its marshaled
type isoMarshalerBitmap struct {
	MarshalerBitmap
	tags
}

// newMarshalerMessage returns a empty message
func newMarshalerMessage() *marshalerMessage {
	return &marshalerMessage{
		Fields:  make([]isoMarshalerField, 0),
		Bitmaps: make([]isoMarshalerBitmap, 0),
	}
}

// addField adds a field to the current message
func (m *marshalerMessage) addField(marsh reflect.Value, tag tags) {
	m.Fields = append(m.Fields, isoMarshalerField{Marshaler: marsh, tags: tag})
}

// addBitmap add a bitmap to the current message
func (m *marshalerMessage) addBitmap(marsh MarshalerBitmap, tag tags) {
	m.Bitmaps = append(m.Bitmaps, isoMarshalerBitmap{MarshalerBitmap: marsh, tags: tag})
}

// Bytes returns the actually ISO8583 formatted message.
func (m *marshalerMessage) Bytes() (messageBytes []byte, returnErr error) {
	var (
		mtiPresent         bool
		firstBitmapPresent bool
	)

	fields := make([]field, 0)

	// Iterate over all fields that NOT implement marshaler bitmap
	for _, f := range m.Fields {
		b, err := resolveMarshalFieldValue(f.Marshaler, f.tags)
		if err != nil {
			return nil, err
		}

		if len(b) > 0 {
			if f.Field == _tagMTI {
				mtiPresent = true
			}

			if f.Field == _tagBITMAP {
				firstBitmapPresent = true
			}

			fields = append(fields, field{name: f.Field, bytes: b})
		}
	}

	// Resolve bitmaps; starting from last to first.
	firstBmapInMarshaler, err := m.resolveMarshalerBitmaps(&fields)
	if err != nil {
		return nil, err
	}

	firstBitmapPresent = firstBmapInMarshaler || firstBitmapPresent

	// Sort fields
	sortFieldsStable(fields, func(index int) string { return fields[index].name })

	// Build message
	messageBytes = make([]byte, 0)
	for _, f := range fields {
		messageBytes = append(messageBytes, f.bytes...)
	}

	// Valdiate MTI and first bitmap presence
	if !firstBitmapPresent {
		returnErr = errors.New("iso8583.marshal: no first bitmap was generated")
	}

	if !mtiPresent {
		returnErr = errors.New("iso8583.marshal: no MTI was generated")
	}

	return messageBytes, returnErr
}

// Resolve bitmaps marshal values.
// Reads bitmaps from "bitmaps" parameter and save them in "fields" and "firstBitmap" variables.
func (m *marshalerMessage) resolveMarshalerBitmaps(fields *[]field) (bool, error) {
	var firstBitmapPresent bool

	sortFieldsStable(m.Bitmaps, func(index int) string { return m.Bitmaps[index].Field })

	// If bitmap length is not indicated: we assume its 64.
	for n := 0; len(m.Bitmaps) > n; n++ {
		if m.Bitmaps[n].tags.Length == 0 {
			m.Bitmaps[n].tags.Length = 64
		}
	}

	// Resolve all bitmaps starting from last one to first,
	// this is because every bitmaps indicates the presence of the next one.
	for n := len(m.Bitmaps) - 1; n >= 0; n-- {
		// Marshal bitmap...

		b, err := m.Bitmaps[n].MarshalerBitmap.MarshalISO8583Bitmap(m.createBitmapMarshalerInput(*fields, n), m.Bitmaps[n].tags.Encoding)
		if err != nil {
			return false, fmt.Errorf("iso8583.marshal: field %s cant be marshaled: %w", m.Bitmaps[n].tags.Field, err)
		}

		if len(b) > 0 {
			if m.Bitmaps[n].Field == _tagBITMAP {
				firstBitmapPresent = true
			}

			*fields = append(*fields, field{name: m.Bitmaps[n].Field, bytes: b})
		}
	}

	return firstBitmapPresent, nil
}

func (m *marshalerMessage) createBitmapMarshalerInput(fields []field, bitmapIndex int) map[int]bool {
	startPosition := 1
	for nn := bitmapIndex - 1; nn >= 0; nn-- {
		startPosition += m.Bitmaps[nn].tags.Length
	}

	// Create new map and only add elements that apply to the current bitmap.
	present := make(map[int]bool)
	for _, f := range fields {
		if f.name == _tagMTI {
			continue
		}

		// Already checked
		numericName, _ := strconv.Atoi(f.name)

		if numericName >= startPosition && numericName < startPosition+m.Bitmaps[bitmapIndex].Length {
			present[numericName-startPosition+1] = true
		}
	}

	// All none selected bits under the capacity are setted false.
	for nn := 1; nn <= m.Bitmaps[bitmapIndex].Length; nn++ {
		if _, exist := present[nn]; !exist {
			present[nn] = false
		}
	}

	return present
}

// sortFieldsStable sorts message fields in the following order: MTI -> BITMAP -> field 0 -> field n-1 -> field n
func sortFieldsStable(obj interface{}, getFieldName func(index int) string) {
	sort.SliceStable(obj, func(i, j int) bool {
		// first bitmap must be always element 0.
		if getFieldName(i) == _tagMTI || getFieldName(j) == _tagMTI {
			return getFieldName(i) == _tagMTI
		}

		// first bitmap must be always be second.
		if getFieldName(i) == _tagBITMAP || getFieldName(j) == _tagBITMAP {
			return getFieldName(i) == _tagBITMAP
		}

		// These errors are not accessible by the public api of this package
		ii, _ := strconv.Atoi(getFieldName(i))
		ji, _ := strconv.Atoi(getFieldName(j))

		return ji > ii
	})
}
