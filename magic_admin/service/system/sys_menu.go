package system

import (
	"github.com/demdxx/gocast"
	"github.com/gin-gonic/gin"
	"go_server/base/core"
	"go_server/model/common/response"
	model "go_server/model/system"
	"go_server/service/base"
	"strings"
)

type MenuService struct {
	base.SysCommonService
}

func (s *MenuService) Tree(c *gin.Context) {
	roleId := c.GetInt64("roleId")
	sysMenu, _ := model.NewMenus().SysTree(s.DB())
	role, ok := base.GetOne[model.Role](s.DB(), "id", roleId)
	if !ok {
		core.Log.Infof("roleId is not found:%+v", roleId)
		response.Resp(c, "roleId is not found")
		return
	}
	personMenu := sysMenu
	if role.Menus != "*" {
		core.Log.Infof("role.Menus:%+v", role.Menus)
		menuSplit := strings.Split(role.Menus, ",")
		personMenus, err := model.NewMenus().GetUserTree(s.DB(), menuSplit)
		if err != nil {
			response.Resp(c, err.Error())
			return
		}
		personMenu = personMenus
	}
	res := make(map[string]interface{})
	res["sysMenu"] = sysMenu       // 系统总树
	res["personMenu"] = personMenu // 当前角色权限[]
	response.Resp(c, res)
	return
}

// 设置用户

func (s *MenuService) Set(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if userId != model.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		Id          interface{} `json:"id"  validate:"required"`
		ParentId    interface{} `json:"parentId"`
		Name        string      `json:"name" validate:"required"`
		Icon        string      `json:"icon"`
		Router      string      `json:"router"`
		IsHide      interface{} `json:"isHide"`
		IsFull      interface{} `json:"isFull"`
		IsAffix     interface{} `json:"isAffix"`
		IsKeepAlive interface{} `json:"isKeepAlive"`
		PageName    string      `json:"pageName" gorm:"type:varchar(200);comment:页面名称"`
		Component   string      `json:"component" gorm:"type:varchar(200);comment:组件"`
		Sort        interface{} `json:"sort"`
		Enable      interface{} `json:"enable" gorm:"comment:是否启用"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if req.Id == 0 {
		response.Resp(c, "菜单不存在")
		return
	}
	menu, ok := base.GetOne[model.Menus](s.DB(), "id", req.Id)
	if !ok {
		response.Resp(c, "菜单不存在")
		return
	}
	parentId := gocast.ToInt64(req.ParentId)
	if parentId != 0 {
		_, ok := base.GetOne[model.Menus](s.DB(), "id", parentId)
		if !ok {
			response.Resp(c, "上级菜单不存在")
			return
		}
	}
	// 限制主菜单不能修改为子菜单
	if menu.ParentId == 0 && parentId != 0 {
		response.Resp(c, "限制主菜单不能修改为子菜单")
		return
	}
	if req.Router != "" {
		menu.Router = req.Router
	}
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if parentId != 0 {
		menu.ParentId = parentId
	}
	if req.Sort != nil {
		menu.ParentId = parentId
	}
	if req.PageName != "" {
		menu.PageName = req.PageName
	}
	if req.Component != "" {
		menu.Component = req.Component
	}
	if req.Enable != nil {
		menu.Enable = gocast.ToBool(req.Enable)
	}
	if req.IsHide != nil {
		menu.IsHide = gocast.ToBool(req.IsHide)
	}
	if req.IsFull != nil {
		menu.IsFull = gocast.ToBool(req.IsFull)
	}
	if req.IsAffix != nil {
		menu.IsAffix = gocast.ToBool(req.IsAffix)
	}
	if req.IsKeepAlive != nil {
		menu.IsKeepAlive = gocast.ToBool(req.IsKeepAlive)
	}
	if err := s.DB().Model(&model.Menus{}).
		Where("id", req.Id).
		Updates(model.Menus{
			ParentId:    menu.ParentId,
			Name:        menu.Name,
			Icon:        menu.Icon,
			Router:      menu.Router,
			IsHide:      menu.IsHide,
			IsFull:      menu.IsFull,
			IsAffix:     menu.IsAffix,
			IsKeepAlive: menu.IsKeepAlive,
			PageName:    menu.PageName,
			Component:   menu.Component,
			Sort:        gocast.ToInt(req.Sort),
			Enable:      menu.Enable, //     interface{} `json:"enable" gorm:"comment:是否启用"`

		}).Error; err != nil {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	response.Resp(c)
	return
}

// 创建菜单

func (s *MenuService) Create(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if userId != model.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		ParentId    interface{} `json:"parentId"`
		Name        string      `json:"name" validate:"required"`
		Icon        string      `json:"icon"`
		Router      string      `json:"router"`
		IsHide      interface{} `json:"isHide"`
		IsFull      interface{} `json:"isFull"`
		IsAffix     interface{} `json:"isAffix"`
		IsKeepAlive interface{} `json:"isKeepAlive"`
		PageName    string      `json:"pageName" gorm:"type:varchar(200);comment:页面名称"`
		Component   string      `json:"component" gorm:"type:varchar(200);comment:组件"`
		Sort        interface{} `json:"sort"`
		Enable      interface{} `json:"enable" gorm:"comment:是否启用"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	parentId := gocast.ToInt64(req.ParentId)
	if req.Name == "" {
		response.Resp(c, "菜单名称未设置")
		return
	}
	if parentId != 0 {
		if _, ok := base.GetOne[model.Menus](s.DB(), "id", parentId); !ok {
			response.Resp(c, "上级菜单不存在")
			return
		}
	}

	if base.CountRow[model.Menus](s.DB(), "router", req.Router) > 0 {
		response.Resp(c, "路由不可重复")
		return
	}
	row := &model.Menus{
		ParentId:    parentId,
		Name:        req.Name,
		Icon:        req.Icon,
		Router:      req.Router,
		IsHide:      gocast.ToBool(req.IsHide),
		IsFull:      gocast.ToBool(req.IsFull),
		IsAffix:     gocast.ToBool(req.IsAffix),
		IsKeepAlive: gocast.ToBool(req.IsKeepAlive),
		PageName:    req.PageName,
		Component:   req.Component,
		Sort:        gocast.ToInt(req.Sort),
		Enable:      gocast.ToBool(req.Enable),
	}
	if err := s.DB().Create(&row).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c)
	return
}

// 删除菜单 同时删除下级菜单

func (s *MenuService) Del(c *gin.Context) {
	adminID := c.GetInt64("userId")
	roleID := c.GetInt64("roleId")
	core.Log.Infof("操作用户ID:%d 角色ID:%d", adminID, roleID)
	if adminID != model.AdminId {
		response.Resp(c, "非技术人员禁止操作")
		return
	}
	id, ok := c.GetQuery("id")
	if !ok {
		response.Resp(c, "未填写ID")
		return
	}
	rowId := gocast.ToInt64(id)
	if rowId == 0 {
		response.Resp(c, "id is zero")
		return
	}
	// 先通过Id获取记录'
	if err := s.DB().First(&model.Menus{}, rowId).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	exCmd := s.DB().Where("id", id).Delete(&model.Menus{})
	if exCmd.RowsAffected != 1 {
		response.Resp(c, "delete fail")
		return
	}
	if exCmd.Error != nil {
		response.Resp(c, exCmd.Error.Error())
		return
	}
	err := s.DB().Where("parent_id", id).Delete(&model.Menus{}).Error
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c)
}

func (s *MenuService) Find(c *gin.Context) {
	type request[T any] struct {
		base.ListRequest[T]
		Id       *interface{} `form:"id"`
		ParentId interface{}  `json:"parentId" form:"parentId"`
	}
	req := new(request[model.Menus])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	db := s.DB()
	if req.Id != nil && *req.Id != 0 {
		db = db.Where("id", req.Id)
	}
	if req.ParentId != nil && gocast.ToInt64(req.ParentId) != 0 {
		db = db.Where("parent_id", req.ParentId)
	}
	resp, err := base.NewQueryBaseHandler(model.NewMenus()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)
}

func (s *MenuService) Get(c *gin.Context) {
	base.NewBaseHandler(model.NewMenus()).Get(c, s.DB())
}
