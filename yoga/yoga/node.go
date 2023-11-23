package yoga

import "unsafe"

type YGNodeFlags struct {
	hasNewLayout        bool  `bitfield:"0"`
	isReferenceBaseline bool  `bitfield:"1"`
	isDirty             bool  `bitfield:"2"`
	nodeType            uint8 `bitfield:"3"`
	measureUsesContext  bool  `bitfield:"4"`
	baselineUsesContext bool  `bitfield:"5"`
	printUsesContext    bool  `bitfield:"6"`
}

type YGNode struct {
	context_ unsafe.Pointer
	flags_   YGNodeFlags
	measure_ struct {
		noContext   YGMeasureFunc
		withContext MeasureWithContextFn
	}
	baseline_ struct {
		noContext   YGBaselineFunc
		withContext BaselineWithContextFn
	}
	print_ struct {
		noContext   YGPrintFunc
		withContext PrintWithContextFn
	}
	dirtied_            YGDirtiedFunc
	style_              YGStyle
	layout_             YGLayout
	lineIndex_          uint32
	owner_              *YGNode
	children_           []*YGNode
	config_             *YGConfig
	resolvedDimensions_ [2]YGValue
}
