package ui

import (
	"MyPicViu/internal/db"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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
		fmt.Printf("节点被点击: ID=%s, 名称=%s, 类型=%s\n",
			node.ID, node.Name, nodeType)

		updateContent(id)
	}

	return container.New(
		layout.NewStackLayout(),
		tree,
	)
}

func updateContent(id string) {
	contentTitle.Text = "当前选择: " + id

	// 根据选择的目录项显示不同内容
	switch id {
	case "a":
		contentTitle.Text = "这是A分类的详细内容\n包含子项: a1, a2, a3"
	case "a1":
		contentTitle.Text = "A1子项的具体信息"
	case "a2":
		contentTitle.Text = "A2子项的具体信息"
	case "a3":
		contentTitle.Text = "A3子项的具体信息"
	case "b":
		contentTitle.Text = "这是B分类的详细内容\n包含子项: b1"
	case "b1":
		contentTitle.Text = "B1子项的具体信息"
	case "c":
		contentTitle.Text = "这是C分类的详细内容"
	default:
		contentTitle.Text = "请从左侧选择一个目录项查看详情"
	}
}
