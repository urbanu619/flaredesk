package {{ .Module }}

import (
	"github.com/gin-gonic/gin"
	"{{ .Server }}/service"
)

// 路由器命名规则: Model+Router 首字母大些

type {{ .ModelUpper }}Router struct {
}

// 路由器命名规则: model 首字母小写

func ({{ .ModelUpper }}Router) Route() string {
	return "/{{ .Model }}"
}

// 控制层与实现层 合二为一 让同一个业务 尽量在一个文件中实现与暴露
// 变量声明 bizUserService moduleModelService 首字母小写
// 服务命名规则:
// service 服务层 -- 固定值
// RealizationLayer 实现层 -- 固定值
// BizServiceGroup 服务组 首字母大些 moduleServiceGroup
// UserService 模型服务 modelService

var {{ .RouterServiceName }}Service = service.RealizationLayer.{{ .ModuleUpper }}ServiceGroup.{{ .ModelUpper }}Service


func (h {{ .ModelUpper }}Router) Register(group *gin.RouterGroup) {
	group.GET("get", {{ .RouterServiceName }}Service.Get)
	group.GET("find", {{ .RouterServiceName }}Service.Find)
	group.GET("comment", {{ .RouterServiceName }}Service.Comment)
}
