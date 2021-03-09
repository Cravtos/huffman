package node

// Node represents node in encoding tree.
type Node struct {
	Value       byte
	Weight      uint
	Left, Right *Node
	Next, Prev  *Node
}
