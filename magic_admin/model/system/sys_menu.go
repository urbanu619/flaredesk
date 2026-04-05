package system

import (
	"fmt"
	"go_server/model/common"
	"gorm.io/gorm"
)

type Menus struct {
	common.GormBaseModel
	ParentId    int64  `json:"parentId" gorm:"column:parent_id;comment:上级菜单Id"`
	Name        string `json:"name" gorm:"column:name;type:varchar(50);comment:菜单名"`
	Icon        string `json:"icon" gorm:"column:icon;type:varchar(200);comment:图标"`
	Router      string `json:"router" gorm:"column:router;unique;type:varchar(80);comment:前端路由"`
	IsHide      bool   `json:"isHide" gorm:"comment:是否隐藏"`
	IsFull      bool   `json:"isFull" gorm:"comment:是否全部"`
	IsAffix     bool   `json:"isAffix" gorm:"comment:是否前缀"`
	IsKeepAlive bool   `json:"isKeepAlive" gorm:"comment:是否保持活跃"`
	PageName    string `json:"pageName" gorm:"type:varchar(200);comment:页面名称"`
	Component   string `json:"component" gorm:"type:varchar(200);comment:组件"`
	Sort        int    `json:"sort" gorm:"comment:排序"`
	Enable      bool   `json:"enable" gorm:"comment:是否启用"`
}

func (*Menus) TableName() string {
	return common.ModelPrefix + "menus"
}

func NewMenus() *Menus {
	return &Menus{}
}

func (*Menus) Comment() string {
	return "前端菜单表"
}

type MenuNode struct {
	Id          int64       `json:"id"`
	ParentId    int64       `json:"parent_id"`
	Name        string      `json:"name"`
	Icon        string      `json:"icon"`
	Router      string      `json:"router"`
	Enable      bool        `json:"enable"`
	IsHide      bool        `json:"isHide" gorm:"comment:是否隐藏"`
	IsFull      bool        `json:"isFull" gorm:"comment:是否全部"`
	IsAffix     bool        `json:"isAffix" gorm:"comment:是否前缀"`
	IsKeepAlive bool        `json:"isKeepAlive" gorm:"comment:是否保持活跃"`
	PageName    string      `json:"pageName" gorm:"type:varchar(200);comment:页面名称"`
	Component   string      `json:"component" gorm:"type:varchar(200);comment:组件"`
	Children    []*MenuNode `json:"children"`
}

func (s *Menus) GetUserTree(db *gorm.DB, userPermissions []string) ([]MenuNode, error) {
	menusNodeTree := make([]MenuNode, 0) // 根节点
	// 所有节点
	allMenus := make([]Menus, 0)
	if err := db.Model(&Menus{}).
		Where("id in (?)", userPermissions).
		Order("sort").
		Find(&allMenus).Error; err != nil {
		return menusNodeTree, err
	}
	for _, m := range allMenus {
		fmt.Println("Component:", m.Component)

		childMenus := make([]*MenuNode, 0)
		if m.ParentId == 0 {
			rootNode := MenuNode{
				Id:          m.ID,
				ParentId:    m.ParentId,
				Name:        m.Name,
				Icon:        m.Icon,
				Router:      m.Router,
				IsHide:      m.IsHide,
				IsFull:      m.IsFull,
				IsAffix:     m.IsAffix,
				IsKeepAlive: m.IsKeepAlive,
				PageName:    m.PageName,
				Component:   m.Component,
				Enable:      m.Enable,
				Children:    childMenus,
			}
			menusNodeTree = append(menusNodeTree, rootNode)
		}
	}
	for i, _ := range menusNodeTree {
		s.walk(allMenus, &menusNodeTree[i])
	}
	return menusNodeTree, nil
}

// MenuTree 系统递归树[满足无限层级递归]

func (s *Menus) SysTree(db *gorm.DB) ([]MenuNode, error) {
	menusNodeTree := make([]MenuNode, 0) // 根节点
	// 所有节点
	allMenus := make([]Menus, 0)
	if err := db.Model(&Menus{}).Where("1=1").
		Order("sort").
		Find(&allMenus).Error; err != nil {
		return menusNodeTree, err
	}
	for _, m := range allMenus {
		childMenus := make([]*MenuNode, 0)
		if m.ParentId == 0 {
			rootNode := MenuNode{
				Id:          m.ID,
				ParentId:    m.ParentId,
				Name:        m.Name,
				Icon:        m.Icon,
				Router:      m.Router,
				IsHide:      m.IsHide,
				IsFull:      m.IsFull,
				IsAffix:     m.IsAffix,
				Enable:      m.Enable,
				IsKeepAlive: m.IsKeepAlive,
				PageName:    m.PageName,
				Component:   m.Component,
				Children:    childMenus,
			}
			menusNodeTree = append(menusNodeTree, rootNode)
		}
	}
	for i, _ := range menusNodeTree {
		s.walk(allMenus, &menusNodeTree[i])
	}
	return menusNodeTree, nil
}

// 递归组装菜单树（根节点）
func (s *Menus) walk(allMenus []Menus, rootNode *MenuNode) {
	// 列出所有下级子目录
	nodes := s.childrenList(allMenus, rootNode.Id)
	if len(nodes) == 0 {
		return
	}
	// 遍历这些文件
	for _, node := range nodes {
		fmt.Println("Component:", node.Component)
		newNode := MenuNode{
			Id:          node.Id,
			ParentId:    node.ParentId,
			Name:        node.Name,
			Icon:        node.Icon,
			Router:      node.Router,
			Enable:      node.Enable,
			IsHide:      node.IsHide,
			IsFull:      node.IsFull,
			IsAffix:     node.IsAffix,
			IsKeepAlive: node.IsKeepAlive,
			PageName:    node.PageName,
			Component:   node.Component,
			Children:    make([]*MenuNode, 0),
		}
		s.walk(allMenus, &newNode)
		rootNode.Children = append(rootNode.Children, &newNode)
	}
	return
}

// 获得子节点列表
func (s *Menus) childrenList(allMenus []Menus, pId int64) (menusNodeTree []MenuNode) {
	for _, m := range allMenus {
		if m.ParentId == pId {
			rootNode := MenuNode{
				Id:          m.ID,
				ParentId:    m.ParentId,
				Name:        m.Name,
				Icon:        m.Icon,
				Router:      m.Router,
				Enable:      m.Enable,
				IsHide:      m.IsHide,
				IsFull:      m.IsFull,
				IsAffix:     m.IsAffix,
				IsKeepAlive: m.IsKeepAlive,
				PageName:    m.PageName,
				Component:   m.Component,
				Children:    make([]*MenuNode, 0),
			}
			menusNodeTree = append(menusNodeTree, rootNode)
		}
	}
	return
}
