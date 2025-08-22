package ui

import (
	"MyPicViu/common/logger"
	"MyPicViu/internal/img"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image"
	"sort"
)

type AdjustType string

const (
	AdjustTypeDirection        AdjustType = "Direction"        // 调整图片方向，镜像，旋转
	AdjustTypeHorizontalMirror AdjustType = "HorizontalMirror" // 水平镜像
	AdjustTypeVerticalMirror   AdjustType = "VerticalMirror"   // 垂直镜像
	AdjustTypeReduction        AdjustType = "Reduction"        // 缩放图片
	AdjustTypeSaturation       AdjustType = "Saturation"       // 调整图片饱和度
	AdjustTypeBrightness       AdjustType = "Brightness"       // 调整图片亮度值
	AdjustTypeContrast         AdjustType = "Contrast"         // 图片对比度值
	AdjustTypeSharpness        AdjustType = "Sharpness"        // 图片锐度值
	AdjustTypeExposure         AdjustType = "Exposure"         // 图片曝光度值
	AdjustTypeTemperature      AdjustType = "Temperature"      // 图片色温值
	AdjustTypeHue              AdjustType = "Hue"              // 图片色调值
	AdjustTypeNoise            AdjustType = "Noise"            // 图片噪点值
)

type ImgLayer struct {
	Type     AdjustType  // 图层调整类型
	Sequence int         // 图层序列
	Value    interface{} // 多个值
}

type ReductionValue struct {
	Width  int
	Height int
}

func NewImgLayerDirection(value int) {
	newImgLayer := &ImgLayer{
		Type:  AdjustTypeDirection,
		Value: value,
	}
	ImgLayerAdd(newImgLayer)
	ImgLayerRefresh()
}

func NewImgLayerHorizontalMirror() {
	newImgLayer := &ImgLayer{
		Type:  AdjustTypeHorizontalMirror,
		Value: true,
	}
	ImgLayerAdd(newImgLayer)
	ImgLayerRefresh()
}

func NewImgLayerVerticalMirror() {
	newImgLayer := &ImgLayer{
		Type:  AdjustTypeVerticalMirror,
		Value: true,
	}
	ImgLayerAdd(newImgLayer)
	ImgLayerRefresh()
}

func NewImgLayerSaturation(value float64) {
	newImgLayer := &ImgLayer{
		Type:  AdjustTypeSaturation,
		Value: value,
	}
	ImgLayerAdd(newImgLayer)
	ImgLayerRefresh()
}

func NewImgLayerReduction(width, height int) {
	newImgLayer := &ImgLayer{
		Type:  AdjustTypeReduction,
		Value: ReductionValue{width, height},
	}
	ImgLayerAdd(newImgLayer)
	ImgLayerRefresh()
}

func NewImgLayerBrightness(value float64) {
	newImgLayer := &ImgLayer{
		Type:  AdjustTypeBrightness,
		Value: value,
	}
	ImgLayerAdd(newImgLayer)
	ImgLayerRefresh()
}

func NewImgLayerContrast(value float64) {
	newImgLayer := &ImgLayer{
		Type:  AdjustTypeContrast,
		Value: value,
	}
	ImgLayerAdd(newImgLayer)
	ImgLayerRefresh()
}

func NewImgLayerSharpness(value float64) {
	newImgLayer := &ImgLayer{
		Type:  AdjustTypeSharpness,
		Value: value,
	}
	ImgLayerAdd(newImgLayer)
	ImgLayerRefresh()
}

func NewImgLayerExposure(value float64) {
	newImgLayer := &ImgLayer{
		Type:  AdjustTypeExposure,
		Value: value,
	}
	ImgLayerAdd(newImgLayer)
	ImgLayerRefresh()
}

func NewImgLayerTemperature(value float64) {
	newImgLayer := &ImgLayer{
		Type:  AdjustTypeTemperature,
		Value: value,
	}
	ImgLayerAdd(newImgLayer)
	ImgLayerRefresh()
}

func NewImgLayerHue(value float64) {
	newImgLayer := &ImgLayer{
		Type:  AdjustTypeHue,
		Value: value,
	}
	ImgLayerAdd(newImgLayer)
	ImgLayerRefresh()
}

func NewImgLayerNoise(value float64) {
	newImgLayer := &ImgLayer{
		Type:  AdjustTypeNoise,
		Value: value,
	}
	ImgLayerAdd(newImgLayer)
	ImgLayerRefresh()
}

var NewImgLayer []*ImgLayer

func NewImgLayerReset() {
	NewImgLayer = make([]*ImgLayer, 0)
}

func ImgLayerAdd(layer *ImgLayer) {
	if NewImgLayer == nil {
		NewImgLayer = make([]*ImgLayer, 0)
	}

	for _, v := range NewImgLayer {
		if v.Type == layer.Type {
			switch v.Type {
			case AdjustTypeHorizontalMirror, AdjustTypeVerticalMirror:
				v.Value = !v.Value.(bool)
			case AdjustTypeDirection:
				rse := v.Value.(int) + layer.Value.(int)
				if rse >= 360 {
					rse -= 360
				}
				v.Value = rse
			default:
				v.Value = layer.Value
			}
			return
		}
	}

	layer.Sequence = len(NewImgLayer) + 1
	NewImgLayer = append(NewImgLayer, layer)
}

var showImgData image.Image

func ImgLayerRefresh() {

	if NewImgLayer == nil || len(NewImgLayer) == 0 {
		logger.Debug("没有图层更新")
		return
	}

	logger.Debug("ImgLayer 数量: ", len(NewImgLayer))

	sort.Slice(NewImgLayer, func(i, j int) bool {
		if NewImgLayer[i].Sequence > NewImgLayer[j].Sequence {
			return true
		}
		return false
	})

	NowImgEdit()

	showImgData = NowImgData

	for _, v := range NewImgLayer {
		switch v.Type {
		case AdjustTypeDirection:
			value, ok := v.Value.(int)
			if !ok {
				logger.Panic("参数解析错误")
			}
			if value == 90 {
				showImgData = img.Rotate90Clockwise(showImgData)
			} else if value == 180 {
				showImgData = img.Rotate180Clockwise(showImgData)
			} else if value == 270 {
				showImgData = img.Rotate270Clockwise(showImgData)
			}

		case AdjustTypeHorizontalMirror:
			value, ok := v.Value.(bool)
			if !ok {
				logger.Panic("参数解析错误")
			}
			if value {
				showImgData = img.HorizontalMirror(showImgData)
			}

		case AdjustTypeVerticalMirror:
			value, ok := v.Value.(bool)
			if !ok {
				logger.Panic("参数解析错误")
			}
			if value {
				showImgData = img.VerticalMirror(showImgData)
			}

		case AdjustTypeReduction:
			value, ok := v.Value.(ReductionValue)
			if !ok {
				logger.Panic("参数解析错误")
			}
			showImgData = img.DrawCatmullRom(showImgData, value.Width, value.Height)

		case AdjustTypeSaturation:
			value, ok := v.Value.(float64)
			if !ok {
				logger.Panic("参数解析错误")
			}
			showImgData = img.SetImageSaturation(showImgData, value)

		case AdjustTypeBrightness:
			value, ok := v.Value.(float64)
			if !ok {
				logger.Panic("参数解析错误")
			}
			showImgData = img.SetImageBrightness(showImgData, value)

		case AdjustTypeContrast:
			value, ok := v.Value.(float64)
			if !ok {
				logger.Panic("参数解析错误")
			}
			showImgData = img.SetImageContrast(showImgData, value)

		case AdjustTypeSharpness:
			value, ok := v.Value.(float64)
			if !ok {
				logger.Panic("参数解析错误")
			}
			showImgData = img.SetImageSharpness(showImgData, value)

		case AdjustTypeExposure:
			value, ok := v.Value.(float64)
			if !ok {
				logger.Panic("参数解析错误")
			}
			showImgData = img.SetImageExposure(showImgData, value)

		case AdjustTypeTemperature:
			value, ok := v.Value.(float64)
			if !ok {
				logger.Panic("参数解析错误")
			}
			showImgData = img.SetImageTemperature(showImgData, value)

		case AdjustTypeHue:
			value, ok := v.Value.(float64)
			if !ok {
				logger.Panic("参数解析错误")
			}
			showImgData = img.SetImageHue(showImgData, value)

		case AdjustTypeNoise:
			value, ok := v.Value.(float64)
			if !ok {
				logger.Panic("参数解析错误")
			}
			showImgData = img.SetImageNoise(showImgData, value)

		}
	}

	imgObj := canvas.NewImageFromImage(showImgData)
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
