package ui

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"image/color"
)

func MiddleContainer() *container.Split {

	//// 图片视图区 在左边
	//imgViewContent := canvas.NewText("图片视图区 在左边", color.Gray{100})
	//imgViewContainer := container.NewVBox(
	//	container.NewPadded(imgViewContent),
	//	layout.NewSpacer(), // 底部留白
	//)

	// 图片视图区 1 图片显示 上
	imgView1Container := container.NewVBox(
		canvas.NewText("图片视图区 1 图片显示 上", color.Gray{100}),
		layout.NewSpacer(), // 底部留白
	)
	// 图片色值分布区  2 下
	imgView2Container := container.NewVBox(
		canvas.NewText("图片色值分布区  2 下", color.Gray{100}),
		layout.NewSpacer(), // 底部留白
	)
	imgViewContainer := container.NewVSplit(imgView1Container, imgView2Container)
	imgViewContainer.SetOffset(0.9)

	// 图片信息与交互区  在右边
	imgOperatContent := canvas.NewText("图片信息与交互区  在右边", color.Gray{100})
	imgOperatContainer := container.NewVBox(
		container.NewPadded(imgOperatContent),
		layout.NewSpacer(), // 底部留白
	)

	middleContainer := container.NewHSplit(imgViewContainer, imgOperatContainer)
	middleContainer.SetOffset(0.7) // 左侧占比25%

	//contentContent := canvas.NewText("", color.Gray{100})
	//// 中间显示区域：垂直布局，添加边距
	//middleContainer := container.NewVBox(
	//	container.NewPadded(contentTitle),
	//	canvas.NewLine(color.Gray{200}), // 分隔线
	//	container.NewPadded(contentContent),
	//	layout.NewSpacer(), // 底部留白
	//)
	return middleContainer
}
