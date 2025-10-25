package v2

type App interface {
	Create() error
	Destroy() error
}

func NewApp[T any](t T) *T {
	return &t
}

func Run(settings *Settings, app App) {
	if err := app.Create(); err != nil {
		panic(err)
	}
	defer app.Destroy()
}
