package login

import (
	"github.com/gin-gonic/gin"
	"go_server/base/core"
	"go_server/model/common/response"
	"go_server/service/base"
)

type AppService struct {
	base.SysCommonService
}

func (s *AppService) GenerateCaptcha(c *gin.Context) {
	bs64, cid := core.CapEngine().Generate()
	data := map[string]any{
		"cid":     cid,
		"captcha": bs64,
	}
	response.Resp(c, data)
}
