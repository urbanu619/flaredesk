package app

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

type TgGroupRouter struct{}

func (TgGroupRouter) Route() string {
	return "/tg_group"
}

var tgGroupService = service.RealizationLayer.AppServiceGroup.TgGroupService

func (h TgGroupRouter) Register(group *gin.RouterGroup) {
	group.GET("find", tgGroupService.Find)
	group.GET("get", tgGroupService.Get)
	group.POST("create", tgGroupService.Create)
	group.POST("update", tgGroupService.Update)
	group.GET("delete", tgGroupService.Delete)
	group.POST("syncFromBot", tgGroupService.SyncFromBot)
}
