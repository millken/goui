package component

import v2 "github.com/millken/goui/v2"

// LabelComponent 是一个简单的只渲染文本的组件
type LabelComponent struct {
	text  string
	rect  v2.Rect
	color v2.Color
	shown bool
}

func NewLabel(text string) *LabelComponent {
	return &LabelComponent{text: text, color: v2.Color{R: 0, G: 0, B: 0, A: 255}, shown: true}
}

func (l *LabelComponent) Paint(c v2.Canvas) {
	if !l.shown {
		return
	}
	c.Save()
	c.SetColor(l.color)
	// draw text at top-left of rect
	c.DrawText(l.text, l.rect.X, l.rect.Y)
	c.Restore()
}

func (l *LabelComponent) HandleEvent(e v2.Event) bool {
	// label 不主动处理事件
	return false
}

func (l *LabelComponent) SetText(s string)    { l.text = s }
func (l *LabelComponent) SetColor(c v2.Color) { l.color = c }
func (l *LabelComponent) SetBounds(x, y, w, h int) {
	l.rect = v2.Rect{X: float32(x), Y: float32(y), Width: float32(w), Height: float32(h)}
}
func (l *LabelComponent) IsVisible() bool   { return l.shown }
func (l *LabelComponent) SetVisible(v bool) { l.shown = v }

// SetPosition sets the label position (implements v2.Component).
func (l *LabelComponent) SetPosition(p v2.Point) { l.rect.X = p.X; l.rect.Y = p.Y }

// SetSize sets the label size (implements v2.Component).
func (l *LabelComponent) SetSize(s v2.Size) { l.rect.Width = s.Width; l.rect.Height = s.Height }

// GetBounds returns the label bounds (implements v2.Component).
func (l *LabelComponent) GetBounds() v2.Rect { return l.rect }

// AddChild is a no-op for Label (implements v2.Component).
func (l *LabelComponent) AddChild(child v2.Component) { /* no-op */ }

// RemoveChild is a no-op for Label (implements v2.Component).
func (l *LabelComponent) RemoveChild(id string) { /* no-op */ }

// GetChild always returns nil for Label (implements v2.Component).
func (l *LabelComponent) GetChild(id string) v2.Component { return nil }
