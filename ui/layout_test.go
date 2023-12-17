package ui

import (
	"testing"

	"github.com/millken/yoga"
)

func TestLayout1(t *testing.T) {
	cfg := yoga.ConfigNew()
	cfg.SetShouldPrintTree(true)
	root := yoga.NewNodeWithConfig(cfg)
	root.StyleSetFlexDirection(yoga.FlexDirectionRow)
	root.StyleSetWidth(800)
	root.StyleSetHeight(600)
	node1 := yoga.NewNode()
	node1.StyleSetWidthPercent(50)
	root.InsertChild(node1, 0)
	node2 := yoga.NewNode()
	node2.StyleSetWidthPercent(50)
	root.InsertChild(node2, 1)
	//root.StyleSetPositionType(yoga.PositionTypeAbsolute)
	yoga.CalculateLayout(root, yoga.Undefined, yoga.Undefined, yoga.DirectionLTR)
}
