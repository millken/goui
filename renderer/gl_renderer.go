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
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/millken/thorvg"
)

type GLRenderer struct {
	canvas        *thorvg.Canvas
	width, height uint32
	lastCtx       unsafe.Pointer
}

func NewGLRenderer() (*GLRenderer, error) {
	// 初始化 OpenGL 上下文和资源
	glErr := gl.Init()
	if glErr != nil {
		return nil, glErr
	}
	// 确保当前线程已有有效的 GL 上下文
	if C.getCurrentGLContext() == nil {
		return nil, fmt.Errorf("no current OpenGL context; create window and make context current before NewGLRenderer")
	}
	canvas := thorvg.NewGlCanvas()
	ctx := C.getCurrentGLContext()

	return &GLRenderer{
		canvas:  canvas,
		lastCtx: unsafe.Pointer(ctx),
	}, nil
}

func (r *GLRenderer) Cleanup() {
	// 清理 ThorVG 画布等资源
	if r.canvas != nil {
		r.canvas.Destroy()
		r.canvas = nil
	}
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
	r.lastCtx = unsafe.Pointer(ctx)
	return r.canvas.GlCanvasSetTarget(r.lastCtx, 0, r.width, r.height, thorvg.ColorspaceABGR8888S)
}

// BeginFrame 准备一帧绘制并返回 ThorVG 画布

func (r *GLRenderer) BeginFrame(width, height int) (*thorvg.Canvas, error) {
	// 若 GL 上下文已切换（例如在多窗口会话之间），需要重新绑定目标
	cur := unsafe.Pointer(C.getCurrentGLContext())
	if cur != nil && cur != r.lastCtx {
		// 若尺寸未知，使用当前保存的尺寸；否则按传入尺寸
		if uint32(width) != 0 && uint32(height) != 0 {
			r.width, r.height = uint32(width), uint32(height)
		}
		r.lastCtx = cur
		if err := r.canvas.GlCanvasSetTarget(r.lastCtx, 0, r.width, r.height, thorvg.ColorspaceABGR8888S); err != nil {
			return nil, err
		}
	}
	if uint32(width) != r.width || uint32(height) != r.height {
		if err := r.resize(uint32(width), uint32(height)); err != nil {
			return nil, err
		}
	}
	// 确保每帧都按当前目标尺寸设置 viewport（viewport 是上下文局部状态）
	if r.width > 0 && r.height > 0 {
		gl.Viewport(0, 0, int32(r.width), int32(r.height))
	}
	// 在绘制开始前清理默认帧缓冲，避免上一帧残留
	gl.ClearColor(0.12, 0.12, 0.12, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	return r.canvas, nil
}

func (r *GLRenderer) EndFrame() error {
	var err error
	if err = r.canvas.Update(); err != nil {
		return err
	}
	if err = r.canvas.Draw(true); err != nil {
		return err
	}
	if err = r.canvas.Sync(); err != nil {
		return err
	}
	return nil
}

// ResetScene 清空 ThorVG 场景，常用于窗口尺寸变更或需要彻底重建场景时
func (r *GLRenderer) ResetScene(width, height int) error {
	if r.canvas != nil {
		r.canvas.Destroy()
	}
	r.canvas = thorvg.NewGlCanvas()
	// 重新保存尺寸并绑定当前 GL 上下文和目标尺寸
	if width > 0 && height > 0 {
		r.width, r.height = uint32(width), uint32(height)
	}
	ctx := C.getCurrentGLContext()
	r.lastCtx = unsafe.Pointer(ctx)
	return r.canvas.GlCanvasSetTarget(r.lastCtx, 0, r.width, r.height, thorvg.ColorspaceABGR8888S)
}
