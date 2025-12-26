package component

import (
	"sync"
	"sync/atomic"
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

// Global focus management - only one input can be focused at a time
// Using atomic.Pointer for lock-free access
var focusedInput atomic.Pointer[Input]

// Input represents a text input field
type Input struct {
	*BaseComponent
	theme       InputTheme
	text        string
	placeholder string
	focused     atomic.Bool // Use atomic to ensure consistent focus state
	disabled    bool
	mu          sync.RWMutex

	// Cursor
	cursorPosition int // rune position, not byte position
	cursorVisible  bool
	lastBlinkTime  time.Time
	scrollOffset   float32 // horizontal scroll offset in pixels

	// Selection
	selectionStart  int // start of selection (rune position), -1 if no selection
	selectionEnd    int // end of selection (rune position)
	selectionAnchor int // anchor point for selection expansion
}

// NewInput creates a new input field
func NewInput(id ...string) *Input {
	input := &Input{
		BaseComponent:  NewBaseComponent(id...),
		theme:          DefaultInputTheme,
		text:           "",
		placeholder:    "",
		disabled:       false,
		cursorPosition: 0,
		cursorVisible:  false, // Only show cursor when focused
		lastBlinkTime:  time.Now(),
		scrollOffset:   0,
	}
	// Initialize atomic focused to false
	input.focused.Store(false)

	// Set default size
	input.SetWidth(200)
	input.SetHeight(40)

	// Register mouse event handler for focus and click positioning
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

		// On mouse down inside input
		if me.Down && inside {
			// Calculate cursor position from click first (before locking)
			newPos := input.positionFromX(me.X)

			input.mu.Lock()
			wasFocused := input.focused.Load()
			input.mu.Unlock()

			// Use global focus management
			input.requestFocus()

			input.mu.Lock()
			input.cursorPosition = newPos
			// Start selection - set anchor to current position
			input.selectionAnchor = newPos
			input.selectionStart = newPos
			input.selectionEnd = newPos
			input.mu.Unlock()

			// Scroll to make cursor visible
			input.scrollCursorIntoView()

			if !wasFocused {
				input.MarkDirty(false)
			} else {
				input.MarkDirty(false)
			}
			return true
		}

		// On mouse click outside, remove focus and clear selection
		if me.Down && me.Count == 1 && !inside {
			input.blur()
		}

		return false
	}))

	// Register mouse drag handler for selection
	input.RegisterEventListener(event.EventMouseDrag, event.EventListenerFunc(func(e event.Event) bool {
		if input.disabled || !input.focused.Load() {
			return false
		}

		me, ok := e.(*event.MouseDragEvent)
		if !ok {
			return false
		}

		r := input.Bounds()
		inside := me.X >= float32(r.Min.X) && me.X < float32(r.Max.X) &&
			me.Y >= float32(r.Min.Y) && me.Y < float32(r.Max.Y)

		if inside {
			newPos := input.positionFromX(me.X)
			input.mu.Lock()
			// Update cursor and selection based on anchor
			input.cursorPosition = newPos
			if newPos < input.selectionAnchor {
				input.selectionStart = newPos
				input.selectionEnd = input.selectionAnchor
			} else {
				input.selectionStart = input.selectionAnchor
				input.selectionEnd = newPos
			}
			input.mu.Unlock()
			input.scrollCursorIntoView()
			input.MarkDirty(false)
			return true
		}

		return false
	}))

	// Register key event handler for text editing
	input.RegisterEventListener(event.EventKey, event.EventListenerFunc(func(e event.Event) bool {
		if !input.focused.Load() || input.disabled {
			return false
		}

		ke, ok := e.(*event.KeyEvent)
		if !ok {
			return false
		}

		// Handle special keys
		if ke.Down {
			// Check for Shift key for selection
			shiftPressed := ke.Mods&event.ModShift != 0
			ctrlPressed := ke.Mods&event.ModControl != 0

			switch ke.Key {
			case event.KeyCodeDeleteBackspace:
				input.handleBackspace()
			case event.KeyCodeReturnEnter:
				// Handle enter (could trigger callback in future)
				return true
			case event.KeyCodeLeftArrow:
				if shiftPressed {
					input.expandSelection(-1)
				} else {
					input.moveCursor(-1)
				}
			case event.KeyCodeRightArrow:
				if shiftPressed {
					input.expandSelection(1)
				} else {
					input.moveCursor(1)
				}
			case event.KeyCodeHome:
				if shiftPressed {
					input.expandSelectionTo(0)
				} else {
					input.mu.Lock()
					input.cursorPosition = 0
					input.scrollOffset = 0
					input.mu.Unlock()
					input.MarkDirty(false)
				}
			case event.KeyCodeEnd:
				if shiftPressed {
					input.expandSelectionTo(utf8.RuneCountInString(input.text))
				} else {
					input.mu.Lock()
					input.cursorPosition = utf8.RuneCountInString(input.text)
					input.scrollCursorIntoViewNoLock()
					input.mu.Unlock()
					input.MarkDirty(false)
				}
			case event.KeyCodeEscape:
				input.blur()
			case event.KeyCodeC:
				// Ctrl+C: Copy
				if ctrlPressed {
					input.copySelection()
				} else {
					return false // Let CharEvent handle
				}
			case event.KeyCodeV:
				// Ctrl+V: Paste
				if ctrlPressed {
					input.pasteFromClipboard()
				} else {
					return false // Let CharEvent handle
				}
			case event.KeyCodeA:
				// Ctrl+A: Select all
				if ctrlPressed {
					input.selectAll()
				} else {
					return false // Let CharEvent handle
				}
			case event.KeyCodeX:
				// Ctrl+X: Cut
				if ctrlPressed {
					input.cutSelection()
				} else {
					return false // Let CharEvent handle
				}
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
		if !input.focused.Load() || input.disabled {
			return false
		}

		ce, ok := e.(*event.CharEvent)
		if !ok {
			return false
		}

		// Filter out control characters only, allow all printable chars including Chinese
		if ce.Rune < 32 {
			return false
		}

		// Insert character at cursor position using rune slice for UTF-8 support
		input.mu.Lock()
		runes := []rune(input.text)
		input.text = string(runes[:input.cursorPosition]) + string(ce.Rune) + string(runes[input.cursorPosition:])
		input.cursorPosition++
		input.scrollCursorIntoViewNoLock()
		input.mu.Unlock()

		// Mark dirty to trigger redraw
		input.MarkDirty(false)
		return true
	}))

	// Register focus event handler - only update cursor style, actual focus state is managed by requestFocus/blur
	input.RegisterEventListener(event.EventFocus, event.EventListenerFunc(func(e event.Event) bool {
		if input.disabled {
			return false
		}

		_, ok := e.(*event.FocusEvent)
		if !ok {
			return false
		}

		// Only update cursor style, don't change focus state (managed by requestFocus/blur)
		if rc, ok := input.Root().(RootComponent); ok && rc != nil {
			if w := rc.GetWindow(); w != nil {
				// Check if we are the globally focused input
				if focusedInput.Load() == input {
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
	inp.scrollOffset = 0
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
			inp.focused.Store(false)
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

	// If already focused, do nothing
	if focusedInput.Load() == inp {
		return
	}

	// Remove focus from previously focused input
	if prev := focusedInput.Load(); prev != nil {
		prev.mu.Lock()
		wasFocused := prev.focused.Load()
		prev.focused.Store(false)
		prev.cursorVisible = false
		prev.mu.Unlock()
		if wasFocused {
			prev.MarkDirty(false)
		}
	}

	// Set focus to this input
	inp.mu.Lock()
	inp.focused.Store(true)
	inp.cursorVisible = true
	inp.lastBlinkTime = time.Now()
	inp.mu.Unlock()

	// Atomically swap to set new focus
	focusedInput.Store(inp)
	inp.MarkDirty(false)
}

// blur removes focus from input
func (inp *Input) blur() {
	// Try to swap out focus - only succeeds if we are the focused input
	if focusedInput.CompareAndSwap(inp, nil) {
		inp.mu.Lock()
		wasFocused := inp.focused.Load()
		inp.focused.Store(false)
		inp.cursorVisible = false
		inp.mu.Unlock()
		if wasFocused {
			inp.MarkDirty(false)
		}
	}
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
		// Scroll to make cursor visible (handles case when text becomes short)
		inp.scrollCursorIntoViewNoLock()
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
	inp.scrollCursorIntoViewNoLock()
	inp.MarkDirty(false)
}

// scrollCursorIntoViewNoLock scrolls to make cursor visible (must be called with lock held)
func (inp *Input) scrollCursorIntoViewNoLock() {
	r := inp.Bounds()
	contentWidth := float32(r.Max.X-r.Min.X) - inp.theme.PaddingX*2

	// Measure text before cursor
	runes := []rune(inp.text)
	var cursorText string
	if inp.cursorPosition > 0 && inp.cursorPosition <= len(runes) {
		cursorText = string(runes[:inp.cursorPosition])
	} else {
		cursorText = ""
	}
	tm := gouifont.MeasureText(cursorText, gouifont.WithFontSize(inp.theme.FontSize))
	cursorPixelPos := float32(tm.Width)

	// Calculate total text width
	totalTextWidth := gouifont.MeasureText(inp.text, gouifont.WithFontSize(inp.theme.FontSize)).Width

	// If text fits, no scroll needed
	if float32(totalTextWidth) <= contentWidth {
		inp.scrollOffset = 0
		return
	}

	// Calculate maximum scroll offset (text width - content width)
	maxScrollOffset := float32(totalTextWidth) - contentWidth
	if maxScrollOffset < 0 {
		maxScrollOffset = 0
	}

	// Calculate cursor visible position
	cursorVisiblePos := cursorPixelPos - inp.scrollOffset

	// Scroll if cursor is near right edge (> contentWidth - 20)
	if cursorVisiblePos > contentWidth-20 {
		inp.scrollOffset = cursorPixelPos - (contentWidth - 20)
	}
	// Scroll left if cursor is past left edge (< 20)
	if cursorVisiblePos < 20 {
		inp.scrollOffset = cursorPixelPos - 20
	}
	// Don't scroll past beginning
	if inp.scrollOffset < 0 {
		inp.scrollOffset = 0
	}
	// Don't scroll past end
	if inp.scrollOffset > maxScrollOffset {
		inp.scrollOffset = maxScrollOffset
	}
}

// Clear clears input text
func (inp *Input) Clear() {
	inp.SetText("")
}

// positionFromX converts X coordinate to cursor position
func (inp *Input) positionFromX(x float32) int {
	inp.mu.RLock()
	defer inp.mu.RUnlock()

	r := inp.Bounds()
	contentX := float32(r.Min.X) + inp.theme.PaddingX

	// Adjust x to be relative to content area and scroll offset
	adjustedX := x - contentX + inp.scrollOffset

	if inp.text == "" || adjustedX <= 0 {
		return 0
	}

	// Find the closest cursor position
	runes := []rune(inp.text)
	minDist := float32(999999)
	bestPos := 0

	for i := 0; i <= len(runes); i++ {
		testText := string(runes[:i])
		tm := gouifont.MeasureText(testText, gouifont.WithFontSize(inp.theme.FontSize))
		textWidth := float32(tm.Width)
		dist := adjustedX - textWidth
		if dist < 0 {
			dist = -dist
		}
		if dist < minDist {
			minDist = dist
			bestPos = i
		}
	}

	return bestPos
}

// scrollCursorIntoView scrolls to make cursor visible
func (inp *Input) scrollCursorIntoView() {
	inp.mu.Lock()
	inp.scrollCursorIntoViewNoLock()
	inp.mu.Unlock()
}

// expandSelection expands selection by moving the cursor
func (inp *Input) expandSelection(offset int) {
	inp.mu.Lock()
	defer inp.mu.Unlock()

	newPos := inp.cursorPosition + offset
	if newPos < 0 {
		newPos = 0
	} else if newPos > utf8.RuneCountInString(inp.text) {
		newPos = utf8.RuneCountInString(inp.text)
	}

	// If no selection exists, start from cursor
	if inp.selectionStart == -1 {
		inp.selectionStart = inp.cursorPosition
	}
	inp.selectionEnd = newPos
	inp.cursorPosition = newPos
	inp.scrollCursorIntoViewNoLock()
	inp.MarkDirty(false)
}

// expandSelectionTo expands selection to specific position
func (inp *Input) expandSelectionTo(pos int) {
	inp.mu.Lock()
	defer inp.mu.Unlock()

	// Clamp position
	if pos < 0 {
		pos = 0
	} else if pos > utf8.RuneCountInString(inp.text) {
		pos = utf8.RuneCountInString(inp.text)
	}

	// If no selection exists, start from cursor
	if inp.selectionStart == -1 {
		inp.selectionStart = inp.cursorPosition
	}
	inp.selectionEnd = pos
	inp.cursorPosition = pos
	inp.scrollCursorIntoViewNoLock()
	inp.MarkDirty(false)
}

// selectWordAt selects the word at the given position
func (inp *Input) selectWordAt(pos int) {
	inp.mu.Lock()
	defer inp.mu.Unlock()

	runes := []rune(inp.text)
	length := len(runes)

	if length == 0 {
		inp.selectionStart = 0
		inp.selectionEnd = 0
		return
	}

	// Clamp position
	if pos < 0 {
		pos = 0
	} else if pos > length {
		pos = length
	}

	// Find word boundaries
	start := pos
	end := pos

	// Helper to check if a rune is part of a word
	isWordChar := func(r rune) bool {
		// Word characters: letters, digits, underscore, and CJK characters
		if r >= 0x4E00 && r <= 0x9FFF { // CJK Unified Ideographs
			return true
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}
		if r >= '0' && r <= '9' {
			return true
		}
		return r == '_'
	}

	// Check if current position is a word character
	if pos < length && isWordChar(runes[pos]) {
		// Expand left
		for start > 0 && isWordChar(runes[start-1]) {
			start--
		}
		// Expand right
		for end < length && isWordChar(runes[end]) {
			end++
		}
	} else {
		// Non-word character: select it
		if pos < length {
			start = pos
			end = pos + 1
		}
	}

	inp.selectionAnchor = start
	inp.selectionStart = start
	inp.selectionEnd = end
	inp.cursorPosition = end
	inp.scrollCursorIntoViewNoLock()
	inp.MarkDirty(false)
}

// selectAll selects all text
func (inp *Input) selectAll() {
	inp.mu.Lock()
	defer inp.mu.Unlock()

	length := utf8.RuneCountInString(inp.text)
	inp.selectionAnchor = 0
	inp.selectionStart = 0
	inp.selectionEnd = length
	inp.cursorPosition = length
	inp.scrollCursorIntoViewNoLock()
	inp.MarkDirty(false)
}

// copySelection copies selected text to clipboard
func (inp *Input) copySelection() {
	inp.mu.Lock()
	defer inp.mu.Unlock()

	// Check if there's a selection
	if inp.selectionStart == -1 || inp.selectionEnd == -1 || inp.selectionStart == inp.selectionEnd {
		return
	}

	// Normalize selection
	start := inp.selectionStart
	end := inp.selectionEnd
	if start > end {
		start, end = end, start
	}

	// Get selected text
	runes := []rune(inp.text)
	if start >= len(runes) {
		return
	}
	if end > len(runes) {
		end = len(runes)
	}
	selectedText := string(runes[start:end])

	// Set clipboard
	if rc, ok := inp.Root().(RootComponent); ok && rc != nil {
		if w := rc.GetWindow(); w != nil {
			_ = w.SetClipboardText(selectedText)
		}
	}
}

// cutSelection cuts selected text to clipboard
func (inp *Input) cutSelection() {
	inp.mu.Lock()
	defer inp.mu.Unlock()

	// Check if there's a selection
	if inp.selectionStart == -1 || inp.selectionEnd == -1 || inp.selectionStart == inp.selectionEnd {
		return
	}

	// Normalize selection
	start := inp.selectionStart
	end := inp.selectionEnd
	if start > end {
		start, end = end, start
	}

	// Get selected text
	runes := []rune(inp.text)
	if start >= len(runes) {
		return
	}
	if end > len(runes) {
		end = len(runes)
	}
	selectedText := string(runes[start:end])

	// Delete selected text
	inp.text = string(runes[:start]) + string(runes[end:])
	inp.cursorPosition = start
	inp.selectionStart = -1
	inp.selectionEnd = -1
	inp.selectionAnchor = start
	inp.scrollCursorIntoViewNoLock()
	inp.MarkDirty(false)

	// Set clipboard
	if rc, ok := inp.Root().(RootComponent); ok && rc != nil {
		if w := rc.GetWindow(); w != nil {
			_ = w.SetClipboardText(selectedText)
		}
	}
}

// pasteFromClipboard pastes text from clipboard
func (inp *Input) pasteFromClipboard() {
	// First clear any existing selection
	inp.mu.Lock()
	if inp.selectionStart != -1 && inp.selectionEnd != -1 && inp.selectionStart != inp.selectionEnd {
		// Delete selected text
		start := inp.selectionStart
		end := inp.selectionEnd
		if start > end {
			start, end = end, start
		}
		runes := []rune(inp.text)
		if start < len(runes) && end <= len(runes) {
			inp.text = string(runes[:start]) + string(runes[end:])
			inp.cursorPosition = start
		}
		inp.selectionStart = -1
		inp.selectionEnd = -1
		inp.selectionAnchor = inp.cursorPosition
	}
	inp.mu.Unlock()

	// Get clipboard text
	var clipboardText string
	if rc, ok := inp.Root().(RootComponent); ok && rc != nil {
		if w := rc.GetWindow(); w != nil {
			text, err := w.GetClipboardText()
			if err == nil {
				clipboardText = text
			}
		}
	}

	// Insert clipboard text at cursor position
	if clipboardText != "" {
		inp.mu.Lock()
		runes := []rune(inp.text)
		clipboardRunes := []rune(clipboardText)
		newText := string(runes[:inp.cursorPosition]) + clipboardText + string(runes[inp.cursorPosition:])
		inp.text = newText
		inp.cursorPosition += len(clipboardRunes)
		inp.scrollCursorIntoViewNoLock()
		inp.MarkDirty(false)
		inp.mu.Unlock()
	}
}

// Paint draws input field
func (inp *Input) Paint(vg *thorvg.Canvas) bool {
	inp.mu.RLock()
	theme := inp.theme
	text := inp.text
	placeholder := inp.placeholder
	focused := inp.focused.Load()
	disabled := inp.disabled
	cursorPos := inp.cursorPosition
	cursorVis := inp.cursorVisible
	scrollOffset := inp.scrollOffset
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
	borderShape := thorvg.NewShape()
	if borderShape != nil {
		bdr, gd, bd, _ := borderColor.RGBA()
		borderShape.AppendRect(x, y, float32(r.Max.X-r.Min.X), float32(r.Max.Y-r.Min.Y), theme.CornerRadius, theme.CornerRadius, false)
		borderShape.SetStrokeColor(bdr, gd, bd, 255)
		borderShape.SetStrokeWidth(theme.BorderWidth)
		vg.Push(&borderShape.Paint)
		borderShape.Release()
	}

	// Determine text to display
	displayText := text
	if text == "" && placeholder != "" {
		displayText = placeholder
	}

	// Create clip shape for content area
	clipShape := thorvg.NewShape()
	if clipShape != nil {
		contentWidth := float32(r.Max.X-r.Min.X) - theme.PaddingX*2
		clipShape.AppendRect(contentX, contentY, contentWidth, contentHeight, 0, 0, false)
	}

	// Draw selection highlight
	if inp.selectionStart != -1 && inp.selectionEnd != -1 && inp.selectionStart != inp.selectionEnd {
		inp.mu.Lock()
		selStart := inp.selectionStart
		selEnd := inp.selectionEnd
		inp.mu.Unlock()

		// Normalize selection
		if selStart > selEnd {
			selStart, selEnd = selEnd, selStart
		}

		// Calculate selection positions in pixels
		runes := []rune(text)
		if selStart < len(runes) {
			startText := string(runes[:selStart])
			tmStart := gouifont.MeasureText(startText, gouifont.WithFontSize(theme.FontSize))
			selX := contentX + float32(tmStart.Width) - scrollOffset

			endText := string(runes[:selEnd])
			tmEnd := gouifont.MeasureText(endText, gouifont.WithFontSize(theme.FontSize))
			selWidth := float32(tmEnd.Width) - float32(tmStart.Width)

			// Draw selection background
			selBg := thorvg.NewShape()
			if selBg != nil {
				// Set clip on selection
				if clipShape != nil {
					selBg.SetClip(&clipShape.Paint)
				}

				selBg.AppendRect(selX, contentY, selWidth, contentHeight, 0, 0, false)
				selBg.SetFillColor(33, 150, 243, 100) // Light blue with alpha
				vg.Push(&selBg.Paint)
				selBg.Release()
			}
		}
	}

	// Draw text with scroll support and clip
	if displayText != "" {
		txt := thorvg.NewText()
		if txt != nil {
			// Set clip on text
			if clipShape != nil {
				txt.SetClip(&clipShape.Paint)
			}

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

			// Apply scroll offset (subtract to move left)
			txt.Translate(contentX-scrollOffset, textY)
			vg.Push(&txt.Paint)
			txt.Release()
		}
	}

	// Draw cursor if focused and cursor is visible
	if focused && cursorVis && !disabled {
		cursor := thorvg.NewShape()
		if cursor != nil {
			// Set clip on cursor
			if clipShape != nil {
				cursor.SetClip(&clipShape.Paint)
			}

			cr, cg, cb, _ := theme.FocusColor.RGBA()

			// Calculate cursor position using rune slice for UTF-8 support
			cursorX := contentX
			if cursorPos > 0 && len(text) > 0 {
				inp.mu.RLock()
				runes := []rune(inp.text)
				// Ensure cursorPos is within bounds
				if cursorPos <= len(runes) {
					cursorText := string(runes[:cursorPos])
					tm := gouifont.MeasureText(cursorText, gouifont.WithFontSize(inp.theme.FontSize))
					cursorX += float32(tm.Width)
				}
				inp.mu.RUnlock()
			}
			// Apply scroll offset
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

	// Release clip shape
	if clipShape != nil {
		clipShape.Release()
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
