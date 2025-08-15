package ui

import (
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

	// 图片信息与交互区  在右边
	rich := widget.NewRichTextFromMarkdown(`
# 图片信息
 

---

## 文件信息

* 图片文件名 
<br>
* 图片文件大小
\n\n
* 图片文件md5

* 图片文件路径

* 权限

* 最后修改时间

* 系统原生信息


---

## 图片基础信息

* 图片宽
* 图片高
* 图片dpi
* 图片是多少位的

 
---

## 图片EXIF

* 图片MIME类型
* 图片方向
* 图片相机制造商
* 相机型号
* 拍摄时间
* 焦距
* iso
* 光圈
* 快门

 
---

## 图片色彩属性

* 图片饱和度值
* 图片亮度值
* 图片对比度值
* 图片锐度值
* 图片曝光度值
* 图片色温值
* 图片色调值
* 图片噪点值

 
---

## 图片指纹

* 图片的差异哈希值获取图像的指纹字符串
* 感知哈希算法获取图像的指纹字符串
* 均值哈希值获取图像的指纹字符串

`)
	rich.Resize(fyne.NewSize(260, 0))
	//imgOperateContainer.Add(rich)

	// 文件信息
	imgFileInfoDetail := container.NewWithoutLayout(widget.NewRichTextFromMarkdown(
		"- 文件名\n\n- 文件大小\n\n- 文件md5\n\n- 文件路径\n\n- 权限\n\n- 最后修改时间\n\n- 系统原生信息\n\n\n\n"))
	imgFileInfoDetail.Objects[0].Move(fyne.NewPos(20, 0))
	// 基础信息
	imgBaseInfoDetail := container.NewWithoutLayout(widget.NewRichTextFromMarkdown(`
- 图片宽
- 图片高
- 图片dpi
- 图片是多少位的

`))
	imgBaseInfoDetail.Objects[0].Move(fyne.NewPos(20, 0))
	//拍摄参数
	imgExifInfoDetail := container.NewWithoutLayout(widget.NewRichTextFromMarkdown(`
- 图片MIME类型
- 图片方向
- 图片相机制造商
- 相机型号
- 拍摄时间
- 焦距
- iso
- 光圈
- 快门

`))
	imgExifInfoDetail.Objects[0].Move(fyne.NewPos(20, 0))
	// 色彩属性
	imgColorInfoDetail := container.NewWithoutLayout(widget.NewRichTextFromMarkdown(`
* 图片饱和度值
* 图片亮度值
* 图片对比度值
* 图片锐度值
* 图片曝光度值
* 图片色温值
* 图片色调值
* 图片噪点值
`))
	imgColorInfoDetail.Objects[0].Move(fyne.NewPos(20, 0))
	// 指纹
	imgFingerprintInfoDetail := container.NewWithoutLayout(widget.NewRichTextFromMarkdown(`
* 图片的差异哈希值获取图像的指纹字符串
* 感知哈希算法获取图像的指纹字符串
* 均值哈希值获取图像的指纹字符串
`))
	imgFingerprintInfoDetail.Objects[0].Move(fyne.NewPos(20, 0))

	ac := widget.NewAccordion(

		&widget.AccordionItem{
			Title:  "文件信息",
			Detail: imgFileInfoDetail,
			Open:   true,
		},

		&widget.AccordionItem{
			Title:  "基础信息",
			Detail: imgBaseInfoDetail,
			Open:   true,
		},

		&widget.AccordionItem{
			Title:  "📷 拍摄参数",
			Detail: imgExifInfoDetail,
			Open:   true,
		},
		&widget.AccordionItem{
			Title:  "色彩属性",
			Detail: imgColorInfoDetail,
		},

		&widget.AccordionItem{
			Title:  "指纹",
			Detail: imgFingerprintInfoDetail,
		},
	)
	ac.MultiOpen = true
	//ac.Resize(fyne.NewSize(260, 760))
	scrollContainer := container.NewScroll(ac)
	scrollContainer.SetMinSize(fyne.NewSize(0, 720))

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("图片信息", theme.FileImageIcon(), scrollContainer),
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
