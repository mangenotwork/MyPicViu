package ui

import (
	"MyPicViu/common/logger"
	"MyPicViu/common/utils"
	"MyPicViu/internal/img"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
)

// 图片视图区 1 图片显示 上
var imgView1Container = container.NewStack(
//layout.NewSpacer(), // 底部留白
)

// 图片色值分布区  2 下
var imgView2Container = container.NewVBox(
	canvas.NewText("图片色值分布区  2 下", color.Gray{100}),
	layout.NewSpacer(), // 底部留白
)

// 图片信息与交互区  在右边
var imgOperateContainer = container.NewVBox()

// 图片信息文本区
var ImgInfoTextContainer = widget.NewAccordion()

func MiddleContainer() *container.Split {

	//// 图片视图区 在左边

	// 图片视图区 1 图片显示 上
	//imgView1Container := container.NewVBox(
	//	canvas.NewText("图片视图区 1 图片显示 上", color.Gray{100}),
	//	layout.NewSpacer(), // 底部留白
	//)
	background := canvas.NewRectangle(color.Black)
	//background.SetMinSize(fyne.NewSize(0, 0))
	imgView1Container.Add(background)
	//// 图片色值分布区  2 下
	//imgView2Container := container.NewVBox(
	//	canvas.NewText("图片色值分布区  2 下", color.Gray{100}),
	//	layout.NewSpacer(), // 底部留白
	//)

	imgViewContainer := container.NewVSplit(imgView1Container, imgView2Container)
	imgViewContainer.SetOffset(0.9)

	ImgInfoTextContainer.Append(&widget.AccordionItem{
		Title:  "文件信息",
		Detail: imgFileInfoDetail,
		Open:   true,
	})
	ImgInfoTextContainer.Append(&widget.AccordionItem{
		Title:  "基础信息",
		Detail: imgBaseInfoDetail,
		Open:   true,
	})
	ImgInfoTextContainer.Append(&widget.AccordionItem{
		Title:  "色彩属性",
		Detail: imgColorInfoDetail,
		Open:   true,
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

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("图片信息", theme.FileImageIcon(), ImgInfoTextScrollContainer),
		container.NewTabItemWithIcon("图片编辑", theme.ColorPaletteIcon(), widget.NewLabel("操作 todo")),
	)

	imgOperateContainer.Add(tabs)

	middleContainer := container.NewHSplit(imgViewContainer, imgOperateContainer)
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

var imgFileInfoDetail = container.NewWithoutLayout()
var imgBaseInfoDetail = container.NewWithoutLayout()
var imgExifInfoDetail = container.NewWithoutLayout()
var imgColorInfoDetail = container.NewWithoutLayout()
var imgFingerprintInfoDetail = container.NewWithoutLayout()

func setImgInfoText(info *img.ImgInfo) {

	logger.Debug("图片信息: %+v", info)

	// 文件信息
	imgFileInfoDetail.RemoveAll()
	imgFileInfoDetailStr := fmt.Sprintf(
		"- 文件名: **%s**\n\n- 大小: **%s**\n\n- md5: **%s**\n\n- 路径: **%s**\n\n- 权限: **%s**\n\n- 最后修改时间: **%s**\n\n\n\n",
		info.FileName, utils.SizeFormat(info.FileSize), info.FileMd5, info.FilePath, info.FileMode, info.FileModTime)
	imgFileInfoDetail.Add(widget.NewRichTextFromMarkdown(imgFileInfoDetailStr))
	imgFileInfoDetail.Objects[0].Move(fyne.NewPos(20, 0))
	imgFileInfoDetail.Refresh()

	// 基础信息
	imgBaseInfoDetail.RemoveAll()
	imgBaseInfoDetailStr := fmt.Sprintf("- 宽高: **%dx%d**\n\n- DPI: **todo**\n\n- 位: **todo**\n\n",
		info.Width, info.Height)
	imgBaseInfoDetail.Add(widget.NewRichTextFromMarkdown(imgBaseInfoDetailStr))
	imgBaseInfoDetail.Objects[0].Move(fyne.NewPos(20, 0))
	imgBaseInfoDetail.Refresh()

	//拍摄参数
	imgExifInfoDetail.RemoveAll()
	imgExifInfoDetail.Add(widget.NewRichTextFromMarkdown(`
- MIME类型: **todo**
- 方向: **todo**
- 相机制造商: **todo**
- 相机型号: **todo**
- 拍摄时间: **todo**
- 焦距: **todo**
- iso: **todo**
- 光圈: **todo**
- 快门: **todo**
`))
	imgExifInfoDetail.Objects[0].Move(fyne.NewPos(20, 0))
	imgExifInfoDetail.Refresh()

	// 色彩属性
	imgColorInfoDetail.RemoveAll()
	imgColorInfoDetailStr := fmt.Sprintf(
		"- 饱和度值: **%f**\n\n- 亮度值: **%f**\n\n- 对比度值: **%f**\n\n- 锐度值: **%f**\n\n- 曝光度值: **%v**\n\n- 色温值: **%f**\n\n- 色调值: **%f**\n\n- 噪点值: **%f**\n\n\n\n",
		info.Saturation, info.Brightness, info.Contrast, info.Sharpness, info.Exposure, info.Temperature, info.Hue, info.Noise,
	)
	imgColorInfoDetail.Add(widget.NewRichTextFromMarkdown(imgColorInfoDetailStr))
	imgColorInfoDetail.Objects[0].Move(fyne.NewPos(20, 0))
	imgColorInfoDetail.Refresh()

	// 指纹
	imgFingerprintInfoDetail.RemoveAll()
	imgFingerprintInfoDetailStr := fmt.Sprintf("- 差异哈希值: **%s**\n\n- 感知哈希算: **%s**\n\n- 均值哈希: **%s**\n\n",
		info.DifferenceHash, info.PHash, info.AverageHash)
	imgFingerprintInfoDetail.Add(widget.NewRichTextFromMarkdown(imgFingerprintInfoDetailStr))
	imgFingerprintInfoDetail.Objects[0].Move(fyne.NewPos(20, 0))
	imgFingerprintInfoDetail.Refresh()

}
