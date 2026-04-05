package app

import (
	"github.com/gin-gonic/gin"
	"go_server/service"
)

// 路由器命名规则: Model+Router 首字母大些

type SysSignConfigRouter struct {
}

// 路由器命名规则: model 首字母小写

func (SysSignConfigRouter) Route() string {
	return "/sys_sign_config"
}

// 控制层与实现层 合二为一 让同一个业务 尽量在一个文件中实现与暴露
// 变量声明 bizUserService moduleModelService 首字母小写
// 服务命名规则:
// service 服务层 -- 固定值
// RealizationLayer 实现层 -- 固定值
// BizServiceGroup 服务组 首字母大些 moduleServiceGroup
// UserService 模型服务 modelService

var appSysSignConfigService = service.RealizationLayer.AppServiceGroup.SysSignConfigService


func (h SysSignConfigRouter) Register(group *gin.RouterGroup) {
	group.GET("get", appSysSignConfigService.Get)
	group.GET("find", appSysSignConfigService.Find)
	group.GET("comment", appSysSignConfigService.Comment)
}
