package art

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode4Next(t *testing.T) {
	n := newNode4()

	e := n.next(b("a"))
	assert.Nil(t, e)

	n = newNode4()
	n.keys[0] = b("a")
	n.keys[1] = b("c")
	n.edges[0] = newNode4()
	n.edges[1] = newNode4()
	n.children = 2

	e = n.next(b("a"))
	assert.NotNil(t, e)

	e = n.next(b("b"))
	assert.Nil(t, e)

	e = n.next(b("c"))
	assert.NotNil(t, e)
}

func TestNode4SetNext(t *testing.T) {
	n := newNode4()
	n.keys[0] = b("a")
	n.keys[1] = b("c")
	n.edges[0] = newNode4()
	n.edges[1] = newNode4()
	n.children = uint8(2)

	n.setNext(b("b"), newNode4())

	assert.Equal(t, b("a"), n.keys[0])
	assert.Equal(t, b("c"), n.keys[1])
	assert.Equal(t, b("b"), n.keys[2])

	assert.NotNil(t, n.edges[0])
	assert.NotNil(t, n.edges[1])
	assert.NotNil(t, n.edges[2])
	assert.Nil(t, n.edges[3])

	assert.Equal(t, uint8(3), n.children)
}

func TestNode16Next(t *testing.T) {
	n := newNode16()

	e := n.next(b("a"))
	assert.Nil(t, e)

	n = newNode16()

	n.keys[0] = b("a")
	n.keys[1] = b("c")
	n.edges[0] = newNode4()
	n.edges[1] = newNode4()
	n.children = uint8(2)

	e = n.next(b("a"))
	assert.NotNil(t, e)

	e = n.next(b("b"))
	assert.Nil(t, e)

	e = n.next(b("c"))
	assert.NotNil(t, e)
}

func TestNode16SetNext(t *testing.T) {
	n := newNode16()

	n.keys[0] = b("a")
	n.keys[1] = b("c")
	n.edges[0] = newNode4()
	n.edges[1] = newNode4()
	n.children = uint8(2)

	n.setNext(b("b"), newNode4())

	assert.Equal(t, b("a"), n.keys[0])
	assert.Equal(t, b("b"), n.keys[1])
	assert.Equal(t, b("c"), n.keys[2])

	assert.NotNil(t, n.edges[0])
	assert.NotNil(t, n.edges[1])
	assert.NotNil(t, n.edges[2])
	assert.Nil(t, n.edges[3])

	assert.Equal(t, uint8(3), n.children)
}

func TestNode48Next(t *testing.T) {
	n := newNode48()

	e := n.next(b("a"))
	assert.Nil(t, e)

	n = newNode48()
	n.keys[0] = 1
	n.keys[1] = 0
	n.keys[2] = 2
	n.edges[0] = newNode4()
	n.edges[1] = newNode4()
	n.children = 2

	e = n.next(byte(0))
	assert.NotNil(t, e)

	e = n.next(byte(1))
	assert.Nil(t, e)

	e = n.next(byte(2))
	assert.NotNil(t, e)
}

func TestNode48SetNext(t *testing.T) {
	n := newNode48()
	n.keys[97] = 1
	n.keys[99] = 2
	n.edges[0] = newNode4()
	n.edges[1] = newNode4()
	n.children = 2

	n.setNext(b("b"), newNode4())

	assert.Equal(t, uint8(1), n.keys[97])
	assert.Equal(t, uint8(3), n.keys[98])
	assert.Equal(t, uint8(2), n.keys[99])

	assert.NotNil(t, n.edges[0])
	assert.NotNil(t, n.edges[1])
	assert.NotNil(t, n.edges[2])
	assert.Nil(t, n.edges[3])

	assert.Equal(t, uint8(3), n.children)
}

func TestNode256Next(t *testing.T) {
	n := newNode256()

	e := n.next(b("a"))
	assert.Nil(t, e)

	n = newNode256()
	n.edges[0] = newNode4()
	n.edges[2] = newNode4()

	e = n.next(byte(0))
	assert.NotNil(t, e)

	e = n.next(byte(1))
	assert.Nil(t, e)

	e = n.next(byte(2))
	assert.NotNil(t, e)
}

func TestNode256SetNext(t *testing.T) {
	n := newNode256()
	n.edges[97] = newNode4()
	n.edges[99] = newNode4()
	n.children = 2

	n.setNext(b("b"), newNode4())

	assert.NotNil(t, n.edges[97])
	assert.NotNil(t, n.edges[98])
	assert.NotNil(t, n.edges[99])
	assert.Nil(t, n.edges[3])

	assert.Equal(t, uint8(3), n.children)
}

func b(s string) byte {
	return byte(s[0])
}
