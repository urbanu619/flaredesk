package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go_server/base/core"
	model "go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
	"go_server/utils"
)

type TgRedPacketSendService struct {
	base.BizCommonService
}

// SendManual 手动发送红包（基于配置ID）
func (s *TgRedPacketSendService) SendManual(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		ConfigId int64 `json:"configId" binding:"required"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	// 查询配置
	var config model.TgRedPacketConfig
	if err := s.DB().Where("id = ? AND status = 1", req.ConfigId).First(&config).Error; err != nil {
		response.Resp(c, "配置不存在或已禁用")
		return
	}

	// 执行发送
	packetNo, err := s.ExecuteSendRedPacket(&config)
	if err != nil {
		core.Log.Errorf("发送红包失败: %v", err)
		response.Resp(c, fmt.Sprintf("发送失败: %v", err))
		return
	}

	// 更新配置的执行统计
	now := time.Now()
	s.DB().Model(&model.TgRedPacketConfig{}).Where("id = ?", config.Id).Updates(map[string]interface{}{
		"last_exec_time": now,
		"exec_count":     config.ExecCount + 1,
	})

	response.Resp(c, map[string]interface{}{
		"message":  "红包发送成功",
		"packetNo": packetNo,
	})
}

// SendDirect 直接发送红包（无需配置，直接传参数）
func (s *TgRedPacketSendService) SendDirect(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		GroupId       int64           `json:"groupId" binding:"required"`
		GroupName     string          `json:"groupName"`
		PacketType    int8            `json:"packetType" binding:"required"`
		TotalAmount   decimal.Decimal `json:"totalAmount" binding:"required"`
		TotalCount    int             `json:"totalCount" binding:"required"`
		Symbol        string          `json:"symbol"`
		BlessingWords string          `json:"blessingWords"`
		ExpireMinutes int             `json:"expireMinutes"`    // 过期时间(分钟)，默认10
		MaxGrabAmount decimal.Decimal `json:"maxGrabAmount"`    // 单人最大可领取金额(0=不限制)
		Lang          string          `json:"lang"`             // 消息语言: vi/id/en/zh
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	// 验证
	if req.PacketType != 1 && req.PacketType != 2 {
		response.Resp(c, "红包类型错误")
		return
	}
	if req.TotalCount <= 0 {
		response.Resp(c, "红包个数必须大于0")
		return
	}
	if req.TotalAmount.LessThanOrEqual(decimal.Zero) {
		response.Resp(c, "红包金额必须大于0")
		return
	}

	// 默认语言
	if req.Lang == "" {
		req.Lang = "vi"
	}
	// 默认币种（根据语言联动）
	if req.Symbol == "" {
		langSymbol := map[string]string{"vi": "VND", "id": "IDR", "en": "USD", "zh": "CNY"}
		if s, ok := langSymbol[req.Lang]; ok {
			req.Symbol = s
		} else {
			req.Symbol = "VND"
		}
	}
	if req.BlessingWords == "" {
		if req.PacketType == 1 {
			req.BlessingWords = "恭喜发财，大吉大利"
		} else {
			req.BlessingWords = "拼手气红包，快来抢！"
		}
	}

	// 默认过期时间
	if req.ExpireMinutes <= 0 {
		req.ExpireMinutes = 10
	}

	// 构建临时配置
	config := &model.TgRedPacketConfig{
		GroupId:       fmt.Sprintf("%d", req.GroupId),
		GroupName:     req.GroupName,
		PacketType:    req.PacketType,
		TotalAmount:   req.TotalAmount,
		TotalCount:    req.TotalCount,
		Symbol:        req.Symbol,
		BlessingWords: req.BlessingWords,
		ExpireMinutes: req.ExpireMinutes,
		MaxGrabAmount: req.MaxGrabAmount,
		Lang:          req.Lang,
	}

	// 执行发送
	packetNo, err := s.ExecuteSendRedPacket(config)
	if err != nil {
		core.Log.Errorf("发送红包失败: %v", err)
		response.Resp(c, fmt.Sprintf("发送失败: %v", err))
		return
	}

	response.Resp(c, map[string]interface{}{
		"message":  "红包发送成功",
		"packetNo": packetNo,
	})
}

// ExecuteSendRedPacket 执行发送红包的核心逻辑
func (s *TgRedPacketSendService) ExecuteSendRedPacket(config *model.TgRedPacketConfig) (string, error) {
	// 调用 magic_server API 创建红包并发送到 Telegram
	packetNo, err := s.callMagicServerSendRedPacket(config)
	if err != nil {
		core.Log.Errorf("创建并发送红包失败: %v", err)
		return "", err
	}

	core.Log.Infof("成功发送红包: 群组=%s, 类型=%d, 金额=%v, 个数=%d, 编号=%s",
		config.GroupId, config.PacketType, config.TotalAmount, config.TotalCount, packetNo)

	return packetNo, nil
}

// callMagicServerSendRedPacket 调用 magic_server API 发送红包
func (s *TgRedPacketSendService) callMagicServerSendRedPacket(config *model.TgRedPacketConfig) (string, error) {
	// magic_server API 地址
	apiURL := "http://localhost:2011/api/v1/adi/telegram/redpacket/send"

	// 构造请求参数
	payload := map[string]interface{}{
		"groupId":       config.GroupId,
		"totalAmount":   config.TotalAmount,
		"totalCount":    config.TotalCount,
		"packetType":    config.PacketType,
		"symbol":        config.Symbol,
		"blessingWords": config.BlessingWords,
		"expireMinutes": config.ExpireMinutes,
		"maxGrabAmount": config.MaxGrabAmount,
		"lang":          config.Lang,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("序列化请求参数失败: %v", err)
	}

	// 发送 HTTP 请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建 HTTP 请求失败: %v", err)
	}

	signMessage, _ := utils.BuildSignMessage()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("sign", signMessage)
	// TODO: 添加认证 Token
	// req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送 HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			PacketNo string `json:"packetNo"`
			PacketId int64  `json:"packetId"`
		} `json:"data"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 记录原始响应，用于调试
	core.Log.Infof("magic_server 原始响应: %s", string(body))

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v (响应内容: %s)", err, string(body))
	}

	if result.Code != 0 && result.Code != 200 {
		return "", fmt.Errorf("API 返回错误: %s", result.Message)
	}

	if result.Data.PacketNo == "" {
		return "", fmt.Errorf("红包编号为空")
	}

	return result.Data.PacketNo, nil
}

// GetGroups 获取机器人所在的群组列表
func (s *TgRedPacketSendService) GetGroups(c *gin.Context) {
	s.SetDbAlias("app")
	// 从数据库中获取已记录的群组
	var groups []model.TgGroup
	if err := s.DB().Where("status = 1").Find(&groups).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}

	response.Resp(c, map[string]interface{}{
		"message": "获取群组列表成功",
		"groups":  groups,
	})
}
