package ui

import (
	"MyPicViu/common/logger"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
)

var MainApp fyne.App
var MainWindow fyne.Window

var tree *widget.Tree            // 目录树的UI组件
var dataManager *TreeDataManager // 目录树的数据管理

func InitUI() {
	MainApp = app.NewWithID("MyPicViu.2025.0814")
	MainWindow = MainApp.NewWindow("MyPicViu")
	logger.Debug("初始化UI")
	TreeData = buildTreeData()
	tree, dataManager = CreateCustomTree()
	tree.OnSelected = TreeOnSelected()
	tree.OnBranchOpened = TreeOnBranch()
	tree.OnBranchClosed = TreeOnBranch()
}

func MainContent() *container.Split {
	LeftContainerInit()
	// 主布局：左右分割
	mainContent := container.NewHSplit(LeftContainer, MiddleContainer())
	mainContent.SetOffset(0.2) // 左侧占比20%
	return mainContent
}

// LeftContainer 最左边的目录文件列表树
var LeftContainer = container.New(layout.NewStackLayout())

// ImgViewContainer 图片视图区 图片显示 上
var ImgViewContainer = container.NewStack()

// ImgColorClustersViewContainer 图片色值分布区 图片显示 下
var ImgColorClustersViewContainer = container.NewStack()

// OpenInitContainer 初始化打开图片的按钮
var OpenInitContainer = container.New(layout.NewVBoxLayout())

// ImgOperateContainer 图片信息与交互区  在右边
var ImgOperateContainer = container.NewVBox()

// ImgInfoTextContainer 图片信息文本区
var ImgInfoTextContainer = widget.NewAccordion()

var ImgOperateImgOperateAbilityContainer = container.New(layout.NewVBoxLayout())

// 创建弹出菜单
var popupMenu = fyne.NewMenu("操作",
	fyne.NewMenuItem("新建", func() {
		log.Println("执行新建操作")
		// todo ...
	}),
	fyne.NewMenuItemSeparator(), // 分隔线
	fyne.NewMenuItem("打开", func() {
		log.Println("执行打开操作")
		// todo ...
	}),
	fyne.NewMenuItemSeparator(), // 分隔线
	fyne.NewMenuItem("保存", func() {
		log.Println("执行保存操作")
		// todo ...
	}),
)

// ImgCanvasObject 创建包含图片和控制按钮的内容
func ImgCanvasObject(img *canvas.Image, scale *float64, originalSize *fyne.Size) fyne.CanvasObject {

	// 放大按钮
	ImgZoomInBtn := &widget.Button{
		Icon: theme.ZoomInIcon(),
		OnTapped: func() {
			if img == nil {
				return
			}
			*scale += 0.2
			updateImageSize(img, scale, originalSize)
		},
	}

	// 缩小按钮
	ImgZoomOutBtn := &widget.Button{
		Icon: theme.ZoomOutIcon(),
		OnTapped: func() {
			if img == nil {
				return
			}
			*scale -= 0.2
			if *scale < 0.2 { // 限制最小缩放
				*scale = 0.2
			}
			updateImageSize(img, scale, originalSize)
		},
	}

	// 重置按钮
	ImgZoomResetBtn := &widget.Button{
		Icon: theme.ViewRefreshIcon(),
		OnTapped: func() {
			if img == nil {
				return
			}
			*scale = 1.0
			updateImageSize(img, scale, originalSize)
		},
	}

	PrevBtn := &widget.Button{
		Icon: theme.NavigateBackIcon(),
		OnTapped: func() {
			logger.Debug("上一个")
			// todo ...
		},
	}

	NextBtn := &widget.Button{
		Icon: theme.NavigateNextIcon(),
		OnTapped: func() {
			logger.Debug("下一个")
			// todo ...
		},
	}

	// 按钮容器
	controls := container.NewHBox(
		layout.NewSpacer(),
		PrevBtn,
		ImgZoomInBtn,
		ImgZoomOutBtn,
		ImgZoomResetBtn,
		NextBtn,
		layout.NewSpacer(),
	)

	controls2 := container.NewVBox(
		layout.NewSpacer(),
		controls,
	)

	// 图片容器（使用滚动容器，方便查看大图）
	scrollContainer := container.NewScroll(img)
	scrollContainer.SetMinSize(fyne.NewSize(originalSize.Width, originalSize.Height))

	// 主容器
	return container.NewVBox(
		scrollContainer,
		controls2,
	)
}

// 更新图片大小
func updateImageSize(img *canvas.Image, scale *float64, originalSize *fyne.Size) {
	newWidth := float32(*scale) * originalSize.Width
	newHeight := float32(*scale) * originalSize.Height
	img.SetMinSize(fyne.NewSize(newWidth, newHeight))
	img.Refresh()
}
