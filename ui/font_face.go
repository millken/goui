package ui

import (
	"runtime"

	"github.com/millken/goui/internal/skia"
)

// FontFace holds the immutable portions of a font description.
type FontFace struct {
	face skia.TypeFace
}

func newFace(face skia.TypeFace) *FontFace {
	if face == nil {
		return nil
	}
	f := &FontFace{face: face}
	runtime.SetFinalizer(f, func(obj *FontFace) {
		skia.TypeFaceUnref(obj.face)
	})
	return f
}
