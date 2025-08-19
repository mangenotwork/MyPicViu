package ui

import (
	"MyPicViu/common/logger"
	"MyPicViu/internal/img"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
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
		//return theme.Color(theme.ColorNameBackground)
		return color.RGBA{R: 128, G: 128, B: 128, A: 255}
	}

	// return theme.Color(theme.ColorNameShadow)
	return color.RGBA{R: 192, G: 192, B: 192, A: 255}
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

	ImgOperateImgOperateAbilityContainer.Add(layout.NewSpacer())
	SaturationOperateAbility()
	SaturationOperateAbility() // todo
	SaturationOperateAbility() // todo
	SaturationOperateAbility() // todo
	SaturationOperateAbility() // todo
	SaturationOperateAbility() // todo

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

// SaturationOperateAbility 调整图片饱和度
func SaturationOperateAbility() {
	// -1.0  1.0
	value := 0.2
	data := binding.BindFloat(&value)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "调整饱和度: %0.2f"))
	entry := widget.NewEntryWithData(binding.FloatToString(data))
	floats := container.NewGridWithColumns(2, label, entry)
	slide := widget.NewSliderWithData(-1, 1, data)
	slide.Step = 0.01
	item := container.NewVBox(floats, slide, widget.NewSeparator())

	data.AddListener(binding.NewDataListener(func() {
		// 获取当前值
		currentVal, err := data.Get()
		if err != nil {
			logger.Error("获取值失败, err = ", err)
			return // 处理错误（如需要）
		}

		// 在这里添加值变化后的操作（例如更新图像饱和度）
		logger.Debug("值的变化 ： ", currentVal)
		if NowImgData != nil {
			NowImgData = img.SetImageSaturation(NowImgData, currentVal)

			imgObj := canvas.NewImageFromImage(NowImgData)
			imgObj.FillMode = canvas.ImageFillContain // 保持比例显示
			dx := ImgViewContainer.Size().Width
			originalSize := fyne.NewSize(
				dx,
				700,
			)
			// 重置缩放
			imgObj.SetMinSize(originalSize)
			scale := 1.0
			ImgViewContainer.RemoveAll()
			background := canvas.NewRasterWithPixels(checkerPattern)
			background.SetMinSize(fyne.NewSize(280, 280))
			ImgViewContainer.Add(background)
			ImgViewContainer.Add(ImgCanvasObject(imgObj, &scale, &originalSize))
			ImgViewContainer.Refresh()

		}

	}))

	ImgOperateImgOperateAbilityContainer.Add(item)
}
