package goui

type BaseComponent struct {
	FrameChangeCallback func()
	frame               Rect
	border              Border
	scale               float32
	NeedsLayout         bool
}

// Scale returns the scale that has been applied to this panel. This will be automatically applied, transforming the
// graphics and mouse events.
func (p *BaseComponent) Scale() float32 {
	if p.scale <= 0 { // This happens if not explicitly set. 0 or less isn't valid, so make it 1
		p.scale = 1
	}
	return p.scale
}

// SetFrameRect sets the location and size of the panel in its parent's coordinate system.
func (p *BaseComponent) SetFrameRect(rect Rect) {
	scale := p.Scale()
	rect.Width /= scale
	rect.Height /= scale
	moved := p.frame.X != rect.X || p.frame.Y != rect.Y
	resized := p.frame.Width != rect.Width || p.frame.Height != rect.Height
	if moved || resized {
		if moved {
			p.frame.Point = rect.Point
		}
		if resized {
			p.frame.Size = rect.Size
			p.NeedsLayout = true
		}
		if p.FrameChangeCallback != nil {
			p.FrameChangeCallback()
		}

		//p.MarkForRedraw()
	}
}
