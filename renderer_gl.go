package goui

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/millken/go/thorvg"
)

// GLRenderer 是基于 RGFW OpenGL 上下文的 ThorVG 画布渲染器
type GLRenderer struct {
	win           *glfw.Window
	canvas        *thorvg.Canvas
	width, height uint32
}

// NewGLRenderer 创建一个 GLRenderer，并将 ThorVG 画布绑定到窗口的 OpenGL 上下文
func NewGLRenderer(win *glfw.Window, width, height uint32) (*GLRenderer, error) {
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
	// 更新 OpenGL 视口
	gl.Viewport(0, 0, int32(w), int32(h))
	// 在 GLFW 中，ThorVG OpenGL 目标绑定通常使用当前上下文，无需显式从窗口获取上下文句柄。
	// 传入 0 作为上下文句柄（取决于 thorvg 绑定实现），若需要可用平台特定 API 获取。
	return r.canvas.GlCanvasSetTarget(nil, 0, r.width, r.height, thorvg.ColorspaceABGR8888S)
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
func (r *GLRenderer) Present() { r.win.SwapBuffers() }

// SetViewport 设置视口（由具体后续实现扩展）
func (r *GLRenderer) SetViewport(x, y, width, height int) { /* noop for now */ }

// Cleanup 释放资源
func (r *GLRenderer) Cleanup() {
	if r.canvas != nil {
		r.canvas.Destroy()
		r.canvas = nil
	}
}
