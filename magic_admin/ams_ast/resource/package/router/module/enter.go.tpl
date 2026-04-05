package {{ .Module }}

import (
	"github.com/gin-gonic/gin"
	"{{ .Server }}/global"
)

// All routes to be registered

var (
	allRouters = []global.ContextInterface{
	}
)

type RouterGroup struct {
}

func (RouterGroup) Route() string {
	return "/{{ .Module }}"
}

func (h RouterGroup) Register(group *gin.RouterGroup) {
	for _, item := range allRouters {
		global.RegisterRouter(group, item)
	}
}
