package impl

import (
	"github.com/millken/go/plutovg"
	v2 "github.com/millken/goui/v2"
)

var _ v2.Canvas = (*PlutovgCanvas)(nil)

type PlutovgCanvas struct {
	surface *plutovg.Surface
	canvas  *plutovg.Canvas
	lineW   float32
}

func NewPlutovgCanvas(width, height int) *PlutovgCanvas {
	surface := plutovg.SurfaceCreate(width, height)
	return &PlutovgCanvas{
		surface: surface,
		canvas:  plutovg.CanvasCreate(surface),
	}
}

func (m *PlutovgCanvas) Save() {
	m.canvas.Save()
}

func (m *PlutovgCanvas) Restore() {
	m.canvas.Restore()
}

func (m *PlutovgCanvas) Translate(tx, ty float32) {
	m.canvas.Translate(tx, ty)
}

func (m *PlutovgCanvas) Scale(sx, sy float32) {
	m.canvas.Scale(sx, sy)
}

func (m *PlutovgCanvas) SetColor(color v2.Color) {
	m.canvas.SetColor(plutovg.MakeColor(float32(color.R), float32(color.G), float32(color.B), float32(color.A)))
}

func (m *PlutovgCanvas) SetFont(font v2.Font) {
	fontface := plutovg.FontFaceLoadFromFile("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 0)
	if fontface == nil {
		panic("Font not found")
	}
	m.canvas.SetFont(fontface, 16)
}

func (m *PlutovgCanvas) DrawRect(rect v2.Rect) {
	m.canvas.Rect(rect.X, rect.Y, rect.Width, rect.Height)
	if m.lineW > 0 {
		m.canvas.SetLineWidth(m.lineW)
		m.canvas.Stroke()
	} else {
		m.canvas.Fill()
	}
}

func (m *PlutovgCanvas) DrawImage(img v2.Image, dest v2.Rect) {
	// best-effort: if img implements a helper to get raw bytes or io.Reader, use it
	// For now, tdewolff/canvas supports embedding images via Image object; leave as TODO
	// Placeholder: draw an empty rect to indicate image area
	// m.ctx.DrawPath(dest.X, dest.Y, canvas.Rectangle(dest.Width, dest.Height))
	// m.ctx.Stroke()
}

func (m *PlutovgCanvas) DrawText(textStr string, x, y float32) {
	m.canvas.FillText(textStr, plutovg.TextEncoding(0), x, y)
}

func (m *PlutovgCanvas) DrawLine(x0, y0, x1, y1 float32) {
	m.canvas.MoveTo(x0, y0)
	m.canvas.LineTo(x1, y1)
	m.canvas.Stroke()
}
