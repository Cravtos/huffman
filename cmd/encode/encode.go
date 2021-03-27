package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cravtos/huffman/internal/pkg/bitio"
	"github.com/cravtos/huffman/internal/pkg/helpers"
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

	// Open file to write compressed data
	outFilePath := inFilePath + ".huff"
	log.Println("creating file", outFilePath)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't create file %s\n", outFilePath)
		return
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
		return
	}

	// Make encoding table
	log.Println("making encoding table")
	table := root.NewEncodingTable()

	// Start reading from begin
	if _, err = inFile.Seek(0, 0); err != nil {
		fmt.Fprintf(os.Stderr, "got error while seeking to begining of file: %v\n", err)
	}

	// Encode file
	log.Println("encoding file", inFilePath, "to file", outFilePath)
	v, err := r.ReadByte()
	for err == nil {
		if err = w.WriteBits(table[v].Code, table[v].Len); err != nil {
			fmt.Fprintf(os.Stderr, "got error while writing data: %v\n", err)
			return
		}
		v, err = r.ReadByte()
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
	ratio := float32(inStat.Size()) / float32(outStat.Size())
	log.Printf("input size: %v, output size: %v, ratio: %v\n", inSize, outSize, ratio)
}
