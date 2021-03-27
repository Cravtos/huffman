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
		fmt.Fprintf(os.Stderr, "usage: %s [file]", filepath.Base(os.Args[0]))
		return
	}

	// Open file to read data
	inFilePath := filepath.Clean(os.Args[1])
	log.Println("opening file", inFilePath)
	inFile, err := os.Open(inFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s", inFilePath)
		return
	}
	defer inFile.Close()

	// Open file to write compressed data
	outFilePath := inFilePath + ".huff"
	log.Println("creating file", outFilePath)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't create file %s", outFilePath)
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
		fmt.Fprintf(os.Stderr, "got error while writing header: %v", err)
		return
	}

	// Make encoding table
	log.Println("making encoding table")
	table := root.NewEncodingTable()

	// Start reading from begin
	if _, err = inFile.Seek(0, 0); err != nil {
		fmt.Fprintf(os.Stderr, "got error while seeking to begining of file: %v", err)
	}

	// Encode file
	log.Println("encoding file", inFilePath, "to file", outFilePath)
	v, err := r.ReadByte()
	for err == nil {
		if err = w.WriteBits(table[v].Code, table[v].Len); err != nil {
			fmt.Fprintf(os.Stderr, "got error while writing data: %v", err)
			return
		}
		v, err = r.ReadByte()
	}

	// Flush everything to file
	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "got error while flushing: %v", err)
		return
	}

	log.Println("finished. see", outFilePath)
}
