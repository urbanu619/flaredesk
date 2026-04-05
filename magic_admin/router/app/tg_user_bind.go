package app

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

type TgUserBindRouter struct{}

func (TgUserBindRouter) Route() string {
	return "/tg_user_bind"
}

var tgUserBindService = service.RealizationLayer.AppServiceGroup.TgUserBindService

func (h TgUserBindRouter) Register(group *gin.RouterGroup) {
	group.GET("find", tgUserBindService.Find)
	group.GET("get", tgUserBindService.Get)
	group.POST("create", tgUserBindService.Create)
	group.POST("update", tgUserBindService.Update)
	group.GET("delete", tgUserBindService.Delete)
}
