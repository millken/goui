package component

import "github.com/millken/go/thorvg"

type Component interface {
	AddChild(Component)
	// RemoveChild(Component)

	Children() []Component

	Render(*thorvg.Canvas)
}
