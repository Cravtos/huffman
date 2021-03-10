package tree

import (
	"github.com/cravtos/huffman/internal/pkg/code"
)

// Node represents node in encoding tree.
type Node struct {
	value       byte
	weight      uint
	left, right *Node
	next, prev  *Node
}

// NewEncodingTree constructs encoding tree from byte frequencies.
// Returns root node.
func NewEncodingTree(freq []uint) *Node {
	var head Node // Fictitious head

	for i, v := range freq {
		if v == 0 {
			continue
		}
		node := &Node{
			value:  byte(i),
			weight: v,
		}
		head.insert(node)
	}

	for head.next != nil && head.next.next != nil {
		l := head.pop()
		r := head.pop()

		node := head.join(l, r)
		head.insert(node)
	}

	// Fictitious head point to tree root
	return head.next
}

// NewEncodingTable fills table with new encoding values.
// E.g.:
// table['a'] = {code: 0b10, len: 2}
// table['b'] = {code: 0b110, len: 3}
func (head *Node) NewEncodingTable() code.Table {
	table := make(code.Table)
	head.fillTable(table, code.Code{})
	return table
}

// fillTable recursively fills table with encoding values.
func (head *Node) fillTable(table code.Table, c code.Code) {
	if head.left == nil && head.right == nil {
		table[head.value] = c
		return
	}

	if head.left != nil {
		lc := code.Code{
			Code: c.Code << 1,
			Len:  c.Len + 1,
		}
		head.left.fillTable(table, lc)
	}

	if head.right != nil {
		rc := code.Code{
			Code: (c.Code << 1) | 1,
			Len:  c.Len + 1,
		}
		head.right.fillTable(table, rc)
	}
}

// WriteHeader writes header which can be used to construct encoding tree.
//
// To store the tree, we use a post-order traversal, writing each node visited.
// When you encounter a leaf node, you write a 1 followed by the byte value of the leaf node.
// When you encounter a non-leaf node, you write a 0.
// To indicate the end of the Huffman coding tree, we write another 0.
//
// For the string "streets are stone stars are not",
// the header information is "1t1a1r001n1o01 01e1s0000", followed by the encoded text.
//func (head *Node) WriteHeader(w *bitio.Writer) error {
//	
//}

// insert puts a node in a list so that the list remains sorted.
func (head *Node) insert(node *Node) {
	after := head
	for after.next != nil && node.weight >= after.next.weight {
		after = after.next
	}

	node.prev = after
	node.next = after.next
	if after.next != nil {
		after.next.prev = node
	}
	after.next = node
}

// pop removes first node after head and returns it.
// If head is the only node, nil is returned.
func (head *Node) pop() *Node {
	node := head.next

	if node != nil {
		head.next = nil
		if node.next != nil {
			node.next.prev = head
			head.next = node.next
		}

		node.next = nil
		node.prev = nil
	}

	return node
}

// join returns node with left and right leaves set to l and r.
// Returned node weight is sum of l and r weights.
func (head *Node) join(l, r *Node) *Node {
	var node Node
	if l != nil {
		node.weight += l.weight
		node.left = l
	}

	if r != nil {
		node.weight += r.weight
		node.right = r
	}

	return &node
}
