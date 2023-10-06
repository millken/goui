package rmlui

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

func NewRenterInterface(renderer *sdl.Renderer) SwigcptrRenderInterface {
	a := unsafe.Pointer(renderer)
	return SwigcptrRenderInterface(a)
}

func NewSystemInterface(window *sdl.Window) SwigcptrSystemInterface {
	return SwigcptrSystemInterface(unsafe.Pointer(window))
}
