package system

import (
	"go_server/model/common"
	"gorm.io/gorm"
	"strings"
)

type Apis struct {
	common.GormFullModel
	ParentId int64  `json:"parentId" gorm:"column:parent_id;comment:上级接口Id"`
	Group    string `json:"group" gorm:"column:group;type:varchar(50);comment:分组名称"`
	Name     string `json:"name" gorm:"column:name;type:varchar(50);comment:接口名称"`
	Method   string `json:"method" gorm:"column:method;type:varchar(20);comment:请求方法"`
	Path     string `json:"path" gorm:"column:path;type:varchar(80);comment:接口全路径;unique"`
}

func (*Apis) TableName() string {
	return common.ModelPrefix + "apis"
}

func NewApis() *Apis {
	return &Apis{}
}

func (*Apis) Comment() string {
	return "api权限表"
}

func (s *Apis) FindOrCreateAuth(db *gorm.DB, path, method string) (*Apis, error) {
	if err := db.Where("path", path).First(&s).Error; err != nil {
		if !strings.Contains(path, "*") {
			pathList := strings.Split(path, "/")
			if len(pathList) > 2 {
				s.Group = pathList[2]
			}
			s.ParentId = 1
			s.Path = path
			s.Method = method
			s.Name = "待设置"
			if err := db.Create(&s).Error; err != nil {
				return s, err
			}
		}
		return s, err
	}
	return s, nil
}

type ApiNode struct {
	Id       int64      `json:"id"`
	ParentId int64      `json:"parentId"`
	Group    string     `json:"group"`
	Name     string     `json:"name"`
	Method   string     `json:"method"`
	Path     string     `json:"path"`
	Children []*ApiNode `json:"children"`
}

func (s *Apis) GetUserTree(db *gorm.DB, ids []string) ([]ApiNode, error) {
	menusNodeTree := make([]ApiNode, 0) // 根节点
	// 所有节点
	allMenus := make([]Apis, 0)
	if err := db.Model(&Apis{}).
		Where("id in (?)", ids).
		Find(&allMenus).Error; err != nil {
		return menusNodeTree, err
	}
	for _, m := range allMenus {
		childMenus := make([]*ApiNode, 0)
		if m.ParentId == 0 {
			rootNode := ApiNode{
				Id:       m.ID,
				ParentId: m.ParentId,
				Group:    m.Group,
				Name:     m.Name,
				Method:   m.Method,
				Path:     m.Path,
				Children: childMenus,
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

func (s *Apis) SysTree(db *gorm.DB) ([]ApiNode, error) {
	menusNodeTree := make([]ApiNode, 0) // 根节点
	// 所有节点
	allMenus := make([]Apis, 0)
	if err := db.Model(&Apis{}).Where("1=1").Find(&allMenus).Error; err != nil {
		return menusNodeTree, err
	}
	for _, m := range allMenus {
		childMenus := make([]*ApiNode, 0)
		if m.ParentId == 0 {
			rootNode := ApiNode{
				Id:       m.ID,
				ParentId: m.ParentId,
				Name:     m.Name,
				Group:    m.Group,
				Method:   m.Method,
				Path:     m.Path,
				Children: childMenus,
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
func (s *Apis) walk(allMenus []Apis, rootNode *ApiNode) {
	// 列出所有下级子目录
	nodes := s.childrenList(allMenus, rootNode.Id)
	if len(nodes) == 0 {
		return
	}
	// 遍历这些文件
	for _, m := range nodes {
		newNode := ApiNode{
			Id:       m.Id,
			ParentId: m.ParentId,
			Group:    m.Group,
			Name:     m.Name,
			Method:   m.Method,
			Path:     m.Path,
			Children: nil,
		}
		s.walk(allMenus, &newNode)
		rootNode.Children = append(rootNode.Children, &newNode)
	}
	return
}

// 获得子节点列表
func (s *Apis) childrenList(allMenus []Apis, pId int64) (menusNodeTree []ApiNode) {
	for _, m := range allMenus {
		if m.ParentId == pId {
			rootNode := ApiNode{
				Id:       m.ID,
				ParentId: m.ParentId,
				Name:     m.Name,
				Group:    m.Group,
				Method:   m.Method,
				Path:     m.Path,
				Children: nil,
			}
			menusNodeTree = append(menusNodeTree, rootNode)
		}
	}
	return
}
