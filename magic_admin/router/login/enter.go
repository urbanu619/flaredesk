package login

import (
	"github.com/gin-gonic/gin"
	"go_server/global"
)

var (
	allRouters = []global.ContextInterface{
		LoginRouter{},
	}
)

type RouterGroup struct {
}

func (RouterGroup) Route() string {
	return "/app"
}

func (h RouterGroup) Register(group *gin.RouterGroup) {
	for _, item := range allRouters {
		global.RegisterRouter(group, item)
	}
}
