package goui

import "github.com/millken/yoga"

// ApplyStyle 将 Style 的部分属性映射到 Yoga 节点。
// 说明：当前绑定仅使用已在工程中验证存在的 API（SetWidth/SetHeight/SetDisplay/SetFlexDirection/SetAlignItems/SetJustifyContent）。
// 百分比与 auto 的设置预留，待确认 yoga 绑定方法名（通常为 SetWidthPercent/SetWidthAuto 等）后再补充。
func ApplyStyle(n *yoga.Node, s Style) {
	if n == nil {
		return
	}
	// 尺寸
	if s.Width != nil {
		switch s.Width.Unit {
		case UnitPoint:
			n.SetWidth(s.Width.Value)
		case UnitPercent:
			// Yoga 常见 API：SetWidthPercent
			n.SetWidthPercent(s.Width.Value)
		case UnitAuto:
			// Yoga 常见 API：SetWidthAuto
			n.SetWidthAuto()
		}
	}
	if s.Height != nil {
		switch s.Height.Unit {
		case UnitPoint:
			n.SetHeight(s.Height.Value)
		case UnitPercent:
			n.SetHeightPercent(s.Height.Value)
		case UnitAuto:
			n.SetHeightAuto()
		}
	}

	// Flex 容器设置
	if s.Display != nil {
		n.SetDisplay(*s.Display)
	}
	if s.FlexDirection != nil {
		n.SetFlexDirection(*s.FlexDirection)
	}
	if s.AlignItems != nil {
		n.SetAlignItems(*s.AlignItems)
	}
	if s.JustifyContent != nil {
		n.SetJustifyContent(*s.JustifyContent)
	}

	if s.Flex != nil {
		n.SetFlex(*s.Flex)
	}
	if s.FlexGrow != nil {
		n.SetFlexGrow(*s.FlexGrow)
	}
	if s.FlexShrink != nil {
		n.SetFlexShrink(*s.FlexShrink)
	}

	// Flex basis
	if s.FlexBasis != nil {
		switch s.FlexBasis.Unit {
		case UnitPoint:
			n.SetFlexBasis(s.FlexBasis.Value)
		case UnitPercent:
			n.SetFlexBasisPercent(s.FlexBasis.Value)
		case UnitAuto:
			n.SetFlexBasisAuto()
		}
	}

	// 位置类型与偏移
	if s.Position != nil {
		n.SetPositionType(*s.Position)
	}
	applyPosition(n, yoga.EdgeLeft, s.Left)
	applyPosition(n, yoga.EdgeTop, s.Top)
	applyPosition(n, yoga.EdgeRight, s.Right)
	applyPosition(n, yoga.EdgeBottom, s.Bottom)

	// 边距与内边距
	applyEdges(n, s.Margin, true)
	applyEdges(n, s.Padding, false)
}

// Edges 构造器
func EdgesAll(d *Dimension) *Edges           { return &Edges{All: d} }
func EdgesHV(h, v *Dimension) *Edges         { return &Edges{Horizontal: h, Vertical: v} }
func EdgesLTRB(l, t, r, b *Dimension) *Edges { return &Edges{Left: l, Top: t, Right: r, Bottom: b} }

func applyPosition(n *yoga.Node, edge yoga.Edge, d *Dimension) {
	if d == nil {
		return
	}
	switch d.Unit {
	case UnitPoint:
		n.SetPosition(edge, d.Value)
	case UnitPercent:
		n.SetPositionPercent(edge, d.Value)
	}
}

func applyEdges(n *yoga.Node, e *Edges, isMargin bool) {
	if e == nil {
		return
	}
	left := pickDim(e.Left, e.Horizontal, e.All)
	right := pickDim(e.Right, e.Horizontal, e.All)
	top := pickDim(e.Top, e.Vertical, e.All)
	bottom := pickDim(e.Bottom, e.Vertical, e.All)

	if isMargin {
		applyEdgeMargin(n, yoga.EdgeLeft, left)
		applyEdgeMargin(n, yoga.EdgeRight, right)
		applyEdgeMargin(n, yoga.EdgeTop, top)
		applyEdgeMargin(n, yoga.EdgeBottom, bottom)
	} else {
		applyEdgePadding(n, yoga.EdgeLeft, left)
		applyEdgePadding(n, yoga.EdgeRight, right)
		applyEdgePadding(n, yoga.EdgeTop, top)
		applyEdgePadding(n, yoga.EdgeBottom, bottom)
	}
}

func pickDim(primary, axis, all *Dimension) *Dimension {
	if primary != nil {
		return primary
	}
	if axis != nil {
		return axis
	}
	return all
}

func applyEdgeMargin(n *yoga.Node, edge yoga.Edge, d *Dimension) {
	if d == nil {
		return
	}
	switch d.Unit {
	case UnitPoint:
		n.SetMargin(edge, d.Value)
	case UnitPercent:
		n.SetMarginPercent(edge, d.Value)
	case UnitAuto:
		n.SetMarginAuto(edge)
	}
}

func applyEdgePadding(n *yoga.Node, edge yoga.Edge, d *Dimension) {
	if d == nil {
		return
	}
	switch d.Unit {
	case UnitPoint:
		n.SetPadding(edge, d.Value)
	case UnitPercent:
		n.SetPaddingPercent(edge, d.Value)
	}
}
