package app

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

type CfDnsRouter struct{}

func (CfDnsRouter) Route() string {
	return "/cf_dns"
}

var cfDnsService = service.RealizationLayer.AppServiceGroup.CfDnsService

func (h CfDnsRouter) Register(group *gin.RouterGroup) {
	group.GET("zones", cfDnsService.Zones)
	group.GET("records", cfDnsService.Records)
	group.POST("record/create", cfDnsService.CreateRecord)
	group.POST("record/batch_create", cfDnsService.BatchCreateRecord)
	group.POST("record/update", cfDnsService.UpdateRecord)
	group.GET("record/delete", cfDnsService.DeleteRecord)
	group.POST("record/toggle_proxy", cfDnsService.ToggleProxy)
	group.POST("cross_zone/delete", cfDnsService.CrossZoneDeleteRecords)
	group.POST("cross_zone/toggle_proxy", cfDnsService.CrossZoneToggleProxy)
}
