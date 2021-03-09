package node

type Node struct {
	Value       byte
	Weight      int
	Left, Right *Node
	Next, Prev  *Node
}
