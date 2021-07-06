package helpers

import (
	"fmt"
	"io"
	"os"
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

// PrintRatio prints compression ratio for two files.
func PrintRatio(f *os.File, s *os.File) error {
	inStat, err := f.Stat()
	if err != nil {
		return err
	}

	outStat, err := s.Stat()
	if err != nil {
		return err
	}

	inSize := inStat.Size()
	outSize := outStat.Size()
	ratio := float32(inStat.Size()) / float32(outStat.Size())

	fmt.Printf("input size: %v\noutput size: %v\nratio: %v bytes\n", inSize, outSize, ratio)
	return nil
}
