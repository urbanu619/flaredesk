package system

import (
	"github.com/demdxx/gocast"
	"github.com/gin-gonic/gin"
	"go_server/model/common/response"
	system2 "go_server/model/system"
	"go_server/service/base"
	"strings"
)

type ApisService struct {
	base.SysCommonService
}

func (s *ApisService) Tree(c *gin.Context) {
	roleId := c.GetInt64("roleId")
	sysApis, _ := system2.NewApis().SysTree(s.DB())
	role, ok := base.GetOne[system2.Role](s.DB(), "id", roleId)
	if !ok {
	}
	personApis := sysApis
	if role.Apis != "*" {
		apisSplit := strings.Split(role.Apis, ",")
		personApis, _ = system2.NewApis().GetUserTree(s.DB(), apisSplit)
	}
	res := make(map[string]interface{})
	res["sysApis"] = sysApis       // 系统总树
	res["personApis"] = personApis // 当前角色权限[]
	response.Resp(c, res)
	return
}

// 设置

func (s *ApisService) Set(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if userId != system2.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		Id       interface{} `json:"id"  validate:"required"`
		ParentId interface{} `json:"parentId"`
		Group    string      `json:"group" gorm:"column:group;type:varchar(50);comment:分组名称"`
		Name     string      `json:"name" gorm:"column:name;type:varchar(50);comment:接口名称"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if req.Id == 0 {
		response.Resp(c, "API不存在")
		return
	}
	menu, ok := base.GetOne[system2.Apis](s.DB(), "id", req.Id)
	if !ok {
		response.Resp(c, "API不存在")
		return
	}
	parentId := gocast.ToInt64(req.ParentId)
	if parentId != 0 {
		_, ok := base.GetOne[system2.Apis](s.DB(), "id", parentId)
		if !ok {
			response.Resp(c, "上级API不存在")
			return
		}
	}
	if req.Group != "" {
		menu.Group = req.Group
	}
	if req.Name != "" {
		menu.Name = req.Name
	}
	if parentId != menu.ParentId {
		menu.ParentId = parentId
	}
	if err := s.DB().Model(&system2.Apis{}).
		Where("id", req.Id).
		Updates(system2.Apis{
			ParentId: menu.ParentId,
			Name:     menu.Name,
			Group:    menu.Group,
		}).Error; err != nil {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	response.Resp(c)
	return
}

// 创建分组

func (s *ApisService) Create(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if userId != system2.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		ParentId interface{} `json:"parentId"`
		Group    string      `json:"group" gorm:"column:group;type:varchar(50);comment:分组名称"`
		Name     string      `json:"name" gorm:"column:name;type:varchar(50);comment:接口名称"`
		Path     string      `json:"path" gorm:"column:path;type:varchar(50);comment:分组路径"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if req.Name == "" || req.Group == "" || req.Path == "" {
		response.Resp(c, "基础信息填写不完整")
		return
	}
	_, ok := base.GetOne[system2.Apis](s.DB(), "path", req.Path)
	if ok {
		response.Resp(c, "分组已存在")
		return
	}
	parentId := gocast.ToInt64(req.ParentId)
	if parentId != 0 {
		_, ok := base.GetOne[system2.Apis](s.DB(), "id", parentId)
		if !ok {
			response.Resp(c, "上级API不存在")
			return
		}
	}
	row := &system2.Apis{
		Name:     req.Name,
		Group:    req.Group,
		Path:     req.Path,
		ParentId: parentId,
	}
	if err := s.DB().Create(&row).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c)
	return
}

func (s *ApisService) Find(c *gin.Context) {
	type request[T any] struct {
		base.ListRequest[T]
		Id       *interface{} `form:"id"`
		ParentId interface{}  `json:"parentId"`
	}
	req := new(request[system2.Apis])
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
	resp, err := base.NewQueryBaseHandler(system2.NewApis()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)
	//base.NewBaseHandler(model.NewApis()).List(c, s.DB())
}

func (s *ApisService) Get(c *gin.Context) {
	base.NewBaseHandler(system2.NewApis()).Get(c, s.DB())
}

func (s *ApisService) Del(c *gin.Context) {
	base.NewBaseHandler(system2.NewApis()).DeleteOne(c, s.DB())
}
