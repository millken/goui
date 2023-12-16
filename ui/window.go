package ui

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	redrawSet = sync.Map{}
)

func LoadRedrawWindows() []*Window {
	var set []*Window
	redrawSet.Range(func(key, value interface{}) bool {
		set = append(set, key.(*Window))
		return true
	})
	return set
}

func MarkWindowForRedraw(w *Window) {
	if _, exists := redrawSet.Load(w); !exists {
		redrawSet.Store(w, struct{}{})
	}
}

func DeleteWindowFromRedraw(w *Window) {
	redrawSet.Delete(w)
}

// WindowOption holds an option for window creation.
type WindowOption func(*Window) error
type Window struct {
	title            string
	width, height    int32
	wnd              *sdl.Window
	surface          *surface
	layout           *Layout
	valid            bool
	lastDrawDuration time.Duration
}

func NewWindow(opts ...WindowOption) (*Window, error) {
	w := &Window{
		wnd:     nil,
		surface: &surface{},
	}
	for _, opt := range opts {
		if err := opt(w); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func (w *Window) SetTitle(title string) {
	w.title = title
}

func (w *Window) SetLayout(layout *Layout) {
	w.layout = layout
}

func (w *Window) SetSize(width, height int32) {
	w.width = width
	w.height = height
}

func (w *Window) Size() (width, height int32) {
	if w.width == 0 || w.height == 0 {
		return 800, 600
	}
	return w.width, w.height
}

func (w *Window) Show() {
	var err error
	if w.wnd == nil {
		width, height := w.Size()
		w.wnd, err = sdl.CreateWindow(w.title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			width, height, windowFlags)
		if err != nil {
			panic(err)
		}
	}
	w.valid = true
	w.ReDraw()

}

func (w *Window) Hide() {
	//w.wnd.Hide()
}

func (w *Window) OnClose(callback func(context.Context)) {
	//w.wnd.OnClose(callback)
}

func (w *Window) Destroy() {
	//w.wnd.Destroy()
}

func (w *Window) Draw(c *Canvas) {
	paint := BackgroundColor.Paint(c, w.LocalContentRect(), Stroke)
	for i := 0; i < 100; i++ {
		paint.SetColor(RGB(rand.Intn(255), rand.Intn(255), rand.Intn(255)))
		c.DrawRect(Rect{Point: Point{X: float32(rand.Intn(800)), Y: float32(rand.Intn(600))}, Size: Size{Width: float32(100), Height: float32(100)}}, paint)
		c.Save()
	}
	//c.DrawPaint(paint)
	//c.Save()
}

func (w *Window) PerformLayout() {
	//w.wnd.PerformLayout()
}

// BackingScale returns the scale of the backing store for this window.
func (w *Window) BackingScale() (x, y float32) {
	displayIndex, _ := w.wnd.GetDisplayIndex()
	ddpi, _, _, _ := sdl.GetDisplayDPI(displayIndex)
	fmt.Println("ddpi:", ddpi, "scale: ", float32(ddpi)/96)
	return float32(ddpi) / 96, float32(ddpi) / 96
	//return 1, 1
}

// IsVisible returns true if the window is currently being shown.
func (w *Window) IsVisible() bool {
	if w.IsValid() {
		return true
	}
	return false
}

// IsValid returns true if the window is still valid (i.e. hasn't been disposed).
func (w *Window) IsValid() bool {
	return w != nil && w.valid
}

// ContentRect returns the boundaries in display coordinates of the window's content area.
func (w *Window) ContentRect() Rect {
	if w.IsValid() {
		x, y := w.wnd.GetPosition()
		width, height := w.wnd.GetSize()
		return Rect{
			Point: Point{X: float32(x), Y: float32(y)},
			Size:  Size{Width: float32(width), Height: float32(height)},
		}
	}
	return Rect{Size: Size{Width: 1, Height: 1}}
}

// LocalContentRect returns the boundaries in local coordinates of the window's content area.
func (w *Window) LocalContentRect() Rect {
	r := w.ContentRect()
	r.Point.X = 0
	r.Point.Y = 0
	return r
}

// ReDraw redraws the window.
func (w *Window) ReDraw() {
	if w.IsValid() {
		sx, sy := w.BackingScale()
		ctx, err := w.wnd.GLCreateContext()
		if err != nil {
			slog.Error(err.Error())
		}
		if err := w.wnd.GLMakeCurrent(ctx); err != nil {
			slog.Error(err.Error())
		}
		c, err := w.surface.prepareCanvas(w.ContentRect().Size, w.LocalContentRect(), sx, sy)
		if err != nil {
			slog.Error(err.Error(), "size", w.ContentRect().Size, "rect", w.LocalContentRect(), "scaleX", sx, "scaleY", sy)
			return
		}
		start := time.Now()
		c.Save()
		w.Draw(c)
		c.Restore()
		c.Flush()
		w.lastDrawDuration = time.Since(start)
		w.wnd.GLSwap()
	}
}

// MarkForRedraw marks this window for drawing at the next update.
func (w *Window) MarkForRedraw() {
	MarkWindowForRedraw(w)
}
