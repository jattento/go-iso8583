package iso8583

import (
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
	Length    string
}

const _tagBITMAP = "bitmap"
const _tagMTI = "mti"

// If a nil pointer to tag struct is returned
// the current field should be omitted.
func readTags(tag reflect.StructTag) *tags {
	var output tags
	rawTags, exist := tag.Lookup("iso8583")
	if !exist {
		// Only fields with iso8583 tags are considerate.
		return nil
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
			output.Length = strings.TrimPrefix(tagBlock, "length:")
			continue
		}

		if strings.HasPrefix(tagBlock, "encoding") && len(strings.Split(tagBlock, ":")) == 2 {
			output.Encoding = strings.TrimPrefix(tagBlock, "encoding:")
			continue
		}

		output.Field = tagBlock
	}

	return &output
}

func (t tags) lenINT() (int, error) {
	l, err := strconv.Atoi(t.Length)
	if err != nil {
		return 0, fmt.Errorf("iso8583.marshal: field %s does not have a valid length", t.Field)
	}
	return l, nil
}

// fieldINT does not considerate the possibility of field names "bitmap" and "mti", which should be checked outside.
func (t tags) fieldINT() (int, error) {
	f, err := strconv.Atoi(t.Field)
	if err != nil {
		return 0, fmt.Errorf("iso8583.marshal: field %s does not have a valid field name", t.Field)
	}
	return f, nil
}
