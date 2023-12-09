package ui

import (
	"runtime"

	"github.com/pkg/errors"

	"gitee.com/millken/goui/internal/skia"
)

// Direction holds the direction of a path.
type Direction byte

// Possible values for Direction.
const (
	Clockwise Direction = iota
	CounterClockwise
)

// Path holds geometry.
type Path struct {
	path skia.Path
}

func newPath(path skia.Path) *Path {
	p := &Path{path: path}
	runtime.SetFinalizer(p, func(obj *Path) {
		skia.PathDelete(obj.path)
	})
	return p
}

// NewPath creates a new, empty path.
func NewPath() *Path {
	return newPath(skia.PathNew())
}

// NewPathFromSVGString attempts to create a path from the given SVG string.
func NewPathFromSVGString(svg string) (*Path, error) {
	p := NewPath()
	if !skia.PathParseSVGString(p.path, svg) {
		return nil, errors.New("unable to parse SVG string into path")
	}
	return p, nil
}

// LineTo appends a straight line segment.
func (p *Path) LineTo(x, y float32) {
	skia.PathLineTo(p.path, x, y)
}

// MoveTo begins a new contour at the specified point.
func (p *Path) MoveTo(x, y float32) {
	skia.PathMoveTo(p.path, x, y)
}

// Rect adds a rectangle to the path with a clockwise direction. The rectangle is a complete contour, i.e. it starts
// with a MoveTo and ends with a Close operation.
func (p *Path) Rect(bounds Rect) {
	skia.PathAddRect(p.path, bounds, skia.Direction(Clockwise))
}

// RoundedRect adds a rectangle with curved corners to the path with a clockwise direction. The rectangle is a complete
// contour, i.e. it starts with a MoveTo and ends with a Close operation.
func (p *Path) RoundedRect(bounds Rect, radiusX, radiusY float32) {
	skia.PathAddRoundedRect(p.path, bounds, radiusX, radiusY, skia.Direction(Clockwise))
}

// Clone this path.
func (p *Path) Clone() *Path {
	return newPath(skia.PathClone(p.path))
}

// Close the current contour.
func (p *Path) Close() {
	skia.PathClose(p.path)
}
