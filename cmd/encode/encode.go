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

	cur := head.Next
	for cur != nil {
		fmt.Printf("(v: %d, w: %d)\n", cur.Value, cur.Weight)
		cur = cur.Next
	}
}
