package system

import (
	"github.com/gin-gonic/gin"
	config2 "go_server/base/config"
	"go_server/global"
)

var (
	allRouters = []global.ContextInterface{
		UserRouter{},
		RoleRouter{},
		MenuRouter{},
		ApisRouter{},
		DictionaryRouter{},
		FileRouter{},
		SignRouter{},
	}
)

type RouterGroup struct {
}

func (RouterGroup) Route() string {
	return "/sys"
}

func (h RouterGroup) Register(group *gin.RouterGroup) {
	// 正式环境不开启
	if config2.AppConf().Mod != config2.ModEnvProd {
		allRouters = append(allRouters, DbRouter{}, AutoRouter{})
	}
	for _, item := range allRouters {
		global.RegisterRouter(group, item)
	}
}
