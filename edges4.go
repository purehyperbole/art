package art

import "bytes"

type edges4 struct {
	edges    [4]*node
	keys     [4]byte
	children uint8
}

func newEdges4() *edges4 {
	return &edges4{}
}

func newEdges4p() *edges {
	var e edges
	e = &edges4{}
	return &e
}

func (e *edges4) ntype() uint8 {
	return Node4
}

func (e *edges4) next(b byte) *node {
	i := bytes.IndexByte(e.keys[:], b)
	if i < 0 {
		return nil
	}

	return e.edges[i]
}

func (e *edges4) setNext(b byte, next *node) {
	i := e.search(b)

	if e.keys[i] == b {
		e.edges[i] = next
		return
	}

	copy(e.keys[i+1:], e.keys[i:])
	copy(e.edges[i+1:], e.edges[i:])

	e.keys[i] = b
	e.edges[i] = next
	e.children++
}

func (e *edges4) search(b byte) uint8 {
	for i := uint8(0); i < uint8(len(e.keys)); i++ {
		if e.keys[i] >= b {
			return i
		}
	}
	return e.children
}

func (e *edges4) copy() edges {
	ne := &edges4{
		children: e.children,
	}

	ne.keys = e.keys
	ne.edges = e.edges

	return ne
}

func (e *edges4) upgrade() edges {
	newEdges := newEdges16()

	for i := 0; i < 4; i++ {
		newEdges.setNext(e.keys[i], e.edges[i])
	}

	return newEdges
}

func (e *edges4) full() bool {
	return e.children == 4
}
