package art

import (
	"bytes"
)

type edges16 struct {
	keys     [16]byte
	edges    [16]*node
	children uint8
}

func newEdges16() *edges16 {
	return &edges16{}
}

func (e *edges16) ntype() uint8 {
	return Node16
}

func (e *edges16) next(b byte) *node {
	i := bytes.IndexByte(e.keys[:], b)
	if i < 0 {
		return nil
	}

	return e.edges[i]
}

func (e *edges16) setNext(b byte, next *node) {
	p := e.search(b)

	if e.keys[p] == b {
		e.edges[p] = next
		return
	}

	copy(e.keys[p+1:], e.keys[p:])
	copy(e.edges[p+1:], e.edges[p:])

	e.keys[p] = b
	e.edges[p] = next
	e.children++
}

func (e *edges16) search(b byte) uint8 {
	for i := uint8(0); i < uint8(len(e.keys)); i++ {
		if e.keys[i] >= b {
			return i
		}
	}
	return e.children
}

func (e *edges16) copy() edges {
	ne := &edges16{
		children: e.children,
	}

	ne.keys = e.keys
	ne.edges = e.edges

	return ne
}

func (e *edges16) upgrade() edges {
	newEdges := newEdges48()

	for i := uint8(0); i < e.children; i++ {
		newEdges.keys[e.keys[i]] = i + 1
		newEdges.edges[i] = e.edges[i]
	}

	return newEdges
}

func (e *edges16) full() bool {
	return e.children == 16
}
