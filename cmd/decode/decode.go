package main

import (
	"fmt"
	"github.com/cravtos/huffman/internal/pkg/bitio"
	"log"
	"os"
	"path/filepath"
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

	// Open file to write decompressed data
	outFilePath := inFilePath + ".decoded"
	log.Println("creating file", outFilePath)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't create file %s", outFilePath)
		return
	}
	defer outFile.Close()

	// Read header and construct encoding tree
	log.Println("reading header")
	r := bitio.NewReader(inFile)
	// nEncoded, root, err := tree.ReadHeader(r)

	// Decoding file
	log.Println("decoding file")
	w := bitio.NewWriter(outFile)


	u, err := r.ReadBits(64)
	for err != nil {
		panic(err)
	}
	fmt.Printf("number of encoded symbols: %d | in bits: %b\n", u, u)

	u, err = r.ReadBits(8)
	for err != nil {
		panic(err)
	}
	fmt.Printf("number of symbols in tree: %d | in bits: %b\n", byte(u), byte(u))

	u, err = r.ReadBits(13)
	for err != nil {
		panic(err)
	}
	fmt.Printf("some bits after: %b\n", u)

	u, err = r.ReadBits(37)
	for err != nil {
		panic(err)
	}
	fmt.Printf("some bits after: %b\n", u)

	// code, err := root.DecodeNext(r)
	//for i := 0; i != nEncoded && err == nil; i++ {
	//	if err = w.WriteByte(code); err != nil {
	//		fmt.Fprintf(os.Stderr, "got error while writing data: %v", err)
	//		return
	//	}
	//
	//	code, err = root.DecodeNext(r)
	//}

	//// Flush everything to file
	//if err := w.Flush(); err != nil {
	//	fmt.Fprintf(os.Stderr, "got error while flushing: %v", err)
	//	return
	//}

	_ = r
	_ = w
	log.Println("finished. see", outFilePath)
}
