package goui

// Settings contains configuration flags and thresholds used by PolyUI.
// This is a Go translation of org.polyfrost.polyui.Settings (Kotlin).
type Settings struct {
	Renderer                   Renderer
	Debug                      bool
	PrintPerfInfo              bool
	EnableDebugKeybind         bool
	MaxFPS                     int
	UnfocusedFPS               int
	EnableVSync                bool
	DefaultLocale              string
	ResourcePolicy             ResourcePolicy
	IsMasterFrameBuffer        bool
	MinDrawablesForFramebuffer int
	FramebuffersEnabled        bool
	RenderPausingEnabled       bool
	ComboMaxIntervalMs         int64
	MaxComboSize               int
	ClearComboWhenMaxed        bool
	DragThreshold              float32
	ScrollMultiplierX          float32
	ScrollMultiplierY          float32
	NaturalScrolling           bool
	MinimumWidth               float32
	MinimumHeight              float32
	MaximumWidth               float32
	MaximumHeight              float32
	AspectRatioX               float32
	AspectRatioY               float32
	CleanupAfterInit           bool
	CleanupOnShutdown          bool
	LoadTranslationsOnInit     bool
	SmoothThemeChanges         bool
	RoundedCorners             bool
}

type ResourcePolicy int

const (
	ResourcePolicyCrash ResourcePolicy = iota
	ResourcePolicyWarn
)

// NewSettings returns a Settings instance initialized with defaults matching the Kotlin implementation.
func NewSettings() *Settings {
	s := &Settings{
		Debug:                      false,
		PrintPerfInfo:              false,
		EnableDebugKeybind:         true,
		MaxFPS:                     0,
		UnfocusedFPS:               10,
		EnableVSync:                true,
		DefaultLocale:              "en_default",
		ResourcePolicy:             ResourcePolicyWarn,
		IsMasterFrameBuffer:        false,
		MinDrawablesForFramebuffer: 50,
		FramebuffersEnabled:        false,
		RenderPausingEnabled:       true,
		ComboMaxIntervalMs:         200,
		MaxComboSize:               4,
		ClearComboWhenMaxed:        false,
		DragThreshold:              3.0,
		ScrollMultiplierX:          4.0,
		ScrollMultiplierY:          4.0,
		NaturalScrolling:           false,
		MinimumWidth:               100.0,
		MinimumHeight:              100.0,
		MaximumWidth:               -1.0,
		MaximumHeight:              -1.0,
		AspectRatioX:               -1.0,
		AspectRatioY:               -1.0,
		CleanupAfterInit:           true,
		CleanupOnShutdown:          false,
		LoadTranslationsOnInit:     true,
		SmoothThemeChanges:         true,
		RoundedCorners:             true,
	}
	return s
}

// SetRenderer sets the Renderer to be used.
func (s *Settings) SetRenderer(r Renderer) { s.Renderer = r }

// SetAspectRatio sets the aspect ratio (x:y). Use -1 to disable an axis.
func (s *Settings) SetAspectRatio(x, y float32) {
	s.AspectRatioX = x
	s.AspectRatioY = y
}

// MinimumSize returns the minimum width and height set in the settings.
func (s *Settings) MinimumSize() (float32, float32) {
	return s.MinimumWidth, s.MinimumHeight
}

// MaximumSize returns the maximum width and height set in the settings.
func (s *Settings) MaximumSize() (float32, float32) {
	return s.MaximumWidth, s.MaximumHeight
}

// GetDragThreshold returns the drag threshold used to start dragging.
func (s *Settings) GetDragThreshold() float32 { return s.DragThreshold }

// GetComboMaxIntervalMs returns the maximum time window (in ms) for combo clicks.
func (s *Settings) GetComboMaxIntervalMs() int64 { return s.ComboMaxIntervalMs }

// GetMaxComboSize returns the maximum combo size allowed.
func (s *Settings) GetMaxComboSize() int { return s.MaxComboSize }

// GetClearComboWhenMaxed returns whether combos clear when maxed.
func (s *Settings) GetClearComboWhenMaxed() bool { return s.ClearComboWhenMaxed }

// GetNaturalScrolling returns the natural scrolling flag.
func (s *Settings) GetNaturalScrolling() bool { return s.NaturalScrolling }

// GetScrollMultiplierX returns the horizontal scroll multiplier.
func (s *Settings) GetScrollMultiplierX() float32 { return s.ScrollMultiplierX }

// GetScrollMultiplierY returns the vertical scroll multiplier.
func (s *Settings) GetScrollMultiplierY() float32 { return s.ScrollMultiplierY }
