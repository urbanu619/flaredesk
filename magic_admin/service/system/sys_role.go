package system

import (
	"github.com/demdxx/gocast"
	"github.com/gin-gonic/gin"
	"go_server/base/core"
	"go_server/model/common/response"
	system2 "go_server/model/system"
	"go_server/service/base"
	"strings"
)

type RoleService struct {
	base.SysCommonService
}

// 设置用户

func (s *RoleService) Set(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if userId != system2.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		Id    interface{}   `json:"id"`
		Name  string        `json:"name"  validate:"required"`
		Desc  string        `json:"desc"  validate:"required"`
		Apis  []interface{} `json:"apis"`
		Menus []interface{} `json:"menus"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	id := gocast.ToInt64(req.Id)
	if id == system2.AdminId {
		response.Resp(c, "超人角色 不许乱改")
		return
	}
	if req.Name == "" || len(req.Apis) == 0 || len(req.Menus) == 0 || id == 0 {
		response.Resp(c, "用户基础信息填写不完整")
		return
	}
	_, ok := base.GetOne[system2.Role](s.DB(), "id", id)
	if !ok {
		response.Resp(c, "角色不存在")
		return
	}
	apis := make([]string, 0)
	if len(req.Apis) > 0 {
		for _, item := range req.Apis {
			apis = append(apis, gocast.ToString(item))
		}
	}
	menus := make([]string, 0)
	if len(req.Menus) > 0 {
		for _, item := range req.Menus {
			menus = append(menus, gocast.ToString(item))
		}
	}
	if err := s.DB().Model(&system2.Role{}).
		Where("id", id).
		Updates(map[string]interface{}{
			"name":  req.Name, // 头像设置
			"desc":  req.Desc,
			"apis":  strings.Join(apis, ","),
			"menus": strings.Join(menus, ","),
		}).Error; err != nil {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	response.Resp(c)
	return
}

// 创建用户

func (s *RoleService) Create(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if userId != system2.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		Name string `json:"name"  validate:"required"`
		Desc string `json:"desc"  validate:"required"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if req.Name == "" {
		response.Resp(c, "用户基础信息填写不完整")
		return
	}
	_, ok := base.GetOne[system2.Role](s.DB(), "name", req.Name)
	if ok {
		response.Resp(c, "角色已存在")
		return
	}
	row := &system2.Role{
		Name: req.Name,
		Desc: req.Desc,
	}
	if err := s.DB().Create(&row).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c)
	return
}

// 删除

func (s *RoleService) Del(c *gin.Context) {
	adminID := c.GetInt64("userId")
	roleID := c.GetInt64("roleId")
	core.Log.Infof("操作用户ID:%d 角色ID:%d", adminID, roleID)
	if adminID != system2.AdminId {
		response.Resp(c, "非开发人员禁止操作")
		return
	}
	type request struct {
		Id int64 `json:"id" form:"id"`
	}
	req := new(request)
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	roleId := gocast.ToInt64(req.Id)
	if roleId == 0 {
		response.Resp(c, "id is zero")
		return
	}
	if roleId == system2.AdminId {
		response.Resp(c, "管理员角色禁止删除")
		return
	}
	// 如果存在该类型角色 不可删除
	if count := base.CountRow[system2.Role](s.DB(), "role_id", roleId); count > 0 {
		response.Resp(c, "存在该角色用户 不允许删除")
		return
	}

	// 先通过Id获取记录'
	if err := s.DB().Debug().First(&system2.Role{}, req.Id).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	exCmd := s.DB().Where("id", req.Id).Delete(&system2.Role{})
	if exCmd.RowsAffected != 1 {
		response.Resp(c, "delete fail")
		return
	}
	if exCmd.Error != nil {
		response.Resp(c, exCmd.Error.Error())
		return
	}
	response.Resp(c)
}

func (s *RoleService) Find(c *gin.Context) {
	//type request[T any] struct {
	//	base.ListRequest[T]
	//}
	//req := new(request[model.Role])
	//if err := c.BindQuery(req); err != nil {
	//	response.Resp(c, err.Error())
	//	return
	//}
	db := s.DB()
	roles, err := base.GetMore[system2.Role](db)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	type RoleInfo struct {
		ID    int64    `json:"id" gorm:"primarykey;comment:id"`
		Name  string   `json:"name" gorm:"column:name;type:varchar(100);comment:角色名称"`
		Apis  []string `json:"auths" gorm:"column:apis;type:text;comment:Api权限"`
		Menus []string `json:"menus" gorm:"column:menus;type:text;comment:menu权限"`
		Desc  string   `json:"desc" gorm:"column:desc;type:varchar(200);comment:角色权限描述;"`
	}
	list := make([]*RoleInfo, 0)
	for _, item := range roles {
		list = append(list, &RoleInfo{
			ID:    item.ID,
			Name:  item.Name,
			Apis:  strings.Split(item.Apis, ","),
			Menus: strings.Split(item.Menus, ","),
			Desc:  item.Desc,
		})
	}
	//paging := base.NewPagination()
	response.Resp(c, map[string]interface{}{
		"list": list,
		//"paging": paging,
	})
}

func (s *RoleService) Get(c *gin.Context) {
	base.NewBaseHandler(system2.NewRole()).Get(c, s.DB())
}
