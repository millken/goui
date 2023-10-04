// SPDX-License-Identifier: GNU Lesser General Public License v2.1

//go:build !android && !ios && !js

package goui

import (
	"context"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/richardwilkes/toolbox"
)

type window struct {
	wnd   *glfw.Window
	panel *Panel
	valid bool
}

func NewWindow() (Window, error) {
	win, err := glfw.CreateWindow(1, 1, "goui", nil, nil)
	if err != nil {
		return nil, err
	}
	return &window{
		wnd:   win,
		valid: true,
	}, nil
}

func (w *window) SetSize(width, height int) {
	w.wnd.SetSize(width, height)
}

func (w *window) Size() (width, height int) {
	return w.wnd.GetSize()
}

func (w *window) Show() {
	w.wnd.Show()
}

func (w *window) Hide() {
	w.wnd.Hide()
}

func (w *window) SetPanel(panel *Panel) {
	w.panel = panel
}

func (w *window) OnClose(callback func(ctx context.Context)) {
	w.wnd.SetCloseCallback(func(win *glfw.Window) {
		ctx := WithGLFWWindow(context.Background(), win)
		callback(ctx)
	})
}

func (w *window) Draw(c *Canvas) {
	if w.wnd != nil {
		toolbox.Call(func() {
			//w.
		})
	}
}
func (w *window) IsValid() bool {
	return w != nil && w.valid
}

// performs any layout that needs to be run by this window or its children.
func (w *window) PerformLayout() {
	rect := w.ContentRect()
	rect.X = 0
	rect.Y = 0
}

func (w *window) ContentRect() Rect {
	if w.IsValid() {
		x, y := w.wnd.GetPos()
		width, height := w.wnd.GetSize()
		return Rect{
			Point: Point{X: float32(x), Y: float32(y)},
			Size:  Size{Width: float32(width), Height: float32(height)},
		}
	}
	return Rect{Size: Size{Width: 1, Height: 1}}
}

func (w *window) Destroy() {
	w.wnd.Destroy()
}
