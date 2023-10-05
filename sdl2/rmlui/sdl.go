package rmlui

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type RenterInterface uintptr

func NewRenterInterface(renderer *sdl.Renderer) SwigcptrRml_RenderInterface {
	a := unsafe.Pointer(renderer)
	return SwigcptrRml_RenderInterface(a)
}

func NewSystemInterface(window *sdl.Window) SwigcptrSystemInterface {
	return SwigcptrSystemInterface(unsafe.Pointer(window))
}
