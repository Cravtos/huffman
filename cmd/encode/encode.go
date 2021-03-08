package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cravtos/huffman/internal/pkg/helpers"
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

	for i, v := range freq {
		if v != 0 {
			fmt.Printf("byte %d has freq %d\n", i, v)
		}
	}
}