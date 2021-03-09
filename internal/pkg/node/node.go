package node

// NewTree constructs encoding tree from byte frequencies.
// Returns fictitious tree head.
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

	// for head.Next != nil {
	// 	l := head.pop()
	// 	r := head.pop()
	//
	// }

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

// Pop removes first node after head and returns it.
// If head is the only node, nil is returned.
func (head *Node) Pop() *Node {
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

// Todo: function to join two nodes (Node.join). should return one node with left right != nil.
// Todo: function that makes a tree from sorted list (NewTree).
