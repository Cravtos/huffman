package node

import "github.com/cravtos/huffman/internal/pkg/code"

// NewTree constructs encoding tree from byte frequencies.
// Returns root node.
func NewTree(freq []uint) *Node {
	var head Node // Fictitious head

	for i, v := range freq {
		if v == 0 {
			continue
		}
		node := &Node{
			Value:  byte(i),
			Weight: v,
		}
		head.insert(node)
	}

	for head.Next != nil && head.Next.Next != nil {
		l := head.pop()
		r := head.pop()

		node := head.join(l, r)
		head.insert(node)
	}

	return head.Next
}

// FillTable fills table with new encoding values.
// E.g.:
// table['a'] = {code: 0b10, len: 2}
// table['b'] = {code: 0b110, len: 3}
func (head *Node) FillTable(table map[byte]code.Code) {
	var c code.Code
	head.fillTable(table, c)
}

// fillTable recursively fills table with encoding values.
func (head *Node) fillTable(table map[byte]code.Code, c code.Code) {
	if head.Left == nil && head.Right == nil {
		table[head.Value] = c
		return
	}

	if head.Left != nil {
		lc := code.Code{
			Code: c.Code << 1,
			Len:  c.Len + 1,
		}
		head.Left.fillTable(table, lc)
	}

	if head.Right != nil {
		rc := code.Code{
			Code: (c.Code << 1) | 1,
			Len:  c.Len + 1,
		}
		head.Right.fillTable(table, rc)
	}
}

// insert puts a node in a list so that the list remains sorted.
func (head *Node) insert(node *Node) {
	after := head
	for after.Next != nil && node.Weight >= after.Next.Weight {
		after = after.Next
	}

	node.Prev = after
	node.Next = after.Next
	if after.Next != nil {
		after.Next.Prev = node
	}
	after.Next = node
}

// pop removes first node after head and returns it.
// If head is the only node, nil is returned.
func (head *Node) pop() *Node {
	node := head.Next

	if node != nil {
		head.Next = nil
		if node.Next != nil {
			node.Next.Prev = head
			head.Next = node.Next
		}

		node.Next = nil
		node.Prev = nil
	}

	return node
}

// join returns node with left and right leaves set to l and r.
// Returned node weight is sum of l and r weights.
func (head *Node) join(l, r *Node) *Node {
	var node Node
	if l != nil {
		node.Weight += l.Weight
		node.Left = l
	}

	if r != nil {
		node.Weight += r.Weight
		node.Right = r
	}

	return &node
}
