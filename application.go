package goui

import (
	"fmt"
	"runtime"
	"time"

	"github.com/millken/go/thorvg"
	"github.com/millken/goui/component"
	"github.com/millken/goui/event"
	rend "github.com/millken/goui/renderer"
	"github.com/millken/goui/window"
	glfwwin "github.com/millken/goui/window/glfw"
	"github.com/millken/yoga"
)

// Application 管理多个窗口会话。移动端可限制为单会话。
type Application struct {
	sessions []*Session
	// 默认配置用于新建会话
	defaultCfg *AppConfig
	// 渲染器工厂，默认使用 GLRenderer；测试可注入 fake renderer
	rendererFactory func() (rend.Renderer, error)
	// 窗口工厂，默认使用 GLFW 的 NewWindow；可注入其他后端
	windowFactory func(title string, width, height int) (window.Window, error)
	// 在主线程（Run 所在线程）执行的任务队列
	tasks chan func()
}

type Session struct {
	win      window.Window
	root     component.Component
	renderer rend.Renderer
	em       *event.EventManager
	cfg      *AppConfig

	lastSize [2]int
	closed   bool
}

type applicationOptions struct {
	thorvgThreads   int
	appConfig       *AppConfig
	rendererFactory func() (rend.Renderer, error)
	windowFactory   func(title string, width, height int) (window.Window, error)
}

type ApplicationOption func(*applicationOptions)

func WithThorvgThreads(n int) ApplicationOption {
	return func(opts *applicationOptions) {
		if n > 0 {
			opts.thorvgThreads = n
		}
	}
}

var defaultApplicationOptions = &applicationOptions{
	thorvgThreads:   1,
	appConfig:       &AppConfig{TargetFPS: 60, EnableVSync: true},
	rendererFactory: func() (rend.Renderer, error) { return rend.NewGLRenderer() },
	windowFactory:   func(title string, w, h int) (window.Window, error) { return glfwwin.NewWindow(title, w, h) },
}

// NewApplication 创建 Application 实例

func NewApplication(opts ...ApplicationOption) *Application {
	options := defaultApplicationOptions
	for _, opt := range opts {
		opt(options)
	}
	// 初始化 ThorVG 全局一次
	if err := thorvg.Init(uint(options.thorvgThreads)); err != nil {
		panic(err)
	}
	return &Application{
		defaultCfg:      options.appConfig,
		rendererFactory: options.rendererFactory,
		windowFactory:   options.windowFactory,
		tasks:           make(chan func(), 64),
	}
}

// SetRendererFactory 允许自定义渲染器工厂（用于测试或替换后端）
func (app *Application) SetRendererFactory(f func() (rend.Renderer, error)) {
	if f != nil {
		app.rendererFactory = f
	}
}

// NewSession 创建一个窗口会话
func (app *Application) NewSession(win window.Window, root component.Component) (*Session, error) {
	return app.NewSessionWithAppConfig(win, root, app.defaultCfg)
}

// NewSessionWithAppConfig 创建一个窗口会话；移动端可在外层限制数量
func (app *Application) NewSessionWithAppConfig(win window.Window, root component.Component, cfg *AppConfig) (*Session, error) {
	if root == nil || win == nil {
		return nil, fmt.Errorf("window/root cannot be nil")
	}
	// 绑定根
	root.SetRoot(root)
	// 渲染器需在窗口上下文 current 后创建
	win.MakeContextCurrent()
	r, err := app.rendererFactory()
	if err != nil {
		return nil, err
	}
	s := &Session{
		win:      win,
		root:     root,
		renderer: r,
		em:       event.NewEventManager(),
		cfg:      cfg,
	}
	// 注入窗口与事件管理器
	if base, ok := root.(component.RootComponent); ok {
		base.SetWindow(win)
		base.SetEventManager(s.em)
	} else {
		return nil, fmt.Errorf("root must implement RootComponent")
	}
	if s.cfg == nil {
		s.cfg = app.defaultCfg
	}
	win.SetVSync(s.cfg.EnableVSync)
	if err := s.em.Start(); err != nil {
		return nil, err
	}
	app.sessions = append(app.sessions, s)
	return s, nil
}

func (app *Application) Run() {
	runtime.LockOSThread()

	statsTicker := time.NewTicker(3 * time.Second)
	defer statsTicker.Stop()

	for app.anyOpen() {
		frameStart := time.Now()

		// 先执行主线程任务（如 SpawnWindow）
		for {
			select {
			case f := <-app.tasks:
				if f != nil {
					f()
				}
			default:
				goto drained
			}
		}
	drained:

		// 全局轮询以触发平台事件回调
		// 各后端可在各自窗口包装内部处理；这里依赖每会话 FetchEvents

		// 遍历会话：事件→布局→渲染
		renderedAny := false
		hadEvents := false
		for _, s := range app.sessions {
			if s == nil || !s.win.IsOpen() {
				// 窗口请求关闭时，主动清理并销毁窗口，避免界面上悬挂不消失
				if s != nil && !s.closed {
					s.closed = true
					if s.em != nil {
						s.em.Stop()
					}
					if s.renderer != nil {
						s.renderer.Cleanup()
					}
					// 销毁平台窗口
					_ = s.win.Close()
				}
				continue
			}
			// 批量取事件
			s.pullEvents(&hadEvents)
			// 尺寸/布局
			w, h := s.win.GetSize()
			resized := s.lastSize != [2]int{w, h}
			if resized {
				s.lastSize = [2]int{w, h}
				s.win.MakeContextCurrent()
				_ = s.renderer.ResetScene(w, h)
				s.root.MarkDirty(true)
			}
			// 布局
			layoutNeeded := resized
			if checker, ok := s.root.(interface{ CheckLayoutDirtyRecursive() bool }); ok {
				if checker.CheckLayoutDirtyRecursive() {
					layoutNeeded = true
				}
			}
			if layoutNeeded {
				s.root.Layout(float32(w), float32(h), yoga.DirectionLTR)
				if clearer, ok := s.root.(interface{ ClearLayoutDirtyRecursive() }); ok {
					clearer.ClearLayoutDirtyRecursive()
				}
			}
			// 渲染
			s.win.MakeContextCurrent()
			canvas, err := s.renderer.BeginFrame(w, h)
			if err != nil {
				continue
			}
			if s.root.Paint(canvas) {
				_ = s.renderer.EndFrame()
				s.win.SwapBuffers()
				renderedAny = true
			}
		}

		// 空闲帧休眠
		budget := time.Second / 60
		if len(app.sessions) > 0 {
			if cfg := app.sessions[0].cfg; cfg != nil && cfg.TargetFPS > 0 {
				budget = time.Second / time.Duration(cfg.TargetFPS)
			}
		}
		if !hadEvents && !renderedAny {
			if dt := time.Since(frameStart); dt < budget {
				time.Sleep(budget - dt)
			}
		}

		select {
		case <-statsTicker.C:
			// 可扩展统计输出
		default:
		}
	}
}

// pullEvents drains window events and posts to the session EventManager
func (s *Session) pullEvents(had *bool) {
	if s == nil || s.em == nil || s.win == nil {
		return
	}
	events := s.win.FetchEvents(0)
	for _, ev := range events {
		if ev != nil {
			_ = s.em.PostEvent(ev)
			if had != nil {
				*had = true
			}
		}
	}
}

// SetWindowFactory 允许外部注入窗口工厂（例如在移动端或自定义后端）
func (app *Application) SetWindowFactory(f func(title string, width, height int) (window.Window, error)) {
	if f != nil {
		app.windowFactory = f
	}
}

// Post 将函数投递到 Application.Run 所在的主线程执行
func (app *Application) Post(f func()) {
	if f == nil {
		return
	}
	select {
	case app.tasks <- f:
	default:
		// 队列满时，尝试阻塞投递，避免丢任务
		app.tasks <- f
	}
}

// SpawnWindow 在主线程创建新窗口并建立 Session；可从任意 goroutine 调用
// buildRoot 用于基于新窗口构建根组件
func (app *Application) SpawnWindow(title string, width, height int, buildRoot func(win window.Window) component.Component) (*Session, error) {
	if app == nil || app.windowFactory == nil || buildRoot == nil {
		return nil, fmt.Errorf("invalid application/windowFactory/buildRoot")
	}
	var (
		result *Session
		retErr error
		done   = make(chan struct{})
	)
	app.Post(func() {
		defer close(done)
		win, err := app.windowFactory(title, width, height)
		if err != nil {
			retErr = err
			return
		}
		root := buildRoot(win)
		sess, err := app.NewSession(win, root)
		if err != nil {
			retErr = err
			return
		}
		result = sess
	})
	<-done
	return result, retErr
}
func (app *Application) anyOpen() bool {
	for _, s := range app.sessions {
		if s != nil && s.win.IsOpen() {
			return true
		}
	}
	return false
}

func (app *Application) Close() {
	for _, s := range app.sessions {
		if s == nil {
			continue
		}
		if s.em != nil {
			s.em.Stop()
		}
		if s.renderer != nil {
			s.renderer.Cleanup()
		}
	}
	thorvg.Term()
}
