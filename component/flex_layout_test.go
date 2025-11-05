package component

import (
	"testing"

	"github.com/dnsoa/go/assert"
	"github.com/millken/yoga"
)

func TestFlex_PaddingAndGaps_AffectLayout(t *testing.T) {
	r := assert.New(t)

	col := Column("col")
	col.SetPaddingEdges(10, 20, 10, 20) // left, top, right, bottom
	col.SetRowGap(5)

	// create three fixed-size boxes using Button with explicit sizes
	b1 := NewButton("b1")
	b1.SetWidth(50)
	b1.SetMinHeight(20)
	b2 := NewButton("b2")
	b2.SetWidth(50)
	b2.SetMinHeight(20)
	b3 := NewButton("b3")
	b3.SetWidth(50)
	b3.SetMinHeight(20)

	col.AddChild(b1)
	col.AddChild(b2)
	col.AddChild(b3)

	// layout under a fixed root size
	col.Layout(200, 200, yoga.DirectionLTR)

	// check top positions respect padding and gap
	y1 := b1.Bounds().Min.Y
	y2 := b2.Bounds().Min.Y
	y3 := b3.Bounds().Min.Y

	r.Equal(20, y1)      // first child starts at top padding
	r.Equal(y1+20+5, y2) // height 20 + rowGap 5
	r.Equal(y2+20+5, y3)
}
