package iso8583

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type tags struct {
	Field     string
	OmitEmpty bool
	Disesteem bool
	Encoding  string
	Length    int
}

const _tagBITMAP = "bitmap"
const _tagMTI = "mti"

var (
	errUnexportedField = errors.New("unexported field")
	errAnonymousField  = errors.New("anonymous field")
	errTagsNotFound    = errors.New("ISO8583 tags not found")
)

func searchTags(field reflect.StructField) (tags, error) {
	// Only exported and non-anonymous field are considered.
	if field.PkgPath != "" {
		return tags{}, errUnexportedField
	}

	if field.Anonymous {
		return tags{}, errAnonymousField
	}

	return readTags(field.Tag)
}

// If a nil pointer is returned the ISO8583 tag is not present.
func readTags(tag reflect.StructTag) (tags, error) {
	var output tags
	var returnErr error
	rawTags, exist := tag.Lookup("iso8583")
	if !exist {
		// Only fields with iso8583 tags are considerate.
		return tags{}, errTagsNotFound
	}

	// Search for provided tags...
	for _, tagBlock := range strings.Split(rawTags, ",") {

		if tagBlock == "-" {
			output.Disesteem = true
			continue
		}

		if tagBlock == "omitempty" {
			output.OmitEmpty = true
			continue
		}

		if strings.HasPrefix(tagBlock, "length") && len(strings.Split(tagBlock, ":")) == 2 {
			strLenght := strings.TrimPrefix(tagBlock, "length:")

			l, err := strconv.Atoi(strLenght)
			if err != nil {
				returnErr = fmt.Errorf("invalid length: %w", err)
			}

			output.Length = l

			continue
		}

		if strings.HasPrefix(tagBlock, "encoding") && len(strings.Split(tagBlock, ":")) == 2 {
			output.Encoding = strings.TrimPrefix(tagBlock, "encoding:")
			continue
		}

		output.Field = tagBlock
	}

	if output.Field == "" {
		return output, errors.New("exported struct field contains ISO8583 tag but no field name")
	}

	if returnErr != nil {
		return output, fmt.Errorf("field %s: %w", output.Field, returnErr)
	}

	return output, nil
}
