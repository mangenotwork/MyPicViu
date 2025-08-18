package ui

import (
	"MyPicViu/common/logger"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io"
	"net/url"
)

// 菜单
func MakeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	newItem := fyne.NewMenuItem("打开", nil)
	fileItem := fyne.NewMenuItem("打开图片", func() {
		logger.Debug("Menu New->File")
		openImgFile(w)
	})
	fileItem.Icon = theme.FileIcon()
	dirItem := fyne.NewMenuItem("打开目录", func() {
		logger.Debug("Menu New->Directory")
		openFile(w)
	})
	dirItem.Icon = theme.FolderIcon()
	newItem.ChildMenu = fyne.NewMenu("",
		dirItem,
		fileItem,
	)

	openSettings := func() {
		w := a.NewWindow("设置")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(440, 520))
		w.Show()
	}
	showAbout := func() {
		w := a.NewWindow("关于")
		w.SetContent(widget.NewLabel("MyPicViu\n我的图片查看器，pc软件，查看图片详细信息，显示颜色分布，图像处理等功能，图片文件分类管理，主打支持linux系统的桌面软件。"))
		w.Show()
	}
	aboutItem := fyne.NewMenuItem("关于", showAbout)
	settingsItem := fyne.NewMenuItem("设置", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut
	w.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		openSettings()
	})

	helpMenu := fyne.NewMenu("帮助",
		fyne.NewMenuItem("使用文档", func() {
			u, _ := url.Parse("https://github.com/mangenotwork/MyPicViu")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItem("项目地址", func() {
			u, _ := url.Parse("https://github.com/mangenotwork/MyPicViu")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("联系作者", func() {
			u, _ := url.Parse("https://github.com/mangenotwork/MyPicViu")
			_ = a.OpenURL(u)
		}))

	updateMenu := fyne.NewMenu("检查更新")

	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("文件", newItem)
	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	file.Items = append(file.Items, aboutItem)
	return fyne.NewMainMenu(
		file,
		helpMenu,
		updateMenu,
	)
}

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
		dataManager.AddRootFileNode(reader.URI().Name(), reader.URI().Path())

	}, win)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
	fd.Show()
}

func showImage(f fyne.URIReadCloser) {
	img := loadImage(f)
	if img == nil {
		return
	}
	img.FillMode = canvas.ImageFillOriginal

	w := fyne.CurrentApp().NewWindow(f.URI().Name())
	w.SetContent(container.NewScroll(img))
	w.Resize(fyne.NewSize(320, 240))
	w.Show()
}

func loadImage(f fyne.URIReadCloser) *canvas.Image {
	data, err := io.ReadAll(f)
	if err != nil {
		fyne.LogError("Failed to load image data", err)
		return nil
	}
	res := fyne.NewStaticResource(f.URI().Name(), data)

	return canvas.NewImageFromResource(res)
}
