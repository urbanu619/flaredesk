package app

import (
	"time"

	"github.com/demdxx/gocast"
	"github.com/gin-gonic/gin"
	"go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
)

type CfDnsTemplateService struct {
	base.BizCommonService
}

// Find 查询模板列表
func (s *CfDnsTemplateService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		Name string `form:"name"`
	}
	req := new(request[app.CfDnsTemplate])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	db := s.DB()
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	db = db.Order("created_at DESC")
	resp, err := base.NewQueryBaseHandler(app.NewCfDnsTemplate()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)
}

// Get 获取单个模板
func (s *CfDnsTemplateService) Get(c *gin.Context) {
	s.SetDbAlias("app")
	base.NewBaseHandler(app.NewCfDnsTemplate()).Get(c, s.DB())
}

// Create 创建模板
func (s *CfDnsTemplateService) Create(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		Name    string `json:"name" binding:"required"`
		Remark  string `json:"remark"`
		Records string `json:"records" binding:"required"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	now := time.Now().Unix()
	tpl := &app.CfDnsTemplate{
		Name:      req.Name,
		Remark:    req.Remark,
		Records:   req.Records,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.DB().Create(tpl).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, tpl)
}

// Update 更新模板
func (s *CfDnsTemplateService) Update(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		Id      int64  `json:"id" binding:"required"`
		Name    string `json:"name"`
		Remark  string `json:"remark"`
		Records string `json:"records"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	tpl, ok := base.GetOne[app.CfDnsTemplate](s.DB(), "id", req.Id)
	if !ok {
		response.Resp(c, "模板不存在")
		return
	}
	updates := map[string]interface{}{"updated_at": time.Now().Unix()}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if req.Records != "" {
		updates["records"] = req.Records
	}
	if err := s.DB().Model(tpl).Updates(updates).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c)
}

// Delete 删除模板
func (s *CfDnsTemplateService) Delete(c *gin.Context) {
	s.SetDbAlias("app")
	id := gocast.ToInt64(c.Query("id"))
	if id == 0 {
		response.Resp(c, "id 必填")
		return
	}
	tpl, ok := base.GetOne[app.CfDnsTemplate](s.DB(), "id", id)
	if !ok {
		response.Resp(c, "模板不存在")
		return
	}
	if err := s.DB().Delete(tpl).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c)
}
