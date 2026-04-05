package app

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

type CfDnsTemplateRouter struct{}

func (CfDnsTemplateRouter) Route() string {
	return "/cf_dns_template"
}

var cfDnsTemplateService = service.RealizationLayer.AppServiceGroup.CfDnsTemplateService

func (h CfDnsTemplateRouter) Register(group *gin.RouterGroup) {
	group.GET("find", cfDnsTemplateService.Find)
	group.GET("get", cfDnsTemplateService.Get)
	group.POST("create", cfDnsTemplateService.Create)
	group.POST("update", cfDnsTemplateService.Update)
	group.GET("delete", cfDnsTemplateService.Delete)
}
