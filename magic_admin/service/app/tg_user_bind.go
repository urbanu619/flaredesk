package app

import (
	"time"

	"github.com/demdxx/gocast"
	"github.com/gin-gonic/gin"
	"go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
)

type TgUserBindService struct {
	base.BizCommonService
}

// Find 查询用户绑定列表
func (s *TgUserBindService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		UserId           *int64  `form:"user_id" json:"user_id"`
		TelegramId       *int64  `form:"telegram_id" json:"telegram_id"`
		TelegramUsername string  `form:"telegram_username" json:"telegram_username"`
		BindStatus       *int    `form:"bind_status" json:"bind_status"`
	}
	req := new(request[app.TgUserBind])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	db := s.DB()

	// 条件筛选
	if req.UserId != nil && *req.UserId != 0 {
		db = db.Where("user_id = ?", *req.UserId)
	}
	if req.TelegramId != nil && *req.TelegramId != 0 {
		db = db.Where("telegram_id = ?", *req.TelegramId)
	}
	if req.TelegramUsername != "" {
		db = db.Where("telegram_username LIKE ?", "%"+req.TelegramUsername+"%")
	}
	if req.BindStatus != nil {
		db = db.Where("bind_status = ?", *req.BindStatus)
	}

	// 按创建时间倒序
	db = db.Order("created_at DESC")

	resp, err := base.NewQueryBaseHandler(app.NewTgUserBind()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)
}

// Get 获取用户绑定详情
func (s *TgUserBindService) Get(c *gin.Context) {
	s.SetDbAlias("app")
	base.NewBaseHandler(app.NewTgUserBind()).Get(c, s.DB())
}

// Create 创建用户绑定
func (s *TgUserBindService) Create(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		UserId            interface{} `json:"userId"`
		TelegramId        interface{} `json:"telegramId"`
		TelegramUsername  string      `json:"telegramUsername"`
		TelegramFirstName string      `json:"telegramFirstName"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	userId := gocast.ToInt64(req.UserId)
	telegramId := gocast.ToInt64(req.TelegramId)

	// 验证必填字段
	if userId == 0 {
		response.Resp(c, "平台用户ID不能为空")
		return
	}
	if req.TelegramUsername == "" {
		response.Resp(c, "Telegram用户名不能为空")
		return
	}

	// 检查是否已绑定
	var count int64
	s.DB().Model(&app.TgUserBind{}).Where("telegram_username = ? AND bind_status = 1", req.TelegramUsername).Count(&count)
	if count > 0 {
		response.Resp(c, "该Telegram用户名已被绑定")
		return
	}

	// 检查 UserId 是否已绑定
	s.DB().Model(&app.TgUserBind{}).Where("user_id = ? AND bind_status = 1", userId).Count(&count)
	if count > 0 {
		response.Resp(c, "该平台用户已绑定其他Telegram账号")
		return
	}

	now := time.Now()
	bind := &app.TgUserBind{
		UserId:            userId,
		TelegramId:        telegramId,
		TelegramUsername:  req.TelegramUsername,
		TelegramFirstName: req.TelegramFirstName,
		BindStatus:        1,
		BindTime:          now,
		CreatedAt:         now.Unix(),
		UpdatedAt:         now.Unix(),
	}

	if err := s.DB().Create(bind).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}

	response.Resp(c, bind)
}

// Update 更新用户绑定
func (s *TgUserBindService) Update(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		Id                int64       `json:"id" binding:"required"`
		UserId            interface{} `json:"userId"`
		TelegramId        interface{} `json:"telegramId"`
		TelegramUsername  string      `json:"telegramUsername"`
		TelegramFirstName string      `json:"telegramFirstName"`
		BindStatus        *int        `json:"bindStatus"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	bind, ok := base.GetOne[app.TgUserBind](s.DB(), "id", req.Id)
	if !ok {
		response.Resp(c, "绑定记录不存在")
		return
	}

	userId := gocast.ToInt64(req.UserId)
	telegramId := gocast.ToInt64(req.TelegramId)

	// 更新字段
	updates := make(map[string]interface{})
	if userId != 0 {
		updates["user_id"] = userId
	}
	if telegramId != 0 {
		updates["telegram_id"] = telegramId
	}
	if req.TelegramUsername != "" {
		updates["telegram_username"] = req.TelegramUsername
	}
	if req.TelegramFirstName != "" {
		updates["telegram_first_name"] = req.TelegramFirstName
	}
	if req.BindStatus != nil {
		updates["bind_status"] = *req.BindStatus
	}
	updates["updated_at"] = time.Now().Unix()

	if err := s.DB().Model(&bind).Updates(updates).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}

	response.Resp(c)
}

// Delete 删除用户绑定（软删除：修改状态为已解绑）
func (s *TgUserBindService) Delete(c *gin.Context) {
	s.SetDbAlias("app")
	id, ok := c.GetQuery("id")
	if !ok {
		response.Resp(c, "未填写ID")
		return
	}

	bindId := gocast.ToInt64(id)
	if bindId == 0 {
		response.Resp(c, "ID无效")
		return
	}

	bind, ok := base.GetOne[app.TgUserBind](s.DB(), "id", bindId)
	if !ok {
		response.Resp(c, "绑定记录不存在")
		return
	}

	// 软删除：状态改为0（已解绑）
	updates := map[string]interface{}{
		"bind_status": 0,
		"updated_at":  time.Now().Unix(),
	}
	if err := s.DB().Model(&bind).Updates(updates).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}

	response.Resp(c)
}
