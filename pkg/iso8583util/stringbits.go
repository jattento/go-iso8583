package iso8583util

import (
	"fmt"
	"strings"
)

const _bitsInByte = 8

// StringBitsToBytes converts a bits string into []byte, is a character is different to space,
// 0 or 1 the function panics. If the amount of 0 and 1 %8 != 0, 0 are append up to the condition is false
func StringBitsToBytes(str string) []byte {
	s := strings.ReplaceAll(str, " ", "")

	for n := range s {
		if s[n] != '0' && s[n] != '1' {
			panic(fmt.Errorf("input character isn't 1 or 0, '%s' found", string(s[n])))
		}
	}

	ss := make([]string, 0)
	for n := 0; n <= len(s) || n%_bitsInByte != 1; n++ {
		if len(s) < n {
			s += "0"
		}

		if n%_bitsInByte == 0 && n != 0 {
			ss = append(ss, s[n-_bitsInByte:n])
		}
	}

	b := make([]byte, 0)
	for _, byteString := range ss {
		var byt byte
		for bitPosition, bitValue := range byteString {
			if bitValue == '1' {
				byt = setBit(byt, uint(_bitsInByte-1-bitPosition))
			}
		}
		b = append(b, byt)
	}

	return b
}

// BytesToStringBits returns the string bit representation of byte slice.
func BytesToStringBits(b []byte, spaceSeparator bool) string {
	str := strings.TrimSuffix(strings.TrimPrefix(fmt.Sprintf("%08b", b), "["), "]")

	if !spaceSeparator {
		return strings.ReplaceAll(str, " ", "")
	}

	return str
}

// Sets the bit at pos in the integer n.
func setBit(n byte, pos uint) byte {
	n |= 1 << pos
	return n
}
