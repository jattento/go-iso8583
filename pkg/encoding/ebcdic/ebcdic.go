package ebcdic

type version struct {
	entries []entry

	// Maps are cached for better performance
	encoding map[rune]byte
	decoding map[byte]rune
}

type entry struct {
	EBCDIC         byte
	ASCII          byte
	representation rune
	description    string
}

const NULL = 0x0

// FromGoString converts a string to ebcdic bytes.
func (v *version) FromGoString(s string) []byte {
	if len(v.encoding) == 0 {
		v.generateEncodingMap()
	}

	rs := []rune(s)
	output := make([]byte, 0)
	for _, r := range rs {
		func() {
			for k, v := range v.encoding {
				if k == r {
					output = append(output, v)
					return
				}
			}
			output = append(output, NULL)
		}()
	}
	return output
}

// FromGoString converts a ebcdic bytes to a string.
func (v *version) ToGoString(b []byte) string {
	if len(v.decoding) == 0 {
		v.generateDecodingMap()
	}

	output := make([]rune, 0)
	for _, byt := range b {
			for k, v := range v.decoding {
				if k == byt {
					output = append(output, v)
					break
				}
			}
	}
	return string(output)
}

func (v *version) generateEncodingMap() {
	v.encoding = make(map[rune]byte)

	for _, kv := range v.entries {
		if kv.representation != 0 {
			v.encoding[kv.representation] = kv.EBCDIC
			continue
		}
		v.encoding[rune(kv.ASCII)] = kv.EBCDIC
	}
}

func (v *version) generateDecodingMap() {
	v.decoding = make(map[byte]rune)

	for _, kv := range v.entries {
		if kv.representation != 0 {
			v.decoding[kv.EBCDIC] = kv.representation
			continue
		}
		v.decoding[kv.EBCDIC] = rune(kv.ASCII)
	}
}
