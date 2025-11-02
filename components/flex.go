package components

import "github.com/millken/goui"

var _ goui.Component = (*Flex)(nil)

type Flex struct {
	Base
}

func NewFlex() *Flex {
	return &Flex{}
}
