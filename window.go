package goui

import "context"

type Window interface {
	// SetSize sets the size of the window.
	SetSize(width, height int)
	// Size returns the size of the window.
	Size() (width, height int)
	//Show shows the window.
	Show()
	// Hide hides the window.
	Hide()
	//SetPanel sets the panel of the window.
	SetPanel(panel *Panel)
	//OnClose sets the callback for when the window is closed.
	OnClose(callback func(context.Context))
	//Destroy destroys the window.
	Destroy()

	Draw(*Canvas)
	PerformLayout()
}
