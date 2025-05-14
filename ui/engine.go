package ui

import (
	"reflect"

	"github.com/millken/goui/event"
)

type Engine struct {
	app       App
	state     *event.State
	isPaused  bool
	isStopped bool
	window    *Window
	windows   []*Window
	listeners map[event.Kind][]event.Listener
}

func NewEngine(app App) *Engine {
	return &Engine{
		app:       app,
		state:     event.NewState(),
		listeners: make(map[event.Kind][]event.Listener),
	}
}

func (eng *Engine) SetPaused(paused bool) {
	eng.isPaused = paused
}

func (eng *Engine) SetStopped(stopped bool) {
	eng.isStopped = stopped
}

func (eng *Engine) IsStopped() bool {
	return eng.isStopped
}

func (eng *Engine) IsPaused() bool {
	return eng.isPaused
}

func (eng *Engine) Stop() {
	if !eng.isPaused {
		eng.isPaused = true
	}
	eng.isStopped = true
	eng.app.Stop()
}

func (eng *Engine) Subscribe(kind event.Kind, listener event.Listener) {
	if _, ok := eng.listeners[kind]; !ok {
		eng.listeners[kind] = make([]event.Listener, 0)
	}
	eng.listeners[kind] = append(eng.listeners[kind], listener)
}

func (eng *Engine) Unsubscribe(kind event.Kind, listener event.Listener) {
	if _, ok := eng.listeners[kind]; ok {
		for i, l := range eng.listeners[kind] {
			if reflect.ValueOf(l).Pointer() == reflect.ValueOf(listener).Pointer() {
				eng.listeners[kind] = append(eng.listeners[kind][:i], eng.listeners[kind][i+1:]...)
				break
			}
		}
	}
}

func (eng *Engine) Publish(event event.Event) {
	if _, ok := eng.listeners[event.Kind]; ok {
		for _, listener := range eng.listeners[event.Kind] {
			if listener(event) {
				break
			}
		}
	}
}
