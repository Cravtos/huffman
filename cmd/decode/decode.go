package main

import (
	"flag"
	"fmt"
	"github.com/cravtos/huffman/internal/pkg/bitio"
	"github.com/cravtos/huffman/internal/pkg/tree"
	"log"
	"os"
)

func main() {
	inPath := flag.String("input", "", "File to decode.")
	outPath := flag.String("output", "", "Output file.")

	flag.Parse()

	// Check if file is specified as argument
	if *inPath == "" || *outPath == "" {
		fmt.Fprintf(os.Stderr, "specify both input and output files path!")
		flag.Usage()
		os.Exit(1)
	}

	// Open file to read data
	log.Println("opening file", *inPath)
	inFile, err := os.Open(*inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s\n", *inPath)
		os.Exit(1)
	}
	defer inFile.Close()

	// Open file to write decompressed data
	log.Println("creating file", *outPath)
	outFile, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't create file %s\n", *outPath)
		os.Exit(1)
	}
	defer outFile.Close()

	// Read header and construct encoding tree
	log.Println("reading header")
	r := bitio.NewReader(inFile)
	nEncoded, root, err := tree.DecodeHeader(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't decode header information: %v\n", err)
		os.Exit(1)
	}

	// Decoding file
	log.Println("decoding file")
	w := bitio.NewWriter(outFile)

	// Decoding code by code
	code, err := root.DecodeNext(r)
	var i uint32
	for i = 0; i != nEncoded && err == nil; i++ {
		if err = w.WriteByte(code); err != nil {
			fmt.Fprintf(os.Stderr, "got error while writing data: %v\n", err)
			os.Exit(1)
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
		os.Exit(1)
	}

	log.Println("finished. see", *outPath)

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
