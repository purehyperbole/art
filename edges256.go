package art

type edges256 struct {
	edges    [256]*node
	children uint8
}

func newEdges256() *edges256 {
	return &edges256{}
}

func (e *edges256) ntype() uint8 {
	return Node256
}

func (e *edges256) next(b byte) *node {
	return e.edges[b]
}

func (e *edges256) setNext(b byte, next *node) {
	if e.edges[b] == nil {
		e.children++
	}
	e.edges[b] = next
}

func (e *edges256) copy() edges {
	ne := &edges256{
		children: e.children,
	}

	ne.edges = e.edges

	return ne
}

func (e *edges256) upgrade() edges {
	return e
}

func (e *edges256) full() bool {
	return false
}
