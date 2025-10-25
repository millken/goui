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

	// 统一渲染驱动（极简版本）：
	r := settings.Renderer
	if r == nil {
		return
	}
	// 对所有注册的窗口执行一帧渲染（演示）。真实实现应有事件循环与重绘调度。
	for _, w := range listWindows() {
		if err := r.Initialize(w); err != nil {
			continue
		}
		w.PreRender(r)
		width, height := w.GetSize()
		canvas, err := r.BeginFrame(width, height)
		if err == nil && canvas != nil {
			if root := getRoot(w); root != nil {
				root.Paint(canvas)
			}
		}
		r.EndFrame()
		r.Present()
		w.PostRenderer(r)
		r.Cleanup()
	}
}

// --- minimal window registry for centralized driving ---
var (
	windows []Window
	roots   = map[Window]Component{}
)

// RegisterWindow registers a window to be centrally driven by Run.
func RegisterWindow(w Window) { windows = append(windows, w) }

// BindRoot binds a root component to a window for centralized painting.
func BindRoot(w Window, c Component) { roots[w] = c }

func listWindows() []Window      { return windows }
func getRoot(w Window) Component { return roots[w] }
