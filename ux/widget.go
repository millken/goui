package ux

// Widget 所有 UI 组件必须实现的接口
// 它负责告诉 Context：我要怎么被布局和渲染
type Widget interface {
	Mount(*Context) // 组件被挂载时调用，用于声明布局行为
}
