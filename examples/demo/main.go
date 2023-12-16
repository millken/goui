package main

import (
	"gitee.com/millken/goui/app"
	"gitee.com/millken/goui/ui"
)

type App struct {
}

func (a *App) Create() error {
	win, err := ui.NewWindow()
	if err != nil {
		return err
	}
	win.SetTitle("Hello World")
	win.SetSize(800, 600)
	layout := ui.NewLayout()
	label := ui.NewLabel("Hello World")
	layout.InsertWidget(label, 0)
	win.SetLayout(layout)
	win.Show()
	return nil
}

func (a *App) Destroy() error {
	return nil
}
func main() {

	simpleApp := &App{}
	app.Run(simpleApp)
}

// type app struct {
// 	rootWindow goui.Window
// }

// func (a *app) Startup(_ context.Context) error {

// 	panel := goui.NewPanel()
// 	layout := goui.NewLayout(1, 1)
// 	label := goui.NewLabel("Hello World")
// 	layout.SetLabel(label)
// 	panel.SetLayout(layout)
// 	a.rootWindow.SetPanel(panel)
// 	a.rootWindow.Show()
// 	return nil
// }

// func (a *app) Shutdown(_ context.Context) error {
// 	golog.Info("shutdown")
// 	os.Exit(0)
// 	return nil
// }

// func (a *app) Window() goui.Window {
// 	return a.rootWindow
// }

// func NewApp() *app {
// 	demoWin, err := goui.NewWindow()
// 	if err != nil {
// 		panic(err)
// 	}
// 	demoWin.SetSize(800, 600)
// 	return &app{
// 		rootWindow: demoWin,
// 	}
// }

// func main() {
// 	goui.Run(NewApp())
// }
