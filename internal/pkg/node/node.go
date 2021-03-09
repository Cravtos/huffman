package node

// NewTree constructs encoding tree from byte frequencies.
// Returns fictitious head node.
func NewTree(freq []int) *Node {
	var head Node

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

	return &head
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
