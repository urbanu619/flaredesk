package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_server/model/common/response"
)

// Recovery 全局错误恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印错误堆栈信息
				response.Resp(c, fmt.Sprintf("recovery error:%v", err))
				c.Abort()
			}
		}()
		c.Next()
	}
}
