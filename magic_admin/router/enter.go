package router

import (
	"github.com/gin-gonic/gin"
	"go_server/global"
	"go_server/router/login"
	"go_server/router/system"
	"go_server/router/app"
)

var (
	pubRouters	= []global.ContextInterface{
		login.RouterGroup{},
	}

	priRouters	= []global.ContextInterface{
		system.RouterGroup{}, app.RouterGroup{},
	}
)

type PubRouterGroupApp struct {
}

func (PubRouterGroupApp) Route() string {
	return "/api"
}

func (h PubRouterGroupApp) Register(group *gin.RouterGroup) {
	for _, item := range pubRouters {
		global.RegisterRouter(group, item)
	}
}

type PriRouterGroupApp struct {
}

func (PriRouterGroupApp) Route() string {
	return "/api"
}

func (h PriRouterGroupApp) Register(group *gin.RouterGroup) {
	for _, item := range priRouters {
		global.RegisterRouter(group, item)
	}
}
