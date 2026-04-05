package login

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

type LoginRouter struct {
}

func (LoginRouter) Route() string {
	return "/login"
}

// 控制层与实现层 合二为一 让同一个业务 尽量在一个文件中实现与暴露
var appService = service.RealizationLayer.LoginServiceGroup

func (h LoginRouter) Register(group *gin.RouterGroup) {
	group.GET("/generateCaptcha", appService.GenerateCaptcha)
	group.POST("check", appService.Check)
	group.POST("in", appService.In)
	group.GET("wallet/info", appService.WalletInfo) // 新增公共接口
}
