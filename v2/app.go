package v2

import (
	"time"
)

type App interface {
	Create(ctx *Context) error
	Destroy() error
}

func NewApp[T any](t T) *T {
	return &t
}

func Run(settings *Settings, app App) {
	ctx := NewContext()
	wm := NewWindowManager()
	ctx.SetWindowManager(wm)
	if err := app.Create(ctx); err != nil {
		panic(err)
	}
	defer app.Destroy()

	// 主循环：基于 Context，并在所有窗口关闭后退出
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	for {
		select {
		case ev := <-wm.Events():
			// 收集来自窗口实现的事件
			ctx.PushWindowEvent(ev)
		case <-ticker.C:
			now := time.Now()
			ctx.tick(now)

			// 驱动平台事件循环（如 GLFW 的 PollEvents）
			wm.CallPump()

			// 先处理收集到的窗口事件（输入、关闭等）并分发到组件树
			evs := ctx.WindowEvents()
			if len(evs) > 0 {
				for _, e := range evs {
					if e.Event != nil {
						if root := wm.GetRoot(e.Window); root != nil {
							handled := root.HandleEvent(e.Event)
							if handled {
								wm.MarkDirty(e.Window)
							}
						}
					}
				}
			}

			// 处理窗口关闭事件与渲染
			windows := wm.ListWindows()
			if len(windows) == 0 {
				return
			}

			// 仅渲染 dirty 窗口
			dirty := wm.PopDirtyWindows()
			if len(dirty) == 0 {
				continue
			}
			for _, w := range dirty {
				r := settings.Renderer
				if r == nil {
					continue
				}
				if err := r.Initialize(w); err != nil {
					continue
				}
				w.PreRender(r)
				width, height := w.GetSize()
				canvas, err := r.BeginFrame(width, height)
				if err == nil && canvas != nil {
					if root := wm.GetRoot(w); root != nil {
						root.Paint(canvas)
					}
				}
				r.EndFrame()
				r.Present()
				w.PostRenderer(r)
				r.Cleanup()
			}
		}
	}
}
