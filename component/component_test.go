package component

import (
	"testing"

	"github.com/dnsoa/go/assert"
)

func TestRoot(t *testing.T) {
	r := assert.New(t)
	root := NewFlex()
	root.SetRoot(root)
	r.Equal(root, root.Root())

	child1 := NewFlex()
	root.AddChild(child1)
	r.Equal(root, child1.Root())

	child2 := NewFlex()
	child1.AddChild(child2)
	r.Equal(root, child2.Root())

	child3 := NewFlex()
	r.NotEqual(root, child3.Root())
}
