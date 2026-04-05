package login

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_server/base/config"
	"go_server/base/core"
	"go_server/model/common/response"
	"go_server/model/system"
	"go_server/service/base"
	"go_server/utils"
)

func (s *AppService) Check(c *gin.Context) {
	type request struct {
		Username string `json:"username"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	allowed, err := core.NewRateLimiter(s.Redis()).CanExecuteMethod(context.Background(), "Check", c.RemoteIP())
	if err != nil {
		core.Log.Error(err.Error())
		response.Resp(c, "请求频繁")
		return
	}
	if !allowed {
		response.Resp(c, "请求频繁")
		return
	}
	resp := map[string]interface{}{
		"isSetGoogleAuth": false,
		"isNeedCaptcha":   true,
		"t":               0,
	}

	t, ok := core.CapEngine().CapCheck(s.Redis(), c.RemoteIP())
	resp["isNeedCaptcha"] = ok
	resp["t"] = t
	if req.Username == "" {
		response.Resp(c, resp)
		return
	}
	user, ok := base.GetOne[system.Administrator](s.DB(), "username", req.Username)
	if !ok {
		response.Resp(c, resp)
		return
	}
	if user.Lock {
		response.Resp(c, resp)
		return
	}
	if user.GoogleKey != "" {
		resp["isSetGoogleAuth"] = true
	}
	response.Resp(c, resp)
}

func (s *AppService) WalletUrl() string {
	if config.AppConf().Mod == config.ModEnvProd {
		return "http://10.0.0.52:1432"
	}
	return "http://127.0.0.1:1432"
}

func (s *AppService) WalletInfo(c *gin.Context) {
	url := fmt.Sprintf("%s/api/wallet/info", s.WalletUrl())

	resp, err := utils.Get[interface{}](url)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)

}
