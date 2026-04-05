package app

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

type TgRedPacketRecordRouter struct{}

func (TgRedPacketRecordRouter) Route() string {
	return "/tg_red_packet"
}

var tgRedPacketRecordService = service.RealizationLayer.AppServiceGroup.TgRedPacketRecordService

func (h TgRedPacketRecordRouter) Register(group *gin.RouterGroup) {
	group.GET("find", tgRedPacketRecordService.Find)
	group.GET("get", tgRedPacketRecordService.Get)
	group.GET("findGrabRecords", tgRedPacketRecordService.FindGrabRecords)
}
