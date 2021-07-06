package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/cravtos/huffman/internal/pkg/helpers"
	"github.com/cravtos/huffman/internal/pkg/tree"
	"github.com/icza/bitio"
)

func main() {
	inPath := flag.String("input", "", "File to encode.")
	outPath := flag.String("output", "", "Output file.")
	printRatio := flag.Bool("pr", false, "Print compression ratio.")

	flag.Parse()

	// Check if file is specified as argument
	if *inPath == "" || *outPath == "" {
		fmt.Fprintln(os.Stderr, "specify both input and output files path!")
		flag.Usage()
		os.Exit(1)
	}

	// Open file to read data
	inFile, err := os.Open(*inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s\n", *inPath)
		os.Exit(1)
	}
	defer inFile.Close()

	// Open file to write compressed data
	outFile, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't create file %s\n", *outPath)
		os.Exit(1)
	}
	defer outFile.Close()

	// Calculate byte frequencies
	r := bufio.NewReader(inFile)
	freq := helpers.CalcFreq(r)

	// Construct encoding tree
	root := tree.NewEncodingTree(freq)

	// Write header information:
	// 4 bytes - number of encoded symbols
	// 1 byte - number of leaf in encoding tree
	// else - encoding tree in post order traversal
	w := bitio.NewWriter(outFile)
	if err = root.WriteHeader(w, freq); err != nil {
		fmt.Fprintf(os.Stderr, "got error while writing header: %v\n", err)
		os.Exit(1)
	}

	// Make encoding table
	table := root.NewEncodingTable()

	// Start reading from begin
	if _, err = inFile.Seek(0, 0); err != nil {
		fmt.Fprintf(os.Stderr, "got error while seeking to begining of file: %v\n", err)
		os.Exit(1)
	}

	// Encode file
	v, err := r.ReadByte()
	for err == nil {
		if err = w.WriteBits(table[v].Code, table[v].Len); err != nil {
			fmt.Fprintf(os.Stderr, "got error while writing data: %v\n", err)
			os.Exit(1)
		}
		v, err = r.ReadByte()
	}

	// Flush everything to file
	if _, err := w.Align(); err != nil {
		fmt.Fprintf(os.Stderr, "got error while flushing: %v\n", err)
		os.Exit(1)
	}

	if *printRatio == true {
		if err := helpers.PrintRatio(inFile, outFile); err != nil {
			fmt.Fprintf(os.Stderr, "got error while getting compression ratio: %v\n", err)
			os.Exit(2)
		}
	}
}
