package huffman

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/cravtos/huffman/internal/pkg/helpers"
	"github.com/cravtos/huffman/internal/pkg/tree"
	"github.com/icza/bitio"
)

// Encode do huffman encoding of os.File to os.File.
func Encode(in *os.File, out *os.File) (err error) {
	r := bufio.NewReader(in)
	w := bitio.NewWriter(out)

	// Calculate byte frequencies
	freq := helpers.CalcFreq(r)

	// Construct encoding tree
	root := tree.NewEncodingTree(freq)

	// Write header information
	if err = root.WriteHeader(w, freq); err != nil {
		return err
	}

	// Make encoding table
	table := root.NewEncodingTable()

	// Start reading from beginning
	if _, err = in.Seek(0, 0); err != nil {
		return err
	}
	r.Reset(in)

	// Encode file
	v, err := r.ReadByte()
	for err == nil {
		if err = w.WriteBits(table[v].Code, table[v].Len); err != nil {
			return err
		}
		v, err = r.ReadByte()
	}

	// Close writer and flush everything to file
	return w.Close()
}

// Decode do huffman decoding of os.File to os.File.
func Decode(in *os.File, out *os.File) (err error) {
	r := bitio.NewReader(in)
	w := bitio.NewWriter(out)

	// Read header and construct encoding tree
	nEncoded, root, err := tree.DecodeHeader(r)
	if err != nil {
		return err
	}

	// Decoding file code by code
	var i uint32
	code, err := root.DecodeNext(r)
	for i = 0; i != nEncoded && err == nil; i++ {
		err = w.WriteByte(code)
		if err != nil {
			return err
		}

		code, err = root.DecodeNext(r)
	}

	// Check if number of decoded symbols equal number of symbols in header
	if i != nEncoded {
		errMsg := fmt.Sprintf( "number of decoded symbols not equal to number of symbols from header: %d != %d\n", i, nEncoded)
		return errors.New(errMsg)
	}

	// Flush everything to file
	return w.Close()
}