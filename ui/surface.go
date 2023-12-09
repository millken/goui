package ui

import "gitee.com/millken/goui/internal/skia"

type surface struct {
	context skia.DirectContext
	backend skia.BackendRenderTarget
	surface skia.Surface
	size    Size
	scaleX  float32
	scaleY  float32
}

func (s *surface) flush(syncCPU bool) {
	if s != nil && s.surface != nil {
		skia.ContextFlushAndSubmit(s.context, syncCPU)
	}
}
