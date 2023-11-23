package goyoga

func isRow(flexDirection YGFlexDirection) bool {
	return flexDirection == YGFlexDirectionRow || flexDirection == YGFlexDirectionRowReverse
}

func isColumn(flexDirection YGFlexDirection) bool {
	return flexDirection == YGFlexDirectionColumn || flexDirection == YGFlexDirectionColumnReverse
}

func flexStartEdge(flexDirection YGFlexDirection) YGEdge {
	switch flexDirection {
	case YGFlexDirectionColumn:
		return YGEdgeTop
	case YGFlexDirectionColumnReverse:
		return YGEdgeBottom
	case YGFlexDirectionRow:
		return YGEdgeLeft
	case YGFlexDirectionRowReverse:
		return YGEdgeRight
	}
	panic("Invalid FlexDirection")
}

func flexEndEdge(flexDirection YGFlexDirection) YGEdge {
	switch flexDirection {
	case YGFlexDirectionColumn:
		return YGEdgeBottom
	case YGFlexDirectionColumnReverse:
		return YGEdgeTop
	case YGFlexDirectionRow:
		return YGEdgeRight
	case YGFlexDirectionRowReverse:
		return YGEdgeLeft
	}
	panic("Invalid FlexDirection")
}
