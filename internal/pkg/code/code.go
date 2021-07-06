package code

// Table maps byte to its code.
type Table map[byte]Code

// Code represents code in boolean vector.
type Code struct {
	Code uint64 // Vector containing code
	Len  uint8  // Length of code
}
