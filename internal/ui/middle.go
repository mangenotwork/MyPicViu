package ui

import (
	"MyPicViu/common/logger"
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

	ImgInfoTextScrollContainer := container.NewScroll(ImgInfoTextContainer)

	logger.Debug("右边框高度: ", MainWindow.Content().Size().Height-40)

	ImgInfoTextScrollContainer.SetMinSize(fyne.NewSize(0, 760))

	ImgOperateImgOperateAbilityContainer.Append(DirectionOperateAbility())
	ReductionOperateAbility()
	ImgOperateImgOperateAbilityContainer.Append(&widget.AccordionItem{
		Title:  "缩放",
		Detail: ReductionOperateAbilityContainer,
		Open:   true,
	})
	ImgOperateImgOperateAbilityContainer.Append(SaturationOperateAbility())
	ImgOperateImgOperateAbilityContainer.Append(BrightnessOperateAbility())
	ImgOperateImgOperateAbilityContainer.Append(ContrastOperateAbility())
	ImgOperateImgOperateAbilityContainer.Append(SharpnessOperateAbility())
	ImgOperateImgOperateAbilityContainer.Append(ExposureOperateAbility())
	ImgOperateImgOperateAbilityContainer.Append(TemperatureOperateAbility())
	ImgOperateImgOperateAbilityContainer.Append(HueOperateAbility())
	ImgOperateImgOperateAbilityContainer.Append(NoiseOperateAbility())

	ImgOperateAbilityScrollContainer := container.NewVScroll(ImgOperateImgOperateAbilityContainer)

	ImgOperateAbilityScrollContainerM := container.NewVSplit(ImgOperateAbilityScrollContainer, ImgEditSaveButton)
	ImgOperateAbilityScrollContainerM.SetOffset(0.93)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("图片信息", theme.FileImageIcon(), ImgInfoTextScrollContainer),
		container.NewTabItemWithIcon("图片编辑", theme.ColorPaletteIcon(), ImgOperateAbilityScrollContainerM),
	)

	ImgOperateContainer.Add(tabs)

	middleContainer := container.NewHSplit(imgViewContainer, ImgOperateContainer)
	middleContainer.SetOffset(0.75) // 左侧占比25%

	return middleContainer
}

// DirectionOperateAbility 调整图片方向，镜像，旋转
func DirectionOperateAbility() *widget.AccordionItem {

	Rotate90ClockwiseBtn := widget.NewButton("旋转90度", func() {
		logger.Debug("旋转90度")
		if NowImgData != nil {
			NewImgLayerDirection(90)
		}
	})

	Rotate180ClockwiseBtn := widget.NewButton("旋转180度", func() {
		logger.Debug("旋转180度")
		if NowImgData != nil {
			NewImgLayerDirection(180)
		}
	})

	Rotate270ClockwiseBtn := widget.NewButton("旋转270度", func() {
		logger.Debug("旋转270度")
		if NowImgData != nil {
			NewImgLayerDirection(270)
		}
	})

	HorizontalMirrorBtn := widget.NewButton("水平镜像", func() {
		logger.Debug("水平镜像")
		if NowImgData != nil {
			NewImgLayerHorizontalMirror()
		}
	})

	VerticalMirrorBtn := widget.NewButton("垂直镜像", func() {
		logger.Debug("垂直镜像")
		if NowImgData != nil {
			NewImgLayerVerticalMirror()
		}
	})

	item := container.New(layout.NewVBoxLayout())
	item.Add(Rotate90ClockwiseBtn)
	item.Add(Rotate180ClockwiseBtn)
	item.Add(Rotate270ClockwiseBtn)
	item.Add(HorizontalMirrorBtn)
	item.Add(VerticalMirrorBtn)
	item.Add(layout.NewSpacer())

	return &widget.AccordionItem{
		Title:  "旋转方向,镜像",
		Detail: item,
		Open:   true,
	}
}

var ReductionOperateAbilityContainer = container.New(layout.NewVBoxLayout())

// ReductionOperateAbility 缩放图片
func ReductionOperateAbility() {
	ReductionOperateAbilityContainer.RemoveAll()
	WidthEntry = widget.NewEntryWithData(binding.IntToString(NowImgWidth))
	widthEntryContainer := widget.NewForm(widget.NewFormItem("宽度:", WidthEntry))
	ReductionOperateAbilityContainer.Add(widthEntryContainer)

	HeightEntry = widget.NewEntryWithData(binding.IntToString(NowImgHeight))
	heightEntryContainer := widget.NewForm(widget.NewFormItem("高度:", HeightEntry))
	ReductionOperateAbilityContainer.Add(heightEntryContainer)

	ReductionOperateSubmitBtn := widget.NewButton("执行调整", func() {
		logger.Debug("执行调整")
		widthVal, err := NowImgWidth.Get()
		if err != nil {
			logger.Error(err)
		}
		logger.Debug("widthVal = ", widthVal)
		heightVal, err := NowImgHeight.Get()
		if err != nil {
			logger.Error(err)
		}
		logger.Debug("heightVal = ", heightVal)

		if NowImgData != nil {
			NewImgLayerReduction(widthVal, heightVal)
		}
	})
	ReductionOperateAbilityContainer.Add(ReductionOperateSubmitBtn)
	ReductionOperateAbilityContainer.Add(layout.NewSpacer())
	ReductionOperateAbilityContainer.Refresh()
}

// SaturationOperateAbility 调整图片饱和度
func SaturationOperateAbility() *widget.AccordionItem {
	// -1.0  1.0
	value := 0.0
	data := binding.BindFloat(&value)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "饱和度: %0.2f"))
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
			NewImgLayerSaturation(currentVal)
		}
	}))

	return &widget.AccordionItem{
		Title:  "调整饱和度",
		Detail: item,
		Open:   true,
	}
}

// BrightnessOperateAbility 调整图片亮度值
func BrightnessOperateAbility() *widget.AccordionItem {
	// -1.0  1.0
	value := 0.0
	data := binding.BindFloat(&value)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "亮度: %0.2f"))
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
		logger.Debug("currentVal = ", currentVal)
		if NowImgData != nil {
			NewImgLayerBrightness(currentVal)
		}
	}))

	return &widget.AccordionItem{
		Title:  "调整亮度",
		Detail: item,
		Open:   true,
	}
}

// ContrastOperateAbility 图片对比度值
func ContrastOperateAbility() *widget.AccordionItem {

	// -1.0  1.0
	value := 0.0
	data := binding.BindFloat(&value)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "对比度: %0.2f"))
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
		logger.Debug("currentVal = ", currentVal)
		if NowImgData != nil {
			NewImgLayerContrast(currentVal)
		}
	}))

	return &widget.AccordionItem{
		Title:  "调整对比度",
		Detail: item,
		Open:   true,
	}
}

// SharpnessOperateAbility 图片锐度值
func SharpnessOperateAbility() *widget.AccordionItem {

	// -1.0  1.0
	value := 0.0
	data := binding.BindFloat(&value)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "锐度: %0.2f"))
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
		logger.Debug("currentVal = ", currentVal)
		if NowImgData != nil {
			NewImgLayerSharpness(currentVal)
		}
	}))

	return &widget.AccordionItem{
		Title:  "调整锐度",
		Detail: item,
		Open:   true,
	}
}

// ExposureOperateAbility 图片曝光度值
func ExposureOperateAbility() *widget.AccordionItem {

	// -1.0  1.0
	value := 0.0
	data := binding.BindFloat(&value)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "曝光度: %0.2f"))
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
		logger.Debug("currentVal = ", currentVal)
		if NowImgData != nil {
			NewImgLayerExposure(currentVal)
		}
	}))

	return &widget.AccordionItem{
		Title:  "调整曝光度",
		Detail: item,
		Open:   true,
	}
}

// TemperatureOperateAbility 图片色温值
func TemperatureOperateAbility() *widget.AccordionItem {

	// -1.0  1.0
	value := 0.0
	data := binding.BindFloat(&value)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "色温: %0.2f"))
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
		logger.Debug("currentVal = ", currentVal)
		if NowImgData != nil {
			NewImgLayerTemperature(currentVal)
		}
	}))

	return &widget.AccordionItem{
		Title:  "调整色温",
		Detail: item,
		Open:   true,
	}
}

// HueOperateAbility 图片色调值
func HueOperateAbility() *widget.AccordionItem {

	// -1.0  1.0
	value := 0.0
	data := binding.BindFloat(&value)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "色调: %0.2f"))
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
		logger.Debug("currentVal = ", currentVal)
		if NowImgData != nil {
			NewImgLayerHue(currentVal)
		}
	}))

	return &widget.AccordionItem{
		Title:  "调整色调",
		Detail: item,
		Open:   true,
	}
}

// NoiseOperateAbility 图片噪点值
func NoiseOperateAbility() *widget.AccordionItem {

	// -1.0  1.0s
	value := 0.0
	data := binding.BindFloat(&value)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "噪点: %0.2f"))
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
		logger.Debug("currentVal = ", currentVal)
		if NowImgData != nil {
			NewImgLayerNoise(currentVal)
		}
	}))

	return &widget.AccordionItem{
		Title:  "调整噪点",
		Detail: item,
		Open:   true,
	}
}
