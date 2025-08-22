package ui

import (
	"MyPicViu/common/logger"
	"MyPicViu/common/utils"
	"MyPicViu/internal/img"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"time"
)

// 目录树的组件事件

// TreeOnSelected 树节点点击事件
func TreeOnSelected() func(uid widget.TreeNodeID) {
	return func(id widget.TreeNodeID) {
		node := dataManager.FindNode(id)
		if node == nil {
			logger.ErrorF("节点被点击: ID=%s (未知节点)\n", id)
			return
		}

		nodeType := "文件"
		if !node.IsFile {
			nodeType = "文件夹"
		}

		logger.DebugF("节点被点击: ID=%s, 名称=%s, 类型=%s 文件路径=%s\n",
			node.ID, node.Name, nodeType, node.FilePath)

		if node.IsFile {
			ShowImg(node.FilePath)
		}

	}

}

func TreeOnBranch() func(uid widget.TreeNodeID) {
	return func(uid widget.TreeNodeID) {
		logger.Debug("点击了箭头打开 uid = ", uid)
		node := FindNode(TreeData, uid)
		if node != nil && !node.IsFile {
			if !node.Expanded {
				node.Expanded = true
			} else {
				node.Expanded = false
			}
		}
	}
}

//// 自定义可点击标签
//type ClickableLabel struct {
//	widget.Label
//}
//
//func (c *ClickableLabel) Tapped(event *fyne.PointEvent) {
//	// 左键点击逻辑（Tapped 只响应左键）
//	println("左键点击了标签")
//	println(event)
//}
//
//func (c *ClickableLabel) TappedSecondary(event *fyne.PointEvent) {
//	println("右键点击了标签")
//	println(event)
//}
//
//// 初始化时必须调用
//func NewClickableLabel(text string) *ClickableLabel {
//	l := &ClickableLabel{}
//	l.SetText(text)
//	l.ExtendBaseWidget(l) // 关键：初始化基础组件
//	return l
//}

var NowImgIsEdit bool = false

func NowImgEdit() {
	NowImgIsEdit = true
	ImgEditSaveButton.RemoveAll()
	ImgEditSaveButton.Add(widget.NewButton("保存", func() {
		logger.Debug("保存")
		SaveImgShow()
	}))
	ImgEditSaveButton.Refresh()
}

func NowImgEditReset() {
	NowImgIsEdit = false
	NewImgLayerReset()
	ImgEditSaveButton.RemoveAll()
	ImgEditSaveButton.Refresh()
}

var NowImgSuffix string = "png"
var NowImgPath string = ""
var NowImgData image.Image
var NowImgWidth binding.ExternalInt
var NowImgHeight binding.ExternalInt
var WidthEntry *widget.Entry
var HeightEntry *widget.Entry

func ImgOriginalSize() fyne.Size {
	dx := ImgViewContainer.Size().Width
	return fyne.NewSize(
		dx,
		700,
	)
}

// ShowImg 打开并显示图片
func ShowImg(filePath string) {

	logger.Debug("show img")

	if NowImgIsEdit {
		logger.Debug("当前正常编辑 ", NowImgPath)

		msg := fmt.Sprintf("当前编辑了图片 %s , 是否另保存", NowImgPath)

		editSaveDialog := dialog.NewConfirm(fmt.Sprintf("当前编辑了图片%s", NowImgPath), // 标题
			msg, // 提示内容
			func(confirmed bool) { // 回调函数
				if confirmed {
					logger.Debug("用户选择：是")
					SaveImgShow()

				} else {
					logger.Debug("用户选择：否")
					// 执行取消逻辑
					NowImgEditReset()
					ShowImg(filePath)
				}
			},
			MainWindow, // 父窗口（必填参数）
		)
		editSaveDialog.Resize(fyne.NewSize(300, 150))
		editSaveDialog.Show()
		return
	}

	// 解码图片
	reader, err := os.Open(filePath)
	if err != nil {
		logger.Error("打开图片文件失败", err)
		return
	}

	defer func() {
		_ = reader.Close()
	}()

	NowImgPath = filePath
	suffix, err := utils.GetFileSuffix(filePath)
	if err == nil {
		NowImgSuffix = suffix
	}
	NowImgData, _, err = image.Decode(reader)
	if err != nil {
		logger.Error("读取图片文件失败", err)
		return
	}
	width := NowImgData.Bounds().Dx()
	NowImgWidth = binding.BindInt(&width)
	height := NowImgData.Bounds().Dy()
	NowImgHeight = binding.BindInt(&height)

	ReductionOperateAbility()

	// 创建Fyne图片对象
	imgObj := canvas.NewImageFromImage(NowImgData)
	imgObj.FillMode = canvas.ImageFillContain // 保持比例显示

	// todo 计算图片的适合位置，不要拉伸图
	//dy := imgView2Container.Size().Height
	//
	//if dx > float32(imgData.Bounds().Dx()) {
	//	dx = float32(imgData.Bounds().Dx())
	//}
	//
	//log.Println("比值: ", imgView1Container.Size().Width/dx)
	//
	//dy = float32(imgData.Bounds().Dy()) * (imgView1Container.Size().Width / dx)
	//
	//if dy > 700 {
	//	dx = float32(imgData.Bounds().Dx()) * (700 / dy)
	//	dy = 700
	//}

	originalSize := ImgOriginalSize()
	// 重置缩放
	imgObj.SetMinSize(originalSize)
	scale := 1.0
	ImgViewContainer.RemoveAll()
	background := canvas.NewRasterWithPixels(checkerPattern)
	background.SetMinSize(fyne.NewSize(280, 280))
	ImgViewContainer.Add(background)
	ImgViewContainer.Add(ImgCanvasObject(imgObj, &scale, &originalSize))
	ImgViewContainer.Refresh()

	// 计算颜色分布
	go func() {
		colorData := img.GetClusters(NowImgData)

		// 创建颜色分布条图片（尺寸为800x100）
		barHeight := ImgColorClustersViewContainer.Size().Height
		totalWidth := ImgColorClustersViewContainer.Size().Width
		barImage := img.CreateColorBarImage(colorData, int(totalWidth), int(barHeight))

		// 转换为Fyne可显示的图片对象
		fyne.Do(func() {
			barImageContainer := canvas.NewImageFromImage(barImage)
			barImageContainer.FillMode = canvas.ImageFillOriginal // 保持原始尺寸

			ImgColorClustersViewContainer.RemoveAll()
			ImgColorClustersViewContainer.Add(barImageContainer)
			ImgColorClustersViewContainer.Refresh()
		})
	}()

	// 计算图片信息
	go func() {
		info := &img.ImgInfo{
			FilePath: filePath,
			Width:    NowImgData.Bounds().Dx(),
			Height:   NowImgData.Bounds().Dy(),
		}
		info.GetFileInfo()
		info.GetImgInfo()

		fyne.Do(func() {
			setImgInfoText(info)
		})

	}()

}

func fileSaved(f fyne.URIWriteCloser, w fyne.Window) {
	defer f.Close()

	if showImgData == nil {
		dialog.ShowError(fmt.Errorf("没有可保存的图片数据"), w)
		return
	}

	var err error

	switch NowImgSuffix {
	case "png":
		// 保存  NowImgData 为png
		err = png.Encode(f, showImgData)
	case "jpg":
		// 保存  NowImgData 为jpeg
		err = jpeg.Encode(f, showImgData, &jpeg.Options{Quality: 90})
	}

	// 处理编码错误
	if err != nil {
		dialog.ShowError(fmt.Errorf("保存图片失败: %v", err), w)
		return
	}

	dialog.ShowInformation("保存成功", fmt.Sprintf("图片已保存至:\n%s", f.URI()), w)
	log.Println("Saved to...", f.URI())
}

func SaveImgShow() {
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		if writer == nil {
			log.Println("Cancelled")
			return
		}

		fileSaved(writer, MainWindow)

		NowImgEditReset()

	}, MainWindow)

	dir, nameWithoutExt, _, ext := utils.ParsePath(NowImgPath)
	saveDialog.SetTitleText("另存图片")

	if uri, err := storage.ParseURI(dir); err == nil {
		// 获取目录URI
		dirURI, _ := storage.ListerForURI(uri)
		if dirURI != nil {
			saveDialog.SetLocation(dirURI) // 设置默认打开的目录
		}
	}

	saveDialog.SetFileName(fmt.Sprintf("%s_%s%s", nameWithoutExt, time.Now().Format(utils.TimeNumberTemplate), ext))
	saveDialog.Show()
}
