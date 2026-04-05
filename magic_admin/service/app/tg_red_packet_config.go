package app

import (
	"go_server/base/core"
	model "go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TgRedPacketConfigService struct {
	base.BizCommonService
}

// Get 获取单条记录
func (s *TgRedPacketConfigService) Get(c *gin.Context) {
	s.SetDbAlias("app")
	base.NewBaseHandler(model.NewTgRedPacketConfig()).Get(c, s.DB())
}

// Find 查询列表
func (s *TgRedPacketConfigService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		Id          *int64  `form:"id"`
		IsExport    *bool   `form:"isExport"`
		Fields      *string `form:"fields"` // 指定返回字段 , 分割
		ConfigName  *string `form:"configName"`
		ConfigType  *int8   `form:"configType"`
		GroupId     *string `form:"groupId"`
		PacketType  *int8   `form:"packetType"`
		Status      *int8   `form:"status"`
		BeginTime   *int64  `form:"beginTime"`
		EndTime     *int64  `form:"endTime"`
	}
	req := new(request[model.TgRedPacketConfig])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	db := s.DB()
	if req.Id != nil && *req.Id != 0 {
		db = db.Where("id", req.Id)
	}
	if req.ConfigName != nil && *req.ConfigName != "" {
		db = db.Where("config_name LIKE ?", "%"+*req.ConfigName+"%")
	}
	if req.ConfigType != nil {
		db = db.Where("config_type", req.ConfigType)
	}
	if req.GroupId != nil && *req.GroupId != "" {
		db = db.Where("group_id", req.GroupId)
	}
	if req.PacketType != nil {
		db = db.Where("packet_type", req.PacketType)
	}
	if req.Status != nil {
		db = db.Where("status", req.Status)
	} else {
		// 默认不显示已删除的
		db = db.Where("status != ?", 3)
	}
	if req.BeginTime != nil && *req.BeginTime > 0 && req.EndTime != nil && *req.EndTime > 0 {
		db = db.Where("created_at between ? and ?", *req.BeginTime, *req.EndTime)
	}

	// 导出处理
	fields := make([]string, 0)
	if req.Fields != nil {
		fields = strings.Split(*req.Fields, ",")
	}
	var url string
	colInfo := s.GetColumnCommentFromStruct(model.TgRedPacketConfig{})
	var err error
	if req.IsExport != nil && *req.IsExport {
		if len(fields) == 0 {
			for _, col := range colInfo {
				fields = append(fields, col.Field)
			}
		}
		core.Log.Infof("导出的字段:%s", fields)
		url, err = base.ExportCsv[model.TgRedPacketConfig](db, fields, colInfo)
		if err != nil {
			response.Resp(c, err.Error())
			return
		}
	}
	resp, err := base.NewQueryBaseHandler(model.NewTgRedPacketConfig()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}

	// 批量填充群组名称：从 magic_tg_group 表查 title
	chatIds := make([]int64, 0, len(resp.List))
	for _, item := range resp.List {
		if item.GroupId != "" {
			if chatId, err := strconv.ParseInt(item.GroupId, 10, 64); err == nil {
				chatIds = append(chatIds, chatId)
			}
		}
	}
	groupNameMap := make(map[string]string)
	if len(chatIds) > 0 {
		var groups []model.TgGroup
		s.DB().Where("chat_id IN ?", chatIds).Find(&groups)
		for _, g := range groups {
			groupNameMap[strconv.FormatInt(g.ChatId, 10)] = g.Title
		}
	}
	// 填充 group_name
	for i := range resp.List {
		if resp.List[i].GroupName == "" {
			if title, ok := groupNameMap[resp.List[i].GroupId]; ok {
				resp.List[i].GroupName = title
			}
		}
	}

	response.Resp(c, map[string]interface{}{
		"url":    url,
		"cols":   colInfo,
		"list":   resp.List,
		"paging": resp.Paging,
	})
}

// Comment 获取表字段注释
func (s *TgRedPacketConfigService) Comment(c *gin.Context) {
	s.SetDbAlias("app")
	dbs, err := s.GetColumnComment("app", model.NewTgRedPacketConfig().TableName())
	if err != nil {
		response.Resp(c, "获取失败")
		return
	} else {
		response.Resp(c, gin.H{"dbs": dbs})
		return
	}
}

// Create 创建红包配置
func (s *TgRedPacketConfigService) Create(c *gin.Context) {
	s.SetDbAlias("app")
	var req model.TgRedPacketConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	// 验证必填字段
	if req.ConfigName == "" {
		response.Resp(c, "配置名称不能为空")
		return
	}
	if req.GroupId == "" {
		response.Resp(c, "群组ID不能为空")
		return
	}
	if req.ConfigType != 1 && req.ConfigType != 2 {
		response.Resp(c, "配置类型错误")
		return
	}
	if req.PacketType != 1 && req.PacketType != 2 {
		response.Resp(c, "红包类型错误")
		return
	}
	if req.TotalCount <= 0 {
		response.Resp(c, "红包个数必须大于0")
		return
	}

	// 定时红包需要 cron 表达式
	if req.ConfigType == 1 && req.CronExpr == "" {
		response.Resp(c, "定时红包必须配置Cron表达式")
		return
	}

	// 设置默认值
	if req.Status == 0 {
		req.Status = 1 // 默认启用
	}
	if req.Symbol == "" {
		req.Symbol = "VND"
	}
	if req.TimeZone == "" {
		req.TimeZone = "Asia/Ho_Chi_Minh"
	}

	if err := s.DB().Create(&req).Error; err != nil {
		core.Log.Errorf("创建红包配置失败: %v", err)
		response.Resp(c, "创建失败")
		return
	}
	response.Resp(c)
}

// Update 更新红包配置
func (s *TgRedPacketConfigService) Update(c *gin.Context) {
	s.SetDbAlias("app")
	var req model.TgRedPacketConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	if req.Id == 0 {
		response.Resp(c, "ID不能为空")
		return
	}

	// 验证配置是否存在
	var exist model.TgRedPacketConfig
	if err := s.DB().Where("id = ?", req.Id).First(&exist).Error; err != nil {
		response.Resp(c, "配置不存在")
		return
	}

	// 使用 map 更新，这样 nil 值也能正确清空日期字段
	updates := map[string]interface{}{
		"config_name":     req.ConfigName,
		"config_type":     req.ConfigType,
		"group_id":        req.GroupId,
		"group_name":      req.GroupName,
		"packet_type":     req.PacketType,
		"total_amount":    req.TotalAmount,
		"total_count":     req.TotalCount,
		"symbol":          req.Symbol,
		"expire_minutes":  req.ExpireMinutes,
		"max_grab_amount": req.MaxGrabAmount,
		"lang":            req.Lang,
		"blessing_words":  req.BlessingWords,
		"cron_expr":       req.CronExpr,
		"time_zone":       req.TimeZone,
		"start_date":      req.StartDate,
		"end_date":        req.EndDate,
		"remark":          req.Remark,
	}

	if err := s.DB().Model(&model.TgRedPacketConfig{}).Where("id = ?", req.Id).Updates(updates).Error; err != nil {
		core.Log.Errorf("更新红包配置失败: %v", err)
		response.Resp(c, "更新失败")
		return
	}
	response.Resp(c)
}

// Delete 删除红包配置（软删除）
func (s *TgRedPacketConfigService) Delete(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		Id int64 `json:"id" binding:"required"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	// 软删除：将状态改为3
	if err := s.DB().Model(&model.TgRedPacketConfig{}).Where("id = ?", req.Id).Update("status", 3).Error; err != nil {
		core.Log.Errorf("删除红包配置失败: %v", err)
		response.Resp(c, "删除失败")
		return
	}
	response.Resp(c)
}

// ToggleStatus 切换状态（启用/禁用）
func (s *TgRedPacketConfigService) ToggleStatus(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		Id     int64 `json:"id" binding:"required"`
		Status int8  `json:"status" binding:"required"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	if req.Status != 1 && req.Status != 2 {
		response.Resp(c, "状态值错误")
		return
	}

	if err := s.DB().Model(&model.TgRedPacketConfig{}).Where("id = ?", req.Id).Update("status", req.Status).Error; err != nil {
		core.Log.Errorf("切换状态失败: %v", err)
		response.Resp(c, "操作失败")
		return
	}
	response.Resp(c)
}
