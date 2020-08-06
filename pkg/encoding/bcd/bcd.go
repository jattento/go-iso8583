// Package bcd provides functions for converting integers to BCD byte array and vice versa.
package bcd

func pow100(power byte) uint64 {
	res := uint64(1)
	for i := byte(0); i < power; i++ {
		res *= 100
	}
	return res
}

func FromUint(value uint64, size int) []byte {
	buf := make([]byte, size)
	if value > 0 {
		remainder := value
		for pos := size - 1; pos >= 0 && remainder > 0; pos-- {
			tail := byte(remainder % 100)
			hi, lo := tail/10, tail%10
			buf[pos] = byte(hi<<4 + lo)
			remainder = remainder / 100
		}
	}
	return buf
}

// Returns uint8 value in BCD format.
//
// If value > 99, function returns value for last two digits of source value
// (Example: uint8(123) = uint8(0x23)).
func FromUint8(value uint8) byte {
	return FromUint(uint64(value), 1)[0]
}

// Returns two-bytes array with uint16 value in BCD format
//
// If value > 9999, function returns value for last two digits of source value
// (Example: uint8(12345) = []byte{0x23, 0x45}).
func FromUint16(value uint16) []byte {
	return FromUint(uint64(value), 2)
}

// Returns four-bytes array with uint32 value in BCD format
//
// If value > 99999999, function returns value for last two digits of source value
// (Example: uint8(1234567890) = []byte{0x23, 0x45, 0x67, 0x89}).
func FromUint32(value uint32) []byte {
	return FromUint(uint64(value), 4)
}

// Returns eight-bytes array with uint64 value in BCD format
//
// If value > 9999999999999999, function returns value for last two digits of source value
// (Example: uint8(12233445566778899) = []byte{0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99}).
func FromUint64(value uint64) []byte {
	return FromUint(value, 8)
}

func toUint(value []byte, size int) uint64 {
	vlen := len(value)
	if vlen > size {
		value = value[vlen-size:]
	}
	res := uint64(0)
	for i, b := range value {
		hi, lo := b>>4, b&0x0f
		if hi > 9 || lo > 9 {
			return 0
		}
		res += uint64(hi*10+lo) * pow100(byte(vlen-i)-1)
	}
	return res
}

// Returns uint8 value converted from bcd byte.
//
// If byte is not BCD (e.g. 0x1A), function returns zero.
func ToUint8(value byte) uint8 {
	return uint8(toUint([]byte{value}, 1))
}

// Return uint16 value converted from at most last two bytes of bcd bytes array.
//
// If any byte of used array part is not BCD (e.g 0x1A), function returns zero.
func ToUint16(value []byte) uint16 {
	return uint16(toUint(value, 2))
}

// Return uint32 value converted from at most last four bytes of bcd bytes array.
//
// If any byte of used array part is not BCD (e.g 0x1A), function returns zero.
func ToUint32(value []byte) uint32 {
	return uint32(toUint(value, 4))
}

// Return uint64 value converted from at most last eight bytes of bcd bytes array.
//
// If any byte of used array part is not BCD (e.g 0x1A), function returns zero.
func ToUint64(value []byte) uint64 {
	return toUint(value, 8)
}
