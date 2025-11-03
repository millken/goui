package goui

import (
	"fmt"

	"github.com/millken/yoga"
)

// App 应用程序主结构
type App struct {
	root Component
}

// NewApp 创建应用程序
func NewApp(title string, width, height int, root Component) (*App, error) {
	if root == nil {
		return nil, fmt.Errorf("root component cannot be nil")
	}

	app := &App{
		root: root,
	}

	return app, nil
}

// Run 运行应用程序
func (a *App) Run() {
	// 初始化根组件
	if rootBase, ok := a.root.(*BaseComponent); ok {
		rootBase.mount()

		// 执行初始布局
		rootBase.layout(800, 600, yoga.DirectionLTR)
	}

	fmt.Println("GOUI 应用程序正在运行...")
	fmt.Println("注意：这是一个简化版本，完整的GUI渲染需要集成渲染引擎")

	// 简单的模拟循环
	// 在实际实现中，这里会是完整的事件循环和渲染循环
}

// Close 关闭应用程序
func (a *App) Close() {
	if a.root != nil {
		if rootBase, ok := a.root.(*BaseComponent); ok {
			rootBase.unmount()
			rootBase.destroy()
		}
	}
}