package app

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

type TgRedPacketConfigRouter struct {
}

func (TgRedPacketConfigRouter) Route() string {
	return "/tg_red_packet_config"
}

var appTgRedPacketConfigService = service.RealizationLayer.AppServiceGroup.TgRedPacketConfigService

func (h TgRedPacketConfigRouter) Register(group *gin.RouterGroup) {
	group.GET("get", appTgRedPacketConfigService.Get)
	group.GET("find", appTgRedPacketConfigService.Find)
	group.GET("comment", appTgRedPacketConfigService.Comment)
	group.POST("create", appTgRedPacketConfigService.Create)
	group.POST("update", appTgRedPacketConfigService.Update)
	group.POST("delete", appTgRedPacketConfigService.Delete)
	group.POST("toggleStatus", appTgRedPacketConfigService.ToggleStatus)
}
