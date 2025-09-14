package canvas

import (
	"log"
	"os"
	"runtime"

	"github.com/millken/goui/ux"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

var _ ux.Renderer = (*CanvasRenderer)(nil)

var defaultFontFamily = "Default"
var defaultFontSize float64 = 16

// var defaultFontColor = glay.NewColorRGBA8(255, 255, 255, 255)

// CanvasRenderer 基于 tdewolff/canvas，将 RenderCommand 渲染到内存画布并可落盘。
// 非并发安全，测试中请在单 goroutine 使用。
type CanvasRenderer struct {
	c   *canvas.Canvas
	ctx *canvas.Context
	w   int
	h   int

	// 字体缓存，避免每次绘制重复构建；未找到字体时禁用文本绘制。
	fontFamily *canvas.FontFamily
	hasFont    bool
}

func NewCanvasRenderer(width, height int) (ux.Renderer, error) {
	c := canvas.New(float64(width), float64(height))
	r := &CanvasRenderer{
		c:   c,
		ctx: canvas.NewContext(c),
		w:   width,
		h:   height,
	}
	// 初始化默认字体，失败时仅关闭文本绘制，避免 panic
	r.initDefaultFont()
	return r, nil
}
func (m *CanvasRenderer) Resize(width, height int) {}

// 初始化默认字体：优先环境变量，其次常见系统路径；失败则禁用文本绘制
func (m *CanvasRenderer) initDefaultFont() {
	m.fontFamily = canvas.NewFontFamily(defaultFontFamily)
	candidates := []string{}
	if p := os.Getenv("CANVAS_DEFAULT_FONT"); p != "" {
		candidates = append(candidates, p)
	}
	switch runtime.GOOS {
	case "linux":
		candidates = append(candidates,
			"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
			"/usr/share/fonts/truetype/freefont/FreeSans.ttf",
		)
	case "darwin":
		candidates = append(candidates,
			"/Library/Fonts/Arial.ttf",
			"/System/Library/Fonts/Supplemental/Arial.ttf",
		)
	case "windows":
		candidates = append(candidates,
			`C:\\Windows\\Fonts\\arial.ttf`,
		)
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			if err := m.fontFamily.LoadFontFile(p, canvas.FontRegular); err == nil {
				m.hasFont = true
				return
			}
		}
	}
	log.Printf("CanvasRenderer: no default font found, text rendering will be skipped")
}

// ProcessRenderCommands 将命令绘制到画布。
func (m *CanvasRenderer) ProcessRenderCommands(commands []ux.RenderCmd) error {
	for _, cmd := range commands {
		switch cmd.Type {
		case ux.CmdFill:
			m.ctx.SetFillColor(cmd.Color)
			x := float64(cmd.Rect.X)
			y := float64(cmd.Rect.Y)
			w := float64(cmd.Rect.Width)
			h := float64(cmd.Rect.Height)
			// Flip Y to use top-left origin (screen coordinates)
			drawY := float64(m.h) - (y + h)
			m.ctx.DrawPath(x, drawY, canvas.Rectangle(w, h))
			m.ctx.Fill()
		case ux.CmdText:
			// 无可用字体时跳过文本绘制以避免 panic
			if !m.hasFont {
				continue
			}
			m.ctx.SetFillColor(cmd.Color)
			x := float64(cmd.Rect.X)
			y := float64(cmd.Rect.Y)
			// w := float64(cmd.Rect.Width)
			h := float64(cmd.Rect.Height)
			// Flip Y to use top-left origin (screen coordinates)
			drawY := float64(m.h) - (y + h)

			face := m.fontFamily.Face(defaultFontSize, canvas.Black, canvas.FontRegular)
			text := canvas.NewTextLine(face, cmd.Text, canvas.Left)
			var offsetX float64

			m.ctx.DrawText(x+offsetX, drawY+h-defaultFontSize, text)
			//m.ctx.DrawPath(x, drawY, canvas.Rectangle(w, h)
		}
	}
	return nil
}
func (m *CanvasRenderer) Cleanup() {}

// Reset 清空记录的渲染命令（重置上下文）。
func (m *CanvasRenderer) Reset() { m.ctx = canvas.NewContext(m.c) }

func (m *CanvasRenderer) WriteToPNG(filename string) error {
	return renderers.Write(filename, m.c)
}
