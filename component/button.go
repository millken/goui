package component

type Button struct {
	*BaseComponent
	text string
}

func NewButton(id ...string) *Button {
	btn := &Button{
		BaseComponent: NewBaseComponent(id...),
	}
	// 这里可以添加更多按钮初始化逻辑，比如设置标签文本等
	return btn
}

func (b *Button) SetText(text string) {
	b.text = text
	b.MarkDirty(false) // 文本改变时标记为脏
}
