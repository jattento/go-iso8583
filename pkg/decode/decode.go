package decode

import (
	"errors"
	"github.com/iso-lib/pkg/bitmap"
)

const (
	_mtiKey    = "MTI"
	_mtiLength = 4
)

func Decode(b []byte) (map[string][]byte, error) {
	const minLength = 12

	content := make(map[string][]byte)

	if len(b) < minLength {
		return nil, errors.New("message to short")
	}

	content[_mtiKey] = b[:_mtiLength]

	bmap,nextBitmap,err := bitmap.ISO8583FromBytes(b[4:12],1)
	if err != nil{
		return nil, err
	}

	if nextBitmap{
		secondBitmap,_,err := bitmap.ISO8583FromBytes(b[12:20],2)
		if err != nil{
			return nil, err
		}

		for k, v := range secondBitmap {
			bmap[k] = v
		}
	}

	return content, nil
}
