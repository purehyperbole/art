package art

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeSetNext(t *testing.T) {
	n := newNode(Node4, nil, nil)

	for i := 255; i >= 0; i-- {
		n.setNext(byte(i), newNode(Node4, nil, nil))
	}

	for i := 0; i < 256; i++ {
		assert.NotNil(t, n.next(byte(i)))
	}
}

func TestNode4Next(t *testing.T) {
	n := newNode(Node4, nil, nil)

	e := n.next(b("a"))
	assert.Nil(t, e)

	n = newNode(Node4, nil, nil)
	n.getEdges().(*edges4).keys[0] = b("a")
	n.getEdges().(*edges4).keys[1] = b("c")
	n.getEdges().(*edges4).edges[0] = newNode(Node4, nil, nil)
	n.getEdges().(*edges4).edges[1] = newNode(Node4, nil, nil)
	n.getEdges().(*edges4).children = 2

	e = n.next(b("a"))
	assert.NotNil(t, e)

	e = n.next(b("b"))
	assert.Nil(t, e)

	e = n.next(b("c"))
	assert.NotNil(t, e)
}

func TestNode4SetNext(t *testing.T) {
	n := newNode(Node4, nil, nil)
	n.getEdges().(*edges4).keys[0] = b("a")
	n.getEdges().(*edges4).keys[1] = b("c")
	n.getEdges().(*edges4).edges[0] = newNode(Node4, nil, nil)
	n.getEdges().(*edges4).edges[1] = newNode(Node4, nil, nil)
	n.getEdges().(*edges4).children = uint8(2)

	n.setNext(b("b"), newNode(Node4, nil, nil))

	assert.Equal(t, b("a"), n.getEdges().(*edges4).keys[0])
	assert.Equal(t, b("b"), n.getEdges().(*edges4).keys[1])
	assert.Equal(t, b("c"), n.getEdges().(*edges4).keys[2])

	assert.NotNil(t, n.getEdges().(*edges4).edges[0])
	assert.NotNil(t, n.getEdges().(*edges4).edges[1])
	assert.NotNil(t, n.getEdges().(*edges4).edges[2])
	assert.Nil(t, n.getEdges().(*edges4).edges[3])

	assert.Equal(t, uint8(3), n.getEdges().(*edges4).children)
}

func TestNode16Next(t *testing.T) {
	n := newNode(Node16, nil, nil)

	e := n.next(b("a"))
	assert.Nil(t, e)

	n = newNode(Node16, nil, nil)

	n.getEdges().(*edges16).keys[0] = b("a")
	n.getEdges().(*edges16).keys[1] = b("c")
	n.getEdges().(*edges16).edges[0] = newNode(Node4, nil, nil)
	n.getEdges().(*edges16).edges[1] = newNode(Node4, nil, nil)
	n.getEdges().(*edges16).children = uint8(2)

	e = n.next(b("a"))
	assert.NotNil(t, e)

	e = n.next(b("b"))
	assert.Nil(t, e)

	e = n.next(b("c"))
	assert.NotNil(t, e)
}

func TestNode16SetNext(t *testing.T) {
	n := newNode(Node16, nil, nil)

	n.getEdges().(*edges16).keys[0] = b("a")
	n.getEdges().(*edges16).keys[1] = b("c")
	n.getEdges().(*edges16).edges[0] = newNode(Node4, nil, nil)
	n.getEdges().(*edges16).edges[1] = newNode(Node4, nil, nil)
	n.getEdges().(*edges16).children = uint8(2)

	n.setNext(b("b"), newNode(Node4, nil, nil))

	assert.Equal(t, b("a"), n.getEdges().(*edges16).keys[0])
	assert.Equal(t, b("b"), n.getEdges().(*edges16).keys[1])
	assert.Equal(t, b("c"), n.getEdges().(*edges16).keys[2])

	assert.NotNil(t, n.getEdges().(*edges16).edges[0])
	assert.NotNil(t, n.getEdges().(*edges16).edges[1])
	assert.NotNil(t, n.getEdges().(*edges16).edges[2])
	assert.Nil(t, n.getEdges().(*edges16).edges[3])

	assert.Equal(t, uint8(3), n.getEdges().(*edges16).children)

	n = newNode(Node4, nil, nil)

	for x := 0; x < 8; x++ {
		for i := 32; i > 16; i-- {
			n.setNext(byte(i), newNode(Node4, nil, nil))
		}
	}

	assert.Equal(t, uint8(16), n.getEdges().(*edges16).children)
}

func TestNode48Next(t *testing.T) {
	n := newNode(Node48, nil, nil)

	e := n.next(b("a"))
	assert.Nil(t, e)

	n = newNode(Node48, nil, nil)
	n.getEdges().(*edges48).keys[0] = 1
	n.getEdges().(*edges48).keys[1] = 0
	n.getEdges().(*edges48).keys[2] = 2
	n.getEdges().(*edges48).edges[0] = newNode(Node4, nil, nil)
	n.getEdges().(*edges48).edges[1] = newNode(Node4, nil, nil)
	n.getEdges().(*edges48).children = 2

	e = n.next(byte(0))
	assert.NotNil(t, e)

	e = n.next(byte(1))
	assert.Nil(t, e)

	e = n.next(byte(2))
	assert.NotNil(t, e)
}

func TestNode48SetNext(t *testing.T) {
	n := newNode(Node48, nil, nil)
	n.getEdges().(*edges48).keys[97] = 1
	n.getEdges().(*edges48).keys[99] = 2
	n.getEdges().(*edges48).edges[0] = newNode(Node4, nil, nil)
	n.getEdges().(*edges48).edges[1] = newNode(Node4, nil, nil)
	n.getEdges().(*edges48).children = 2

	n.setNext(b("b"), newNode(Node4, nil, nil))

	assert.Equal(t, uint8(1), n.getEdges().(*edges48).keys[97])
	assert.Equal(t, uint8(3), n.getEdges().(*edges48).keys[98])
	assert.Equal(t, uint8(2), n.getEdges().(*edges48).keys[99])

	assert.NotNil(t, n.getEdges().(*edges48).edges[0])
	assert.NotNil(t, n.getEdges().(*edges48).edges[1])
	assert.NotNil(t, n.getEdges().(*edges48).edges[2])
	assert.Nil(t, n.getEdges().(*edges48).edges[3])

	assert.Equal(t, uint8(3), n.getEdges().(*edges48).children)
}

func TestNode256Next(t *testing.T) {
	n := newNode(Node256, nil, nil)

	e := n.next(b("a"))
	assert.Nil(t, e)

	n = newNode(Node256, nil, nil)
	n.getEdges().(*edges256).edges[0] = newNode(Node4, nil, nil)
	n.getEdges().(*edges256).edges[2] = newNode(Node4, nil, nil)

	e = n.next(byte(0))
	assert.NotNil(t, e)

	e = n.next(byte(1))
	assert.Nil(t, e)

	e = n.next(byte(2))
	assert.NotNil(t, e)
}

func TestNode256SetNext(t *testing.T) {
	n := newNode(Node256, nil, nil)
	n.getEdges().(*edges256).edges[97] = newNode(Node4, nil, nil)
	n.getEdges().(*edges256).edges[99] = newNode(Node4, nil, nil)
	n.getEdges().(*edges256).children = 2

	n.setNext(b("b"), newNode(Node4, nil, nil))

	assert.NotNil(t, n.getEdges().(*edges256).edges[97])
	assert.NotNil(t, n.getEdges().(*edges256).edges[98])
	assert.NotNil(t, n.getEdges().(*edges256).edges[99])
	assert.Nil(t, n.getEdges().(*edges256).edges[3])

	assert.Equal(t, uint8(3), n.getEdges().(*edges256).children)
}

func b(s string) byte {
	return byte(s[0])
}
