// SPDX-License-Identifier: GNU Lesser General Public License v2.1

//go:build !android && !ios && !js

package goui

import (
	"context"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/richardwilkes/toolbox"
)

type window struct {
	win *glfw.Window
	panel *Panel
}

func NewWindow() (Window, error) {
	win, err := glfw.CreateWindow(1, 1, "goui", nil, nil)
	if err != nil {
		return nil, err
	}
	return &window{
		win: win,
	}, nil
}

func (w *window) SetSize(width, height int) {
	w.win.SetSize(width, height)
}

func (w *window) Size() (width, height int) {
	return w.win.GetSize()
}

func (w *window) Show() {
	w.win.Show()
}

func (w *window) Hide() {
	w.win.Hide()
}

func (w *window) SetPanel(panel *Panel) {
	w.panel = panel
}

func (w *window) OnClose(callback func(ctx context.Context)) {
	w.win.SetCloseCallback(func(win *glfw.Window) {
		ctx := WithGLFWWindow(context.Background(), win)
		callback(ctx)
	})
}

func (w *window) Draw(c *Canvas) {
	if w.win != nil {
		toolbox.Call(func() {
			w.
		})
	}
}

func (w *window) Destroy() {
	w.win.Destroy()
}
