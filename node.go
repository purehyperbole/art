package art

import (
	"fmt"
	"sort"
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

type edges struct {
	ntype    uint8
	keys     []byte
	edges    []*node
	children uint8
}

type node struct {
	prefix []byte
	edges  *unsafe.Pointer
	value  interface{}
}

func newNode(size int, prefix []byte, value interface{}) *node {
	var e unsafe.Pointer

	switch size {
	case Node4:
		e = unsafe.Pointer(newEdges4())
	case Node16:
		e = unsafe.Pointer(newEdges16())
	case Node48:
		e = unsafe.Pointer(newEdges48())
	case Node256:
		e = unsafe.Pointer(newEdges256())
	}

	return &node{
		prefix: prefix,
		value:  value,
		edges:  &e,
	}
}

func (n *node) next(b byte) *node {
	return (*edges)(atomic.LoadPointer(n.edges)).next(b)
}

func newEdges4() *edges {
	return &edges{
		ntype: Node4,
		keys:  make([]byte, 4, 4),
		edges: make([]*node, 4, 4),
	}
}

func newEdges16() *edges {
	return &edges{
		ntype: Node16,
		keys:  make([]byte, 16),
		edges: make([]*node, 16),
	}
}

func newEdges48() *edges {
	return &edges{
		ntype: Node48,
		keys:  make([]byte, 256),
		edges: make([]*node, 48),
	}
}

func newEdges256() *edges {
	return &edges{
		ntype: Node256,
		edges: make([]*node, 256),
	}
}

func (n *node) swapNext(b byte, existing, next *node) bool {
	e := (*edges)(atomic.LoadPointer(n.edges))

	cn := e.next(b)

	if cn != existing {
		return false
	}

	var ne *edges

	if e.full() && e.next(b) == nil {
		ne = e.upgrade()
	} else {
		ne = e.copy()
	}

	switch ne.ntype {
	case Node4:
		ne.setNext4(b, next)
	case Node16:
		ne.setNext16(b, next)
	case Node48:
		ne.setNext48(b, next)
	case Node256:
		ne.setNext256(b, next)
	}

	return atomic.CompareAndSwapPointer(n.edges, unsafe.Pointer(e), unsafe.Pointer(ne))
}

func (n *node) setNext(b byte, next *node) {
	e := (*edges)(atomic.LoadPointer(n.edges))

	var newEdges *edges

	if e.full() && e.next(b) == nil {
		newEdges = e.upgrade()
	} else {
		newEdges = e
	}

	switch newEdges.ntype {
	case Node4:
		newEdges.setNext4(b, next)
	case Node16:
		newEdges.setNext16(b, next)
	case Node48:
		newEdges.setNext48(b, next)
	case Node256:
		newEdges.setNext256(b, next)
	}

	*e = *newEdges
}

func (e *edges) next(b byte) *node {
	switch e.ntype {
	case Node4:
		return e.next4(b)
	case Node16:
		return e.next16(b)
	case Node48:
		return e.next48(b)
	}

	return e.edges[b]
}

func (e *edges) next4(b byte) *node {
	for i := 0; i < 4; i++ {
		if e.keys[i] == b {
			return e.edges[i]
		}
	}

	return nil
}

func (e *edges) next16(b byte) *node {
	i := sort.Search(int(e.children), func(i int) bool {
		return e.keys[i] >= b
	})

	if i == 16 {
		return nil
	}

	if e.keys[i] != b {
		return nil
	}

	return e.edges[i]
}

func (e *edges) next48(b byte) *node {
	i := uint8(e.keys[b])

	if i == 0 {
		return nil
	}

	return e.edges[i-1]
}

func (e *edges) setNext4(b byte, next *node) {
	for i := 0; i < 4; i++ {
		if e.keys[i] == b {
			e.edges[i] = next
			return
		}
	}

	i := e.children
	e.keys[i] = b
	e.edges[i] = next
	e.children++
}

func (e *edges) setNext16(b byte, next *node) {
	p := e.search(b)

	if e.keys[p] == b {
		e.edges[p] = next
		return
	}

	for i := uint8(15); i > p; i-- {
		e.keys[i] = e.keys[i-1]
		e.edges[i] = e.edges[i-1]
	}

	e.keys[p] = b
	e.edges[p] = next
	e.children++
}

func (e *edges) setNext48(b byte, next *node) {
	if e.keys[b] != 0 {
		e.edges[e.keys[b]-1] = next
		return
	}

	e.keys[b] = e.children + 1
	e.edges[e.children] = next
	e.children++
}

func (e *edges) setNext256(b byte, next *node) {
	if e.edges[b] == nil {
		e.children++
	}
	e.edges[b] = next
}

func (e *edges) upgrade() *edges {
	var newEdges *edges

	switch e.ntype {
	case Node4:
		newEdges = e.upgrade4()
	case Node16:
		newEdges = e.upgrade16()
	case Node48:
		newEdges = e.upgrade48()
	}

	newEdges.children = e.children

	return newEdges
}

func (e *edges) upgrade4() *edges {
	newEdges := newEdges16()

	for i := 0; i < 4; i++ {
		newEdges.setNext16(e.keys[i], e.edges[i])
	}

	return newEdges
}

func (e *edges) upgrade16() *edges {
	newEdges := newEdges48()

	for i := uint8(0); i < e.children; i++ {
		newEdges.keys[e.keys[i]] = i + 1
		newEdges.edges[i] = e.edges[i]
	}

	return newEdges
}

func (e *edges) upgrade48() *edges {
	newEdges := newEdges256()

	for i := 0; i < 256; i++ {
		if e.keys[i] > 0 {
			newEdges.edges[i] = e.edges[e.keys[i]-1]
		}
	}

	return newEdges
}

func (e *edges) search(b byte) uint8 {
	for i := uint8(0); i < uint8(len(e.keys)); i++ {
		if e.keys[i] >= b {
			return i
		}
	}
	return e.children
}

func (e *edges) full() bool {
	switch e.ntype {
	case Node4:
		return e.children == 4
	case Node16:
		return e.children == 16
	case Node48:
		return e.children == 48
	}

	return false
}

func (e *edges) copy() *edges {
	ne := &edges{
		ntype:    e.ntype,
		keys:     make([]byte, len(e.keys)),
		edges:    make([]*node, len(e.edges)),
		children: e.children,
	}

	copy(ne.keys, e.keys)
	copy(ne.edges, e.edges)

	return ne
}

func (n *node) getEdges() *edges {
	return (*edges)(atomic.LoadPointer(n.edges))
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
