package ui

import (
	"MyPicViu/internal/db"
	"MyPicViu/internal/img"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"log"
	"os"
)

var tree, dataManager = db.CreateCustomTree(db.TreeData)

func LeftContainer() *fyne.Container {

	// 树节点点击事件
	tree.OnSelected = func(id widget.TreeNodeID) {
		node := dataManager.FindNode(id)
		if node == nil {
			fmt.Printf("节点被点击: ID=%s (未知节点)\n", id)
			return
		}

		nodeType := "文件"
		if !node.IsFile {
			nodeType = "文件夹"
		}
		fmt.Printf("节点被点击: ID=%s, 名称=%s, 类型=%s 文件路径=%s\n",
			node.ID, node.Name, nodeType, node.FilePath)

		// 解码图片
		reader, err := os.Open(node.FilePath)
		if err != nil {
			log.Println("打开图片文件失败", err)
			return
		}
		imgData, _, err := image.Decode(reader)
		if err != nil {
			log.Println("读取图片文件失败", err)
			return
		}

		// 创建Fyne图片对象
		imgObj := canvas.NewImageFromImage(imgData)
		imgObj.FillMode = canvas.ImageFillContain // 保持比例显示

		dx := imgView1Container.Size().Width

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
		imgView1Container.RemoveAll()
		background := canvas.NewRectangle(color.Black)
		background.SetMinSize(fyne.NewSize(0, 0))
		imgView1Container.Add(background)
		imgView1Container.Add(createContent(imgObj, &scale, &originalSize))
		imgView1Container.Refresh()

		go func() {
			colorData := img.GetClusters(imgData)

			// 创建颜色分布条图片（尺寸为800x100）
			barHeight := imgView2Container.Size().Height
			totalWidth := imgView2Container.Size().Width
			barImage := createColorBarImage(colorData, int(totalWidth), int(barHeight))

			// 转换为Fyne可显示的图片对象
			fyneImage := canvas.NewImageFromImage(barImage)
			fyneImage.FillMode = canvas.ImageFillOriginal // 保持原始尺寸

			imgView2Container.RemoveAll()
			imgView2Container.Add(fyneImage)
			imgView2Container.Refresh()
		}()

	}

	return container.New(
		layout.NewStackLayout(),
		tree,
	)
}

// ColorRatio 定义颜色和对应的比例
type ColorRatio struct {
	Color color.Color
	Ratio float64 // 比例值，总和应接近1
}

// createColorBarImage 生成颜色分布条图片
func createColorBarImage(colors []img.ColorCluster, width, height int) image.Image {
	// 创建一个RGBA图像
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 计算每个颜色块的宽度并绘制
	currentX := 0
	for _, item := range colors {
		log.Println("item.Percent = ", item.Percent/100)
		// 根据比例计算当前颜色块的宽度
		segmentWidth := int(item.Percent / 100 * float64(width))
		if segmentWidth <= 0 {
			continue // 跳过比例为0的颜色
		}

		// 定义当前颜色块的矩形区域
		rect := image.Rect(
			currentX,
			0,
			currentX+segmentWidth,
			height,
		)

		// 填充颜色
		draw.Draw(img, rect, &image.Uniform{item.Color}, image.Point{}, draw.Src)

		// 更新X坐标，准备绘制下一个颜色块
		currentX += segmentWidth
	}

	// 处理最后可能的像素偏差（确保填满整个宽度）
	if currentX < width {
		rect := image.Rect(currentX, 0, width, height)
		draw.Draw(img, rect, &image.Uniform{colors[len(colors)-1].Color}, image.Point{}, draw.Src)
	}

	return img
}

// 创建包含图片和控制按钮的内容
func createContent(img *canvas.Image, scale *float64, originalSize *fyne.Size) fyne.CanvasObject {
	// 放大按钮
	zoomInBtn := widget.NewButton("放大", func() {
		if img == nil {
			return
		}
		*scale += 0.2
		updateImageSize(img, scale, originalSize)
	})

	// 缩小按钮
	zoomOutBtn := widget.NewButton("缩小", func() {
		if img == nil {
			return
		}
		*scale -= 0.2
		if *scale < 0.2 { // 限制最小缩放
			*scale = 0.2
		}
		updateImageSize(img, scale, originalSize)
	})

	// 重置按钮
	resetBtn := widget.NewButton("重置大小", func() {
		if img == nil {
			return
		}
		*scale = 1.0
		updateImageSize(img, scale, originalSize)
	})

	// 按钮容器
	controls := container.NewHBox(
		//openNewBtn,
		zoomInBtn,
		zoomOutBtn,
		resetBtn,
	)

	// 图片容器（使用滚动容器，方便查看大图）
	scrollContainer := container.NewScroll(img)
	scrollContainer.SetMinSize(fyne.NewSize(originalSize.Width, originalSize.Height))

	// 主容器
	return container.NewVBox(
		controls,
		scrollContainer,
	)
}

// 更新图片大小
func updateImageSize(img *canvas.Image, scale *float64, originalSize *fyne.Size) {
	newWidth := float32(*scale) * originalSize.Width
	newHeight := float32(*scale) * originalSize.Height
	img.SetMinSize(fyne.NewSize(newWidth, newHeight))
	img.Refresh()
}

// ColorSegment 定义颜色分布片段
type ColorSegment struct {
	Color color.Color // 片段颜色
	Ratio float64     // 占比(0-1)
	Label string      // 标签
}
