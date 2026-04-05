package system

import (
	"go_server/model/common"
)

type Role struct {
	common.GormFullModel
	Name  string `json:"name" gorm:"column:name;type:varchar(100);comment:角色名称"`
	Apis  string `json:"auths" gorm:"column:apis;type:text;comment:Api权限"`
	Menus string `json:"menus" gorm:"column:menus;type:text;comment:menu权限"`
	Desc  string `json:"desc" gorm:"column:desc;type:varchar(200);comment:角色权限描述;"`
}

func (*Role) TableName() string {
	return common.ModelPrefix + "role"
}

func NewRole() *Role {
	return &Role{}
}

func (*Role) Comment() string {
	return "角色表"
}
