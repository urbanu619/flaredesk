package login

import (
	"context"
	"fmt"
	"github.com/demdxx/gocast"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"go_server/base/core"
	"go_server/base/engine/middleware"
	"go_server/model/common/response"
	"go_server/model/system"
	"go_server/service/base"
	"time"
)

func (s *AppService) In(c *gin.Context) {
	type request struct {
		Username   interface{} `json:"username" validate:"required"`
		Password   interface{} `json:"password"  validate:"required"`
		GoogleCode interface{} `json:"googleCode"` // 谷歌验证码
		Cid        *string     `json:"cid"`        // 验证Id
		Code       *string     `json:"code"`       // 验证码
	}
	var err error
	defer func() {
		if err != nil {
			_ = core.CapEngine().CapAdd(s.Redis(), c.RemoteIP())
		} else {
			_ = core.CapEngine().CapClear(s.Redis(), c.RemoteIP())
			_ = core.NewRateLimiter(s.Redis()).ClearLimit(context.Background(), "LOGIN", c.RemoteIP())
		}
	}()
	// 接口请求频率限制
	allowed, err := core.NewRateLimiter(s.Redis()).CanExecuteMethod(context.Background(), "LOGIN", c.RemoteIP())
	if err != nil {
		response.Resp(c, "请求频繁")
		return
	}
	if !allowed {
		response.Resp(c, "请求频繁")
		return
	}
	req := new(request)
	if err = c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	// Check
	username := gocast.ToString(req.Username)
	password := gocast.ToString(req.Password)
	googleCode := gocast.ToString(req.GoogleCode)
	if username == "" || password == "" {
		response.Resp(c, response.ResponseCodeParamError)
		return
	}
	// 请求频率限制
	//serviceName := "AdminLogin" + c.RemoteIP()
	//checkRep := s.RepeatFilter(serviceName, time.Duration(2)*time.Second)
	//if !checkRep {
	//	response.Resp(c, response.ResponseCodeFrequentOperation)
	//	return
	//}
	t, ok := core.CapEngine().CapCheck(s.Redis(), c.RemoteIP())
	if ok {
		if req.Cid == nil || req.Code == nil {
			err = fmt.Errorf("连续%d次输入错误信息 请输入验证码", t)
			response.Resp(c, err.Error())
			return
		}
		if !core.CapEngine().Verify(*req.Cid, *req.Code) {
			err = fmt.Errorf("连续%d次输入错误信息 请注意区分大小写", t)
			response.Resp(c, err.Error())
			return
		}
	}
	user, ok := base.GetOne[system.Administrator](s.DB(), "username", username)
	if !ok {
		err = fmt.Errorf("连续%d次输入错误信息", t)
		response.Resp(c, err.Error())
		return
	}
	if user.Lock {
		err = fmt.Errorf("连续%d次输入错误信息", t)
		response.Resp(c, err.Error())
		return
	}
	if !user.CheckPassWord(password) {
		err = fmt.Errorf("连续%d次输入错误信息", t)
		response.Resp(c, err.Error())
		return
	}
	if user.GoogleKey != "" {
		if googleCode == "" {
			err = fmt.Errorf("连续%d次输入错误信息:GoogleCodeError", t)
			response.Resp(c, err.Error())
			return
		}
		if ok, _ := core.NewGoogleAuth().VerifyCode(user.GoogleKey, googleCode); !ok {
			err = fmt.Errorf("连续%d次输入错误信息:GoogleCodeError", t)
			response.Resp(c, err.Error())
			return
		}
	}
	tokenString, err := middleware.GenerateJWT(middleware.Member{
		ID:     user.ID,
		RoleId: user.RoleId,
	})
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	// 保存token到mysql
	user.LastLoginIp = c.RemoteIP()
	user.LastLoginTime = lo.ToPtr(time.Now())
	user.Token = tokenString
	if err = s.DB().Save(&user).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["nickname"] = user.Nickname
	data["avatar"] = user.Avatar
	data["token"] = tokenString
	response.Resp(c, data)
	return
}
