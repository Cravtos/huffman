package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/cravtos/huffman/internal/pkg/bitio"
	"github.com/cravtos/huffman/internal/pkg/helpers"
	"github.com/cravtos/huffman/internal/pkg/tree"
	"log"
	"os"
)

func main() {
	inPath := flag.String("input", "", "File to encode.")
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

	// Open file to write compressed data
	log.Println("creating file", *outPath)
	outFile, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't create file %s\n", *outPath)
		os.Exit(1)
	}
	defer outFile.Close()

	// Calculate byte frequencies
	log.Println("calculating frequencies")
	r := bufio.NewReader(inFile)
	freq := helpers.CalcFreq(r)

	// Construct encoding tree
	log.Println("constructing encoding tree")
	root := tree.NewEncodingTree(freq)

	// Write header information:
	// 4 bytes - number of encoded symbols
	// 1 byte - number of leaf in encoding tree
	// else - encoding tree in post order traversal
	log.Println("writing header")
	w := bitio.NewWriter(outFile)
	if err = root.WriteHeader(w, freq); err != nil {
		fmt.Fprintf(os.Stderr, "got error while writing header: %v\n", err)
		os.Exit(1)
	}

	// Make encoding table
	log.Println("making encoding table")
	table := root.NewEncodingTable()

	// Start reading from begin
	if _, err = inFile.Seek(0, 0); err != nil {
		fmt.Fprintf(os.Stderr, "got error while seeking to begining of file: %v\n", err)
		os.Exit(1)
	}

	// Encode file
	log.Println("encoding file", *inPath, "to file", *outPath)
	v, err := r.ReadByte()
	for err == nil {
		if err = w.WriteBits(table[v].Code, table[v].Len); err != nil {
			fmt.Fprintf(os.Stderr, "got error while writing data: %v\n", err)
			os.Exit(1)
		}
		v, err = r.ReadByte()
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
	ratio := float32(inStat.Size()) / float32(outStat.Size())
	log.Printf("input size: %v, output size: %v, ratio: %v\n", inSize, outSize, ratio)
}
