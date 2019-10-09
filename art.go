package art

import "sync"

// ART an adaptive radix tree implementation
type ART struct {
	root  *node
	locks [256]sync.RWMutex
}

// New creates a new radix tree
func New() *ART {
	return &ART{root: newNode256()}
}

// Insert value into the tree
func (t *ART) Insert(key []byte, value interface{}) {
	t.locks[key[0]].Lock()

	parent, current, pos, dv := t.find(key)

	switch {
	case shouldInsert(key, current, parent, pos, dv):
		t.insertNode(key, value, parent, current, pos, dv)
	case shouldUpdate(key, current, parent, pos, dv):
		t.updateNode(key, value, parent, current, pos, dv)
	case shouldSplitThreeWay(key, current, parent, pos, dv):
		t.splitThreeWay(key, value, parent, current, pos, dv)
	case shouldSplitTwoWay(key, current, parent, pos, dv):
		t.splitTwoWay(key, value, parent, current, pos, dv)
	}

	t.locks[key[0]].Unlock()
}

// Lookup a value from the tree
func (t *ART) Lookup(key []byte) interface{} {
	t.locks[key[0]].RLock()
	_, current, pos, _ := t.find(key)
	t.locks[key[0]].RUnlock()

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

func (t *ART) insertNode(key []byte, value interface{}, parent, current *node, pos, dv int) {
	newNode := newNode4()
	newNode.prefix = key[pos+1:]
	newNode.value = value

	parent.setNext(key[pos], newNode)
}

func (t *ART) updateNode(key []byte, value interface{}, parent, current *node, pos, dv int) {
	current.value = value
}

func (t *ART) splitTwoWay(key []byte, value interface{}, parent, current *node, pos, dv int) {
	var pfx []byte

	// fix issue where key is found, but is occupied by another current with prefix
	if len(key) > pos {
		pfx = key[pos : pos+dv]
	}

	n1 := newNode4()
	n1.prefix = pfx
	n1.value = value

	n1.setNext(current.prefix[dv], current)

	current.prefix = current.prefix[dv+1:]

	parent.setNext(key[pos-1], n1)
}

func (t *ART) splitThreeWay(key []byte, value interface{}, parent, current *node, pos, dv int) {
	n1 := newNode4()
	n1.prefix = current.prefix[:dv]

	n3 := newNode4()
	n3.prefix = key[pos+dv+1:]
	n3.value = value

	n1.setNext(current.prefix[dv], current)
	n1.setNext(key[pos+dv], n3)

	current.prefix = current.prefix[dv+1:]

	parent.setNext(key[pos-1], n1)
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
