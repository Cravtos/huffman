package node

// Node represents node in encoding tree.
type Node struct {
	Value       byte
	Weight      int
	Left, Right *Node
	Next, Prev  *Node
}
