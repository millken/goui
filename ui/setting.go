package ui

// EventMask defines mask event for enable/disable events receivers at runtime level
type EventMask int

const (
	// AllEventsDisable disables all input events on App
	AllEventsDisable = 0x00
	// MouseButtonEventEnabled enabled mouse buttons events receiver on App
	MouseButtonEventEnabled = 0x01
	// MouseMotionEventEnabled enabled mouse motion events receiver on App
	MouseMotionEventEnabled = 0x02
	// ScrollEventEnabled enabled scroll event receiver on App
	ScrollEventEnabled = 0x04
	// KeyEventEnabled enabled key event receiver on App
	KeyEventEnabled = 0x08
	// AllEventsEnabled enables all input events on App
	AllEventsEnabled = 0xFFFF
)

type Settings struct {
	Width     int
	Height    int
	EventMask EventMask
}
