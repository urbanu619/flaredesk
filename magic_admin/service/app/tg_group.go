package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/demdxx/gocast"
	"github.com/gin-gonic/gin"
	"go_server/base/core"
	"go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
	"go_server/utils"
)

type TgGroupService struct {
	base.BizCommonService
}

// Find 查询群组列表
func (s *TgGroupService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		ChatId   *int64 `form:"chat_id" json:"chat_id"`
		Title    string `form:"title" json:"title"`
		Status   *int8  `form:"status" json:"status"`
		Username string `form:"username" json:"username"`
	}
	req := new(request[app.TgGroup])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	db := s.DB()

	// 条件筛选
	if req.ChatId != nil && *req.ChatId != 0 {
		db = db.Where("chat_id = ?", *req.ChatId)
	}
	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	// 按创建时间倒序
	db = db.Order("created_at DESC")

	resp, err := base.NewQueryBaseHandler(app.NewTgGroup()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)
}

// Get 获取群组详情
func (s *TgGroupService) Get(c *gin.Context) {
	s.SetDbAlias("app")
	base.NewBaseHandler(app.NewTgGroup()).Get(c, s.DB())
}

// Create 创建群组
func (s *TgGroupService) Create(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		ChatId      interface{} `json:"chatId"` // 接受字符串或数字
		ChatType    string      `json:"chatType"`
		Title       string      `json:"title"`
		Username    string      `json:"username"`
		Description string      `json:"description"`
		MemberCount int         `json:"memberCount"`
		Status      int8        `json:"status"`
		BotJoinedAt int64       `json:"botJoinedAt"`
		Remark      string      `json:"remark"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	// 转换 ChatId
	chatId := gocast.ToInt64(req.ChatId)

	// 验证必填字段
	if chatId == 0 {
		response.Resp(c, "群组ID不能为空")
		return
	}
	if req.Title == "" {
		response.Resp(c, "群组标题不能为空")
		return
	}

	// 检查群组是否已存在
	var count int64
	s.DB().Model(&app.TgGroup{}).Where("chat_id = ?", chatId).Count(&count)
	if count > 0 {
		response.Resp(c, "该群组已存在")
		return
	}

	// 设置默认值
	if req.Status == 0 {
		req.Status = 1
	}
	if req.ChatType == "" {
		req.ChatType = "supergroup"
	}

	now := time.Now().Unix()
	group := &app.TgGroup{
		ChatId:      chatId,
		ChatType:    req.ChatType,
		Title:       req.Title,
		Username:    req.Username,
		Description: req.Description,
		MemberCount: req.MemberCount,
		Status:      req.Status,
		BotJoinedAt: req.BotJoinedAt,
		Remark:      req.Remark,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.DB().Create(group).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}

	response.Resp(c, group)
}

// Update 更新群组
func (s *TgGroupService) Update(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		Id          int64  `json:"id" binding:"required"`
		Title       string `json:"title"`
		Username    string `json:"username"`
		Description string `json:"description"`
		MemberCount int    `json:"memberCount"`
		Status      *int8  `json:"status"`
		Remark      string `json:"remark"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	group, ok := base.GetOne[app.TgGroup](s.DB(), "id", req.Id)
	if !ok {
		response.Resp(c, "群组不存在")
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.MemberCount > 0 {
		updates["member_count"] = req.MemberCount
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if err := s.DB().Model(&group).Updates(updates).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}

	response.Resp(c)
}

// Delete 删除群组（软删除：修改状态为已退出）
func (s *TgGroupService) Delete(c *gin.Context) {
	s.SetDbAlias("app")
	id, ok := c.GetQuery("id")
	if !ok {
		response.Resp(c, "未填写ID")
		return
	}

	groupId := gocast.ToInt64(id)
	if groupId == 0 {
		response.Resp(c, "ID无效")
		return
	}

	group, ok := base.GetOne[app.TgGroup](s.DB(), "id", groupId)
	if !ok {
		response.Resp(c, "群组不存在")
		return
	}

	// 软删除：状态改为3（已退出）
	if err := s.DB().Model(&group).Update("status", 3).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}

	response.Resp(c)
}

// SyncFromBot 从 Telegram Bot 同步群组信息
func (s *TgGroupService) SyncFromBot(c *gin.Context) {
	s.SetDbAlias("app")

	// 调用 magic_server API 同步群组信息
	groups, err := s.callMagicServerSyncGroups()
	if err != nil {
		core.Log.Errorf("同步群组失败: %v", err)
		response.Resp(c, fmt.Sprintf("同步失败: %v", err))
		return
	}

	// 统计新增和更新的数量
	newCount := 0
	updateCount := 0

	// 遍历群组，保存或更新到数据库
	for _, groupData := range groups {
		// 检查群组是否已存在
		var existingGroup app.TgGroup
		err := s.DB().Where("chat_id = ?", groupData.ChatId).First(&existingGroup).Error

		now := time.Now().Unix()

		if err != nil {
			// 群组不存在，创建新记录
			newGroup := &app.TgGroup{
				ChatId:      groupData.ChatId,
				ChatType:    groupData.ChatType,
				Title:       groupData.Title,
				Username:    groupData.Username,
				Description: groupData.Description,
				MemberCount: groupData.MemberCount,
				Status:      1, // 默认正常状态
				BotJoinedAt: groupData.BotJoinedAt,
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			if err := s.DB().Create(newGroup).Error; err != nil {
				core.Log.Errorf("创建群组记录失败: %v", err)
				continue
			}
			newCount++
		} else {
			// 群组已存在，更新信息
			updates := map[string]interface{}{
				"title":        groupData.Title,
				"username":     groupData.Username,
				"description":  groupData.Description,
				"member_count": groupData.MemberCount,
				"updated_at":   now,
			}

			if err := s.DB().Model(&existingGroup).Updates(updates).Error; err != nil {
				core.Log.Errorf("更新群组记录失败: %v", err)
				continue
			}
			updateCount++
		}
	}

	core.Log.Infof("群组同步完成: 新增 %d 个, 更新 %d 个", newCount, updateCount)

	response.Resp(c, map[string]interface{}{
		"message":     "同步成功",
		"totalCount":  len(groups),
		"newCount":    newCount,
		"updateCount": updateCount,
	})
}

// callMagicServerSyncGroups 调用 magic_server API 同步群组
func (s *TgGroupService) callMagicServerSyncGroups() ([]GroupData, error) {
	core.Log.Info("========== 开始同步群组 ==========")

	// magic_server API 地址
	apiURL := "http://localhost:2011/api/v1/adi/telegram/groups/sync"

	// 生成签名
	core.Log.Info("正在生成签名...")
	signMessage, err := utils.BuildSignMessage()
	if err != nil {
		core.Log.Errorf("生成签名失败: %v", err)
		return nil, fmt.Errorf("生成签名失败: %v", err)
	}
	core.Log.Infof("签名生成成功: %s", signMessage[:50]+"...")

	// 发送 HTTP 请求
	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建 HTTP 请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("sign", signMessage)

	core.Log.Infof("开始调用 magic_server API: %s", apiURL)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送 HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	core.Log.Infof("收到响应，状态码: %d", resp.StatusCode)

	// 解析响应
	var result struct {
		Code    int    `json:"code"`
		Msg     string `json:"msg"`     // magic_server 使用 msg 而不是 message
		Message string `json:"message"` // 兼容两种格式
		Data    struct {
			Groups []GroupData `json:"groups"`
			Count  int         `json:"count"`
		} `json:"data"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 打印原始响应用于调试
	core.Log.Infof("API 响应: %s", string(body))

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, 原始响应: %s", err, string(body))
	}

	if result.Code != 0 && result.Code != 200 {
		errMsg := result.Message
		if errMsg == "" {
			errMsg = result.Msg
		}
		return nil, fmt.Errorf("API 返回错误(code=%d): %s", result.Code, errMsg)
	}

	return result.Data.Groups, nil
}

// GroupData 群组数据结构（用于与 magic_server 交互）
type GroupData struct {
	ChatId      int64  `json:"chatId"`
	ChatType    string `json:"chatType"`
	Title       string `json:"title"`
	Username    string `json:"username"`
	Description string `json:"description"`
	MemberCount int    `json:"memberCount"`
	BotJoinedAt int64  `json:"botJoinedAt"`
}
