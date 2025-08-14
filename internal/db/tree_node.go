package db

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"path/filepath"
)

// 自定义数据结构：表示目录树节点
type TreeNode struct {
	ID       string      // 节点唯一标识
	Name     string      // 节点显示名称
	Children []*TreeNode // 子节点
	IsFile   bool        // 区分文件(叶子)和文件夹(分支)
	Parent   *TreeNode   // 父节点引用
}

// 树数据管理器
type TreeDataManager struct {
	rootNodes []*TreeNode
	nodeMap   map[string]*TreeNode // 节点缓存映射
	tree      *widget.Tree         // 关联的树组件
}

// 创建新的数据管理器
func NewTreeDataManager(rootNodes []*TreeNode) *TreeDataManager {
	manager := &TreeDataManager{
		rootNodes: rootNodes,
		nodeMap:   make(map[string]*TreeNode),
	}
	manager.buildNodeMap(rootNodes)
	return manager
}

// 构建节点映射表
func (m *TreeDataManager) buildNodeMap(nodes []*TreeNode) {
	for _, node := range nodes {
		m.nodeMap[node.ID] = node
		if len(node.Children) > 0 {
			m.buildNodeMap(node.Children)
		}
	}
}

// 关联树组件
func (m *TreeDataManager) SetTree(tree *widget.Tree) {
	m.tree = tree
}

// 根据ID查找节点
func (m *TreeDataManager) FindNode(id string) *TreeNode {
	return m.nodeMap[id]
}

// 获取子节点ID列表
func (m *TreeDataManager) GetChildrenIDs(id string) []widget.TreeNodeID {
	if id == "" { // 根节点
		ids := make([]widget.TreeNodeID, len(m.rootNodes))
		for i, node := range m.rootNodes {
			ids[i] = node.ID // 为什么显示id  不能显示name
		}
		return ids
	}

	node := m.FindNode(id)
	if node == nil {
		return []widget.TreeNodeID{}
	}

	ids := make([]widget.TreeNodeID, len(node.Children))
	for i, child := range node.Children {
		ids[i] = child.ID
	}
	return ids
}

// 判断节点是否为分支
func (m *TreeDataManager) IsBranch(id string) bool {
	if id == "" { // 根节点视为分支
		return true
	}

	node := m.FindNode(id)
	return node != nil && !node.IsFile
}

var TreeData = buildTreeData()

//var BrowseTree = createCustomTree(TreeData)

func buildTreeData() []*TreeNode {
	return []*TreeNode{
		{
			ID:   "a",
			Name: "a",
			Children: []*TreeNode{
				{ID: "a1", Name: "a1", IsFile: false, Children: []*TreeNode{{ID: "a3", Name: "a3", IsFile: true}}},
				{ID: "a2", Name: "a2", IsFile: true},
			},
			IsFile: false, // 是文件夹(分支)
		},
		{
			ID:       "b",
			Name:     "b",
			Children: []*TreeNode{},
			IsFile:   false, // 是文件夹(分支)，但没有子节点
		},
		{
			ID:       "c",
			Name:     "c",
			Children: []*TreeNode{},
			IsFile:   true, // 是文件(叶子)
		},
	}
}

// 在根节点动态添加节点
func (m *TreeDataManager) AddRootNode(name string) bool {
	newNode := &TreeNode{
		ID:     name,
		Name:   name,
		IsFile: true,
	}

	// 添加到根节点
	m.rootNodes = append(m.rootNodes, newNode)

	log.Println("%v", newNode)

	// 更新节点映射
	m.nodeMap[newNode.ID] = newNode

	// 刷新树组件
	if m.tree != nil {
		m.tree.Refresh() // 根节点变化，刷新整个树
	}

	return true
}

func (m *TreeDataManager) AddRootDirNode(dir string) bool {
	// 生成唯一ID，避免路径作为ID可能产生的冲突
	//id := fmt.Sprintf("root_%s", filepath.Base(dir))

	newNode := &TreeNode{
		ID:       dir, // 使用更安全的唯一ID
		Name:     dir, // 显示目录名而非完整路径
		IsFile:   false,
		Children: make([]*TreeNode, 0),
		Parent:   nil, // 根节点没有父节点
	}

	log.Println("准备遍历目录 : ", dir)
	err := traverseDir(dir, newNode, m.nodeMap)
	if err != nil {
		log.Println(err)
	}

	//log.Println("%v", newNode)

	// 添加到根节点列表
	m.rootNodes = append(m.rootNodes, newNode)

	// 更新节点映射（使用唯一ID）
	m.nodeMap[newNode.ID] = newNode

	// 刷新树组件
	if m.tree != nil {
		m.tree.Refresh()
	}

	return true
}

// traverseDir 递归遍历目录并创建节点
func traverseDir(path string, parentNode *TreeNode, nodeMap map[string]*TreeNode) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("读取目录 %s 失败: %v", path, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		// 为每个节点生成唯一ID（结合父节点ID避免重名冲突）
		//nodeID := fmt.Sprintf("%s_%s", parentNode.ID, entry.Name())

		if entry.IsDir() {
			fmt.Printf("[目录] %s\n", fullPath)

			dirNode := &TreeNode{
				ID:       fullPath,
				Name:     entry.Name(),
				IsFile:   false,
				Children: make([]*TreeNode, 0),
				Parent:   parentNode, // 正确设置父节点为当前的parentNode
			}
			parentNode.Children = append(parentNode.Children, dirNode)
			// 关键修改：将当前目录节点添加到nodeMap
			nodeMap[dirNode.ID] = dirNode
			// 递归遍历子目录，传入当前目录节点作为父节点
			if err := traverseDir(fullPath, dirNode, nodeMap); err != nil {
				return err
			}
		} else {
			fileInfo, err := entry.Info()
			if err != nil {
				return fmt.Errorf("获取文件信息失败: %v", err)
			}

			fileNode := &TreeNode{
				ID:       fullPath,
				Name:     entry.Name(),
				IsFile:   true,
				Children: make([]*TreeNode, 0),
				Parent:   parentNode, // 正确设置父节点
			}
			parentNode.Children = append(parentNode.Children, fileNode)
			// 关键修改：将当前文件节点添加到nodeMap
			nodeMap[fileNode.ID] = fileNode
			fmt.Printf("[文件] %s (大小: %d bytes)\n", fullPath, fileInfo.Size())
		}
	}

	return nil
}

// 动态删除节点
func (m *TreeDataManager) DeleteNode(id string) bool {
	node := m.FindNode(id)
	if node == nil {
		return false
	}

	// 递归删除所有子节点
	var deleteChildren func(nodes []*TreeNode)
	deleteChildren = func(nodes []*TreeNode) {
		for _, n := range nodes {
			delete(m.nodeMap, n.ID)
			deleteChildren(n.Children)
		}
	}
	deleteChildren(node.Children)

	// 从父节点中移除
	if node.Parent == nil {
		// 根节点
		for i, n := range m.rootNodes {
			if n.ID == id {
				m.rootNodes = append(m.rootNodes[:i], m.rootNodes[i+1:]...)
				break
			}
		}
	} else {
		// 非根节点
		parent := node.Parent
		for i, child := range parent.Children {
			if child.ID == id {
				parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
				break
			}
		}
	}

	// 从映射中删除当前节点
	delete(m.nodeMap, id)

	// 刷新树组件
	if node.Parent == nil {
		m.tree.Refresh() // 根节点变化，刷新整个树
	} else {
		m.tree.RefreshItem(node.Parent.ID) // 刷新父节点
	}
	return true
}

// 动态修改节点名称
func (m *TreeDataManager) RenameNode(id, newName string) bool {
	node := m.FindNode(id)
	if node == nil {
		return false
	}

	node.Name = newName

	// 刷新指定节点
	if m.tree != nil {
		m.tree.RefreshItem(id)
	}

	return true
}

// 创建自定义树
func CreateCustomTree(treeData []*TreeNode) (*widget.Tree, *TreeDataManager) {
	// 创建数据管理器
	dataManager := NewTreeDataManager(treeData)

	// 创建树组件
	tree := widget.NewTree(
		// 节点ID获取函数：根据父节点ID返回子节点ID列表
		//func(id widget.TreeNodeID) []widget.TreeNodeID {
		//	if id == "" { // 根节点
		//		ids := make([]widget.TreeNodeID, len(treeData))
		//		for i, node := range treeData {
		//			ids[i] = node.ID
		//		}
		//		return ids
		//	}
		//
		//	// 查找对应节点并返回其子节点ID
		//	node := FindNode(treeData, id)
		//	if node == nil {
		//		return []widget.TreeNodeID{}
		//	}
		//
		//	ids := make([]widget.TreeNodeID, len(node.Children))
		//	for i, child := range node.Children {
		//		ids[i] = child.ID
		//	}
		//	return ids
		//},
		dataManager.GetChildrenIDs,
		// 分支判断函数：判断节点是否为分支(文件夹)
		//func(id widget.TreeNodeID) bool {
		//	if id == "" { // 根节点视为分支
		//		return true
		//	}
		//
		//	node := FindNode(treeData, id)
		//	return node != nil && !node.IsFile
		//},
		dataManager.IsBranch,
		// 模板创建函数：创建分支和叶子节点的UI模板
		func(branch bool) fyne.CanvasObject {
			if branch {
				return widget.NewLabel("分支模板")
			}
			return widget.NewLabel("叶子模板")
		},
		// 数据绑定函数：将节点数据绑定到UI组件
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			node := FindNode(treeData, id)
			if node == nil {
				o.(*widget.Label).SetText(id)
				return
			}

			text := node.Name
			if branch {
				text += " (分支)"
			}
			o.(*widget.Label).SetText(text)
		},
	)

	// 关联树组件到数据管理器
	dataManager.SetTree(tree)

	// 节点选中事件
	tree.OnSelected = func(id widget.TreeNodeID) {
		node := dataManager.FindNode(id)
		if node != nil {
			nodeType := "文件"
			if !node.IsFile {
				nodeType = "文件夹"
			}
			fmt.Printf("选中: %s (%s)\n", node.Name, nodeType)
		}
	}

	return tree, dataManager
}

// 查找节点函数：根据ID在树中查找对应的节点
func FindNode(rootNodes []*TreeNode, id string) *TreeNode {
	// 遍历根节点查找
	for _, node := range rootNodes {
		if node.ID == id {
			return node
		}

		// 递归查找子节点
		found := FindNode(node.Children, id)
		if found != nil {
			return found
		}
	}
	return nil
}
