package goui

// dispatchEvent 遍历组件树进行命中测试与冒泡分发（不依赖具体组件实现）。
// 策略：
// 1) 鼠标事件：按子组件逆序（后添加在上）做命中测试与递归分发；若未处理，最后交给当前组件。
// 2) 其它事件：从子到父，只要某个组件 HandleEvent 返回 true 即停止。
func dispatchEvent(root Component, ev Event) bool {
	switch e := ev.(type) {
	case MouseEvent:
		return dispatchMouse(root, e, 0, 0)
	default:
		// 非鼠标事件：深度优先，子优先
		for i := len(root.GetChildren()) - 1; i >= 0; i-- {
			if dispatchEvent(root.GetChildren()[i], ev) {
				return true
			}
		}
		return root.HandleEvent(ev)
	}
}

func dispatchMouse(c Component, e MouseEvent, ox, oy float32) bool {
	// 当前组件的绝对位置
	l := c.GetComputedLayout()
	cx := ox + l.Left
	cy := oy + l.Top

	// 先子组件（逆序）
	children := c.GetChildren()
	for i := len(children) - 1; i >= 0; i-- {
		ch := children[i]
		cl := ch.GetComputedLayout()
		chx := cx + cl.Left
		chy := cy + cl.Top
		if pointInRect(e.X, e.Y, chx, chy, cl.Width, cl.Height) {
			if dispatchMouse(ch, e, cx, cy) { // 注意下一层偏移基于当前组件的绝对坐标
				return true
			}
			if ch.HandleEvent(e) {
				return true
			}
		}
	}
	// 再当前组件
	if pointInRect(e.X, e.Y, cx, cy, l.Width, l.Height) {
		if c.HandleEvent(e) {
			return true
		}
	}
	return false
}

func pointInRect(px, py, x, y, w, h float32) bool {
	return px >= x && px <= x+w && py >= y && py <= y+h
}
