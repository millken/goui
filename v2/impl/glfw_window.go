package impl

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
	v2 "github.com/millken/goui/v2"
)

var _ v2.Window = (*GLFWWindow)(nil)

type GLFWWindow struct {
	win     *glfw.Window
	root    v2.Component
	visible bool
}

func NewGLFWWindow(title string, width, height int) (*GLFWWindow, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("glfw init failed: %w", err)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	win, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, fmt.Errorf("create window failed: %w", err)
	}
	return &GLFWWindow{
		win: win,
	}, nil
}

func (w *GLFWWindow) Show() {
	// Implementation to show the window
	w.visible = true
	v2.RegisterWindow(w)
}

func (w *GLFWWindow) Hide() {
	// Implementation to hide the window
	w.visible = false
}

func (w *GLFWWindow) Close() {
	// Implementation to close the window
}

func (w *GLFWWindow) SetTitle(title string) {
	// Implementation to set the window title
}

func (w *GLFWWindow) SetSize(width, height int) {
	// Implementation to set the window size
}

func (w *GLFWWindow) SetRoot(root v2.Component) {
	w.root = root
	v2.BindRoot(w, root)
}

func (w *GLFWWindow) IsVisible() bool {
	// Implementation to check if the window is visible
	return w.visible
}

func (w *GLFWWindow) GetSize() (width, height int) {
	// Implementation to get the window size
	return 0, 0
}

// PreRender is called before the renderer is initialized.
func (w *GLFWWindow) PreRender(r v2.Renderer) {
	// Setup code before renderer initialization
}

// PostRenderer is called after the renderer is cleaned up.
func (w *GLFWWindow) PostRenderer(r v2.Renderer) {
	// Cleanup code after renderer cleanup
}

func (w *GLFWWindow) RequestRedraw() {
	// mark dirty; for stub we do nothing
}
