package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cravtos/huffman/internal/pkg/helpers"
	"github.com/cravtos/huffman/internal/pkg/node"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [file]", filepath.Base(os.Args[0]))
		return
	}

	filePath := filepath.Clean(os.Args[1])
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file %s", filePath)
		return
	}
	defer file.Close()

	r := bufio.NewReader(file)
	freq := helpers.CalcFreq(r)

	head := node.NewTree(freq)

	// Little test
	var freqSum int
	for _, v := range freq {
		freqSum += v
	}

	root := head.Next
	if root.Weight != freqSum {
		fmt.Fprintf(os.Stderr, "root weight is incorrect!\nexp: %d\ngot: %d\n", root.Weight, freqSum)
		return
	}
	fmt.Printf("all ok :)\n")
}
