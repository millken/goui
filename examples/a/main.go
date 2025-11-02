package main

import (
	"github.com/millken/goui"
	"github.com/millken/goui/components"
	"github.com/millken/goui/windows"
)

type App struct {
}

func (app *App) OnStart(ctx *goui.Context) error {
	flex := components.NewFlex()
	// return flex
	win1 := windows.NewRgfwWindow()
	win1.Init(ctx)
	win1.SetTitle("RGFW Window")
	win1.SetSize(800, 600)
	win1.SetRoot(flex)
	return nil
}

func (app *App) OnStop(ctx *goui.Context) error {
	return nil
}

func main() {
	// win1 := windows.NewRgfwWindow()
	// win1.Init(nil)
	// win1.SetTitle("RGFW Window")
	// win1.SetSize(800, 600)
	// win1.Show()
	setting := goui.NewSettings()
	app := goui.NewApp(App{})
	goui.Run(setting, app)
}
