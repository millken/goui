package ui

type App interface {
	Create() error
	Destroy() error
}
