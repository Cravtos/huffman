package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cravtos/huffman/internal/pkg/bitio"
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

	// Open file to write decompressed data
	outFilePath := inFilePath + ".decoded"
	log.Println("creating file", outFilePath)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't create file %s\n", outFilePath)
		return
	}
	defer outFile.Close()

	// Read header and construct encoding tree
	log.Println("reading header")
	r := bitio.NewReader(inFile)
	nEncoded, root, err := tree.DecodeHeader(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't decode header information: %v\n", err)
		return
	}

	// Decoding file
	log.Println("decoding file")
	w := bitio.NewWriter(outFile)

	//code, err := root.DecodeNext(r)
	//for i := 0; i != nEncoded && err == nil; i++ {
	//	if err = w.WriteByte(code); err != nil {
	//		fmt.Fprintf(os.Stderr, "got error while writing data: %v", err)
	//		return
	//	}
	//
	//	code, err = root.DecodeNext(r)
	//}

	// Flush everything to file
	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "got error while flushing: %v\n", err)
		return
	}

	_ = nEncoded
	_ = root

	log.Println("finished. see", outFilePath)
}
