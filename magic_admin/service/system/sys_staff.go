package system

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_server/base/core"
	"go_server/model/common/response"
	system2 "go_server/model/system"
	"go_server/service/base"
	"strings"
)

// 用户登陆时候获取基本信息

func (s *UserService) Info(c *gin.Context) {
	userId := c.GetInt64("userId")
	if userId == 0 {
		response.Resp(c)
		return
	}
	user, ok := base.GetOne[system2.Administrator](s.DB(), "id", userId)
	if !ok {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	role, ok := base.GetOne[system2.Role](s.DB(), "id", user.RoleId)
	if !ok {
		response.Resp(c, "角色信息不存在")
		return
	}
	res := make(map[string]interface{})
	res["avatar"] = user.Avatar
	res["nickName"] = user.Nickname
	res["roleName"] = role.Name
	sysMenus, err := system2.NewMenus().SysTree(s.DB())
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	res["sysMenus"] = sysMenus // 系统菜单
	// 前端菜单
	if role.Menus == "*" {
		res["menus"] = sysMenus // 系统菜单
	} else {
		menuIds := strings.Split(role.Menus, ",")
		menus, err := system2.NewMenus().GetUserTree(s.DB(), menuIds)
		if err != nil {
			response.Resp(c, "菜单信息不存在")
			return
		}
		res["menus"] = menus // 用户菜单
	}
	sysApis, err := system2.NewApis().SysTree(s.DB())
	if err != nil {
		response.Resp(c, "菜单信息不存在")
		return
	}
	res["sysApis"] = sysApis // 系统菜单

	// 个人api
	if role.Apis == "*" {
		res["apis"] = "*" // 开放全部菜单
	} else {
		apiIds := strings.Split(role.Apis, ",")
		apis, err := system2.NewApis().GetUserTree(s.DB(), apiIds)
		if err != nil {
			response.Resp(c, "菜单信息不存在")
			return
		}
		res["apis"] = apis // 用户api权限
	}
	response.Resp(c, res)
}

// 用户信息设置

func (s *UserService) Set(c *gin.Context) {
	type request struct {
		Avatar   string `json:"avatar"`
		Nickname string `json:"nickname"`
		Password string `json:"password"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	userId := c.GetInt64("userId")

	user, ok := base.GetOne[system2.Administrator](s.DB(), "id", userId)
	if !ok {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Password != "" {
		user.Password, user.Salt = user.EncodePassword(req.Password)
	}
	if err := s.DB().Model(&system2.Administrator{}).
		Where("id", userId).
		Updates(map[string]interface{}{
			"avatar":   user.Avatar, // 头像设置
			"nickname": user.Nickname,
			"password": user.Password,
			"salt":     user.Salt,
		}).Error; err != nil {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	response.Resp(c, map[string]interface{}{
		"avatar":   user.Avatar, // 头像设置
		"nickname": user.Nickname,
	})
}

// 获取谷歌密钥

func (s *UserService) GetGoogleKey(c *gin.Context) {
	userId := c.GetInt64("userId")
	googleSecret := core.NewGoogleAuth().GetSecret()[:16]
	res := make(map[string]interface{})
	res["googleKey"] = googleSecret
	res["qrcode"] = core.NewGoogleAuth().GetQrcode(fmt.Sprintf("AmsAdmin-(%d)", userId),
		googleSecret)
	response.Resp(c, res)
}

// 设置/重置谷歌密钥

func (s *UserService) ReplaceGoogleKey(c *gin.Context) {
	type request struct {
		GoogleKey          string `json:"googleKey" validate:"required"`
		GoogleCode         string `json:"googleCode" validate:"required"`
		OriginalGoogleCode string `json:"originalGoogleCode"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	if req.GoogleCode == "" || req.GoogleKey == "" {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	userId := c.GetInt64("userId")
	user, ok := base.GetOne[system2.Administrator](s.DB(), "id", userId)
	if !ok {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	if user.GoogleKey != "" {
		if req.OriginalGoogleCode == "" {
			response.Resp(c, response.ResponseCodeFailure)
			return
		}
		if ok, _ := core.NewGoogleAuth().VerifyCode(user.GoogleKey, req.OriginalGoogleCode); !ok {
			response.Resp(c, "谷歌验证码错误")
			return
		}
	}
	//  保存谷歌秘钥
	if err := s.DB().Model(&system2.Administrator{}).
		Where("id", userId).
		Updates(system2.Administrator{
			GoogleKey: req.GoogleKey,
		}).Error; err != nil {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	response.Resp(c)
}

// 取消谷歌验证码

func (s *UserService) CancelGoogleKey(c *gin.Context) {
	type request struct {
		GoogleCode string `json:"googleCode" validate:"required"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	userId := c.GetInt64("userId")
	user, ok := base.GetOne[system2.Administrator](s.DB(), "id", userId)
	if !ok {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	if user.GoogleKey == "" {
		response.Resp(c, "未设置谷歌密钥")
		return
	}
	if req.GoogleCode == "" {
		response.Resp(c, "请输入谷歌验证码")
		return
	}
	if ok, _ := core.NewGoogleAuth().VerifyCode(user.GoogleKey, req.GoogleCode); !ok {
		response.Resp(c, "验证码错误")
		return
	}
	//  保存谷歌秘钥
	if err := s.DB().Model(&system2.Administrator{}).
		Where("id", userId).
		Updates(system2.Administrator{
			GoogleKey: "",
		}).Error; err != nil {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	response.Resp(c)
}
