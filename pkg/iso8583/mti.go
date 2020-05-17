package iso8583

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jattento/go-iso8583/pkg/mti"
)

// MTI shares behaviour with VAR type.
type MTI struct {
	mti.MTI
}

// MarshalISO8583 wrappes VAR behavior.
func (mtiV MTI) MarshalISO8583(length int, enc string) ([]byte, error) {
	return VAR.MarshalISO8583(VAR(mtiV.String()), length, enc)
}

// UnmarshalISO8583 wrappes VAR behavior and validates that the output is a number to convert it to mti.MTI
func (mtiV *MTI) UnmarshalISO8583(b []byte, length int, enc string) (int, error) {
	if b == nil {
		return 0, errors.New("bytes input is nil")
	}

	var v VAR
	n, err := v.UnmarshalISO8583(b, length, enc)
	if err != nil {
		return n, err
	}

	_, err = strconv.Atoi(string(v))
	if err != nil {
		return n, fmt.Errorf("mti characters arent numbers: %w", err)
	}

	if len(v) != 4 {
		return n, fmt.Errorf("mti isnt 4 characters long, its: %v", len(v))
	}

	*mtiV = MTI{MTI: mti.MTI(v)}

	return n, nil
}
