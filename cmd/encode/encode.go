package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cravtos/huffman/internal/pkg/helpers"
	"github.com/cravtos/huffman/internal/pkg/huffman"
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

	err = huffman.Encode(inFile, outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "got error while encoding: %v\n", err)
		os.Exit(1)
	}

	if *printRatio == true {
		if err := helpers.PrintRatio(inFile, outFile); err != nil {
			fmt.Fprintf(os.Stderr, "got error while getting compression ratio: %v\n", err)
			os.Exit(2)
		}
	}
}
