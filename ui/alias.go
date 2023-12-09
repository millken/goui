package ui

import (
	"gitee.com/millken/goui/internal/geom"
)

// Point is an alias for geom.Point[float32], for convenience.
type Point = geom.Point[float32]

// NewPoint creates a new Point.
func NewPoint(x, y float32) Point {
	return geom.NewPoint(x, y)
}

// Size is an alias for geom.Size[float32], for convenience.
type Size = geom.Size[float32]

// NewSize creates a new Size.
func NewSize(width, height float32) Size {
	return geom.NewSize(width, height)
}

// Rect is an alias for geom.Rect[float32], for convenience.
type Rect = geom.Rect[float32]

// NewRect creates a new Rect.
func NewRect(x, y, width, height float32) Rect {
	return geom.NewRect(x, y, width, height)
}

// Insets is an alias for geom.Insets[float32], for convenience.
type Insets = geom.Insets[float32]

// NewUniformInsets creates a new Insets whose edges all have the same value.
func NewUniformInsets(amount float32) Insets {
	return geom.NewUniformInsets(amount)
}

// NewSymmetricInsets returns an Insets whose edges match their opposite edge.
func NewSymmetricInsets(h, v float32) Insets {
	return geom.NewSymmetricInsets(h, v)
}

// NewHorizontalInsets creates a new Insets whose left and right edges have the specified value.
func NewHorizontalInsets(amount float32) Insets {
	return geom.NewHorizontalInsets(amount)
}

// NewVerticalInsets creates a new Insets whose top and bottom edges have the specified value.
func NewVerticalInsets(amount float32) Insets {
	return geom.NewVerticalInsets(amount)
}

const (
	// DefaultMaxSize is the default size that should be used for a maximum dimension if the target has no real
	// preference and can be expanded beyond its preferred size. This is intentionally not something very large to allow
	// basic math operations an opportunity to succeed when laying out panels. It is perfectly acceptable to use a
	// larger value than this, however, if that makes sense for your specific target.
	DefaultMaxSize = 10000
	// StdHSpacing is the typical spacing between columns.
	StdHSpacing = 8
	// StdVSpacing is the typical spacing between rows.
	StdVSpacing = 4
)

// StdInsets returns insets preset to the standard spacing.
func StdInsets() Insets {
	return Insets{
		Top:    StdVSpacing,
		Left:   StdHSpacing,
		Bottom: StdVSpacing,
		Right:  StdHSpacing,
	}
}

// Matrix is an alias for geom.Matrix2D[float32], for convenience.
type Matrix = geom.Matrix[float32]

// NewIdentityMatrix creates a new identity transformation 2D matrix.
func NewIdentityMatrix() Matrix {
	return geom.NewIdentityMatrix[float32]()
}

// NewTranslationMatrix creates a new 2D matrix that translates by 'tx' and 'ty'.
func NewTranslationMatrix(tx, ty float32) Matrix {
	return geom.NewTranslationMatrix(tx, ty)
}

// NewScaleMatrix creates a new 2D matrix that scales by 'sx' and 'sy'.
func NewScaleMatrix(sx, sy float32) Matrix {
	return geom.NewScaleMatrix(sx, sy)
}

// NewRotationMatrix creates a new 2D matrix that rotates by 'radians'. Positive values are clockwise.
func NewRotationMatrix(radians float32) Matrix {
	return geom.NewRotationMatrix(radians)
}

// NewRotationByDegreesMatrix creates a new 2D matrix that rotates by 'degrees'. Positive values are clockwise.
func NewRotationByDegreesMatrix(degrees float32) Matrix {
	return geom.NewRotationByDegreesMatrix(degrees)
}
