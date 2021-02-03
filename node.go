package art

import (
	"fmt"
	"strings"
	"sync/atomic"
	"unsafe"
)

const (
	Node4 = iota
	Node16
	Node48
	Node256
)

type edges interface {
	ntype() uint8
	next(b byte) *node
	setNext(b byte, next *node)
	copy() edges
	upgrade() edges
	full() bool
}

type node struct {
	prefix []byte
	edges  *unsafe.Pointer
	value  Comparable
}

func newNode(size int, prefix []byte, value Comparable) *node {
	var e unsafe.Pointer
	var ne edges

	switch size {
	case Node4:
		ne = newEdges4()
	case Node16:
		ne = newEdges16()
	case Node48:
		ne = newEdges48()
	case Node256:
		ne = newEdges256()
	}

	e = unsafe.Pointer(&ne)

	return &node{
		prefix: prefix,
		value:  value,
		edges:  &e,
	}
}

func (n *node) next(b byte) *node {
	return (*(*edges)(atomic.LoadPointer(n.edges))).next(b)
}

func (n *node) swapNext(b byte, existing, next *node) bool {
	e := (*edges)(atomic.LoadPointer(n.edges))

	cn := (*e).next(b)

	if cn != existing {
		return false
	}

	var ne edges

	if (*e).full() && (*e).next(b) == nil {
		ne = (*e).upgrade()
	} else {
		ne = (*e).copy()
	}

	ne.setNext(b, next)

	return atomic.CompareAndSwapPointer(n.edges, unsafe.Pointer(e), unsafe.Pointer(&ne))
}

func (n *node) setNext(b byte, next *node) {
	e := (*edges)(atomic.LoadPointer(n.edges))

	var newEdges edges

	if (*e).full() && (*e).next(b) == nil {
		newEdges = (*e).upgrade()
	} else {
		newEdges = *e
	}

	newEdges.setNext(b, next)

	*e = newEdges
}

func (n *node) getEdges() edges {
	return (*(*edges)(atomic.LoadPointer(n.edges)))
}

func (n *node) print() {
	output := []string{"{"}

	output = append(output, fmt.Sprintf("	Prefix Length: %d", len(n.prefix)))
	output = append(output, fmt.Sprintf("	Prefix: %s", string(n.prefix)))
	output = append(output, fmt.Sprintf("	Value: %d", n.value))

	output = append(output, "	Edges: [")

	edges := n.getEdges()

	if edges != nil {
		for i := 0; i < 256; i++ {
			edge := edges.next(byte(i))
			if edge != nil {
				output = append(output, fmt.Sprintf("		%s: %s", string(byte(i)), edge.prefix))
			}
		}
	}

	output = append(output, "	]")
	output = append(output, "}")

	fmt.Println(strings.Join(output, "\n"))
}
