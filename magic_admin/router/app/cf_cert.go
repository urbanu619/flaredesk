package app

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

type CfCertRouter struct{}

func (CfCertRouter) Route() string {
	return "/cf_cert"
}

var cfCertService = service.RealizationLayer.AppServiceGroup.CfCertService

func (h CfCertRouter) Register(group *gin.RouterGroup) {
	group.GET("list", cfCertService.ListCerts)
	group.POST("batch_generate", cfCertService.BatchGenerateCerts)
}
