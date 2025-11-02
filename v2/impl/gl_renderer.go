package impl

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	v2 "github.com/millken/goui/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
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
	// textured drawing program and resources
	texProg uint32
	texVao  uint32
	texVbo  uint32
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
	// 初始化纹理渲染程序
	tprog, err := buildProgram(texVertexSrc, texFragmentSrc)
	if err != nil {
		return err
	}
	r.texProg = tprog
	gl.GenVertexArrays(1, &r.texVao)
	gl.BindVertexArray(r.texVao)
	gl.GenBuffers(1, &r.texVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.texVbo)
	// 每个顶点: pos.xy (2 floats) + tex.xy (2 floats) => 4 floats
	gl.BufferData(gl.ARRAY_BUFFER, 4*4*4, nil, gl.DYNAMIC_DRAW)
	// pos attrib
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	// texcoord attrib
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
	// unbind
	gl.BindVertexArray(0)
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
	if r.texProg != 0 {
		gl.DeleteProgram(r.texProg)
		r.texProg = 0
	}
	if r.vbo != 0 {
		gl.DeleteBuffers(1, &r.vbo)
		r.vbo = 0
	}
	if r.texVbo != 0 {
		gl.DeleteBuffers(1, &r.texVbo)
		r.texVbo = 0
	}
	if r.vao != 0 {
		gl.DeleteVertexArrays(1, &r.vao)
		r.vao = 0
	}
	if r.texVao != 0 {
		gl.DeleteVertexArrays(1, &r.texVao)
		r.texVao = 0
	}
	// 可在此处分离上下文，如需要：glfw.DetachCurrentContext()
	r.inited = false
}

func (r *GLRenderer) BeginFrame(width, height int) (v2.Canvas, error) {
	if !r.inited {
		return nil, fmt.Errorf("renderer not initialized")
	}
	// 使用帧缓冲尺寸（适配 DPI 缩放）
	fbW, fbH := width, height
	if r.win != nil && r.win.win != nil {
		wfw, wfh := r.win.win.GetFramebufferSize()
		if wfw > 0 && wfh > 0 {
			fbW, fbH = wfw, wfh
		}
	}
	// 设置视口并清屏
	gl.Viewport(0, 0, int32(fbW), int32(fbH))
	gl.ClearColor(r.clearR, r.clearG, r.clearB, r.clearA)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	return &glCanvas{r: r, fbW: fbW, fbH: fbH}, nil
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

// (glCanvas 类型在下方有更完整的声明)

// 绘制状态（保持简单）
type drawState struct {
	color v2.Color
	font  v2.Font
	lineW float32
	// 2D 变换（3x3 矩阵）
	tx, ty float32 // 仅支持平移（为简单起见），后续可扩展为完整矩阵
}

type clipRect struct {
	x, y, w, h int32
}

type glCanvas struct {
	r        *GLRenderer
	fbW, fbH int

	// 状态栈
	stack []drawState
	state drawState
	clips []clipRect
}

// helper: 转换像素坐标（屏幕左上）到 NDC
func (c *glCanvas) ndcRect(rect v2.Rect) (x0, y0, x1, y1 float32) {
	if c.r == nil || !c.r.inited || c.r.prog == 0 || c.fbW <= 0 || c.fbH <= 0 {
		return
	}
	x0 = float32((rect.X/float32(c.fbW))*2.0 - 1.0)
	y0 = float32(1.0 - (rect.Y/float32(c.fbH))*2.0)
	x1 = float32(((rect.X+rect.Width)/float32(c.fbW))*2.0 - 1.0)
	y1 = float32(1.0 - ((rect.Y+rect.Height)/float32(c.fbH))*2.0)
	return
}

// Save/Restore
func (c *glCanvas) Save() {
	c.stack = append(c.stack, c.state)
}

func (c *glCanvas) Restore() {
	if len(c.stack) == 0 {
		return
	}
	c.state = c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]
	// 恢复裁剪
	if len(c.clips) > 0 {
		cr := c.clips[len(c.clips)-1]
		gl.Enable(gl.SCISSOR_TEST)
		gl.Scissor(cr.x, cr.y, cr.w, cr.h)
	} else {
		gl.Disable(gl.SCISSOR_TEST)
	}
}

func (c *glCanvas) Translate(tx, ty float32) {
	c.state.tx += float32(tx)
	c.state.ty += float32(ty)
}

func (c *glCanvas) Scale(sx, sy float32) {
	// 简化：不实现缩放，仅记录 TODO
}

func (c *glCanvas) SetColor(color v2.Color) {
	c.state.color = color
}

func (c *glCanvas) SetFont(font v2.Font) {
	c.state.font = font
}

func (c *glCanvas) SetLineWidth(w float64) {
	c.state.lineW = float32(w)
	gl.LineWidth(c.state.lineW)
}

// DrawRect 使用 TRIANGLE_STRIP 渲染矩形（填充）
func (c *glCanvas) DrawRect(rect v2.Rect) {
	if c.r == nil || !c.r.inited || c.r.prog == 0 || c.fbW <= 0 || c.fbH <= 0 {
		return
	}
	// apply translation
	rect.X += float32(c.state.tx)
	rect.Y += float32(c.state.ty)

	x0, y0, x1, y1 := c.ndcRect(rect)
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
	col := c.state.color
	r := float32(col.R) / 255.0
	g := float32(col.G) / 255.0
	b := float32(col.B) / 255.0
	a := float32(col.A) / 255.0
	gl.Uniform4f(c.r.colorLoc, r, g, b, a)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

// DrawLine 使用 GL_LINES
func (c *glCanvas) DrawLine(x0, y0, x1, y1 float32) {
	if c.r == nil || !c.r.inited || c.r.prog == 0 || c.fbW <= 0 || c.fbH <= 0 {
		return
	}
	// apply translation
	x0 += float32(c.state.tx)
	y0 += float32(c.state.ty)
	x1 += float32(c.state.tx)
	y1 += float32(c.state.ty)

	r0 := v2.Rect{X: x0, Y: y0, Width: 1, Height: 1}
	r1 := v2.Rect{X: x1, Y: y1, Width: 1, Height: 1}
	nx0, ny0, _, _ := c.ndcRect(r0)
	nx1, ny1, _, _ := c.ndcRect(r1)
	verts := [...]float32{nx0, ny0, nx1, ny1}
	gl.UseProgram(c.r.prog)
	gl.BindVertexArray(c.r.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, c.r.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(verts)*4, gl.Ptr(&verts[0]))
	col := c.state.color
	r := float32(col.R) / 255.0
	g := float32(col.G) / 255.0
	b := float32(col.B) / 255.0
	a := float32(col.A) / 255.0
	gl.Uniform4f(c.r.colorLoc, r, g, b, a)
	gl.DrawArrays(gl.LINES, 0, 2)
}

// DrawImage/DrawText 是占位，后续可接入字体与纹理支持
func (c *glCanvas) DrawImage(img v2.Image, dest v2.Rect) {
	// TODO: implement image drawing with textures
}

func (c *glCanvas) DrawText(text string, x, y float32) {
	if c.r == nil || !c.r.inited || c.r.prog == 0 || c.fbW <= 0 || c.fbH <= 0 {
		return
	}
	if c.r.win == nil || c.r.win.ctx == nil {
		return
	}
	fm := c.r.win.ctx.GetFontManager()
	// 使用默认字体名 "default" 和字号 16（可改为参数）
	face, err := fm.GetFace("default", 16)
	if err != nil {
		// 无法获得字体，跳过绘制
		fmt.Println("DrawText: GetFace error:", err)
		return
	}

	// 使用 font.Drawer 在 Alpha 图上绘制文本
	d := &font.Drawer{
		Face: face,
	}
	adv := d.MeasureString(text)
	w := int((adv + fixed.Int26_6(63)) >> 6)
	if w <= 0 {
		w = 1
	}
	m := face.Metrics()
	h := int((m.Ascent + m.Descent).Ceil())
	if h <= 0 {
		h = 1
	}

	// debug print
	fmt.Println("DrawText: text=", text, "adv=", adv, "w=", w, "h=", h)
	fmtry := "nil"
	if fm != nil {
		fmtry = "ok"
	}
	fmt.Println("DrawText: fontmanager=", fmtry)
	fmt.Println("DrawText: face metrics:", face.Metrics())

	alpha := image.NewAlpha(image.Rect(0, 0, w, h))
	d.Dst = alpha
	d.Src = image.NewUniform(color.Alpha{A: 255})
	d.Dot = fixed.Point26_6{X: 0, Y: m.Ascent}
	d.DrawString(text)

	// convert alpha to RGBA where RGB=1, A=alpha
	rgba := image.NewRGBA(alpha.Bounds())
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{C: color.RGBA{R: 255, G: 255, B: 255, A: 0}}, image.Point{}, draw.Src)
	// copy alpha channel
	for yy := 0; yy < h; yy++ {
		for xx := 0; xx < w; xx++ {
			a := alpha.AlphaAt(xx, yy).A
			rgba.SetRGBA(xx, yy, color.RGBA{R: 255, G: 255, B: 255, A: a})
		}
	}

	// upload texture
	tex := c.r.uploadRGBAAsTexture(rgba)
	fmt.Println("DrawText: uploaded texture id=", tex, "size=", w, "x", h)
	defer func() { gl.DeleteTextures(1, &tex) }()

	// compute dest rect and NDC
	dest := v2.Rect{X: x, Y: y, Width: float32(w), Height: float32(h)}
	// apply translation
	dest.X += c.state.tx
	dest.Y += c.state.ty
	x0, y0, x1, y1 := c.ndcRect(dest)

	// texture coords
	verts := []float32{
		x0, y0, 0.0, 1.0,
		x1, y0, 1.0, 1.0,
		x0, y1, 0.0, 0.0,
		x1, y1, 1.0, 0.0,
	}

	// use textured program
	gl.UseProgram(c.r.texProg)
	gl.BindVertexArray(c.r.texVao)
	gl.BindBuffer(gl.ARRAY_BUFFER, c.r.texVbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, int(len(verts))*4, gl.Ptr(&verts[0]))

	// set color uniform
	colorLoc := gl.GetUniformLocation(c.r.texProg, gl.Str("uColor\x00"))
	col := c.state.color
	fr := float32(col.R) / 255.0
	fg := float32(col.G) / 255.0
	fb := float32(col.B) / 255.0
	fa := float32(col.A) / 255.0
	gl.Uniform4f(colorLoc, fr, fg, fb, fa)

	// bind texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	texLoc := gl.GetUniformLocation(c.r.texProg, gl.Str("uTex\x00"))
	gl.Uniform1i(texLoc, 0)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)
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

var texVertexSrc = `#version 330 core
layout(location = 0) in vec2 aPos;
layout(location = 1) in vec2 aTex;
out vec2 vTex;
void main(){
	gl_Position = vec4(aPos, 0.0, 1.0);
	vTex = aTex;
}` + "\x00"

var texFragmentSrc = `#version 330 core
in vec2 vTex;
out vec4 FragColor;
uniform sampler2D uTex;
uniform vec4 uColor;
void main(){
	vec4 t = texture(uTex, vTex);
	// assume white glyphs in texture with alpha channel
	FragColor = vec4(uColor.rgb, uColor.a * t.a);
}` + "\x00"

// uploadRGBAAsTexture 上传 RGBA 图像到 GPU 并返回生成的 texture id
func (r *GLRenderer) uploadRGBAAsTexture(img *image.RGBA) uint32 {
	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	w := int32(img.Bounds().Dx())
	h := int32(img.Bounds().Dy())
	if w == 0 || h == 0 {
		return tex
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&img.Pix[0]))
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return tex
}
