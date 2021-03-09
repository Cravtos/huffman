package node

func NewTree(freq []int) *Node {
	var head Node

	prev := &head
	for i, v := range freq {
		if v == 0 {
			continue
		}
		cur := &Node{
			Value:  byte(i),
			Weight: v,
			Prev:   prev,
		}
		prev.Next = cur
		prev = cur
	}

	return &head
}
