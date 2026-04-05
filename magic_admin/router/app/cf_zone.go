package app

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

type CfZoneRouter struct{}

func (CfZoneRouter) Route() string {
	return "/cf_zone"
}

var cfZoneService = service.RealizationLayer.AppServiceGroup.CfZoneService

func (h CfZoneRouter) Register(group *gin.RouterGroup) {
	group.GET("sync", cfZoneService.Sync)
	group.GET("find", cfZoneService.Find)
	group.GET("all", cfZoneService.All)
	group.POST("update_remark", cfZoneService.UpdateRemark)
}
