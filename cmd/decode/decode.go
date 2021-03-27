package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cravtos/huffman/internal/pkg/bitio"
	"github.com/cravtos/huffman/internal/pkg/tree"
)

func main() {
	// Check if file is specified as argument
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [file]\n", filepath.Base(os.Args[0]))
		return
	}

	// Open file to read data
	inFilePath := filepath.Clean(os.Args[1])
	log.Println("opening file", inFilePath)
	inFile, err := os.Open(inFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s\n", inFilePath)
		return
	}
	defer inFile.Close()

	// Open file to write decompressed data
	outFilePath := inFilePath + ".decoded"
	log.Println("creating file", outFilePath)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't create file %s\n", outFilePath)
		return
	}
	defer outFile.Close()

	// Read header and construct encoding tree
	log.Println("reading header")
	r := bitio.NewReader(inFile)
	nEncoded, root, err := tree.DecodeHeader(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't decode header information: %v\n", err)
		return
	}

	// Decoding file
	log.Println("decoding file")
	w := bitio.NewWriter(outFile)

	// Decoding code by code
	code, err := root.DecodeNext(r)
	var i uint32
	for i = 0; i != nEncoded && err == nil; i++ {
		if err = w.WriteByte(code); err != nil {
			fmt.Fprintf(os.Stderr, "got error while writing data: %v", err)
			return
		}

		code, err = root.DecodeNext(r)
	}

	// Check if number of decoded symbols equal number of symbols in header
	if i != nEncoded {
		fmt.Fprintf(os.Stderr, "number of decoded symbols not equal to number of symbols from header: %d != %d\n", i, nEncoded)
	}

	// Flush everything to file
	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "got error while flushing: %v\n", err)
		return
	}

	log.Println("finished. see", outFilePath)

	// Get size statistics
	inStat, err := inFile.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't obtain stat for input file: %v\n", err)
		return
	}

	outStat, err := outFile.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't obtain stat for output file: %v\n", err)
		return
	}

	inSize := inStat.Size()
	outSize := outStat.Size()
	ratio := float32(outStat.Size()) / float32(inStat.Size())
	log.Printf("input size: %v, output size: %v, ratio: %v\n", inSize, outSize, ratio)
}
