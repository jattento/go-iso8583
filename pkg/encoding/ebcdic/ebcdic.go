package ebcdic

type version []entry

type entry struct {
	EBCDIC         byte
	ASCII          byte
	representation string
	description    string
}

const null = 0x0

func (v version) FromASCII(b []byte) []byte {
	output := make([]byte, 0)

	for _, byt := range b {
		func() {
			for _, e := range v {
				if e.ASCII == byt {
					output = append(output, e.EBCDIC)
					return
				}
			}
			output = append(output, null)
		}()
	}

	return output
}

func (v version) ToASCII(b []byte) []byte {
	output := make([]byte, 0)

	for _, byt := range b {
		func() {
			for _, e := range v {
				if e.EBCDIC == byt {
					output = append(output, e.ASCII)
					return
				}
			}
			output = append(output, null)
		}()
	}

	return output
}
