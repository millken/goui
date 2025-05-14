package event

type Kind uint8

const (
	// Mouse events
	KindMouse Kind = iota
	KindResize
	KindKeyboard
)

type Type uint8

const (
	TypeNone Type = 0x00
	TypeDown Type = 0x01
	TypeUp   Type = 0x02
	TypeMove Type = 0x04
)

type Button uint8

const (
	ButtonUndefine    = Button(0x00)
	ButtonLeft        = Button(0x01)
	ButtonMiddle      = Button(0x04)
	ButtonRight       = Button(0x02)
	ButtonTouchFirst  = Button(0x10)
	ButtonTouchSecond = Button(0x20)
	ButtonTouchThird  = Button(0x40)
)

func (b Button) String() string {
	switch b {
	case ButtonUndefine:
		return "button-undefine"
	case ButtonLeft:
		return "button-left"
	case ButtonMiddle:
		return "button-middle"
	case ButtonRight:
		return "button-right"
	case ButtonTouchFirst:
		return "button-touch-first"
	case ButtonTouchSecond:
		return "button-touch-second"
	case ButtonTouchThird:
		return "button-touch-third"
	default:
		return "button-unknown"
	}
}

type Event struct {
	Kind    Kind
	Type    Type
	Button  Button
	X, Y    int32
	W, H    int32
	KeyCode KeyCode
	KeyName string
}

type Listener func(event Event) bool
