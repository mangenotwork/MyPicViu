package ui

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

func MainContent() *container.Split {
	// 主布局：左右分割
	mainContent := container.NewHSplit(LeftContainer(), MiddleContainer())
	mainContent.SetOffset(0.2) // 左侧占比20%
	return mainContent
}

var contentTitle = canvas.NewText("请选择左侧目录项", color.Black)
