package iso8583

import (
	"fmt"
	"strconv"
	"strings"
)

// LengthMarshal receives the expected amount of "L" the content (already encoded) that comes after L and a encoding
// it returns the result bytes after combines the L and the value.
func LengthMarshal(l int, v []byte, enc string) ([]byte, error) {
	varContent := v
	llEncoding := enc

	llValue := strconv.Itoa(len(varContent))
	if len(llValue) > l {
		return nil, fmt.Errorf("content length exceeded the %s limit for %s elements",
			strings.Repeat("9", l), strings.Repeat("L", l))
	}

	for len(llValue) < l {
		llValue = "0" + llValue
	}

	llContent, err := applyEncoding([]byte(llValue), llEncoding, MarshalEncodings)
	if err != nil {
		return nil, err
	}

	return append(llContent, varContent...), nil
}

// LengthUnmarshal receives the amount of "L", the source bytes, the amount of bytes to read, and a encoding;
// it returns the amount of bytes readed, the actually value bytes and a error.
func LengthUnmarshal(l int, b []byte, length int, enc string) (int, []byte, error) {
	if len(b) < length {
		return 0, nil, fmt.Errorf("message remain (%v bytes) is shorter than %s byte length (%v)",
			len(b), strings.Repeat("L", l), length)
	}

	llContent := make([]byte, length)
	copy(llContent, b[:length])

	llEncoding := enc

	llContent, err := applyEncoding(llContent, llEncoding, UnmarshalDecodings)
	if err != nil {
		return 0, nil, err
	}

	llValue, err := strconv.Atoi(string(llContent))
	if err != nil {
		return 0, nil, fmt.Errorf("obtained %s after decoding is not a valid integer: %v",
			strings.Repeat("L", l), string(llContent))
	}

	if len(b)-length < llValue {
		return 0, nil, fmt.Errorf("message remain (%v bytes) is shorter than %s indicated length (%v)",
			len(b)-length, strings.Repeat("L", l), llValue)
	}

	varContent := make([]byte, llValue)
	copy(varContent, b[length:length+llValue])

	return length + llValue, varContent, nil
}

// ReadSplitEncodings returns two copies of str or if it contains a encoding separator "/" it returns
// bots encodings splitted.
func ReadSplitEncodings(str string) (string, string) {
	if strings.Contains(str, "/") {
		splitStr := strings.Split(str, "/")
		return splitStr[0], splitStr[1]
	}

	return str, str
}
