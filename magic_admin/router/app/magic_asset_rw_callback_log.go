package app

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

// 路由器命名规则: Model+Router 首字母大些

type MagicAssetRwCallbackLogRouter struct {
}

// 路由器命名规则: model 首字母小写

func (MagicAssetRwCallbackLogRouter) Route() string {
	return "/magic_asset_rw_callback_log"
}

// 控制层与实现层 合二为一 让同一个业务 尽量在一个文件中实现与暴露
// 变量声明 bizUserService moduleModelService 首字母小写
// 服务命名规则:
// service 服务层 -- 固定值
// RealizationLayer 实现层 -- 固定值
// BizServiceGroup 服务组 首字母大些 moduleServiceGroup
// UserService 模型服务 modelService

var appMagicAssetRwCallbackLogService = service.RealizationLayer.AppServiceGroup.MagicAssetRwCallbackLogService


func (h MagicAssetRwCallbackLogRouter) Register(group *gin.RouterGroup) {
	group.GET("get", appMagicAssetRwCallbackLogService.Get)
	group.GET("find", appMagicAssetRwCallbackLogService.Find)
	group.GET("comment", appMagicAssetRwCallbackLogService.Comment)
}
