package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"go_server/base/core"
	"go_server/model/common/response"
	"go_server/model/system"
	"net/http"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// /proxy/*path
		// /api/app/login/check
		// /api/sys/user/find
		api, err := system.NewApis().FindOrCreateAuth(core.MainDb(), path, c.Request.Method)
		if err == nil {
			// 必须存在用户ID 角色ID
			roleId := c.GetInt64("roleId")
			var role *system.Role
			core.MainDb().Model(&system.Role{}).Where("id", roleId).First(&role)
			if role.Apis != "*" {
				authsList := strings.Split(role.Apis, ",") // 通过缓存获取角色配置的权限
				if !lo.Contains(authsList, fmt.Sprintf("%d", api.ID)) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorObjByCode(response.ResponseCodeInsufficientAuthority))
					return
				}
			}
		}
		// 继续执行
		c.Next()
	}
}
