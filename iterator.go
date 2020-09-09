package art

// Iterate over every key from a given point
func (t *ART) Iterate(from []byte, fn func(key []byte, value Comparable)) {
	var current *node

	if len(from) > 0 {
		_, current, _, _ = t.find(from)
	} else {
		current = t.root
	}

	t.iterate(from, current, fn)
}

func (t *ART) iterate(key []byte, current *node, fn func(key []byte, value Comparable)) {
	if current.edges == nil {
		return
	}

	for i := 0; i < 256; i++ {
		next := current.next(byte(i))
		if next == nil {
			continue
		}

		ckey := make([]byte, len(key))
		copy(ckey, key)

		ckey = append(ckey, byte(i))

		if len(next.prefix) > 0 {
			ckey = append(ckey, next.prefix...)
		}

		if next.value != nil {
			fn(ckey, next.value)
		}

		t.iterate(ckey, next, fn)
	}
}
