package main

import (
	"bufio"
	"fmt"
	"github.com/cravtos/huffman/internal/pkg/bitio"
	"os"
	"path/filepath"

	"github.com/cravtos/huffman/internal/pkg/code"
	"github.com/cravtos/huffman/internal/pkg/helpers"
	"github.com/cravtos/huffman/internal/pkg/node"
)

func main() {
	// Check if file is specified as argument
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [file]", filepath.Base(os.Args[0]))
		return
	}

	// Open file to read data
	inFilePath := filepath.Clean(os.Args[1])
	inFile, err := os.Open(inFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s", inFilePath)
		return
	}
	defer inFile.Close()

	// Open file to write compressed data
	outFilePath := inFilePath + ".huff"
	outFile, err := os.Create(outFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't create file %s", outFilePath)
		return
	}
	defer outFile.Close()

	// Calculate byte frequencies
	r := bufio.NewReader(inFile)
	freq := helpers.CalcFreq(r)

	// Construct encoding tree
	root := node.NewTree(freq)

	// Make encoding table
	table := make(map[byte]code.Code)
	root.FillTable(table)

	w := bitio.NewWriter(outFile)
	for _, c := range table {
		fmt.Printf("writing code %b | len: %d\n", c.Code, c.Len)

		err = w.WriteBits(c.Code, c.Len)
		if err != nil {
			fmt.Fprintf(os.Stderr, "got error while writing: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "got error while flushing: %v", err)
	}
}
