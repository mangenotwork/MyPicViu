package ui

import (
	"MyPicViu/common/logger"
	"MyPicViu/common/utils"
	"MyPicViu/internal/img"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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

	ImgInfoTextContainer.Open(0)
	ImgInfoTextContainer.Open(1)
	ImgInfoTextContainer.Open(2)
}
