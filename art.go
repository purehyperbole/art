package art

import (
	"reflect"
)

// ART an adaptive radix tree implementation
type ART struct {
	root *node
}

// New creates a new radix tree
func New() *ART {
	return &ART{
		root: newNode(Node256, nil, nil),
	}
}

// Insert value into the tree
func (t *ART) Insert(key []byte, value interface{}) bool {
	var success bool

	parent, current, pos, dv := t.find(key)

	for {
		switch {
		case shouldInsert(key, current, parent, pos, dv):
			success = t.insertNode(key, value, parent, current, pos, dv)
		case shouldUpdate(key, current, parent, pos, dv):
			success = t.updateNode(key, value, parent, current, pos, dv)
		case shouldSplitThreeWay(key, current, parent, pos, dv):
			success = t.splitThreeWay(key, value, parent, current, pos, dv)
		case shouldSplitTwoWay(key, current, parent, pos, dv):
			success = t.splitTwoWay(key, value, parent, current, pos, dv)
		}

		if success {
			return true
		}

		parent, current, pos, dv = t.find(key)

		if shouldUpdate(key, current, parent, pos, dv) {
			// someone else updated the same value we did
			return false
		}
	}
}

// Swap atomically swaps a value
func (t *ART) Swap(key []byte, old, new interface{}) bool {
	var success bool

	parent, current, pos, dv := t.find(key)

	// if we didnt find a node and the old value is not empty, fail
	if shouldUpdate(key, current, parent, pos, dv) && old == nil {
		return false
	}

	// if we did find a node, check that the value we have matches or fail
	if current != nil && old != nil {
		// this is probably going to be slow :/
		if !reflect.DeepEqual(current.value, old) {
			return false
		}
	}

	for {
		switch {
		case shouldInsert(key, current, parent, pos, dv):
			success = t.insertNode(key, new, parent, current, pos, dv)
		case shouldUpdate(key, current, parent, pos, dv):
			success = t.updateNode(key, new, parent, current, pos, dv)
		case shouldSplitThreeWay(key, current, parent, pos, dv):
			success = t.splitThreeWay(key, new, parent, current, pos, dv)
		case shouldSplitTwoWay(key, current, parent, pos, dv):
			success = t.splitTwoWay(key, new, parent, current, pos, dv)
		}

		if success {
			return true
		}

		parent, current, pos, dv = t.find(key)

		// if true, someone else updated the same value we did
		if shouldUpdate(key, current, parent, pos, dv) {
			return false
		}
	}
}

// Lookup a value from the tree
func (t *ART) Lookup(key []byte) interface{} {
	_, current, pos, _ := t.find(key)

	if current == nil || len(key) > pos {
		return nil
	}

	return current.value
}

func (t *ART) find(key []byte) (*node, *node, int, int) {
	var pos, dv int
	var current, parent *node

	current = t.root

	for current.next(key[pos]) != nil {
		parent = current
		current = current.next(key[pos])
		pos++

		if len(current.prefix) > 0 {
			dv = divergence(current.prefix, key[pos:])

			if len(current.prefix) > dv {
				return parent, current, pos, dv
			}

			pos = pos + dv
		}

		// if we've found the key, return its parent current so it can be pointed to the new current
		if pos == len(key) {
			return parent, current, pos, dv
		}
	}

	return current, nil, pos, dv
}

func (t *ART) insertNode(key []byte, value interface{}, parent, current *node, pos, dv int) bool {
	n := newNode(Node4, key[pos+1:], value)
	return parent.swapNext(key[pos], nil, n)
}

func (t *ART) updateNode(key []byte, value interface{}, parent, current *node, pos, dv int) bool {
	edgePos := pos - (len(current.prefix) + 1)

	n := newNode(Node4, current.prefix, value)
	n.edges = current.edges

	return parent.swapNext(key[edgePos], current, n)
}

func (t *ART) splitTwoWay(key []byte, value interface{}, parent, current *node, pos, dv int) bool {
	var pfx []byte

	// fix issue where key is found, but is occupied by another current with prefix
	if len(key) > pos {
		pfx = key[pos : pos+dv]
	}

	n1 := newNode(Node4, pfx, value)
	n2 := newNode(Node4, current.prefix[dv+1:], current.value)
	n2.edges = current.edges

	n1.setNext(current.prefix[dv], n2)

	return parent.swapNext(key[pos-1], current, n1)
}

func (t *ART) splitThreeWay(key []byte, value interface{}, parent, current *node, pos, dv int) bool {
	n1 := newNode(Node4, current.prefix[:dv], nil)
	n2 := newNode(Node4, current.prefix[dv+1:], current.value)
	n3 := newNode(Node4, key[pos+dv+1:], value)

	n2.edges = current.edges

	n1.setNext(current.prefix[dv], n2)
	n1.setNext(key[pos+dv], n3)

	return parent.swapNext(key[pos-1], current, n1)
}

func shouldInsert(key []byte, current, parent *node, pos, dv int) bool {
	return pos < len(key) && current == nil
}

func shouldUpdate(key []byte, current, parent *node, pos, dv int) bool {
	return len(key) == pos && dv == len(current.prefix) || len(key) == pos && len(current.prefix) == 0
}

func shouldSplitTwoWay(key []byte, current, parent *node, pos, dv int) bool {
	return (len(key) - (pos + dv)) == 0
}

func shouldSplitThreeWay(key []byte, current, parent *node, pos, dv int) bool {
	return (len(key) - (pos + dv)) > 0
}

// returns shared and divergent characters respectively
func divergence(prefix, key []byte) int {
	var i int

	for i < len(key) && i < len(prefix) {
		if key[i] != prefix[i] {
			break
		}
		i++
	}

	return i

}
