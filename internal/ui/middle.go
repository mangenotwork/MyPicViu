package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
)

func checkerPattern(x, y, _, _ int) color.Color {
	x /= 20
	y /= 20

	if x%2 == y%2 {
		return theme.Color(theme.ColorNameBackground)
	}

	return theme.Color(theme.ColorNameShadow)
}

func MiddleContainer() *container.Split {
	background := canvas.NewRasterWithPixels(checkerPattern)
	background.SetMinSize(fyne.NewSize(280, 280))
	//background := canvas.NewRectangle(color.Black)
	ImgViewContainer.Add(background)
	imgViewContainer := container.NewVSplit(ImgViewContainer, ImgColorClustersViewContainer)
	imgViewContainer.SetOffset(0.9)

	ImgInfoTextContainer.Append(&widget.AccordionItem{
		Title:  "文件信息",
		Detail: imgFileInfoDetail,
		Open:   false,
	})
	ImgInfoTextContainer.Append(&widget.AccordionItem{
		Title:  "基础信息",
		Detail: imgBaseInfoDetail,
		Open:   false,
	})
	ImgInfoTextContainer.Append(&widget.AccordionItem{
		Title:  "色彩属性",
		Detail: imgColorInfoDetail,
		Open:   false,
	})
	ImgInfoTextContainer.Append(&widget.AccordionItem{
		Title:  "📷 拍摄参数",
		Detail: imgExifInfoDetail,
		Open:   false,
	})
	ImgInfoTextContainer.Append(&widget.AccordionItem{
		Title:  "指纹",
		Detail: imgFingerprintInfoDetail,
		Open:   false,
	})
	ImgInfoTextContainer.MultiOpen = true

	//ac.Resize(fyne.NewSize(260, 760))

	ImgInfoTextScrollContainer := container.NewScroll(ImgInfoTextContainer)
	ImgInfoTextScrollContainer.SetMinSize(fyne.NewSize(0, 720))

	ImgOperateImgOperateAbilityContainer.Add(widget.NewLabel("操作 todo"))
	ImgOperateAbilityScrollContainer := container.NewVScroll(ImgOperateImgOperateAbilityContainer)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("图片信息", theme.FileImageIcon(), ImgInfoTextScrollContainer),
		container.NewTabItemWithIcon("图片编辑", theme.ColorPaletteIcon(), ImgOperateAbilityScrollContainer),
	)

	ImgOperateContainer.Add(tabs)

	middleContainer := container.NewHSplit(imgViewContainer, ImgOperateContainer)
	middleContainer.SetOffset(0.75) // 左侧占比25%

	return middleContainer
}
