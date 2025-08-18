package ui

import (
	"MyPicViu/common/logger"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// LeftContainerInit 初始化左边视图
func LeftContainerInit(w fyne.Window) {
	if dataManager.GetNodeMapLen() == 0 {
		OpenInitContainer.Add(layout.NewSpacer())
		OpenInitContainer.Add(container.NewCenter(widget.NewLabel("请选择图片文件或含有图片的目录")))
		openFileBtnContainer := container.NewCenter()
		openFileBtnContainer.Resize(fyne.NewSize(10, 0))
		openFileBtn := widget.NewButton("选择图片", func() {
			logger.Debug("选择图片")
			openImgFile(w)
		})
		openFileBtnContainer.Add(openFileBtn)
		OpenInitContainer.Add(openFileBtnContainer)
		openDirBtnContainer := container.NewCenter()
		openDirBtnContainer.Resize(fyne.NewSize(10, 0))
		openDirBtn := widget.NewButton("选择目录", func() {
			logger.Debug("选择图片")
			openFile(w)
		})
		openDirBtnContainer.Add(openDirBtn)
		OpenInitContainer.Add(openDirBtnContainer)
		OpenInitContainer.Add(layout.NewSpacer())

		LeftContainer.Add(OpenInitContainer)
		return
	}

	LeftContainer.Add(tree)

}
