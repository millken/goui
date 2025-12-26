# CheckBox Component (Refactored - Unison-Inspired)

The CheckBox component is a modern, themeable UI element inspired by the [unison](https://github.com/richardwilkes/unison) framework's design patterns. It allows users to toggle between checked, unchecked, and mixed states with professional visual feedback.

## ✨ Features

- 🎨 **Theme System** - Centralized styling with CheckboxTheme for consistent appearance
- ✅ **State Enum** - Three-state support (On/Off/Mixed) for complex UI scenarios
- 🖱️ **Label Click Support** - Click anywhere on checkbox or label to toggle
- 🎬 **Press Animation** - Visual feedback during mouse interactions
- 🔒 **Thread-Safe** - Mutex-protected state for concurrent access
- 📝 **Optional Label** - Automatic layout with proper spacing
- 🚫 **Disabled State** - Visual feedback and interaction prevention
- 🎯 **Dual Callbacks** - Separate Click and Change callbacks for flexibility
- 📏 **Configurable** - Theme-based customization of all visual aspects
- ♿ **Accessibility** - Proper event handling and clear visual states
- 🔄 **Backward Compatible** - Legacy API methods still work

## Architecture (Unison-Inspired)

This component follows the design patterns from [unison](https://github.com/richardwilkes/unison):

1. **Theme System**: Centralized styling through `CheckboxTheme` struct
2. **State Enum**: Three-state system (CheckOn/CheckOff/CheckMixed)
3. **Callback Separation**: `ClickCallback` for actions, `ChangeCallback` for state changes
4. **Pressed State**: Visual feedback during mouse interactions
5. **Thread Safety**: Mutex protection for all state access
6. **Naming Convention**: Follows unison's `SetTitle()`, `Title()`, `Click()` patterns

## Usage

### Basic Usage

```go
import "github.com/millken/goui/component"
import gouicolor "github.com/millken/goui/color"

// Create a basic checkbox
checkbox := component.NewCheckBox("my-checkbox")
checkbox.SetLabel("Enable notifications")

// Set change callback (receives state)
checkbox.ChangeCallback = func(state component.CheckState) {
    fmt.Printf("Checkbox state changed: %v\n", state)
}

// Set click callback (for action handling)
checkbox.ClickCallback = func() {
    fmt.Println("Checkbox was clicked!")
}
```

### Custom Theming

```go
// Create checkbox with custom theme
checkbox := component.NewCheckBox()
checkbox.SetLabel("Dark mode")

// Define custom theme
customTheme := component.CheckboxTheme{
    EdgeInk:        gouicolor.NewColor(76, 175, 80, 255),   // Green border
    SelectionInk:   gouicolor.NewColor(76, 175, 80, 255),   // Green when hovered/pressed
    OnSelectionInk: gouicolor.NewColor(255, 255, 255, 255), // White on pressed
    ControlInk:     gouicolor.NewColor(255, 255, 255, 255), // White background
    OnControlInk:   gouicolor.NewColor(76, 175, 80, 255),   // Green checkmark
    DisabledInk:    gouicolor.NewColor(200, 200, 200, 255), // Gray when disabled
    BoxSize:        18,    // Larger checkbox
    CornerRadius:   4,     // More rounded
    FontSize:       14.0,
    FontFamily:     "Arial",
    Gap:            10,    // More spacing
    BorderWidth:    2.0,   // Thicker border
    ClickAnimationTime: 150 * time.Millisecond,
}

checkbox.SetTheme(customTheme)
```

### State Management

```go
checkbox := component.NewCheckBox()

// Using boolean API (backward compatible)
checkbox.SetChecked(true)
if checkbox.Checked() {
    fmt.Println("Checkbox is checked")
}

// Using state enum API (new)
checkbox.SetState(component.CheckOn)
switch checkbox.State() {
case component.CheckOn:
    fmt.Println("Checkbox is On")
case component.CheckOff:
    fmt.Println("Checkbox is Off")
case component.CheckMixed:
    fmt.Println("Checkbox is Mixed (for tree views)")
}

// State conversions
boolValue := checkbox.Checked()           // CheckOn -> true, others -> false
state := component.CheckStateFromBool(true) // true -> CheckOn, false -> CheckOff
```

### Mixed State (Tree Selection)

```go
// Mixed state is useful for "select all" in tree views
selectAllCheckbox := component.NewCheckBox()
selectAllCheckbox.SetLabel("Select all")
selectAllCheckbox.SetState(component.CheckMixed)

// Update mixed state based on children
func updateSelectAllState() {
    allChecked := allChildrenChecked()
    noneChecked := noChildrenChecked()

    if allChecked {
        selectAllCheckbox.SetState(component.CheckOn)
    } else if noneChecked {
        selectAllCheckbox.SetState(component.CheckOff)
    } else {
        selectAllCheckbox.SetState(component.CheckMixed)
    }
}
```

### Disabled State

```go
checkbox := component.NewCheckBox()
checkbox.SetLabel("Premium feature")
checkbox.SetDisabled(true)
checkbox.SetChecked(true) // Pre-checked but disabled
```

## API Reference

### Constructors

- `NewCheckBox(id ...string) *CheckBox` - Creates a new checkbox (preferred)
- `NewCheckbox(id ...string) *CheckBox` - Alias for backward compatibility

### State Management

- `SetState(state CheckState)` - Set checkbox state (On/Off/Mixed)
- `State() CheckState` - Get current state
- `SetChecked(checked bool)` - Set checked state (bool → CheckOn/CheckOff)
- `Checked() bool` - Get checked state (CheckOn → true, others → false)

### Label (Unison Naming)

- `SetTitle(text string)` - Set label text (unison-style)
- `Title() string` - Get label text (unison-style)
- `SetLabel(label string)` - Alias for SetTitle (backward compatibility)
- `GetLabel() string` - Alias for Title (backward compatibility)

### Theme

- `SetTheme(theme CheckboxTheme)` - Set complete theme
- `Theme() CheckboxTheme` - Get current theme

### Interaction

- `Click()` - Programmatically click the checkbox
- `SetDisabled(disabled bool)` - Enable or disable
- `IsDisabled() bool` - Check if disabled

### Callbacks

- `ClickCallback func()` - Called when checkbox is clicked
- `ChangeCallback func(state CheckState)` - Called when state changes

### CheckboxTheme Structure

```go
type CheckboxTheme struct {
    EdgeInk        gouicolor.Color  // Border color
    SelectionInk   gouicolor.Color  // Hover/pressed/focus color
    OnSelectionInk gouicolor.Color  // Foreground when pressed
    ControlInk     gouicolor.Color  // Background color
    OnControlInk   gouicolor.Color  // Checkmark/indicator color
    DisabledInk    gouicolor.Color  // Disabled state color

    FontSize   float64
    FontFamily string
    BoxSize    float32  // Checkbox square size
    Gap        float32  // Spacing between box and label
    CornerRadius float32
    BorderWidth  float32

    ClickAnimationTime time.Duration
}
```

### Default Theme Values

- **BoxSize**: 16px
- **Gap**: 8px
- **CornerRadius**: 3px
- **BorderWidth**: 1.5px
- **FontSize**: 14.0
- **FontFamily**: "OPPO Sans 4.0"
- **EdgeInk**: Gray (#969696)
- **SelectionInk**: Blue (#2196F3)
- **OnSelectionInk**: White (#FFFFFF)
- **ControlInk**: White (#FFFFFF)
- **OnControlInk**: Blue (#2196F3)
- **DisabledInk**: Light gray (#C8C8C8)
- **ClickAnimationTime**: 100ms

## CheckState Enum

```go
type CheckState int

const (
    CheckOff CheckState = iota  // Unchecked
    CheckOn                      // Checked
    CheckMixed                   // Mixed state (for tree views)
)
```

### State Methods

- `String() string` - Returns "Off", "On", or "Mixed"
- `ToBool() bool` - CheckOn → true, others → false
- `CheckStateFromBool(b bool) CheckState` - true → CheckOn, false → CheckOff

## Behavior

### Mouse Interaction (Unison Patterns)

1. **Mouse Down**: Sets pressed state, shows pressed visual
2. **Mouse Drag**: Updates pressed state based on cursor position
3. **Mouse Up**: If pressed and inside, toggles state and calls callbacks
4. **Hover**: Shows pointer cursor and hover highlight

### Visual States

- **Normal**: Default appearance with ControlInk background
- **Hovered**: SelectionInk border + semi-transparent background highlight
- **Pressed**: SelectionInk background with OnSelectionInk foreground
- **Disabled**: DisabledInk for all elements, no interaction
- **Checked**: Shows checkmark (CheckOn) or horizontal line (CheckMixed)

### Thread Safety

All state access is protected by `sync.RWMutex` for safe concurrent access.

## Integration

The CheckBox component integrates with:

- **Event System**: Mouse events (down/drag/up) and hover
- **Layout System**: Yoga-based flexible layout
- **Styling System**: Theme-based customization
- **Rendering**: Thorvg GPU-accelerated drawing

## Migration from Old API

The refactored CheckBox maintains backward compatibility:

```go
// Old API (still works)
checkbox := component.NewCheckbox()
checkbox.SetLabel("Text")
checkbox.SetCheckColor(color)
checkbox.SetBorderColor(color)
checkbox.SetSize(20)
checkbox.SetPadding(10)
checkbox.OnChange = func(checked bool) { ... }

// New API (recommended)
checkbox := component.NewCheckBox()
checkbox.SetLabel("Text")
theme := checkbox.Theme()
theme.OnControlInk = color
theme.EdgeInk = color
theme.BoxSize = 20
theme.Gap = 10
checkbox.SetTheme(theme)
checkbox.ChangeCallback = func(state CheckState) { ... }
```

## Examples

See [examples/checkbox_demo](../examples/checkbox_demo/) for a complete working example demonstrating:
- Default themed checkboxes
- Custom themes (green, red, purple)
- Disabled state
- Mixed state (for tree selection)
- State inspection
- Both Click and Change callbacks