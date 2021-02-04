package art

type edges interface {
	ntype() uint8
	next(b byte) *node
	setNext(b byte, next *node)
	copy() edges
	upgrade() edges
	full() bool
}

var leaf = (edges)(&edgesLeaf{})

type edgesLeaf struct{}

func (e *edgesLeaf) ntype() uint8 {
	return NodeLeaf
}

func (e *edgesLeaf) next(b byte) *node {
	return nil
}

func (e *edgesLeaf) setNext(b byte, next *node) {
}

func (e *edgesLeaf) copy() edges {
	return e
}

func (e *edgesLeaf) upgrade() edges {
	return newEdges4()
}

func (e *edgesLeaf) full() bool {
	return true
}
