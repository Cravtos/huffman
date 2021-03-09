package helpers

import (
	"io"
)

// CalcFreq reads everything from ByteReader and returns byte frequencies.
func CalcFreq(br io.ByteReader) []int {
	freq := make([]int, 256)

	v, err := br.ReadByte()
	for err == nil {
		freq[v]++
		v, err = br.ReadByte()
	}

	return freq
}
