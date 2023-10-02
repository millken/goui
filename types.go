package goui

type Component interface {
	Init()
	Dimension()
}

type component struct {
	Count  int
	Type   UIType
	Parent *Cell
}

type FuncDimension func(component any, i uint32, dim0, dim1 float32)
