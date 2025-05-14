package event

// shiftMap maps keys to their Shift versions.
var shiftMap = map[string]string{
	"`": "~",
	"1": "!",
	"2": "@",
	"3": "#",
	"4": "$",
	"5": "%",
	"6": "^",
	"7": "&",
	"8": "*",
	"9": "(",
	"0": ")",
	"-": "_",
	"=": "+",
	"[": "{",
	"]": "}",
	`\`: "|",
	";": ":",
	`'`: `"`,
	",": "<",
	".": ">",
	"/": "?",
}

// State holds the current state of key/mouse events.
type State struct {
	// Mouse buttons.
	Button1 bool // Left
	Button2 bool // Middle
	Button3 bool // Right

	// Special keys
	Escape bool
	Space  bool
	Enter  bool
	Shift  bool
	Ctrl   bool
	Alt    bool
	Up     bool
	Left   bool
	Right  bool
	Down   bool

	// Pressed keys.
	keydown map[string]interface{}

	// Cursor position
	CursorX int
	CursorY int

	// Window resized
	WindowResized bool

	// Touch state
	Touching        bool
	TouchNumFingers int
	TouchCenterX    int
	TouchCenterY    int
	GestureRotated  float64
	GesturePinched  float64
}

// NewState creates a new event.State.
func NewState() *State {
	return &State{
		keydown: map[string]interface{}{},
	}
}
