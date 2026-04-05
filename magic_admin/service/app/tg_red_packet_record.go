package app

import (
	"github.com/gin-gonic/gin"
	"go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
)

type TgRedPacketRecordService struct {
	base.BizCommonService
}

// Find 查询红包发放记录列表（主表）
func (s *TgRedPacketRecordService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		PacketNo string `form:"packet_no" json:"packet_no"`
		GroupId  string `form:"group_id" json:"group_id"`
		Status   *int   `form:"status" json:"status"`
	}
	req := new(request[app.TgRedPacket])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	db := s.DB()

	// 条件筛选
	if req.PacketNo != "" {
		db = db.Where("packet_no = ?", req.PacketNo)
	}
	if req.GroupId != "" {
		db = db.Where("group_id = ?", req.GroupId)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	// 按创建时间倒序
	db = db.Order("created_at DESC")

	resp, err := base.NewQueryBaseHandler(app.NewTgRedPacket()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)
}

// Get 获取红包详情
func (s *TgRedPacketRecordService) Get(c *gin.Context) {
	s.SetDbAlias("app")
	base.NewBaseHandler(app.NewTgRedPacket()).Get(c, s.DB())
}

// FindGrabRecords 查询抢红包记录列表
func (s *TgRedPacketRecordService) FindGrabRecords(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		PacketNo         string `form:"packet_no" json:"packet_no"`
		PacketId         int64  `form:"packet_id" json:"packet_id"`
		TelegramUsername string `form:"telegram_username" json:"telegram_username"`
		UserId           int64  `form:"user_id" json:"user_id"`
	}
	req := new(request[app.TgRedPacketRecord])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	db := s.DB()

	if req.PacketNo != "" {
		db = db.Where("packet_no = ?", req.PacketNo)
	}
	if req.PacketId != 0 {
		db = db.Where("packet_id = ?", req.PacketId)
	}
	if req.TelegramUsername != "" {
		db = db.Where("telegram_username = ?", req.TelegramUsername)
	}
	if req.UserId != 0 {
		db = db.Where("user_id = ?", req.UserId)
	}

	db = db.Order("created_at DESC")

	resp, err := base.NewQueryBaseHandler(app.NewTgRedPacketRecord()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)
}
