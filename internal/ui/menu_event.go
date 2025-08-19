package ui

import (
	"MyPicViu/common/logger"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

// 打开目录
func openFile(win fyne.Window) {
	dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if list == nil {
			logger.Debug("Cancelled")
			return
		}

		children, err := list.List()
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		logger.Debug(list.String())
		out := fmt.Sprintf("打开目录 %s (含%d目录或文件):\n%s", list.Name(), len(children), list.String())
		// 改为确认对话框
		dialog.ShowConfirm(
			"打开目录", // 对话框标题
			out,    // 显示的消息内容
			func(confirmed bool) { // 用户选择后的回调函数
				if confirmed {
					// 用户点击了"确定"按钮，执行相应操作
					logger.Debug("用户确认了操作")
					// 可以在这里添加确认后的逻辑，比如实际打开目录

					if dataManager.GetNodeMapLen() == 0 {
						LeftContainer.Remove(OpenInitContainer)
						LeftContainer.Add(tree)
						LeftContainer.Refresh()
					}

					// 遍历目录
					dataManager.AddRootDirNode(list.Path())
				} else {
					// 用户点击了"取消"按钮
					logger.Debug("用户取消了操作")
				}
			},
			win, // 父窗口
		)

	}, win)
}

// 打开图片
func openImgFile(win fyne.Window) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if reader == nil {
			logger.Debug("Cancelled")
			return
		}

		if reader == nil {
			logger.Debug("Cancelled")
			return
		}
		defer reader.Close()

		logger.Debug(reader.URI().Name())

		if dataManager.GetNodeMapLen() == 0 {
			LeftContainer.Remove(OpenInitContainer)
			LeftContainer.Add(tree)
			LeftContainer.Refresh()
		}

		dataManager.AddRootFileNode(reader.URI().Name(), reader.URI().Path())
		ShowImg(reader.URI().Path())

	}, win)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
	fd.Show()
}
