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
	n.getEdges().keys[0] = b("a")
	n.getEdges().keys[1] = b("c")
	n.getEdges().edges[0] = newNode(Node4, nil, nil)
	n.getEdges().edges[1] = newNode(Node4, nil, nil)
	n.getEdges().children = 2

	e = n.next(b("a"))
	assert.NotNil(t, e)

	e = n.next(b("b"))
	assert.Nil(t, e)

	e = n.next(b("c"))
	assert.NotNil(t, e)
}

func TestNode4SetNext(t *testing.T) {
	n := newNode(Node4, nil, nil)
	n.getEdges().keys[0] = b("a")
	n.getEdges().keys[1] = b("c")
	n.getEdges().edges[0] = newNode(Node4, nil, nil)
	n.getEdges().edges[1] = newNode(Node4, nil, nil)
	n.getEdges().children = uint8(2)

	n.setNext(b("b"), newNode(Node4, nil, nil))

	assert.Equal(t, b("a"), n.getEdges().keys[0])
	assert.Equal(t, b("c"), n.getEdges().keys[1])
	assert.Equal(t, b("b"), n.getEdges().keys[2])

	assert.NotNil(t, n.getEdges().edges[0])
	assert.NotNil(t, n.getEdges().edges[1])
	assert.NotNil(t, n.getEdges().edges[2])
	assert.Nil(t, n.getEdges().edges[3])

	assert.Equal(t, uint8(3), n.getEdges().children)
}

func TestNode16Next(t *testing.T) {
	n := newNode(Node16, nil, nil)

	e := n.next(b("a"))
	assert.Nil(t, e)

	n = newNode(Node16, nil, nil)

	n.getEdges().keys[0] = b("a")
	n.getEdges().keys[1] = b("c")
	n.getEdges().edges[0] = newNode(Node4, nil, nil)
	n.getEdges().edges[1] = newNode(Node4, nil, nil)
	n.getEdges().children = uint8(2)

	e = n.next(b("a"))
	assert.NotNil(t, e)

	e = n.next(b("b"))
	assert.Nil(t, e)

	e = n.next(b("c"))
	assert.NotNil(t, e)
}

func TestNode16SetNext(t *testing.T) {
	n := newNode(Node16, nil, nil)

	n.getEdges().keys[0] = b("a")
	n.getEdges().keys[1] = b("c")
	n.getEdges().edges[0] = newNode(Node4, nil, nil)
	n.getEdges().edges[1] = newNode(Node4, nil, nil)
	n.getEdges().children = uint8(2)

	n.setNext(b("b"), newNode(Node4, nil, nil))

	assert.Equal(t, b("a"), n.getEdges().keys[0])
	assert.Equal(t, b("b"), n.getEdges().keys[1])
	assert.Equal(t, b("c"), n.getEdges().keys[2])

	assert.NotNil(t, n.getEdges().edges[0])
	assert.NotNil(t, n.getEdges().edges[1])
	assert.NotNil(t, n.getEdges().edges[2])
	assert.Nil(t, n.getEdges().edges[3])

	assert.Equal(t, uint8(3), n.getEdges().children)

	n = newNode(Node4, nil, nil)

	for x := 0; x < 8; x++ {
		for i := 32; i > 16; i-- {
			n.setNext(byte(i), newNode(Node4, nil, nil))
		}
	}

	assert.Equal(t, uint8(16), n.getEdges().children)
}

func TestNode48Next(t *testing.T) {
	n := newNode(Node48, nil, nil)

	e := n.next(b("a"))
	assert.Nil(t, e)

	n = newNode(Node48, nil, nil)
	n.getEdges().keys[0] = 1
	n.getEdges().keys[1] = 0
	n.getEdges().keys[2] = 2
	n.getEdges().edges[0] = newNode(Node4, nil, nil)
	n.getEdges().edges[1] = newNode(Node4, nil, nil)
	n.getEdges().children = 2

	e = n.next(byte(0))
	assert.NotNil(t, e)

	e = n.next(byte(1))
	assert.Nil(t, e)

	e = n.next(byte(2))
	assert.NotNil(t, e)
}

func TestNode48SetNext(t *testing.T) {
	n := newNode(Node48, nil, nil)
	n.getEdges().keys[97] = 1
	n.getEdges().keys[99] = 2
	n.getEdges().edges[0] = newNode(Node4, nil, nil)
	n.getEdges().edges[1] = newNode(Node4, nil, nil)
	n.getEdges().children = 2

	n.setNext(b("b"), newNode(Node4, nil, nil))

	assert.Equal(t, uint8(1), n.getEdges().keys[97])
	assert.Equal(t, uint8(3), n.getEdges().keys[98])
	assert.Equal(t, uint8(2), n.getEdges().keys[99])

	assert.NotNil(t, n.getEdges().edges[0])
	assert.NotNil(t, n.getEdges().edges[1])
	assert.NotNil(t, n.getEdges().edges[2])
	assert.Nil(t, n.getEdges().edges[3])

	assert.Equal(t, uint8(3), n.getEdges().children)
}

func TestNode256Next(t *testing.T) {
	n := newNode(Node256, nil, nil)

	e := n.next(b("a"))
	assert.Nil(t, e)

	n = newNode(Node256, nil, nil)
	n.getEdges().edges[0] = newNode(Node4, nil, nil)
	n.getEdges().edges[2] = newNode(Node4, nil, nil)

	e = n.next(byte(0))
	assert.NotNil(t, e)

	e = n.next(byte(1))
	assert.Nil(t, e)

	e = n.next(byte(2))
	assert.NotNil(t, e)
}

func TestNode256SetNext(t *testing.T) {
	n := newNode(Node256, nil, nil)
	n.getEdges().edges[97] = newNode(Node4, nil, nil)
	n.getEdges().edges[99] = newNode(Node4, nil, nil)
	n.getEdges().children = 2

	n.setNext(b("b"), newNode(Node4, nil, nil))

	assert.NotNil(t, n.getEdges().edges[97])
	assert.NotNil(t, n.getEdges().edges[98])
	assert.NotNil(t, n.getEdges().edges[99])
	assert.Nil(t, n.getEdges().edges[3])

	assert.Equal(t, uint8(3), n.getEdges().children)
}

func b(s string) byte {
	return byte(s[0])
}
