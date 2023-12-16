package ui

import "github.com/millken/yoga"

type Widget interface {
	Node() *yoga.Node
}
