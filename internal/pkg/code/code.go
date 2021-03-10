package code

// Table maps byte to its code.
type Table map[byte]Code

// Code contains of byte and its new encoding.
type Code struct {
	Code uint64 // Vector containing code
	Len  uint8  // Length of code
}

