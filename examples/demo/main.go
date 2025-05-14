package main

import (
	"fmt"
	"time"

	"github.com/millken/goui/ui"
	"github.com/millken/yoga"
)

type App struct {
}

func (a *App) Create(settings *ui.Settings) error {
	settings.Width = 800
	return nil
}

func (a *App) Start() error {
	win, err := ui.NewWindow()
	if err != nil {
		return err
	}
	win.SetTitle("Hello World")
	win.SetSize(800, 600)

	yConfig := yoga.ConfigNew()
	yConfig.SetPointScaleFactor(1.445)
	layout := ui.NewLayer(yConfig)
	layout.Node().StyleSetFlexDirection(yoga.FlexDirectionRow)
	layout1 := ui.NewLayer(yConfig)
	layout1.Style().SetBackground("red")
	layout1.Node().StyleSetWidthPercent(50)
	layout.InsertWidget(layout1, 0)

	layout2 := ui.NewLayer(yConfig)
	layout2.Style().SetBackground("green")
	layout2.Node().StyleSetWidthPercent(50)
	layout.InsertWidget(layout2, 1)
	// label := ui.NewLabel("Hello World")
	// layout.InsertWidget(label, 0)
	win.SetLayout(layout)
	win.Show()
	return nil
}

func (a *App) Pause() {
	fmt.Println("pause")
}
func (a *App) Resume() {
	fmt.Println("resume")
}
func (a *App) Destroy() {
	fmt.Println("destroy")
}
func (a *App) Stop() {
	fmt.Println("stop")
}

func (a *App) Tick(elaspedTime time.Duration, syncChan chan<- struct{}) {
}
func main() {

	simpleApp := &App{}
	ui.Run(simpleApp)
}
