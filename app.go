package goui

import (
	"fmt"
	"time"
)

type App interface {
	OnStart(ctx *Context) error
	OnStop(ctx *Context) error
}

func NewApp[T any](t T) *T {
	return &t
}

func Run(settings *Settings, app App) {
	ctx := NewContext()
	// // mark scheduler main thread so Scheduler can execute tasks directly when called from the main loop
	// if ctx.GetScheduler() != nil {
	// 	ctx.GetScheduler().SetMain()
	// }
	wm := NewWindowManager()
	ctx.SetWindowManager(wm)
	if err := app.OnStart(ctx); err != nil {
		panic(err)
	}
	defer app.OnStop(ctx)

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

			// 执行主线程任务（由后端/其它 goroutine 请求在主线程执行的操作）
			if ctx.GetScheduler() != nil {
				ctx.GetScheduler().Poll()
			}

			// 先处理收集到的窗口事件（输入、关闭等）并分发到组件树
			evs := ctx.WindowEvents()
			if len(evs) > 0 {
				for _, e := range evs {
					// handle explicit payload events (mouse/key)
					if e.Event != nil {
						if root := wm.GetRoot(e.Window); root != nil {
							handled := root.HandleEvent(e.Event)
							if handled {
								wm.MarkDirty(e.Window)
							}
						}
						continue
					}
					// handle window-level lifecycle events
					switch e.Type {
					case WindowClosed:
						// call Close on the window to ensure proper cleanup on main thread
						if e.Window != nil {
							e.Window.Close()
						}
					}
				}
			}

			// 处理窗口关闭事件与渲染
			windows := wm.ListWindows()
			if len(windows) == 0 {
				fmt.Println("App.Run: no windows left, exiting")
				return
			}
			// 如果所有窗口都不可见（hidden），也退出
			anyVisible := false
			for _, ww := range windows {
				if ww != nil && ww.IsVisible() {
					anyVisible = true
					break
				}
			}
			if !anyVisible {
				fmt.Println("App.Run: all windows hidden, exiting")
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
