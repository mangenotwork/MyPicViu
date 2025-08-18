package ui

import (
	"MyPicViu/common/logger"
	"MyPicViu/internal/img"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"image"
	"image/color"
	"os"
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

		// 解码图片
		reader, err := os.Open(node.FilePath)
		if err != nil {
			logger.Error("打开图片文件失败", err)
			return
		}

		defer func() {
			_ = reader.Close()
		}()

		imgData, _, err := image.Decode(reader)
		if err != nil {
			logger.Error("读取图片文件失败", err)
			return
		}

		// 创建Fyne图片对象
		imgObj := canvas.NewImageFromImage(imgData)
		imgObj.FillMode = canvas.ImageFillContain // 保持比例显示

		dx := ImgViewContainer.Size().Width

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

		originalSize := fyne.NewSize(
			dx,
			700,
		)

		// 重置缩放
		imgObj.SetMinSize(originalSize)
		scale := 1.0
		ImgViewContainer.RemoveAll()
		background := canvas.NewRectangle(color.Black)
		background.SetMinSize(fyne.NewSize(0, 0))
		ImgViewContainer.Add(background)
		ImgViewContainer.Add(createContent(imgObj, &scale, &originalSize))
		ImgViewContainer.Refresh()

		// 计算颜色分布
		go func() {
			colorData := img.GetClusters(imgData)

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
				FilePath: node.FilePath,
				Width:    imgData.Bounds().Dx(),
				Height:   imgData.Bounds().Dy(),
			}
			info.GetFileInfo()
			info.GetImgInfo()

			fyne.Do(func() {
				setImgInfoText(info)
			})

		}()

	}

}
