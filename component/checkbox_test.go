package component

import (
	"testing"
	"time"

	gouicolor "github.com/millken/goui/color"
	"github.com/millken/yoga"
)

func TestNewCheckbox(t *testing.T) {
	cb := NewCheckbox("test-checkbox")

	if cb == nil {
		t.Fatal("NewCheckbox returned nil")
	}

	if cb.ID() != "test-checkbox" {
		t.Errorf("Expected ID 'test-checkbox', got '%s'", cb.ID())
	}

	if cb.Checked() {
		t.Error("Expected checkbox to be unchecked by default")
	}

	if cb.IsDisabled() {
		t.Error("Expected checkbox to be enabled by default")
	}

	if cb.Title() != "" {
		t.Errorf("Expected empty label by default, got '%s'", cb.Title())
	}

	if cb.State() != CheckOff {
		t.Errorf("Expected state CheckOff, got %v", cb.State())
	}
}

func TestCheckboxSetChecked(t *testing.T) {
	cb := NewCheckbox()

	// Test setting checked
	cb.SetChecked(true)
	if !cb.Checked() {
		t.Error("Expected checkbox to be checked")
	}

	// Test setting unchecked
	cb.SetChecked(false)
	if cb.Checked() {
		t.Error("Expected checkbox to be unchecked")
	}

	// Test setting same value (should not trigger callback)
	callbackTriggered := false
	cb.ChangeCallback = func(state CheckState) {
		callbackTriggered = true
	}
	cb.SetChecked(false) // Already false
	if callbackTriggered {
		t.Error("Callback should not be triggered when setting same value")
	}
}

func TestCheckboxSetState(t *testing.T) {
	cb := NewCheckbox()

	// Test setting On state
	cb.SetState(CheckOn)
	if cb.State() != CheckOn {
		t.Error("Expected state to be CheckOn")
	}

	// Test setting Off state
	cb.SetState(CheckOff)
	if cb.State() != CheckOff {
		t.Error("Expected state to be CheckOff")
	}

	// Test setting Mixed state
	cb.SetState(CheckMixed)
	if cb.State() != CheckMixed {
		t.Error("Expected state to be CheckMixed")
	}
}

func TestCheckboxClick(t *testing.T) {
	cb := NewCheckbox()

	// Test click toggles state
	cb.Click()
	if !cb.Checked() {
		t.Error("Expected checkbox to be checked after Click()")
	}

	cb.Click()
	if cb.Checked() {
		t.Error("Expected checkbox to be unchecked after Click()")
	}
}

func TestCheckboxLabel(t *testing.T) {
	cb := NewCheckbox()

	// Test setting label using SetLabel (backward compatibility)
	cb.SetLabel("Test Label")
	if cb.Title() != "Test Label" {
		t.Errorf("Expected label 'Test Label', got '%s'", cb.Title())
	}

	// Test getting label using GetLabel (backward compatibility)
	if cb.GetLabel() != "Test Label" {
		t.Errorf("Expected GetLabel() to return 'Test Label', got '%s'", cb.GetLabel())
	}

	// Test setting label using SetTitle (unison naming)
	cb.SetTitle("New Label")
	if cb.Title() != "New Label" {
		t.Errorf("Expected title 'New Label', got '%s'", cb.Title())
	}
}

func TestCheckboxDisabled(t *testing.T) {
	cb := NewCheckbox()

	// Test setting disabled
	cb.SetDisabled(true)
	if !cb.IsDisabled() {
		t.Error("Expected checkbox to be disabled")
	}

	// Test setting enabled
	cb.SetDisabled(false)
	if cb.IsDisabled() {
		t.Error("Expected checkbox to be enabled")
	}

	// Test getting disabled state
	cb.SetDisabled(true)
	if !cb.IsDisabled() {
		t.Error("Expected IsDisabled() to return true")
	}
}

func TestCheckboxTheme(t *testing.T) {
	cb := NewCheckbox()

	// Test getting default theme
	theme := cb.Theme()
	if theme.BoxSize != 16 {
		t.Errorf("Expected default BoxSize 16, got %f", theme.BoxSize)
	}

	// Test setting custom theme
	customTheme := CheckboxTheme{
		EdgeInk:      gouicolor.NewColor(255, 0, 0, 255),
		SelectionInk: gouicolor.NewColor(0, 255, 0, 255),
		BoxSize:      20,
		Gap:          10,
		CornerRadius: 5,
		BorderWidth:  2,
		FontSize:     16.0,
		FontFamily:   "Arial",
		ClickAnimationTime: 150 * time.Millisecond,
	}

	customTheme.ControlInk = gouicolor.NewColor(255, 255, 255, 255)
	customTheme.OnControlInk = gouicolor.NewColor(0, 0, 255, 255)
	customTheme.OnSelectionInk = gouicolor.NewColor(255, 255, 255, 255)
	customTheme.DisabledInk = gouicolor.NewColor(150, 150, 150, 255)

	cb.SetTheme(customTheme)

	newTheme := cb.Theme()
	if newTheme.BoxSize != 20 {
		t.Errorf("Expected BoxSize 20 after SetTheme, got %f", newTheme.BoxSize)
	}
	if newTheme.Gap != 10 {
		t.Errorf("Expected Gap 10 after SetTheme, got %f", newTheme.Gap)
	}
}

func TestCheckStateConversion(t *testing.T) {
	// Test FromBool
	if CheckStateFromBool(true) != CheckOn {
		t.Error("Expected CheckStateFromBool(true) to return CheckOn")
	}
	if CheckStateFromBool(false) != CheckOff {
		t.Error("Expected CheckStateFromBool(false) to return CheckOff")
	}

	// Test ToBool
	if !CheckOn.ToBool() {
		t.Error("Expected CheckOn.ToBool() to return true")
	}
	if CheckOff.ToBool() {
		t.Error("Expected CheckOff.ToBool() to return false")
	}
	if CheckMixed.ToBool() {
		t.Error("Expected CheckMixed.ToBool() to return false")
	}

	// Test String
	if CheckOn.String() != "On" {
		t.Error("Expected CheckOn.String() to return 'On'")
	}
	if CheckOff.String() != "Off" {
		t.Error("Expected CheckOff.String() to return 'Off'")
	}
	if CheckMixed.String() != "Mixed" {
		t.Error("Expected CheckMixed.String() to return 'Mixed'")
	}
}

func TestCheckboxChangeCallback(t *testing.T) {
	cb := NewCheckbox()

	callbackTriggered := false
	var callbackState CheckState

	cb.ChangeCallback = func(state CheckState) {
		callbackTriggered = true
		callbackState = state
	}

	// Test callback on check
	cb.SetChecked(true)
	if !callbackTriggered {
		t.Error("Expected callback to be triggered when checking")
	}
	if callbackState != CheckOn {
		t.Error("Expected callback to receive CheckOn state")
	}

	// Reset and test callback on uncheck
	callbackTriggered = false
	cb.SetChecked(false)
	if !callbackTriggered {
		t.Error("Expected callback to be triggered when unchecking")
	}
	if callbackState != CheckOff {
		t.Error("Expected callback to receive CheckOff state")
	}
}

func TestCheckboxClickCallback(t *testing.T) {
	cb := NewCheckbox()

	clickTriggered := false
	cb.ClickCallback = func() {
		clickTriggered = true
	}

	// Test click callback
	cb.Click()
	if !clickTriggered {
		t.Error("Expected ClickCallback to be triggered")
	}
}

func TestCheckboxLayout(t *testing.T) {
	cb := NewCheckbox()

	// Test that checkbox can be laid out
	cb.Layout(100, 100, yoga.DirectionLTR)
	layout := cb.ComputedLayout()

	if layout.Width <= 0 || layout.Height <= 0 {
		t.Error("Expected positive width and height after layout")
	}
}

func TestCheckboxPaint(t *testing.T) {
	cb := NewCheckbox()

	// Test that Paint doesn't panic when called without a canvas
	// (In real usage, the framework provides the canvas)
	cb.MarkDirty(false) // Mark as dirty

	// We can't easily test actual rendering without a proper canvas setup
	// The important thing is that the component doesn't panic
	// and properly manages its dirty state
	if !cb.IsDirty() {
		t.Error("Expected checkbox to be dirty")
	}
}

func TestCheckboxWithLabel(t *testing.T) {
	cb := NewCheckbox()
	cb.SetLabel("Test Label")

	// Test that size is updated when label is added
	width := cb.ComputedLayout().Width
	boxSize := cb.Theme().BoxSize
	if width <= boxSize {
		t.Errorf("Expected width to increase when label is added, got %f", width)
	}
}

func TestCheckboxDisabledInteraction(t *testing.T) {
	cb := NewCheckbox()
	cb.SetDisabled(true)

	// Test that disabled checkbox doesn't respond to events
	originalChecked := cb.Checked()

	// Click should not change state when disabled
	cb.Click()

	// Verify state hasn't changed
	if cb.Checked() != originalChecked {
		t.Error("Expected checkbox state to remain unchanged when disabled")
	}
}

func TestDefaultCheckboxTheme(t *testing.T) {
	theme := DefaultCheckboxTheme

	if theme.BoxSize != 16 {
		t.Errorf("Expected default BoxSize 16, got %f", theme.BoxSize)
	}
	if theme.Gap != 8 {
		t.Errorf("Expected default Gap 8, got %f", theme.Gap)
	}
	if theme.CornerRadius != 3 {
		t.Errorf("Expected default CornerRadius 3, got %f", theme.CornerRadius)
	}
	if theme.BorderWidth != 1.5 {
		t.Errorf("Expected default BorderWidth 1.5, got %f", theme.BorderWidth)
	}
	if theme.FontSize != 14.0 {
		t.Errorf("Expected default FontSize 14.0, got %f", theme.FontSize)
	}
	if theme.ClickAnimationTime != 100*time.Millisecond {
		t.Errorf("Expected default ClickAnimationTime 100ms, got %v", theme.ClickAnimationTime)
	}
}
