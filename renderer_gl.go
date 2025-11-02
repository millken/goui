package goui

import (
	"fmt"

	"github.com/millken/go/rgfw"
	"github.com/millken/go/thorvg"
)

// GLRenderer 是基于 RGFW OpenGL 上下文的 ThorVG 画布渲染器
type GLRenderer struct {
	win           *rgfw.Window
	canvas        *thorvg.Canvas
	width, height uint32
}

// NewGLRenderer 创建一个 GLRenderer，并将 ThorVG 画布绑定到窗口的 OpenGL 上下文
func NewGLRenderer(win *rgfw.Window, width, height uint32) (*GLRenderer, error) {
	c := thorvg.NewGlCanvas()
	if c == nil {
		return nil, fmt.Errorf("thorvg.NewGlCanvas returned nil")
	}
	r := &GLRenderer{win: win, canvas: c}
	if err := r.resize(width, height); err != nil {
		c.Destroy()
		return nil, err
	}
	return r, nil
}

func (r *GLRenderer) resize(w, h uint32) error {
	if w == 0 || h == 0 {
		return nil
	}
	r.width, r.height = w, h
	return r.canvas.GlCanvasSetTarget(r.win.GetOpenGLContext(), 0, r.width, r.height, thorvg.ColorspaceABGR8888S)
}

// BeginFrame 准备一帧绘制并返回 ThorVG 画布
func (r *GLRenderer) BeginFrame(width, height int) (*thorvg.Canvas, error) {
	if uint32(width) != r.width || uint32(height) != r.height {
		if err := r.resize(uint32(width), uint32(height)); err != nil {
			return nil, err
		}
	}
	return r.canvas, nil
}

// EndFrame 完成一帧绘制（调用 ThorVG 的 Update/Draw/Sync）
func (r *GLRenderer) EndFrame() {
	_ = r.canvas.Update()
	_ = r.canvas.Draw(true)
	_ = r.canvas.Sync()
}

// Present 交换缓冲，呈现到屏幕
func (r *GLRenderer) Present() { r.win.SwapBuffersOpenGL() }

// SetViewport 设置视口（由具体后续实现扩展）
func (r *GLRenderer) SetViewport(x, y, width, height int) { /* noop for now */ }

// Cleanup 释放资源
func (r *GLRenderer) Cleanup() {
	if r.canvas != nil {
		r.canvas.Destroy()
		r.canvas = nil
	}
}
