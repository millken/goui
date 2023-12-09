package main

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
