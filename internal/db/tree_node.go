package db

//import (
//	"MyPicViu/common/logger"
//	"MyPicViu/common/utils"
//	"MyPicViu/internal/ui"
//	"fmt"
//	"fyne.io/fyne/v2"
//	"fyne.io/fyne/v2/widget"
//	"os"
//	"path/filepath"
//)
//
//var TreeData []*TreeNode
//
//// TreeNode 目录树节点
//type TreeNode struct {
//	ID       string      `json:"id"`       // 节点唯一标识
//	Name     string      `json:"name"`     // 节点显示名称
//	FilePath string      `json:"filePath"` // 文件路径
//	Children []*TreeNode `json:"children"` // 子节点
//	IsFile   bool        `json:"isFile"`   // 区分文件(叶子)和文件夹(分支)
//	Parent   *TreeNode   `json:"-"`        // 父节点引用
//}
//
//// TreeDataManager 树数据管理器
//type TreeDataManager struct {
//	rootNodes []*TreeNode
//	nodeMap   map[string]*TreeNode // 节点缓存映射
//	tree      *widget.Tree         // 关联的树组件
//}
//
//func NewTreeDataManager(rootNodes []*TreeNode) *TreeDataManager {
//	manager := &TreeDataManager{
//		rootNodes: rootNodes,
//		nodeMap:   make(map[string]*TreeNode),
//	}
//	manager.buildNodeMap(rootNodes)
//	return manager
//}
//
//func (m *TreeDataManager) GetNodeMapLen() int {
//	return len(m.nodeMap)
//}
//
//func (m *TreeDataManager) buildNodeMap(nodes []*TreeNode) {
//	for _, node := range nodes {
//		m.nodeMap[node.ID] = node
//		if len(node.Children) > 0 {
//			m.buildNodeMap(node.Children)
//		}
//	}
//}
//
//// SetTree 关联树组件
//func (m *TreeDataManager) SetTree(tree *widget.Tree) {
//	m.tree = tree
//}
//
//// FindNode 根据ID查找节点
//func (m *TreeDataManager) FindNode(id string) *TreeNode {
//	return m.nodeMap[id]
//}
//
//// GetChildrenIDs 获取子节点ID列表
//func (m *TreeDataManager) GetChildrenIDs(id string) []widget.TreeNodeID {
//	if id == "" {
//		ids := make([]widget.TreeNodeID, len(m.rootNodes))
//		for i, node := range m.rootNodes {
//			ids[i] = node.ID
//		}
//		return ids
//	}
//
//	node := m.FindNode(id)
//	if node == nil {
//		return []widget.TreeNodeID{}
//	}
//
//	ids := make([]widget.TreeNodeID, len(node.Children))
//	for i, child := range node.Children {
//		ids[i] = child.ID
//	}
//	return ids
//}
//
//// IsBranch 判断节点是否为分支
//func (m *TreeDataManager) IsBranch(id string) bool {
//	if id == "" {
//		return true
//	}
//
//	node := m.FindNode(id)
//	return node != nil && !node.IsFile
//}
//
//func PreOrderTraversal(node *TreeNode, handler func(*TreeNode)) {
//	if node == nil {
//		return
//	}
//	// 处理当前节点
//	handler(node)
//	// 递归遍历子节点
//	for _, child := range node.Children {
//		PreOrderTraversal(child, handler)
//	}
//}
//
//// 从数据db中初始化TreeData
//func buildTreeData() []*TreeNode {
//
//	data := make([]*TreeNode, 0)
//	err := DB.Get(TreeNodeTable, TreeNodeKey, &data)
//	if err != nil {
//		logger.Error("获取存储节点数据失败: ", err)
//	}
//
//	for _, node := range data {
//		PreOrderTraversal(node, func(node *TreeNode) {
//			node.Parent = node
//		})
//	}
//
//	return data
//}
//
//// AddRootFileNode 在根节点动态添加节点
//func (m *TreeDataManager) AddRootFileNode(name, filePath string) bool {
//	newNode := &TreeNode{
//		ID:       name,
//		Name:     name,
//		FilePath: filePath,
//		IsFile:   true,
//	}
//
//	n := FindNode(TreeData, newNode.ID)
//	if n != nil {
//		logger.Debug("已存在")
//		return false
//	}
//
//	// 添加到根节点
//	m.rootNodes = append(m.rootNodes, newNode)
//	TreeData = append(TreeData, newNode)
//	err := DB.Set(TreeNodeTable, TreeNodeKey, &TreeData)
//	if err != nil {
//		logger.Error("存储节点数据失败", err)
//		return false
//	}
//
//	// 更新节点映射
//	m.nodeMap[newNode.ID] = newNode
//	// 刷新树组件
//	if m.tree != nil {
//		m.tree.Refresh() // 根节点变化，刷新整个树
//	}
//
//	return true
//}
//
//func (m *TreeDataManager) AddRootDirNode(dir string) bool {
//	newNode := &TreeNode{
//		ID:       dir,
//		Name:     dir,
//		FilePath: dir,
//		IsFile:   false,
//		Children: make([]*TreeNode, 0),
//		Parent:   nil, // 根节点没有父节点
//	}
//	n := FindNode(TreeData, newNode.ID)
//	if n != nil {
//		logger.Debug("已存在")
//		return false
//	}
//
//	m.nodeMap[newNode.ID] = newNode
//
//	err := traverseDir(dir, newNode, m.nodeMap)
//	if err != nil {
//		logger.Error(err)
//	}
//
//	//logger.Debug("\n\n数据拉取后")
//	//TreeNodeBFSPrint(newNode)
//
//	m.rootNodes = append(m.rootNodes, newNode)
//	TreeData = append(TreeData, newNode)
//	err = DB.Set(TreeNodeTable, TreeNodeKey, &TreeData)
//	if err != nil {
//		logger.Error("存储节点数据失败", err)
//		return false
//	}
//	m.tree.Refresh()
//
//	return true
//}
//
//// TreeNodeBFSPrint 广度优先遍历（使用队列）
//func TreeNodeBFSPrint(root *TreeNode) {
//	if root == nil {
//		return
//	}
//
//	// 使用队列存储节点和当前层级
//	queue := []struct {
//		node  *TreeNode
//		level int
//	}{
//		{root, 0},
//	}
//
//	for len(queue) > 0 {
//		current := queue[0]
//		queue = queue[1:]
//		// 打印当前节点
//		fmt.Printf("Level %d: %s (%s)\n", current.level, current.node.Name, current.node.FilePath)
//		for _, child := range current.node.Children {
//			queue = append(queue, struct {
//				node  *TreeNode
//				level int
//			}{child, current.level + 1})
//		}
//	}
//}
//
//// traverseDir 递归遍历目录并创建节点
//func traverseDir(path string, parentNode *TreeNode, nodeMap map[string]*TreeNode) error {
//	entries, err := os.ReadDir(path)
//	if err != nil {
//		return fmt.Errorf("读取目录 %s 失败: %v", path, err)
//	}
//
//	for _, entry := range entries {
//		fullPath := filepath.Join(path, entry.Name())
//		if entry.IsDir() {
//
//			dirNode := &TreeNode{
//				ID:       fullPath,
//				Name:     entry.Name(),
//				FilePath: fullPath,
//				IsFile:   false,
//				Children: make([]*TreeNode, 0),
//				Parent:   parentNode, // 正确设置父节点为当前的parentNode
//			}
//
//			n := FindNode(TreeData, dirNode.ID)
//			if n != nil {
//				logger.Debug("已存在")
//				continue
//			}
//
//			parentNode.Children = append(parentNode.Children, dirNode)
//			// 关键修改：将当前目录节点添加到nodeMap
//			nodeMap[dirNode.ID] = dirNode
//			// 递归遍历子目录，传入当前目录节点作为父节点
//			if err := traverseDir(fullPath, dirNode, nodeMap); err != nil {
//				return err
//			}
//		} else {
//
//			fileInfo, err := entry.Info()
//			if err != nil {
//				return fmt.Errorf("获取文件信息失败: %v", err)
//			}
//
//			// 判断是否是图片类型
//			if !utils.IsImgFile(fullPath) {
//				logger.Debug("不是图片")
//			}
//
//			fileNode := &TreeNode{
//				ID:       fullPath,
//				Name:     entry.Name(),
//				FilePath: fullPath,
//				IsFile:   true,
//				Children: make([]*TreeNode, 0),
//				Parent:   parentNode, // 正确设置父节点
//			}
//
//			//n := FindNode(TreeData, fileNode.ID)
//			//if n != nil {
//			//	logger.Debug("已存在")
//			//	continue
//			//}
//
//			parentNode.Children = append(parentNode.Children, fileNode)
//			// 关键修改：将当前文件节点添加到nodeMap
//			nodeMap[fileNode.ID] = fileNode
//			logger.Debug("[文件] %s (大小: %d bytes)\n", fullPath, fileInfo.Size())
//		}
//	}
//
//	return nil
//}
//
//// 动态删除节点
//func (m *TreeDataManager) DeleteNode(id string) bool {
//	node := m.FindNode(id)
//	if node == nil {
//		return false
//	}
//
//	// 递归删除所有子节点
//	var deleteChildren func(nodes []*TreeNode)
//	deleteChildren = func(nodes []*TreeNode) {
//		for _, n := range nodes {
//			delete(m.nodeMap, n.ID)
//			deleteChildren(n.Children)
//		}
//	}
//	deleteChildren(node.Children)
//
//	// 从父节点中移除
//	if node.Parent == nil {
//		// 根节点
//		for i, n := range m.rootNodes {
//			if n.ID == id {
//				m.rootNodes = append(m.rootNodes[:i], m.rootNodes[i+1:]...)
//				break
//			}
//		}
//	} else {
//		// 非根节点
//		parent := node.Parent
//		for i, child := range parent.Children {
//			if child.ID == id {
//				parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
//				break
//			}
//		}
//	}
//
//	// 从映射中删除当前节点
//	delete(m.nodeMap, id)
//
//	// 刷新树组件
//	if node.Parent == nil {
//		m.tree.Refresh() // 根节点变化，刷新整个树
//	} else {
//		m.tree.RefreshItem(node.Parent.ID) // 刷新父节点
//	}
//	return true
//}
//
//// 动态修改节点名称
//func (m *TreeDataManager) RenameNode(id, newName string) bool {
//	node := m.FindNode(id)
//	if node == nil {
//		return false
//	}
//
//	node.Name = newName
//
//	// 刷新指定节点
//	if m.tree != nil {
//		m.tree.RefreshItem(id)
//	}
//
//	return true
//}
//
//// 创建自定义树
//func CreateCustomTree() (*widget.Tree, *TreeDataManager) {
//	// 创建数据管理器
//	dataManager := NewTreeDataManager(TreeData)
//
//	// 创建树组件
//	tree := widget.NewTree(
//		// 节点ID获取函数：根据父节点ID返回子节点ID列表
//		dataManager.GetChildrenIDs,
//
//		// 分支判断函数：判断节点是否为分支(文件夹)
//		dataManager.IsBranch,
//
//		// 模板创建函数：创建分支和叶子节点的UI模板
//		func(branch bool) fyne.CanvasObject {
//			if branch {
//				return NewClickableLabel("目录", "")
//			}
//			return NewClickableLabel("图片", "")
//		},
//
//		// 数据绑定函数：将节点数据绑定到UI组件
//		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
//			logger.Debug("id = ", id)
//			name := id
//			node := FindNode(TreeData, id)
//			logger.Debug("node = ", node)
//			if node != nil {
//				name = node.Name
//			}
//			//if branch {
//			//	name = "[目录]" + name
//			//}
//			//o.(*widget.Label).SetText(name)
//
//			//logger.Debug(name)
//			//NewClickableLabel(o, name, func(pos fyne.Position) {
//			//	logger.Debug("aaaa")
//			//})
//
//			// 类型转换为自定义标签，并设置文本和节点ID
//			if label, ok := o.(*ClickableLabel); ok {
//				label.SetText(name) // 设置显示文本
//				label.NodeID = id   // 绑定实际节点ID
//			}
//
//		},
//	)
//
//	// 关联树组件到数据管理器
//	dataManager.SetTree(tree)
//
//	return tree, dataManager
//}
//
//// 自定义可点击标签
//type ClickableLabel struct {
//	*widget.Label
//	NodeID string
//}
//
//func (c *ClickableLabel) Tapped(event *fyne.PointEvent) {
//	logger.DebugF("左键点击了节点 [ID: %s, 文本: %s, 位置: %v]",
//		c.NodeID, c.Text, event.Position)
//
//	node := FindNode(TreeData, c.NodeID)
//	if node != nil && node.IsFile {
//		ui.ShowImg(node.FilePath)
//	}
//}
//
//func (c *ClickableLabel) TappedSecondary(event *fyne.PointEvent) {
//	logger.DebugF("右键点击了节点 [ID: %s, 文本: %s, 位置: %v]",
//		c.NodeID, c.Text, event.Position)
//}
//
//// 初始化时必须调用
//func NewClickableLabel(text, nodeID string) *ClickableLabel {
//	l := &ClickableLabel{
//		widget.NewLabel(text),
//		nodeID,
//	}
//	l.ExtendBaseWidget(l) // 关键：初始化基础组件
//	return l
//}
//
//// 查找节点函数：根据ID在树中查找对应的节点
//func FindNode(rootNodes []*TreeNode, id string) *TreeNode {
//	// 遍历根节点查找
//	for _, node := range rootNodes {
//		if node.ID == id {
//			return node
//		}
//
//		// 递归查找子节点
//		found := FindNode(node.Children, id)
//		if found != nil {
//			return found
//		}
//	}
//	return nil
//}
