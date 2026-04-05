package global

import "github.com/gin-gonic/gin"

func RegisterRouter(group *gin.RouterGroup, iCore ContextInterface) {
	iCore.Register(group.Group(iCore.Route()))
}

type ContextInterface interface {
	Route() string
	Register(group *gin.RouterGroup)
}

type CommonHandler struct {
}

func (c CommonHandler) Route() string {
	panic("implement me")
}

func (CommonHandler) Register(group *gin.RouterGroup) {
	panic("implement me")
}
