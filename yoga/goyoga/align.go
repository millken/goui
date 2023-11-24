package goyoga

func resolveChildAlignment(node *YGNode, child *YGNode) YGAlign {
	if child.getStyle().alignSelf() != YGAlignAuto {
		return child.getStyle().alignSelf()
	}
	return node.getStyle().alignItems()
}
