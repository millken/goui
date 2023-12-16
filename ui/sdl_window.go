package ui

import "github.com/veandco/go-sdl2/sdl"

var (
	windowFlags uint32
)

func init() {
	if IsAndroid() {
		windowFlags = sdl.WINDOW_OPENGL | sdl.WINDOW_RESIZABLE |
			sdl.WINDOW_BORDERLESS | sdl.WINDOW_FULLSCREEN_DESKTOP |
			sdl.WINDOW_ALLOW_HIGHDPI
	} else {
		windowFlags = sdl.WINDOW_OPENGL | sdl.WINDOW_RESIZABLE
	}
}
