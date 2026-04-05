package system

import (
	"github.com/samber/lo"
	"gorm.io/gorm"
	"time"
)

func (*Role) DataInit(db *gorm.DB) error {
	// 初始化用户等级
	reqs := []Role{
		{
			Name:  "超级管理员",
			Desc:  "管理所有权限",
			Apis:  "*",
			Menus: "*",
		},
	}
	for _, req := range reqs {
		find := Role{}
		if stat := db.Model(&Role{}).Where("name", req.Name).
			Find(&find).Statement; stat.RowsAffected == 0 {
			if err := db.Create(&req).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (*Apis) DataInit(db *gorm.DB) error {
	// 初始化用户等级
	reqs := []Apis{
		{
			ParentId: 0,
			Name:     "待定义接口",
			Group:    "待分配组",
			Method:   "All",
			Path:     "/",
		},
	}
	for _, req := range reqs {
		find := Apis{}
		if stat := db.Model(&Apis{}).Where("name = ?", req.Name).Find(&find).Statement; stat.RowsAffected == 0 {
			if err := db.Create(&req).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (*Menus) DataInit(db *gorm.DB) error {
	// 初始化菜单
	reqs := []Menus{
		// 红包管理主菜单
		{
			ParentId:    0,
			Name:        "红包管理",
			Icon:        "Wallet",
			Router:      "/redpacket",
			Sort:        1,
			Enable:      true,
			IsHide:      false,
			IsFull:      false,
			IsAffix:     false,
			IsKeepAlive: true,
		},
		// 红包配置
		{
			ParentId:    1,
			Name:        "红包配置",
			Icon:        "Menu",
			Router:      "/redpacket/config",
			Component:   "@/views/redpacket/config/index.vue",
			Sort:        1,
			Enable:      true,
			IsHide:      false,
			IsFull:      false,
			IsAffix:     false,
			IsKeepAlive: true,
		},
		// 手动发送
		{
			ParentId:    1,
			Name:        "手动发送",
			Icon:        "Menu",
			Router:      "/redpacket/send",
			Component:   "@/views/redpacket/send/index.vue",
			Sort:        2,
			Enable:      true,
			IsHide:      false,
			IsFull:      false,
			IsAffix:     false,
			IsKeepAlive: true,
		},
		// 红包记录
		{
			ParentId:    1,
			Name:        "红包记录",
			Icon:        "Menu",
			Router:      "/redpacket/record",
			Component:   "@/views/redpacket/record/index.vue",
			Sort:        3,
			Enable:      true,
			IsHide:      false,
			IsFull:      false,
			IsAffix:     false,
			IsKeepAlive: true,
		},
		// 群组管理
		{
			ParentId:    1,
			Name:        "群组管理",
			Icon:        "Menu",
			Router:      "/redpacket/group",
			Component:   "@/views/redpacket/group/index.vue",
			Sort:        4,
			Enable:      true,
			IsHide:      false,
			IsFull:      false,
			IsAffix:     false,
			IsKeepAlive: true,
		},
		// 系统配置
		{
			ParentId:    0,
			Name:        "系统配置",
			Icon:        "Setting",
			Router:      "/system",
			Sort:        100,
			Enable:      true,
			IsHide:      false,
			IsFull:      false,
			IsAffix:     false,
			IsKeepAlive: true,
		},
	}
	for _, req := range reqs {
		find := Menus{}
		if stat := db.Model(&Menus{}).Where("router = ?", req.Router).Find(&find).Statement; stat.RowsAffected == 0 {
			if err := db.Create(&req).Error; err != nil {
				return err
			}
		}
	}

	// 初始化 系统配置 子菜单
	sysParent := Menus{}
	db.Model(&Menus{}).Where("router = ?", "/system").First(&sysParent)
	if sysParent.ID != 0 {
		sysChildren := []Menus{
			{Name: "账号管理", Icon: "Menu", Router: "/system/accountManage", Component: "@/views/system/accountManage/index.vue", Sort: 1, Enable: true, IsKeepAlive: true},
			{Name: "角色管理", Icon: "Menu", Router: "/system/roleManage", Component: "@/views/system/roleManage/index.vue", Sort: 2, Enable: true, IsKeepAlive: true},
			{Name: "菜单管理", Icon: "Menu", Router: "/system/menuMange", Component: "@/views/system/menuMange/index.vue", Sort: 3, Enable: true, IsKeepAlive: true},
			{Name: "APIs管理", Icon: "Menu", Router: "/system/apisManage", Component: "@/views/system/apisManage/index.vue", Sort: 4, Enable: true, IsKeepAlive: true},
		}
		for _, child := range sysChildren {
			find := Menus{}
			if stat := db.Model(&Menus{}).Where("router = ?", child.Router).Find(&find).Statement; stat.RowsAffected == 0 {
				child.ParentId = sysParent.ID
				if err := db.Create(&child).Error; err != nil {
					return err
				}
			}
		}
	}

	// 初始化 Cloudflare 菜单（查找或创建父菜单，再补充子菜单）
	cfParent := Menus{}
	db.Model(&Menus{}).Where("router = ?", "/cloudflare").First(&cfParent)
	if cfParent.ID == 0 {
		cfParent = Menus{
			ParentId:    0,
			Name:        "Cloudflare",
			Icon:        "Monitor",
			Router:      "/cloudflare",
			Sort:        2,
			Enable:      true,
			IsKeepAlive: true,
		}
		if err := db.Create(&cfParent).Error; err != nil {
			return err
		}
	}
	cfChildren := []Menus{
		{
			Name:        "账号管理",
			Icon:        "Key",
			Router:      "/cloudflare/account",
			Component:   "@/views/cloudflare/account/index.vue",
			Sort:        1,
			Enable:      true,
			IsKeepAlive: true,
		},
		{
			Name:        "Zone 列表",
			Icon:        "List",
			Router:      "/cloudflare/zones",
			Component:   "@/views/cloudflare/zones/index.vue",
			Sort:        2,
			Enable:      true,
			IsKeepAlive: true,
		},
		{
			Name:        "DNS 管理",
			Icon:        "Connection",
			Router:      "/cloudflare/dns",
			Component:   "@/views/cloudflare/dns/index.vue",
			Sort:        3,
			Enable:      true,
			IsKeepAlive: true,
		},
		{
			Name:        "DNS 模板",
			Icon:        "Document",
			Router:      "/cloudflare/template",
			Component:   "@/views/cloudflare/template/index.vue",
			Sort:        4,
			Enable:      true,
			IsKeepAlive: true,
		},
	}
	for _, child := range cfChildren {
		find := Menus{}
		if stat := db.Model(&Menus{}).Where("router = ?", child.Router).Find(&find).Statement; stat.RowsAffected == 0 {
			child.ParentId = cfParent.ID
			if err := db.Create(&child).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

const AdminId int64 = 1

func (r *Administrator) DataInit(db *gorm.DB) error {
	// 初始化用户等级
	password, salt := r.EncodePassword("666666")
	reqs := []Administrator{
		{
			ID:            AdminId,
			Nickname:      "超级管理员",
			Username:      "superman",
			Password:      password,
			Salt:          salt,
			GoogleKey:     "",
			Avatar:        "",
			RoleId:        1,
			LastLoginTime: lo.ToPtr(time.Now()),
		},
	}
	for _, req := range reqs {
		find := NewAdministrator()
		if stat := db.Model(&Administrator{}).Where("id", AdminId).
			Find(&find).Statement; stat.RowsAffected == 0 {
			if err := db.Create(&req).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
