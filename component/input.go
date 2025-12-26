package component

import (
	"sync"
	"time"
	"unicode/utf8"

	gouicolor "github.com/millken/goui/color"
	"github.com/millken/goui/event"
	gouifont "github.com/millken/goui/font"
	gouiwindow "github.com/millken/goui/window"
	"github.com/millken/thorvg"
)

// InputTheme holds theming data for an input field
type InputTheme struct {
	BackgroundColor  gouicolor.Color
	BorderColor      gouicolor.Color
	FocusColor       gouicolor.Color
	TextColor        gouicolor.Color
	PlaceholderColor gouicolor.Color
	DisabledColor    gouicolor.Color

	FontSize     float64
	FontFamily   string
	PaddingX     float32
	PaddingY     float32
	CornerRadius float32
	BorderWidth  float32

	CursorBlinkInterval time.Duration
}

// DefaultInputTheme provides default theme for input fields
var DefaultInputTheme = InputTheme{
	BackgroundColor:     gouicolor.NewColor(255, 255, 255, 255),
	BorderColor:         gouicolor.NewColor(200, 200, 200, 255),
	FocusColor:          gouicolor.NewColor(33, 150, 243, 255),
	TextColor:           gouicolor.NewColor(30, 30, 30, 255),
	PlaceholderColor:    gouicolor.NewColor(180, 180, 180, 255),
	DisabledColor:       gouicolor.NewColor(235, 235, 235, 255),
	FontSize:            14.0,
	FontFamily:          "OPPO Sans 4.0",
	PaddingX:            12,
	PaddingY:            8,
	CornerRadius:        4,
	BorderWidth:         1.5,
	CursorBlinkInterval: 500 * time.Millisecond,
}

// Input represents a text input field
type Input struct {
	*BaseComponent
	theme       InputTheme
	text        string
	placeholder string
	focused     bool
	disabled    bool
	mu          sync.RWMutex

	// Cursor
	cursorPosition int // rune position, not byte position
	cursorVisible  bool
	lastBlinkTime  time.Time
	scrollOffset   float32 // horizontal scroll offset in pixels
}

// NewInput creates a new input field
func NewInput(id ...string) *Input {
	input := &Input{
		BaseComponent:  NewBaseComponent(id...),
		theme:          DefaultInputTheme,
		text:           "",
		placeholder:    "",
		focused:        false,
		disabled:       false,
		cursorPosition: 0,
		cursorVisible:  true,
		lastBlinkTime:  time.Now(),
		scrollOffset:   0,
	}

	// Set default size
	input.SetWidth(200)
	input.SetHeight(40)

	// Register mouse event handler for focus
	input.RegisterEventListener(event.EventMouse, event.EventListenerFunc(func(e event.Event) bool {
		if input.disabled {
			return false
		}

		me, ok := e.(*event.MouseEvent)
		if !ok {
			return false
		}

		r := input.Bounds()
		inside := me.X >= float32(r.Min.X) && me.X < float32(r.Max.X) &&
			me.Y >= float32(r.Min.Y) && me.Y < float32(r.Max.Y)

		// On mouse down inside input, request focus
		if me.Down && me.Count == 1 && inside {
			input.mu.Lock()
			wasFocused := input.focused
			input.focused = true
			input.cursorVisible = true
			input.lastBlinkTime = time.Now()
			input.mu.Unlock()
			if !wasFocused {
				input.MarkDirty(false)
			}
			return true
		}

		// On mouse click outside, remove focus
		if me.Down && me.Count == 1 && !inside && input.focused {
			input.mu.Lock()
			input.focused = false
			input.cursorVisible = false
			input.mu.Unlock()
			input.MarkDirty(false)
		}

		return false
	}))

	// Register key event handler for text editing
	input.RegisterEventListener(event.EventKey, event.EventListenerFunc(func(e event.Event) bool {
		if !input.focused || input.disabled {
			return false
		}

		ke, ok := e.(*event.KeyEvent)
		if !ok {
			return false
		}

		// Handle special keys
		if ke.Down {
			switch ke.Key {
			case event.KeyCodeDeleteBackspace:
				input.handleBackspace()
			case event.KeyCodeReturnEnter:
				// Handle enter (could trigger callback in future)
				return true
			case event.KeyCodeLeftArrow:
				input.moveCursor(-1)
			case event.KeyCodeRightArrow:
				input.moveCursor(1)
			case event.KeyCodeHome:
				input.mu.Lock()
				input.cursorPosition = 0
				input.mu.Unlock()
				input.MarkDirty(false)
			case event.KeyCodeEnd:
				input.mu.Lock()
				input.cursorPosition = utf8.RuneCountInString(input.text)
				input.mu.Unlock()
				input.MarkDirty(false)
			case event.KeyCodeEscape:
				input.blur()
			default:
				// Let CharEvent handle printable characters
				return false
			}
			return true
		}

		return false
	}))

	// Register char event handler for text input
	input.RegisterEventListener(event.EventChar, event.EventListenerFunc(func(e event.Event) bool {
		if !input.focused || input.disabled {
			return false
		}

		ce, ok := e.(*event.CharEvent)
		if !ok {
			return false
		}

		// Handle character input directly here
		input.mu.Lock()
		defer input.mu.Unlock()

		// Filter out control characters only, allow all printable chars including Chinese
		if ce.Rune < 32 {
			return false
		}

		// Insert character at cursor position using rune slice for UTF-8 support
		runes := []rune(input.text)
		input.text = string(runes[:input.cursorPosition]) + string(ce.Rune) + string(runes[input.cursorPosition:])
		input.cursorPosition++

		// Mark dirty to trigger redraw
		input.MarkDirty(false)
		return true
	}))

	// Register focus event handler
	input.RegisterEventListener(event.EventFocus, event.EventListenerFunc(func(e event.Event) bool {
		if input.disabled {
			return false
		}

		fe, ok := e.(*event.FocusEvent)
		if !ok {
			return false
		}

		input.mu.Lock()
		input.focused = fe.Gained
		if input.focused {
			input.cursorVisible = true
			input.lastBlinkTime = time.Now()
		}
		input.mu.Unlock()
		input.MarkDirty(false)

		// Change cursor based on focus state
		if rc, ok := input.Root().(RootComponent); ok && rc != nil {
			if w := rc.GetWindow(); w != nil {
				if input.focused {
					w.SetCursor(gouiwindow.CursorText)
				} else {
					w.SetCursor(gouiwindow.CursorDefault)
				}
			}
		}

		return false
	}))
	return input
}

// SetTheme sets theme for input
func (inp *Input) SetTheme(theme InputTheme) {
	inp.mu.Lock()
	inp.theme = theme
	inp.mu.Unlock()
	inp.MarkDirty(false)
}

// Theme returns current theme
func (inp *Input) Theme() InputTheme {
	inp.mu.RLock()
	defer inp.mu.RUnlock()
	return inp.theme
}

// SetText sets input text
func (inp *Input) SetText(text string) {
	inp.mu.Lock()
	inp.text = text
	inp.cursorPosition = utf8.RuneCountInString(text)
	inp.mu.Unlock()
	inp.MarkDirty(false)
}

// Text returns current text
func (inp *Input) Text() string {
	inp.mu.RLock()
	defer inp.mu.RUnlock()
	return inp.text
}

// SetPlaceholder sets placeholder text
func (inp *Input) SetPlaceholder(placeholder string) {
	inp.mu.Lock()
	inp.placeholder = placeholder
	inp.mu.Unlock()
	inp.MarkDirty(false)
}

// Placeholder returns placeholder text
func (inp *Input) Placeholder() string {
	inp.mu.RLock()
	defer inp.mu.RUnlock()
	return inp.placeholder
}

// SetDisabled sets disabled state
func (inp *Input) SetDisabled(disabled bool) {
	inp.mu.Lock()
	if inp.disabled != disabled {
		inp.disabled = disabled
		if disabled {
			inp.focused = false
		}
		inp.mu.Unlock()
		inp.MarkDirty(false)
	} else {
		inp.mu.Unlock()
	}
}

// IsDisabled returns true if input is disabled
func (inp *Input) IsDisabled() bool {
	inp.mu.RLock()
	defer inp.mu.RUnlock()
	return inp.disabled
}

// Focus sets focus to input
func (inp *Input) Focus() {
	inp.requestFocus()
}

// focusRequest requests focus for input
func (inp *Input) requestFocus() {
	if inp.disabled {
		return
	}
	// Directly set focus state for now
	// In a more complete implementation, this would be managed by a focus manager
	inp.mu.Lock()
	inp.focused = true
	inp.cursorVisible = true
	inp.lastBlinkTime = time.Now()
	inp.mu.Unlock()
	inp.MarkDirty(false)
}

// blur removes focus from input
func (inp *Input) blur() {
	// Directly remove focus for now
	inp.mu.Lock()
	inp.focused = false
	inp.cursorVisible = false
	inp.mu.Unlock()
	inp.MarkDirty(false)
}

// handleChar handles character input (deprecated - logic moved inline)
func (inp *Input) handleChar(r rune) {
	// This is now handled directly in event listener
}

// handleBackspace handles backspace key
func (inp *Input) handleBackspace() {
	inp.mu.Lock()
	defer inp.mu.Unlock()

	if inp.cursorPosition > 0 {
		runes := []rune(inp.text)
		inp.text = string(runes[:inp.cursorPosition-1]) + string(runes[inp.cursorPosition:])
		inp.cursorPosition--
		// Mark dirty to trigger redraw
		inp.MarkDirty(false)
	}
}

// moveCursor moves cursor by given offset
func (inp *Input) moveCursor(offset int) {
	inp.mu.Lock()
	defer inp.mu.Unlock()

	newPos := inp.cursorPosition + offset
	if newPos < 0 {
		newPos = 0
	} else if newPos > utf8.RuneCountInString(inp.text) {
		newPos = utf8.RuneCountInString(inp.text)
	}
	inp.cursorPosition = newPos
	inp.MarkDirty(false)
}

// Clear clears input text
func (inp *Input) Clear() {
	inp.SetText("")
}

// Paint draws input field
func (inp *Input) Paint(vg *thorvg.Canvas) bool {
	inp.mu.RLock()
	theme := inp.theme
	text := inp.text
	placeholder := inp.placeholder
	focused := inp.focused
	disabled := inp.disabled
	cursorPos := inp.cursorPosition
	cursorVis := inp.cursorVisible
	inp.mu.RUnlock()

	// If focused, always redraw to show cursor blinking
	// Otherwise, only redraw if dirty
	if !focused && !inp.IsDirty() {
		return inp.BaseComponent.Paint(vg)
	}

	// Update cursor blink state when focused
	if focused {
		now := time.Now()
		inp.mu.Lock()
		if now.Sub(inp.lastBlinkTime) >= theme.CursorBlinkInterval {
			inp.cursorVisible = !inp.cursorVisible
			inp.lastBlinkTime = now
			cursorVis = inp.cursorVisible
		}
		inp.mu.Unlock()
		// Always request redraw when focused
		inp.BaseComponent.MarkDirty(false)
	}

	// Clear dirty flag if not focused (focused state needs continuous redraw)
	if !focused {
		inp.ClearDirty()
	}

	r := inp.Bounds()
	x, y := float32(r.Min.X), float32(r.Min.Y)

	// Calculate content area
	contentX := x + theme.PaddingX
	contentY := y + theme.PaddingY
	contentHeight := float32(r.Max.Y-r.Min.Y) - theme.PaddingY*2
	contentWidth := float32(r.Max.X-r.Min.X) - theme.PaddingX*2

	// Calculate scroll offset to keep cursor visible
	var scrollOffset float32
	if focused {
		inp.mu.Lock()
		// Measure text before cursor
		runes := []rune(inp.text)
		cursorText := string(runes[:cursorPos])
		tm := gouifont.MeasureText(cursorText, gouifont.WithFontSize(theme.FontSize))
		cursorPixelPos := float32(tm.Width)

		// Scroll if cursor is near right edge
		if cursorPixelPos-inp.scrollOffset > contentWidth-20 {
			inp.scrollOffset = cursorPixelPos - (contentWidth - 20)
		}
		// Scroll left if cursor is past left edge
		if cursorPixelPos-inp.scrollOffset < 20 {
			inp.scrollOffset = cursorPixelPos - 20
		}
		// Don't scroll past beginning
		if inp.scrollOffset < 0 {
			inp.scrollOffset = 0
		}
		scrollOffset = inp.scrollOffset

		// Calculate total text width
		totalTextWidth := gouifont.MeasureText(inp.text, gouifont.WithFontSize(theme.FontSize)).Width

		// If total text fits, no scroll
		if float32(totalTextWidth) <= contentWidth {
			scrollOffset = 0
			inp.scrollOffset = 0
		}
		inp.mu.Unlock()
	} else {
		inp.mu.RLock()
		scrollOffset = inp.scrollOffset
		inp.mu.RUnlock()
	}

	// Determine colors
	var bgColor, borderColor, textColor gouicolor.Color
	if disabled {
		bgColor = theme.DisabledColor
		borderColor = theme.DisabledColor
		textColor = theme.DisabledColor
	} else if focused {
		bgColor = theme.BackgroundColor
		borderColor = theme.FocusColor
		textColor = theme.TextColor
	} else {
		bgColor = theme.BackgroundColor
		borderColor = theme.BorderColor
		textColor = theme.TextColor
	}

	// Force text color to always be visible
	if text == "" && placeholder != "" {
		textColor = theme.PlaceholderColor
	}

	// Draw background
	bg := thorvg.NewShape()
	if bg != nil {
		br, gg, bb, aa := bgColor.RGBA()
		bg.AppendRect(x, y, float32(r.Max.X-r.Min.X), float32(r.Max.Y-r.Min.Y), theme.CornerRadius, theme.CornerRadius, false)
		bg.SetFillColor(br, gg, bb, aa)
		vg.Push(&bg.Paint)
		bg.Release()
	}

	// Draw border
	border := thorvg.NewShape()
	if border != nil {
		bdr, gd, bd, _ := borderColor.RGBA()
		border.AppendRect(x, y, float32(r.Max.X-r.Min.X), float32(r.Max.Y-r.Min.Y), theme.CornerRadius, theme.CornerRadius, false)
		border.SetStrokeColor(bdr, gd, bd, 255)
		border.SetStrokeWidth(theme.BorderWidth)
		vg.Push(&border.Paint)
		border.Release()
	}

	// Determine text to display
	displayText := text
	if text == "" && placeholder != "" {
		displayText = placeholder
	}

	// Draw text with scroll support
	if displayText != "" {
		txt := thorvg.NewText()
		if txt != nil {
			txt.SetText(displayText)
			txt.SetSize(float32(theme.FontSize))

			// Use determined text color
			tr, tg, tb, _ := textColor.RGBA()
			txt.SetColor(tr, tg, tb)

			cs := inp.ComputedStyle()
			if cs.Font != nil && *cs.Font != "" {
				txt.SetFont(*cs.Font)
			} else {
				txt.SetFont(theme.FontFamily)
			}

			// Vertically center text
			tm := gouifont.MeasureText(displayText, gouifont.WithFontSize(theme.FontSize))
			textY := contentY + (contentHeight-float32(tm.Height))/2

			// Apply scroll offset
			txt.Translate(contentX-scrollOffset, textY)
			vg.Push(&txt.Paint)
			txt.Release()
		}
	}

	// Draw cursor if focused and cursor is visible
	if focused && cursorVis && !disabled {
		cursor := thorvg.NewShape()
		if cursor != nil {
			cr, cg, cb, _ := theme.FocusColor.RGBA()

			// Calculate cursor position using rune slice for UTF-8 support
			cursorX := contentX
			if cursorPos > 0 {
				runes := []rune(inp.text)
				cursorText := string(runes[:cursorPos])
				tm := gouifont.MeasureText(cursorText, gouifont.WithFontSize(theme.FontSize))
				cursorX += float32(tm.Width)
			}
			cursorX -= scrollOffset

			tm := gouifont.MeasureText("Hg", gouifont.WithFontSize(theme.FontSize))
			cursorHeight := float32(tm.Height)
			cursorWidth := float32(2.0)

			cursor.AppendRect(cursorX, contentY, cursorWidth, cursorHeight, 0, 0, false)
			cursor.SetFillColor(cr, cg, cb, 255)
			vg.Push(&cursor.Paint)
			cursor.Release()
		}
	}

	childRendered := inp.BaseComponent.Paint(vg)
	// If focused, always return true to ensure continuous redraw
	// Otherwise, only return true if we actually rendered
	if focused {
		return true
	}
	return true || childRendered
}

// SetTextSize sets font size
func (inp *Input) SetTextSize(size float64) {
	inp.mu.Lock()
	if inp.theme.FontSize != size {
		inp.theme.FontSize = size
		inp.mu.Unlock()
		inp.MarkDirty(false)
	} else {
		inp.mu.Unlock()
	}
}

// SetFontFamily sets font family
func (inp *Input) SetFontFamily(name string) {
	inp.mu.Lock()
	if inp.theme.FontFamily != name {
		inp.theme.FontFamily = name
		inp.mu.Unlock()
		inp.MarkDirty(false)
	} else {
		inp.mu.Unlock()
	}
}
