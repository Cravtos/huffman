package helpers

import (
	"io"
)

// CalcFreq reads everything from ByteReader and returns byte frequencies.
func CalcFreq(br io.ByteReader) map[uint8]uint {
	freq := make(map[uint8]uint)

	v, err := br.ReadByte()
	for err == nil {
		freq[v]++
		v, err = br.ReadByte()
	}

	return freq
}
