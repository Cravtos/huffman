package helpers

import (
	"bufio"
)

func CalcFreq(r *bufio.Reader) []int {
	freq := make([]int, 256)

	v, err := r.ReadByte()
	for err == nil {
		freq[v] += 1
		v, err = r.ReadByte()
	}

	return freq
}