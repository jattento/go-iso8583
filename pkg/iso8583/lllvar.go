package iso8583

// LLLVAR: For use of different encoding for 'LLL' and 'VAR' separate both encodings with a slash,
// where first element is the lll encoding and the second the var encoding.
// For Unmarshal length indicate the amount of byte that contain the LLL value
// For example:
// 	`iso8583:"2,length:3,encoding:ascii/ebcdic"`
type LLLVAR string

// MarshalISO8583 allows to use this type in structs and be able tu iso8583.Marshal it.
func (v LLLVAR) MarshalISO8583(length int, enc string) ([]byte, error) {
	return lengthMarshal(3, string(v), enc)
}

// UnmarshalISO8583 allows to use this type in structs and be able tu iso8583.Unmarshal it.
func (v *LLLVAR) UnmarshalISO8583(b []byte, length int, enc string) (int, error) {
	str, n, err := lengthUnmarshal(3, b, length, enc)
	*v = LLLVAR(str)
	return n, err
}
