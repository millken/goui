package goui

type Listener struct {
	event        Event
	eventHandler func(event Event)
}
