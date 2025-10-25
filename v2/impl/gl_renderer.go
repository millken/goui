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
	win      *GLFWWindow
	vao      uint32
	vbo      uint32
	prog     uint32
	colorLoc int32
	inited   bool

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

	// 创建 VBO（4 顶点 * 2 分量 * 4 字节）
	gl.GenBuffers(1, &r.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*2*4, nil, gl.DYNAMIC_DRAW)

	// 构建着色器程序
	prog, err := buildProgram(vertexSrc, fragmentSrc)
	if err != nil {
		return err
	}
	r.prog = prog
	r.colorLoc = gl.GetUniformLocation(r.prog, gl.Str("uColor\x00"))

	gl.UseProgram(r.prog)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	r.inited = true
	return nil
}

func (r *GLRenderer) Cleanup() {
	if !r.inited {
		return
	}
	if r.prog != 0 {
		gl.DeleteProgram(r.prog)
		r.prog = 0
	}
	if r.vbo != 0 {
		gl.DeleteBuffers(1, &r.vbo)
		r.vbo = 0
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
	return &glCanvas{r: r, fbW: width, fbH: height}, nil
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

// glCanvas 是一个最小实现，仅提供 DrawRect(rect,color)
type glCanvas struct {
	r        *GLRenderer
	fbW, fbH int
}

// DrawRect 使用 TRIANGLE_STRIP 渲染矩形
func (c *glCanvas) DrawRect(rect v2.Rect, color v2.Color) {
	if c.r == nil || !c.r.inited || c.r.prog == 0 || c.fbW <= 0 || c.fbH <= 0 {
		return
	}

	// 将屏幕像素坐标（左上角原点）转换到 NDC（左下角原点）
	x0 := float32((rect.X/float64(c.fbW))*2.0 - 1.0)
	y0 := float32(1.0 - (rect.Y/float64(c.fbH))*2.0)
	x1 := float32(((rect.X+rect.Width)/float64(c.fbW))*2.0 - 1.0)
	y1 := float32(1.0 - ((rect.Y+rect.Height)/float64(c.fbH))*2.0)

	verts := [...]float32{
		x0, y0,
		x1, y0,
		x0, y1,
		x1, y1,
	}

	gl.UseProgram(c.r.prog)
	gl.BindVertexArray(c.r.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, c.r.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(verts)*4, gl.Ptr(&verts[0]))

	r := float32(color.R) / 255.0
	g := float32(color.G) / 255.0
	b := float32(color.B) / 255.0
	a := float32(color.A) / 255.0
	gl.Uniform4f(c.r.colorLoc, r, g, b, a)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

// 顶点/片段着色器源码（最小实现）
var vertexSrc = `#version 330 core
layout(location = 0) in vec2 aPos;
void main(){
	gl_Position = vec4(aPos, 0.0, 1.0);
}` + "\x00"

var fragmentSrc = `#version 330 core
out vec4 FragColor;
uniform vec4 uColor;
void main(){
	FragColor = uColor;
}` + "\x00"

func buildProgram(vsSrc, fsSrc string) (uint32, error) {
	vs := gl.CreateShader(gl.VERTEX_SHADER)
	vsrc, vfree := gl.Strs(vsSrc)
	gl.ShaderSource(vs, 1, vsrc, nil)
	vfree()
	gl.CompileShader(vs)
	if err := checkShader(vs); err != nil {
		return 0, fmt.Errorf("vertex shader: %w", err)
	}

	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	fsrc, ffree := gl.Strs(fsSrc)
	gl.ShaderSource(fs, 1, fsrc, nil)
	ffree()
	gl.CompileShader(fs)
	if err := checkShader(fs); err != nil {
		return 0, fmt.Errorf("fragment shader: %w", err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vs)
	gl.AttachShader(prog, fs)
	gl.LinkProgram(prog)
	gl.DeleteShader(vs)
	gl.DeleteShader(fs)
	if err := checkProgram(prog); err != nil {
		return 0, fmt.Errorf("link program: %w", err)
	}
	return prog, nil
}

func checkShader(s uint32) error {
	var status int32
	gl.GetShaderiv(s, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &logLen)
		log := make([]byte, logLen)
		gl.GetShaderInfoLog(s, logLen, nil, &log[0])
		return fmt.Errorf("compile failed: %s", string(log))
	}
	return nil
}

func checkProgram(p uint32) error {
	var status int32
	gl.GetProgramiv(p, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(p, gl.INFO_LOG_LENGTH, &logLen)
		log := make([]byte, logLen)
		gl.GetProgramInfoLog(p, logLen, nil, &log[0])
		return fmt.Errorf("link failed: %s", string(log))
	}
	return nil
}
