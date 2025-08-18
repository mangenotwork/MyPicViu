package main

import (
	"MyPicViu/common/define"
	"MyPicViu/common/logger"
	"MyPicViu/internal/ui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	logger.SetLogFile("./logs/", "MyPicViu", 7)
	InitDB()
}

func main() {

	a := app.NewWithID("MyPicViu.2025.0814")

	makeTray(a)
	logLifecycle(a)
	w := a.NewWindow("MyPicViu")

	icon, _ := fyne.LoadResourceFromPath("./logo.png")
	w.SetIcon(icon)
	w.SetMainMenu(ui.MakeMenu(a, w))
	w.SetMaster()

	w.SetContent(ui.MainContent())

	w.Resize(fyne.NewSize(1600, 900))

	//w.SetFullScreen(true)
	// 初始化弹框通知
	//notice := widget.NewRichTextFromMarkdown(
	//	"This demo has been moved to a new repository.\n\n[Fyne demo on GitHub](https://github.com/fyne-io/demo)")
	//notice.Segments[2].(*widget.HyperlinkSegment).Alignment = fyne.TextAlignCenter
	//dialog.ShowCustom("Fyne Demo Moved", "OK", notice, w)

	w.ShowAndRun()
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		logger.Debug("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		logger.Debug("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		logger.Debug("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		logger.Debug("Lifecycle: Exited Foreground")
	})
}

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("Hello", func() {})
		h.Icon = theme.HomeIcon()
		menu := fyne.NewMenu("Hello World", h)
		h.Action = func() {
			logger.Debug("System tray menu tapped")
			h.Label = "Welcome"
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}

func shortcutFocused(s fyne.Shortcut, cb fyne.Clipboard, f fyne.Focusable) {
	switch sh := s.(type) {
	case *fyne.ShortcutCopy:
		sh.Clipboard = cb
	case *fyne.ShortcutCut:
		sh.Clipboard = cb
	case *fyne.ShortcutPaste:
		sh.Clipboard = cb
	}
	if focused, ok := f.(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func resolveDBPath() string {
	var rawPath string
	switch runtime.GOOS {
	case "windows":
		rawPath = define.DBFilePathWindows
	case "linux":
		rawPath = define.DBFilePathLinux
	case "darwin":
		rawPath = define.DBFilePathMac
	default:
		panic("不支持的操作系统")
	}

	// 解析环境变量和~
	if runtime.GOOS == "windows" {
		return os.ExpandEnv(rawPath)
	}
	return filepath.Join(os.Getenv("HOME"), rawPath[2:]) // 替换~为用户目录
}

func InitDB() {
	dbPath := resolveDBPath()
	dir := filepath.Dir(dbPath)
	logger.Debug("应用数据文件: ", dir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
