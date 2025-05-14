package main

import (
	"fmt"

	"github.com/go-text/draw"
	"github.com/go-text/font"
	"github.com/go-text/shaping"
	"github.com/go-text/typesetting"
)

func main() {
	// 创建一个字体源
	source := font.NewSource()

	// 创建一个形状器
	shaper := shaping.NewShaper(source)

	// 创建一个排版器
	typesetter := typesetting.NewTypesetter(shaper, draw.RGBA{})

	// 添加一些文本
	typesetter.AddText("Hello, world!")

	// 进行排版
	typesetter.Typeset()

	// 输出排版结果
	fmt.Println(typesetter.Result())
}
