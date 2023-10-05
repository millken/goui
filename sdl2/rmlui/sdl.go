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

type SystemInterface uintptr

func NewSystemInterface(window *sdl.Window) SwigcptrRml_SystemInterface {
	return SwigcptrRml_SystemInterface(unsafe.Pointer(window))
}
