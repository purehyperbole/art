package art

import (
	"sort"
)

const (
	Node4 = iota
	Node16
	Node48
	Node256
)

type node struct {
	ntype    uint8
	keys     []byte
	edges    []*node
	prefix   []byte
	children uint8
	value    interface{}
}

func newNode4() *node {
	return &node{
		ntype: Node4,
		keys:  make([]byte, 4, 4),
		edges: make([]*node, 4, 4),
	}
}

func newNode16() *node {
	return &node{
		ntype: Node16,
		keys:  make([]byte, 16),
		edges: make([]*node, 16),
	}
}

func newNode48() *node {
	return &node{
		ntype: Node48,
		keys:  make([]byte, 256),
		edges: make([]*node, 48),
	}
}

func newNode256() *node {
	return &node{
		ntype: Node256,
		edges: make([]*node, 256),
	}
}

func (n *node) next(b byte) *node {
	switch n.ntype {
	case Node4:
		return n.next4(b)
	case Node16:
		return n.next16(b)
	case Node48:
		return n.next48(b)
	}

	return n.edges[b]
}

func (n *node) next4(b byte) *node {
	for i := 0; i < 4; i++ {
		if n.keys[i] == b {
			return n.edges[i]
		}
	}

	return nil
}

func (n *node) next16(b byte) *node {
	i := sort.Search(int(n.children), func(i int) bool {
		return n.keys[i] >= b
	})

	if i == 16 {
		return nil
	}

	if n.keys[i] != b {
		return nil
	}

	return n.edges[i]
}

func (n *node) next48(b byte) *node {
	i := uint8(n.keys[b])

	if i == 0 {
		return nil
	}

	return n.edges[i-1]
}

func (n *node) setNext(b byte, next *node) {
	if n.full() && n.next(b) == nil {
		n.upgrade()
	}

	switch n.ntype {
	case Node4:
		n.setNext4(b, next)
	case Node16:
		n.setNext16(b, next)
	case Node48:
		n.setNext48(b, next)
	case Node256:
		n.setNext256(b, next)
	}
}

func (n *node) setNext4(b byte, next *node) {
	for i := 0; i < 4; i++ {
		if n.keys[i] == b {
			n.edges[i] = next
			return
		}
	}

	i := n.children
	n.keys[i] = b
	n.edges[i] = next
	n.children++
}

func (n *node) setNext16(b byte, next *node) {
	p := n.search(b)

	if n.keys[p] == b {
		n.edges[p] = next
		return
	}

	for i := uint8(15); i > p; i-- {
		n.keys[i] = n.keys[i-1]
		n.edges[i] = n.edges[i-1]
	}

	n.keys[p] = b
	n.edges[p] = next
	n.children++
}

func (n *node) setNext48(b byte, next *node) {
	if n.keys[b] != 0 {
		n.edges[n.keys[b]-1] = next
		return
	}

	n.keys[b] = n.children + 1
	n.edges[n.children] = next
	n.children++
}

func (n *node) setNext256(b byte, next *node) {
	if n.edges[b] == nil {
		n.children++
	}
	n.edges[b] = next
}

func (n *node) upgrade() {
	var newNode *node

	switch n.ntype {
	case Node4:
		newNode = n.upgrade4()
	case Node16:
		newNode = n.upgrade16()
	case Node48:
		newNode = n.upgrade48()
	}

	newNode.prefix = n.prefix
	newNode.children = n.children
	newNode.value = n.value

	*n = *newNode
}

func (n *node) upgrade4() *node {
	newNode := newNode16()

	for i := 0; i < 4; i++ {
		newNode.setNext(n.keys[i], n.edges[i])
	}

	return newNode
}

func (n *node) upgrade16() *node {
	newNode := newNode48()

	for i := uint8(0); i < n.children; i++ {
		newNode.keys[n.keys[i]] = byte(i + 1)
		newNode.edges[i] = n.edges[i]
	}

	return newNode
}

func (n *node) upgrade48() *node {
	newNode := newNode256()

	for i := 0; i < 256; i++ {
		if n.keys[i] > 0 {
			newNode.edges[i] = n.edges[n.keys[i]-1]
		}
	}

	return newNode
}

func (n *node) search(b byte) uint8 {
	for i := uint8(0); i < uint8(len(n.keys)); i++ {
		if n.keys[i] >= b {
			return i
		}
	}
	return n.children
}

func (n *node) full() bool {
	switch n.ntype {
	case Node4:
		return n.children == 4
	case Node16:
		return n.children == 16
	case Node48:
		return n.children == 48
	}

	return false
}
