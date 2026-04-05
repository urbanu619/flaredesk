package app

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

type TgRedPacketSendRouter struct {
}

func (TgRedPacketSendRouter) Route() string {
	return "/tg_red_packet_send"
}

var appTgRedPacketSendService = service.RealizationLayer.AppServiceGroup.TgRedPacketSendService

func (h TgRedPacketSendRouter) Register(group *gin.RouterGroup) {
	group.POST("sendManual", appTgRedPacketSendService.SendManual)       // 基于配置ID发送
	group.POST("sendDirect", appTgRedPacketSendService.SendDirect)       // 直接发送（无需配置）
	group.GET("getGroups", appTgRedPacketSendService.GetGroups)          // 获取群组列表
}
