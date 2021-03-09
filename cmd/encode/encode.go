package main

import (
	"bufio"
	"fmt"
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

	// Open file
	filePath := filepath.Clean(os.Args[1])
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s", filePath)
		return
	}
	defer file.Close()

	// Calculate byte frequencies
	r := bufio.NewReader(file)
	freq := helpers.CalcFreq(r)

	// Construct encoding tree
	root := node.NewTree(freq)

	// Make encoding table
	table := make(map[byte]code.Code)
	root.FillTable(table)

	for i, v := range table {
		fmt.Printf("byte %d:\t%b (len %d)\n", i, v.Code, v.Len)
	}
}
