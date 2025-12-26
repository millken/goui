package component

import (
	"sync"
	"time"

	gouicolor "github.com/millken/goui/color"
	"github.com/millken/goui/event"
	gouifont "github.com/millken/goui/font"
	gouiwindow "github.com/millken/goui/window"
	"github.com/millken/thorvg"
)

// CheckState represents the checkbox state (similar to unison's check.Enum)
type CheckState int

const (
	CheckOff CheckState = iota
	CheckOn
	CheckMixed
)

// String returns the string representation of the state
func (s CheckState) String() string {
	switch s {
	case CheckOn:
		return "On"
	case CheckMixed:
		return "Mixed"
	default:
		return "Off"
	}
}

// FromBool converts a bool to CheckState
func CheckStateFromBool(b bool) CheckState {
	if b {
		return CheckOn
	}
	return CheckOff
}

// ToBool converts CheckState to bool (Mixed returns false)
func (s CheckState) ToBool() bool {
	return s == CheckOn
}

// CheckboxTheme holds theming data for a checkbox (inspired by unison)
type CheckboxTheme struct {
	EdgeInk        gouicolor.Color
	SelectionInk   gouicolor.Color
	OnSelectionInk gouicolor.Color
	ControlInk     gouicolor.Color
	OnControlInk   gouicolor.Color
	DisabledInk    gouicolor.Color

	FontSize   float64
	FontFamily string
	BoxSize    float32
	Gap        float32
	CornerRadius float32
	BorderWidth  float32

	ClickAnimationTime time.Duration
}

// DefaultCheckboxTheme provides the default theme for checkboxes
var DefaultCheckboxTheme = CheckboxTheme{
	EdgeInk:        gouicolor.NewColor(150, 150, 150, 255),
	SelectionInk:   gouicolor.NewColor(33, 150, 243, 255),
	OnSelectionInk: gouicolor.NewColor(255, 255, 255, 255),
	ControlInk:     gouicolor.NewColor(255, 255, 255, 255),
	OnControlInk:   gouicolor.NewColor(33, 150, 243, 255),
	DisabledInk:    gouicolor.NewColor(200, 200, 200, 255),
	FontSize:       14.0,
	FontFamily:     "OPPO Sans 4.0",
	BoxSize:        16,
	Gap:            8,
	CornerRadius:   3,
	BorderWidth:    1.5,
	ClickAnimationTime: 100 * time.Millisecond,
}

// CheckBox represents a clickable checkbox with optional label (refactored following unison patterns)
type CheckBox struct {
	*BaseComponent
	theme    CheckboxTheme
	label    string
	state    CheckState
	disabled bool
	pressed  bool
	hovered  bool
	mu       sync.RWMutex

	// Callbacks
	ClickCallback   func()
	ChangeCallback  func(state CheckState)
}

// NewCheckBox creates a new checkbox (following unison naming convention)
func NewCheckBox(id ...string) *CheckBox {
	cb := &CheckBox{
		BaseComponent: NewBaseComponent(id...),
		theme:         DefaultCheckboxTheme,
		state:         CheckOff,
		disabled:      false,
		pressed:       false,
		hovered:       false,
	}

	// Set initial size
	cb.updateSize()

	// Register mouse event handlers
	cb.RegisterEventListener(event.EventMouse, event.EventListenerFunc(func(e event.Event) bool {
		if cb.disabled {
			return false
		}

		me, ok := e.(*event.MouseEvent)
		if !ok {
			return false
		}

		r := cb.Bounds()
		inside := me.X >= float32(r.Min.X) && me.X < float32(r.Max.X) &&
				 me.Y >= float32(r.Min.Y) && me.Y < float32(r.Max.Y)

		// Mouse down - start press (unison's DefaultMouseDown)
		if me.Down && me.Count == 1 {
			if inside {
				cb.setPressed(true)
				return true
			}
			return false
		}

		// Mouse drag - update press state based on position (unison's DefaultMouseDrag)
		if !me.Down && me.Count == 0 {
			wasPressed := cb.pressed
			if wasPressed {
				cb.setPressed(inside)
			}
			return wasPressed
		}

		// Mouse up - complete click (unison's DefaultMouseUp)
		if !me.Down && me.Count == 1 {
			wasPressed := cb.pressed
			cb.setPressed(false)
			if wasPressed && inside {
				cb.Click()
				return true
			}
			return false
		}

		return false
	}))

	// Handle hover for cursor changes (unison's DefaultUpdateCursor)
	cb.OnHover(func(hover bool) {
		if cb.disabled {
			return
		}
		cb.mu.Lock()
		cb.hovered = hover
		cb.mu.Unlock()
		if rc, ok := cb.Root().(RootComponent); ok && rc != nil {
			if w := rc.GetWindow(); w != nil {
				if hover {
					w.SetCursor(gouiwindow.CursorPointer)
				} else {
					w.SetCursor(gouiwindow.CursorDefault)
				}
			}
		}
		cb.MarkDirty(false)
	})

	return cb
}

// SetTheme sets the theme for the checkbox
func (cb *CheckBox) SetTheme(theme CheckboxTheme) {
	cb.mu.Lock()
	cb.theme = theme
	cb.mu.Unlock()
	cb.updateSize()
	cb.MarkDirty(false)
}

// Theme returns the current theme
func (cb *CheckBox) Theme() CheckboxTheme {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.theme
}

// SetState sets the checkbox state
func (cb *CheckBox) SetState(state CheckState) {
	cb.mu.Lock()
	if cb.state != state {
		cb.state = state
		cb.mu.Unlock()
		cb.MarkDirty(false)
		if cb.ChangeCallback != nil {
			cb.ChangeCallback(state)
		}
	} else {
		cb.mu.Unlock()
	}
}

// State returns the current checkbox state
func (cb *CheckBox) State() CheckState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Checked returns true if the checkbox is in the On state
func (cb *CheckBox) Checked() bool {
	return cb.State().ToBool()
}

// SetChecked sets the checked state (converts bool to CheckState)
func (cb *CheckBox) SetChecked(checked bool) {
	cb.SetState(CheckStateFromBool(checked))
}

// SetTitle sets the label text (following unison naming)
func (cb *CheckBox) SetTitle(text string) {
	cb.mu.Lock()
	if cb.label != text {
		cb.label = text
		cb.mu.Unlock()
		cb.updateSize()
		cb.MarkDirty(false)
	} else {
		cb.mu.Unlock()
	}
}

// Title returns the label text
func (cb *CheckBox) Title() string {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.label
}

// SetLabel is an alias for SetTitle for backward compatibility
func (cb *CheckBox) SetLabel(label string) {
	cb.SetTitle(label)
}

// GetLabel is an alias for Title for backward compatibility
func (cb *CheckBox) GetLabel() string {
	return cb.Title()
}

// SetDisabled sets the disabled state
func (cb *CheckBox) SetDisabled(disabled bool) {
	cb.mu.Lock()
	if cb.disabled != disabled {
		cb.disabled = disabled
		cb.mu.Unlock()
		cb.MarkDirty(false)
	} else {
		cb.mu.Unlock()
	}
}

// IsDisabled returns true if the checkbox is disabled
func (cb *CheckBox) IsDisabled() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.disabled
}

// Click simulates a click on the checkbox (following unison pattern)
func (cb *CheckBox) Click() {
	cb.mu.Lock()
	if cb.disabled {
		cb.mu.Unlock()
		return
	}
	cb.mu.Unlock()

	// Toggle state (unison's updateState)
	cb.updateState()

	// Trigger press animation
	wasDisabled := cb.disabled
	cb.mu.Lock()
	cb.pressed = true
	cb.mu.Unlock()
	cb.MarkDirty(false)

	// Call click callback
	if cb.ClickCallback != nil {
		cb.ClickCallback()
	}

	// Reset pressed state after animation
	// Note: In a real implementation, you'd use a timer for proper animation
	cb.mu.Lock()
	cb.pressed = wasDisabled
	cb.mu.Unlock()
	cb.MarkDirty(false)
}

func (cb *CheckBox) updateState() {
	cb.mu.Lock()
	currentState := cb.state
	if currentState == CheckOn {
		cb.state = CheckOff
	} else {
		cb.state = CheckOn
	}
	newState := cb.state
	cb.mu.Unlock()

	cb.MarkDirty(false)
	if cb.ChangeCallback != nil {
		cb.ChangeCallback(newState)
	}
}

func (cb *CheckBox) setPressed(pressed bool) {
	cb.mu.Lock()
	if cb.pressed != pressed {
		cb.pressed = pressed
		cb.mu.Unlock()
		cb.MarkDirty(false)
	} else {
		cb.mu.Unlock()
	}
}

func (cb *CheckBox) boxSize() float32 {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.theme.BoxSize
}

func (cb *CheckBox) updateSize() {
	cb.mu.RLock()
	theme := cb.theme
	label := cb.label
	cb.mu.RUnlock()

	boxSize := theme.BoxSize

	if label == "" {
		cb.SetWidth(boxSize)
		cb.SetHeight(boxSize)
		return
	}

	tm := gouifont.MeasureText(label, gouifont.WithFontSize(theme.FontSize))
	width := boxSize + theme.Gap + float32(tm.Width)
	height := boxSize
	if float32(tm.Height) > height {
		height = float32(tm.Height)
	}

	cb.SetWidth(width)
	cb.SetHeight(height)
}

func (cb *CheckBox) Paint(vg *thorvg.Canvas) bool {
	if !cb.IsDirty() {
		return cb.BaseComponent.Paint(vg)
	}

	cb.mu.RLock()
	theme := cb.theme
	label := cb.label
	state := cb.state
	disabled := cb.disabled
	pressed := cb.pressed
	hovered := cb.hovered
	cb.mu.RUnlock()

	r := cb.Bounds()
	x, y := float32(r.Min.X), float32(r.Min.Y)
	boxSize := cb.boxSize()

	// Calculate total width for hover effect
	totalWidth := boxSize
	if label != "" {
		tm := gouifont.MeasureText(label, gouifont.WithFontSize(theme.FontSize))
		totalWidth = boxSize + theme.Gap + float32(tm.Width)
	}

	// Draw hover background highlight
	if hovered && !disabled && !pressed {
		hoverBg := thorvg.NewShape()
		if hoverBg != nil {
			hr, hg, hb, _ := theme.SelectionInk.RGBA()
			hoverBg.AppendRect(x-2, y-2, totalWidth+4, boxSize+4, 4, 4, false)
			hoverBg.SetFillColor(hr, hg, hb, 20)
			vg.Push(&hoverBg.Paint)
			hoverBg.Release()
		}
	}

	// Draw checkbox background (unison's DrawRoundedRectBase)
	bg := thorvg.NewShape()
	if bg != nil {
		var bgColor gouicolor.Color
		if disabled {
			bgColor = theme.DisabledInk
		} else if pressed {
			bgColor = theme.SelectionInk
		} else {
			bgColor = theme.ControlInk
		}

		br, gg, bb, aa := bgColor.RGBA()
		bg.AppendRect(x, y, boxSize, boxSize, theme.CornerRadius, theme.CornerRadius, false)
		bg.SetFillColor(br, gg, bb, aa)

		// Draw border
		var borderCol gouicolor.Color
		if disabled {
			borderCol = theme.DisabledInk
		} else if hovered {
			borderCol = theme.SelectionInk
		} else {
			borderCol = theme.EdgeInk
		}

		brdr, ggdr, bbdr, _ := borderCol.RGBA()
		bg.SetStrokeColor(brdr, ggdr, bbdr, 255)
		bg.SetStrokeWidth(theme.BorderWidth)

		vg.Push(&bg.Paint)
		bg.Release()
	}

	// Draw check/mixed indicator (unison's path drawing)
	if state != CheckOff && !disabled {
		checkmark := thorvg.NewShape()
		if checkmark != nil {
			var fgColor gouicolor.Color
			if pressed {
				fgColor = theme.OnSelectionInk
			} else {
				fgColor = theme.OnControlInk
			}

			cr, cg, cb_, _ := fgColor.RGBA()
			checkmark.SetStrokeColor(cr, cg, cb_, 255)
			checkmark.SetStrokeWidth(2.0)
			checkmark.SetStrokeCap(thorvg.StrokeCapRound)
			checkmark.SetStrokeJoin(thorvg.StrokeJoinRound)

			if state == CheckOn {
				// Draw checkmark (unison's check path - centered)
				// Using proportions similar to unison: (0.25, 0.55) -> (0.45, 0.7) -> (0.75, 0.3)
				checkmark.MoveTo(x + boxSize*0.25, y + boxSize*0.55)
				checkmark.LineTo(x + boxSize*0.45, y + boxSize*0.70)
				checkmark.LineTo(x + boxSize*0.75, y + boxSize*0.30)
			} else if state == CheckMixed {
				// Draw horizontal line for mixed state (unison's mixed state - centered)
				lineY := y + boxSize*0.5
				checkmark.MoveTo(x + boxSize*0.25, lineY)
				checkmark.LineTo(x + boxSize*0.70, lineY)
			}

			vg.Push(&checkmark.Paint)
			checkmark.Release()
		}
	}

	// Draw label (vertically centered with checkbox)
	if label != "" {
		text := thorvg.NewText()
		if text != nil {
			text.SetText(label)
			text.SetSize(float32(theme.FontSize))

			var labelCol gouicolor.Color
			if disabled {
				labelCol = theme.DisabledInk
			} else {
				labelCol = theme.EdgeInk
			}

			lr, lg, lb, _ := labelCol.RGBA()
			text.SetColor(lr, lg, lb)

			cs := cb.ComputedStyle()
			if cs.Font != nil && *cs.Font != "" {
				text.SetFont(*cs.Font)
			} else {
				text.SetFont(theme.FontFamily)
			}

			// Vertically center label text with checkbox
			// 使用与 Button 相同的方式：以文本的高度为准，将文本块与复选框垂直居中
			tm := gouifont.MeasureText(label, gouifont.WithFontSize(theme.FontSize))
			textX := x + boxSize + theme.Gap
			textY := y + (boxSize - float32(tm.Height)) / 2

			text.Translate(textX, textY)
			vg.Push(&text.Paint)
			text.Release()
		}
	}

	cb.ClearDirty()
	childRendered := cb.BaseComponent.Paint(vg)
	return true || childRendered
}

// Checkbox is an alias for CheckBox for backward compatibility
type Checkbox = CheckBox

// NewCheckbox is an alias for NewCheckBox for backward compatibility
func NewCheckbox(id ...string) *CheckBox {
	return NewCheckBox(id...)
}