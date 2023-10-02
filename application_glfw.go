// SPDX-License-Identifier: GNU Lesser General Public License v2.1

//go:build !android && !ios && !js

package goui

import (
	"context"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type (
	glfwWindowContextKey struct{}
)

func WithGLFWWindow(ctx context.Context, win *glfw.Window) context.Context {
	return context.WithValue(ctx, glfwWindowContextKey{}, win)
}

func GLFWWindowFromContext(ctx context.Context) (*glfw.Window, bool) {
	win, ok := ctx.Value(glfwWindowContextKey{}).(*glfw.Window)
	return win, ok
}

//	func NewApplication[T any](t T) *T {
//		app := new(T)
//		return (*T)(unsafe.Pointer(app))
//	}
func init() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.Visible, glfwEnabled(false))
	glfw.WindowHint(glfw.Resizable, glfwEnabled(true))
	glfw.WindowHint(glfw.Decorated, glfwEnabled(true))

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	// macOS requires both of these hints to be set
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
}

func Run(app Application) error {

	defer glfw.Terminate()
	if err := app.Startup(context.Background()); err != nil {
		return err
	}

	app.Window().OnClose(func(ctx context.Context) {
		win, ok := GLFWWindowFromContext(ctx)
		if !ok {
			return
		}
		win.Destroy()
		app.Shutdown(ctx)
	})
	for {
		glfw.WaitEvents()
	}
}
func postEmptyEvent() {
	glfw.PostEmptyEvent()
}
