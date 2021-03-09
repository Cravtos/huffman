package node

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

	return &head
}

// insert puts a node in a list so that the list remains sorted
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

// Todo: function to join two nodes (Node.join). should return one node with left right != nil.
// Todo: function that makes a tree from sorted list (NewTree).
