package goui

import (
	"image"
)

// MouseEventDispatcher 鼠标事件分发器
type MouseEventDispatcher struct{}

func (d *MouseEventDispatcher) CanDispatch(event Event) bool {
	_, ok := event.(*MouseEvent)
	return ok
}

func (d *MouseEventDispatcher) Dispatch(event Event, target Component) bool {
	mouseEvent, ok := event.(*MouseEvent)
	if !ok {
		return false
	}

	// 检查组件是否可见和启用
	if !target.IsVisible() || !target.IsEnabled() {
		return false
	}

	// 检查鼠标位置是否在组件边界内
	bounds := target.Bounds()
	point := image.Pt(int(mouseEvent.X), int(mouseEvent.Y))
	if point.X < bounds.Min.X || point.X >= bounds.Max.X || point.Y < bounds.Min.Y || point.Y >= bounds.Max.Y {
		return false
	}

	// 尝试让组件处理事件
	if handler, ok := target.(EventListener); ok {
		return handler.HandleEvent(mouseEvent)
	}

	return false
}

// KeyEventDispatcher 键盘事件分发器
type KeyEventDispatcher struct{}

func (d *KeyEventDispatcher) CanDispatch(event Event) bool {
	_, ok := event.(*KeyEvent)
	return ok
}

func (d *KeyEventDispatcher) Dispatch(event Event, target Component) bool {
	keyEvent, ok := event.(*KeyEvent)
	if !ok {
		return false
	}

	// 检查组件是否可见和启用
	if !target.IsVisible() || !target.IsEnabled() {
		return false
	}

	// 键盘事件通常需要有焦点的组件处理
	if !d.hasFocus(target) {
		return false
	}

	if handler, ok := target.(EventListener); ok {
		return handler.HandleEvent(keyEvent)
	}

	return false
}

func (d *KeyEventDispatcher) hasFocus(component Component) bool {
	// TODO: 实现焦点检查
	// 现在简单返回true，后续实现焦点管理
	return true
}

// WindowEventDispatcher 窗口事件分发器
type WindowEventDispatcher struct{}

func (d *WindowEventDispatcher) CanDispatch(event Event) bool {
	_, ok := event.(*WindowEvent)
	return ok
}

func (d *WindowEventDispatcher) Dispatch(event Event, target Component) bool {
	windowEvent, ok := event.(*WindowEvent)
	if !ok {
		return false
	}

	// 窗口事件通常只对根组件有意义
	if target.Parent() != nil {
		return false
	}

	if handler, ok := target.(EventListener); ok {
		return handler.HandleEvent(windowEvent)
	}

	return false
}

// TouchEventDispatcher 触摸事件分发器
type TouchEventDispatcher struct{}

func (d *TouchEventDispatcher) CanDispatch(event Event) bool {
	_, ok := event.(*TouchEvent)
	return ok
}

func (d *TouchEventDispatcher) Dispatch(event Event, target Component) bool {
	touchEvent, ok := event.(*TouchEvent)
	if !ok {
		return false
	}

	// 检查组件是否可见和启用
	if !target.IsVisible() || !target.IsEnabled() {
		return false
	}

	// 检查触摸位置是否在组件边界内
	bounds := target.Bounds()
	point := image.Pt(int(touchEvent.X), int(touchEvent.Y))
	if point.X < bounds.Min.X || point.X >= bounds.Max.X || point.Y < bounds.Min.Y || point.Y >= bounds.Max.Y {
		return false
	}

	if handler, ok := target.(EventListener); ok {
		return handler.HandleEvent(touchEvent)
	}

	return false
}

// FocusEventDispatcher 焦点事件分发器
type FocusEventDispatcher struct{}

func (d *FocusEventDispatcher) CanDispatch(event Event) bool {
	_, ok := event.(*FocusEvent)
	return ok
}

func (d *FocusEventDispatcher) Dispatch(event Event, target Component) bool {
	focusEvent, ok := event.(*FocusEvent)
	if !ok {
		return false
	}

	// 焦点事件对可见的组件有意义
	if !target.IsVisible() {
		return false
	}

	if handler, ok := target.(EventListener); ok {
		return handler.HandleEvent(focusEvent)
	}

	return false
}

// ScrollEventDispatcher 滚动事件分发器
type ScrollEventDispatcher struct{}

func (d *ScrollEventDispatcher) CanDispatch(event Event) bool {
	_, ok := event.(*ScrollEvent)
	return ok
}

func (d *ScrollEventDispatcher) Dispatch(event Event, target Component) bool {
	scrollEvent, ok := event.(*ScrollEvent)
	if !ok {
		return false
	}

	// 检查组件是否可见和启用
	if !target.IsVisible() || !target.IsEnabled() {
		return false
	}

	// 检查滚动位置是否在组件边界内
	bounds := target.Bounds()
	point := image.Pt(int(scrollEvent.X), int(scrollEvent.Y))
	if point.X < bounds.Min.X || point.X >= bounds.Max.X || point.Y < bounds.Min.Y || point.Y >= bounds.Max.Y {
		return false
	}

	if handler, ok := target.(EventListener); ok {
		return handler.HandleEvent(scrollEvent)
	}

	return false
}