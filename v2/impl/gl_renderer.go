package impl

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	v2 "github.com/millken/goui/v2"
)

var _ v2.Renderer = (*GLRenderer)(nil)

// GLRenderer 是使用 OpenGL 4.6 Core Profile 的简单参考实现。
// 注意：要求在主线程（被 LockOSThread 的线程）进行 GL 上下文操作。
type GLRenderer struct {
	win    *GLFWWindow
	vao    uint32
	inited bool

	clearR, clearG, clearB, clearA float32
}

func NewOpenGLRenderer() *GLRenderer {
	return &GLRenderer{clearA: 1.0}
}

func (r *GLRenderer) Initialize(window v2.Window) error {
	// 仅支持 GLFWWindow
	gw, ok := window.(*GLFWWindow)
	if !ok {
		return fmt.Errorf("GLRenderer requires *GLFWWindow, got %T", window)
	}
	r.win = gw

	if r.win.win == nil {
		return fmt.Errorf("GLFW window is nil")
	}

	// 绑定 OpenGL 上下文到当前线程
	r.win.win.MakeContextCurrent()

	// 初始化 go-gl（需在有当前上下文的线程调用）
	if err := gl.Init(); err != nil {
		return fmt.Errorf("gl.Init failed: %w", err)
	}

	// 开启垂直同步（如需）
	glfw.SwapInterval(1)

	// Core profile 要求至少绑定一个 VAO
	gl.GenVertexArrays(1, &r.vao)
	gl.BindVertexArray(r.vao)

	r.inited = true
	return nil
}

func (r *GLRenderer) Cleanup() {
	if !r.inited {
		return
	}
	if r.vao != 0 {
		gl.DeleteVertexArrays(1, &r.vao)
		r.vao = 0
	}
	// 可在此处分离上下文，如需要：glfw.DetachCurrentContext()
	r.inited = false
}

func (r *GLRenderer) BeginFrame(width, height int) (v2.Canvas, error) {
	if !r.inited {
		return nil, fmt.Errorf("renderer not initialized")
	}
	// 设置视口并清屏
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.ClearColor(r.clearR, r.clearG, r.clearB, r.clearA)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	return v2.NopCanvas{}, nil
}

func (r *GLRenderer) EndFrame() {
	// 本示例渲染路径简单，无需在 EndFrame 做事
}

func (r *GLRenderer) Present() {
	if r.win != nil && r.win.win != nil {
		r.win.win.SwapBuffers()
	}
}

func (r *GLRenderer) SetViewport(x, y, width, height int) {
	if !r.inited {
		return
	}
	gl.Viewport(int32(x), int32(y), int32(width), int32(height))
}
