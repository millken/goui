// Package renderer provides the GLRenderer which bridges OpenGL contexts
// and the ThorVG canvas for GPU-accelerated vector rendering.
package renderer

/*
#cgo linux LDFLAGS: -lGL
#cgo windows LDFLAGS: -lopengl32
#cgo darwin LDFLAGS: -framework OpenGL
#include <stdlib.h>
#ifdef __linux__
#include <GL/glx.h>
typedef GLXContext GLContext;
GLContext getCurrentGLContext() { return glXGetCurrentContext(); }
#elif defined(_WIN32)
#include <windows.h>
typedef HGLRC GLContext;
GLContext getCurrentGLContext() { return wglGetCurrentContext(); }
#elif defined(__APPLE__)
#include <OpenGL/OpenGL.h>
typedef CGLContextObj GLContext;
GLContext getCurrentGLContext() { return CGLGetCurrentContext(); }
#endif
*/
import "C"
import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/millken/go/thorvg"
)

type GLRenderer struct {
	canvas        *thorvg.Canvas
	width, height uint32
}

func NewGLRenderer() (*GLRenderer, error) {
	// 初始化 OpenGL 上下文和资源
	glErr := gl.Init()
	if glErr != nil {
		return nil, glErr
	}
	canvas := thorvg.NewGlCanvas()

	return &GLRenderer{
		canvas: canvas,
	}, nil
}

func (r *GLRenderer) Cleanup() {
	// 清理 OpenGL 资源
}

func (r *GLRenderer) resize(w, h uint32) error {
	if w == 0 || h == 0 {
		return nil
	}
	r.width, r.height = w, h
	// 更新 OpenGL 视口
	gl.Viewport(0, 0, int32(w), int32(h))
	// 获取当前 GL 上下文句柄
	ctx := C.getCurrentGLContext()
	return r.canvas.GlCanvasSetTarget(unsafe.Pointer(ctx), 0, r.width, r.height, thorvg.ColorspaceABGR8888S)
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

func (r *GLRenderer) EndFrame() {
	var err error
	// 渲染前清理背景，仅在确实要绘制时执行
	gl.ClearColor(0.12, 0.12, 0.12, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	if err = r.canvas.Update(); err != nil {
		panic(err)
	}
	if err = r.canvas.Draw(true); err != nil {
		panic(err)
	}
	if err = r.canvas.Sync(); err != nil {
		panic(err)
	}
}

func (r *GLRenderer) Present() {
	// 将渲染结果呈现到屏幕
}
