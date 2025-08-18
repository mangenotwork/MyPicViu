package main

import (
	"MyPicViu/common/define"
	"MyPicViu/common/logger"
	"MyPicViu/internal/db"
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
	ui.InitUI()
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

	w.SetContent(ui.MainContent(w))

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
	switch runtime.GOOS {
	case "windows":
		appDataPath := os.Getenv("APPDATA")
		if appDataPath == "" {
			userProfile := os.Getenv("USERPROFILE")
			appDataPath = filepath.Join(userProfile, "AppData", "Roaming")
		}
		return filepath.Join(appDataPath, define.DBFileDirName, define.DBFileFileName)
	case "linux":
		home := os.Getenv("HOME")
		return filepath.Join(home, ".local", "share", define.DBFileDirName, define.DBFileFileName)
	case "darwin":
		home := os.Getenv("HOME")
		return filepath.Join(home, "Library", "Application Support", define.DBFileDirName, define.DBFileFileName)
	default:
		panic("不支持的操作系统")
	}
}

func InitDB() {
	logger.Debug("初始化DB")
	dbPath := resolveDBPath()
	dir := filepath.Dir(dbPath)
	logger.Debug("应用数据文件: ", dir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	db.InitDB(dbPath)
	stats, err := db.DB.Stats(db.TreeNodeTable)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(stats)

	err = db.DB.Set(db.TreeNodeTable, "data1", "111111111")
	if err != nil {
		logger.Error(err)
	}

	var data1 string = ""
	err = db.DB.Get(db.TreeNodeTable, "data1", &data1)
	if err != nil {
		logger.Error(err)
	}

	logger.Debug(data1)
}
