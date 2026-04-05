package middleware

import (
	"github.com/gin-gonic/gin"
	core2 "go_server/base/core"
	"go_server/model/system"
	"time"
)

// 自定义日志拦截

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			// 在请求结束时记录日志
			elapsed := time.Since(start)
			core2.Log.Infof("GIN|%d|%s|IP:%s|%s|URI:%s|",
				c.Writer.Status(),
				elapsed.Round(time.Millisecond),
				c.ClientIP(),
				c.Request.Method,
				c.Request.URL.Path)

		}()
		c.Next()
		adminId := c.GetInt64("userId")
		core2.MainDb().Create(&system.AdministratorLog{
			AdminId:   adminId,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Ip:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		})
	}
}
