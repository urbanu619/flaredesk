package app

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

type CfAccountRouter struct{}

func (CfAccountRouter) Route() string {
	return "/cf_account"
}

var cfAccountService = service.RealizationLayer.AppServiceGroup.CfAccountService

func (h CfAccountRouter) Register(group *gin.RouterGroup) {
	group.GET("find", cfAccountService.Find)
	group.GET("get", cfAccountService.Get)
	group.POST("create", cfAccountService.Create)
	group.POST("update", cfAccountService.Update)
	group.GET("delete", cfAccountService.Delete)
}
