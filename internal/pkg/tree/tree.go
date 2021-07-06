package tree

import (
	"errors"

	"github.com/cravtos/huffman/internal/pkg/code"
	"github.com/icza/bitio"
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
func NewEncodingTree(freq map[uint8]uint) *Node {
	var head Node // Fictitious head

	for i, v := range freq {
		node := &Node{
			value:  i,
			weight: v,
		}
		head.insert(node)
	}

	for head.next != nil && head.next.next != nil {
		l := head.popFirst()
		r := head.popFirst()

		node := join(l, r)
		head.insert(node)
	}

	// Fictitious head point to tree root
	if head.next != nil {
		head.next.prev = nil
	}
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
	if head == nil {
		return
	}

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
// Header: uint32 (number of encoded symbols in file)
//		   uint8 (number of symbols in tree)
//		   tree in raw bits
//
// To store the tree, we use a post-order traversal, writing each node visited.
// When you encounter a leaf node, you write a 1 followed by the byte value of the leaf node.
// When you encounter a non-leaf node, you write a 0.
// To indicate the end of the Huffman coding tree, we write another 0.
//
// For the string "streets are stone stars are not",
// the header information is "1t1a1r001n1o01 01e1s0000", followed by the encoded text.
// (https://engineering.purdue.edu/ece264/17au/hw/HW13/resources//streetstar.jpg)
func (head *Node) WriteHeader(w *bitio.Writer, freq map[uint8]uint) (err error) {
	var nEncoded uint32
	for _, v := range freq {
		nEncoded += uint32(v)
	}

	// Write total number of encoded symbols
	w.TryWriteBitsUnsafe(uint64(nEncoded), 32)

	// Write total number of symbols in graph
	w.TryWriteBitsUnsafe(uint64(len(freq)), 8)

	// Write encoding tree information
	if err = head.writeHeader(w); err != nil {
		return err
	}
	w.TryWriteBitsUnsafe(0, 1)
	return w.TryError
}

func (head *Node) writeHeader(w *bitio.Writer) (err error) {
	if head == nil {
		return
	}

	if head.left == nil && head.right == nil {
		w.TryWriteBitsUnsafe(1, 1)
		w.TryWriteByte(head.value)
		return w.TryError
	}

	if head.left != nil {
		if err = head.left.writeHeader(w); err != nil {
			return err
		}
	}

	if head.right != nil {
		if err = head.right.writeHeader(w); err != nil {
			return err
		}
	}

	w.TryWriteBitsUnsafe(0, 1)
	return w.TryError
}

// DecodeHeader reads from bitio.Reader total number of encoded symbols,
// number of leaf in tree, the tree itself.
// Returns constructed tree and number of encoded symbols.
func DecodeHeader(r *bitio.Reader) (nEncoded uint32, root *Node, err error) {
	var buf uint64
	buf, err = r.ReadBits(32)
	nEncoded = uint32(buf)
	if err != nil {
		return 0, nil, err
	}

	buf, err = r.ReadBits(8)
	nTree := byte(buf)
	if err != nil {
		return 0, nil, err
	}

	root, err = decodeTree(r, nTree)
	if err != nil {
		return 0, nil, err
	}

	return nEncoded, root, nil
}

// decodeTree constructs tree from header information
func decodeTree(r *bitio.Reader, nTree byte) (root *Node, err error) {
	var head Node
	var nodes byte
	var leaves byte
	var u uint64

	for nodes < nTree {
		u, err = r.ReadBits(1)
		if err != nil {
			return nil, err
		}

		if u == 1 {
			leaves++
			symbol, err := r.ReadBits(8)
			if err != nil {
				return nil, err
			}
			node := &Node{value: byte(symbol)}
			head.pushBack(node)
		}

		if u == 0 {
			nodes++
			if nodes == nTree {
				break
			}
			r := head.popLast()
			l := head.popLast()
			node := join(l, r)
			head.pushBack(node)
		}
	}

	if nodes != leaves {
		err = errors.New("nodes != leaves")
	}

	return head.next, err
}

// insert puts a node to list so that the list remains sorted.
func (head *Node) insert(node *Node) {
	if head == nil {
		return
	}

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

// pushBack puts a node to the end of a list.
func (head *Node) pushBack(node *Node) {
	if head == nil {
		return
	}

	after := head
	for after.next != nil {
		after = after.next
	}

	node.prev = after
	node.next = nil
	after.next = node
}

// popFirst removes first node after head and returns it.
// If head is the only node, nil is returned.
func (head *Node) popFirst() *Node {
	if head == nil {
		return nil
	}

	node := head.next
	if node == nil {
		return nil
	}

	head.next = nil
	if node.next != nil {
		node.next.prev = head
		head.next = node.next
	}
	node.next = nil
	node.prev = nil

	return node
}

// popLast removes first node after head and returns it.
// If head is the only node, nil is returned.
func (head *Node) popLast() *Node {
	if head == nil {
		return nil
	}

	node := head.next
	if node == nil {
		return nil
	}

	for node.next != nil {
		node = node.next
	}

	node.prev.next = nil
	node.prev = nil

	return node
}

// join returns node with left and right leaves set to l and r.
// Returned node weight is sum of l and r weights.
func join(l, r *Node) *Node {
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

// DecodeNext reads bits from Reader until reaching a leaf in encoding tree.
// Returns code corresponding to leaf.
func (head *Node) DecodeNext(r *bitio.Reader) (b byte, err error) {
	if head == nil {
		return
	}

	var u uint64
	node := head

	for node.left != nil && node.right != nil { // or node.value == 0
		u, err = r.ReadBits(1)
		if err != nil {
			return 0, err
		}

		if u == 0 {
			node = node.left
		} else {
			node = node.right
		}
	}

	return node.value, err
}
