package test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cravtos/huffman/internal/pkg/helpers"
	"github.com/cravtos/huffman/internal/pkg/huffman"
)

func encodeAndDecode(t *testing.T, file string) {
	// Create temp directory for resulting files
	tmp := t.TempDir()

	origName := file
	origFile, err := os.Open(origName)
	if err != nil {
		t.Errorf("got error while opening test file %s: %v\n", origName, err)
		t.FailNow()
	}
	defer origFile.Close()

	// Create file for encoded data
	encName := tmp + ".encoded"
	encFile, err := os.Create(encName)
	if err != nil {
		t.Errorf("got error while opening encoded test file %s: %v\n", encName, err)
		t.FailNow()
	}
	defer encFile.Close()
	defer os.Remove(encName)

	err = huffman.Encode(origFile, encFile)
	if err != nil {
		t.Errorf("got error while encoding test file %s: %v\n", origName, err)
		t.FailNow()
	}

	// Seek to beginning (since it will be used again to be decoded)
	_, err = encFile.Seek(0, 0)
	if err != nil {
		t.Errorf("got error while seeking test file %s: %v\n", encName, err)
		t.FailNow()
	}

	// Create file for decoded data
	decName := tmp + ".decoded"
	decFile, err := os.Create(decName)
	if err != nil {
		t.Errorf("got error while opening decoded test file %s: %v\n", decName, err)
		t.FailNow()
	}
	defer decFile.Close()
	defer os.Remove(decName)

	err = huffman.Decode(encFile, decFile)
	if err != nil {
		t.Errorf("got error while decoding test file %s: %v\n", encName, err)
		t.FailNow()
	}

	// Seek files to beginning for comparing
	_, err = origFile.Seek(0, 0)
	if err != nil {
		t.Errorf("got error while seeking test file %s: %v\n", origName, err)
		t.FailNow()
	}

	_, err = decFile.Seek(0, 0)
	if err != nil {
		t.Errorf("got error while seeking test file %s: %v\n", decName, err)
		t.FailNow()
	}

	equal, err := helpers.CompareFiles(origFile, decFile)
	if err != nil {
		t.Errorf("got error while comparing test files %s and %s: %v\n", origName, decName, err)
		t.FailNow()
	}

	if equal == false {
		t.Errorf("original and decoded files are not equal: %s and %s", origName, decName)
	}
}

// TestHuffman encodes and decodes every file in test/testdata, and compares results to originals.
func TestHuffman(t *testing.T) {
	testFiles, err := ioutil.ReadDir("./testdata")
	if err != nil {
		t.Errorf("got error while getting testdata: %v\n", err)
		t.FailNow()
	}

	for _, file := range testFiles {
		file := file
		t.Run(file.Name(), func(t *testing.T) {
			t.Parallel()
			encodeAndDecode(t, "./testdata/"+file.Name())
		})
	}
}
