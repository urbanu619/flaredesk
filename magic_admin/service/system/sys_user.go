package system

import (
	"github.com/demdxx/gocast"
	"github.com/gin-gonic/gin"
	"go_server/base/core"
	"go_server/model/common/response"
	system2 "go_server/model/system"
	"go_server/service/base"
)

type UserService struct {
	base.SysCommonService
}

func (s *UserService) Logs(c *gin.Context) {
	base.NewBaseHandler(system2.NewAdministratorLog()).List(c, s.DB())
}

// 获取用户列表

func (s *UserService) Find(c *gin.Context) {
	db := s.DB()
	users, err := base.GetMore[system2.Administrator](db)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	//type RoleInfo struct {
	//	ID    int64    `json:"id" gorm:"primarykey;comment:id"`
	//	Name  string   `json:"name" gorm:"column:name;type:varchar(100);comment:角色名称"`
	//	Apis  []string `json:"auths" gorm:"column:apis;type:text;comment:Api权限"`
	//	Menus []string `json:"menus" gorm:"column:menus;type:text;comment:menu权限"`
	//	Desc  string   `json:"desc" gorm:"column:desc;type:varchar(200);comment:角色权限描述;"`
	//}
	//list := make([]*RoleInfo, 0)
	//for _, item := range roles {
	//	list = append(list, &RoleInfo{
	//		ID:    item.ID,
	//		Name:  item.Name,
	//		Apis:  strings.Split(item.Apis, ","),
	//		Menus: strings.Split(item.Menus, ","),
	//		Desc:  item.Desc,
	//	})
	//}
	//paging := base.NewPagination()
	response.Resp(c, map[string]interface{}{
		"list": users,
		//"paging": paging,
	})
	//base.NewBaseHandler(model.NewAdministrator()).List(c, s.DB())
}

// 查看用户详情

func (s *UserService) Get(c *gin.Context) {
	base.NewBaseHandler(system2.NewAdministrator()).Get(c, s.DB())
}

// 设置用户

func (s *UserService) SetUser(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if userId != system2.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		UserId         interface{} `json:"userId"`
		Avatar         string      `json:"avatar"`
		Nickname       string      `json:"nickname"`
		Password       string      `json:"password"`
		Lock           interface{} `json:"lock"`
		ResetGoogleKey interface{} `json:"resetGoogleKey"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	user, ok := base.GetOne[system2.Administrator](s.DB(), "id", req.UserId)
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
	resetGoogleKey := gocast.ToBool(req.ResetGoogleKey)
	if user.GoogleKey != "" && resetGoogleKey {
		user.GoogleKey = ""
	}
	if req.Lock != nil {
		user.Lock = gocast.ToBool(req.Lock)
	}
	if err := s.DB().Model(&system2.Administrator{}).
		Where("id", user.ID).
		Updates(map[string]interface{}{
			"avatar":     user.Avatar, // 头像设置
			"nickname":   user.Nickname,
			"password":   user.Password,
			"salt":       user.Salt,
			"google_key": user.GoogleKey,
			"lock":       user.Lock,
		}).Error; err != nil {
		response.Resp(c, response.ResponseCodeFailure)
		return
	}
	response.Resp(c)
	return
}

// 创建用户

func (s *UserService) Create(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if userId != system2.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		Avatar   string      `json:"avatar"`
		Username string      `json:"username"`
		Nickname string      `json:"nickname"`
		Password string      `json:"password"`
		RoleId   interface{} `json:"roleId"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	roleId := gocast.ToInt64(req.RoleId)
	if req.Username == "" || req.Nickname == "" || req.Password == "" || roleId == 0 {
		response.Resp(c, "用户基础信息填写不完整")
		return
	}
	if len(req.Password) < 6 {
		response.Resp(c, "密码长度不得小与6位数")
		return
	}
	_, ok := base.GetOne[system2.Role](s.DB(), "id", roleId)
	if !ok {
		response.Resp(c, "角色信息不存在")
		return
	}
	_, ok = base.GetOne[system2.Administrator](s.DB(), "username", req.Username)
	if ok {
		response.Resp(c, "用户已存在请更好用户名")
		return
	}
	row := &system2.Administrator{
		Nickname: req.Nickname,
		Username: req.Username,
		RoleId:   roleId,
	}
	row.Password, row.Salt = row.EncodePassword(req.Password)
	if err := s.DB().Create(&row).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c)
	return
}

// 删除用户

func (s *UserService) Del(c *gin.Context) {
	adminID := c.GetInt64("userId")
	roleID := c.GetInt64("roleId")
	core.Log.Infof("操作用户ID:%d 角色ID:%d", adminID, roleID)
	if adminID != system2.AdminId {
		response.Resp(c, "非技术人员禁止操作")
		return
	}
	type request struct {
		Id interface{} `json:"id" form:"id"`
	}
	req := new(request)
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	userID := gocast.ToInt64(req.Id)
	if userID == 0 {
		response.Resp(c, "id is zero")
		return
	}
	if userID == system2.AdminId {
		response.Resp(c, "ROOT账号禁止删除")
		return
	}
	// 先通过Id获取记录'
	if err := s.DB().Debug().First(&system2.Administrator{}, req.Id).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	exCmd := s.DB().Where("id", req.Id).Delete(&system2.Administrator{})
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
