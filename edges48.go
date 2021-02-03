package art

type edges48 struct {
	keys     [256]byte
	edges    [48]*node
	children uint8
}

func newEdges48() *edges48 {
	return &edges48{}
}

func (e *edges48) ntype() uint8 {
	return Node48
}

func (e *edges48) next(b byte) *node {
	i := uint8(e.keys[b])

	if i == 0 {
		return nil
	}

	return e.edges[i-1]
}

func (e *edges48) setNext(b byte, next *node) {
	if e.keys[b] != 0 {
		e.edges[e.keys[b]-1] = next
		return
	}

	e.keys[b] = e.children + 1
	e.edges[e.children] = next
	e.children++
}

func (e *edges48) copy() edges {
	ne := &edges48{
		children: e.children,
	}

	ne.keys = e.keys
	ne.edges = e.edges

	return ne
}

func (e *edges48) upgrade() edges {
	newEdges := newEdges256()

	for i := 0; i < 256; i++ {
		if e.keys[i] > 0 {
			newEdges.edges[i] = e.edges[e.keys[i]-1]
		}
	}

	return newEdges
}

func (e *edges48) full() bool {
	return e.children == 48
}
