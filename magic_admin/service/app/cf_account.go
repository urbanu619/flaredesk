package app

import (
	"time"

	"github.com/demdxx/gocast"
	"github.com/gin-gonic/gin"
	"go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
)

type CfAccountService struct {
	base.BizCommonService
}

// Find 查询账号列表
func (s *CfAccountService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		Name   string `form:"name" json:"name"`
		Status *int8  `form:"status" json:"status"`
	}
	req := new(request[app.CfAccount])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	db := s.DB()
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	db = db.Order("created_at DESC")

	resp, err := base.NewQueryBaseHandler(app.NewCfAccount()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)
}

// Get 获取账号详情
func (s *CfAccountService) Get(c *gin.Context) {
	s.SetDbAlias("app")
	base.NewBaseHandler(app.NewCfAccount()).Get(c, s.DB())
}

// Create 创建账号
func (s *CfAccountService) Create(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email"`
		ApiToken string `json:"apiToken" binding:"required"`
		Remark   string `json:"remark"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	now := time.Now().Unix()
	account := &app.CfAccount{
		Name:      req.Name,
		Email:     req.Email,
		ApiToken:  req.ApiToken,
		Status:    1,
		Remark:    req.Remark,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.DB().Create(account).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, account)
}

// Update 更新账号
func (s *CfAccountService) Update(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		Id       int64  `json:"id" binding:"required"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		ApiToken string `json:"apiToken"`
		Status   *int8  `json:"status"`
		Remark   string `json:"remark"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	account, ok := base.GetOne[app.CfAccount](s.DB(), "id", req.Id)
	if !ok {
		response.Resp(c, "账号不存在")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.ApiToken != "" {
		updates["api_token"] = req.ApiToken
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	updates["updated_at"] = time.Now().Unix()

	if err := s.DB().Model(&account).Updates(updates).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c)
}

// Delete 删除账号
func (s *CfAccountService) Delete(c *gin.Context) {
	s.SetDbAlias("app")
	id, ok := c.GetQuery("id")
	if !ok {
		response.Resp(c, "未填写ID")
		return
	}
	accountId := gocast.ToInt64(id)
	if accountId == 0 {
		response.Resp(c, "ID无效")
		return
	}
	account, ok := base.GetOne[app.CfAccount](s.DB(), "id", accountId)
	if !ok {
		response.Resp(c, "账号不存在")
		return
	}
	if err := s.DB().Delete(&account).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c)
}
