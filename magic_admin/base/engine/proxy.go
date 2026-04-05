package engine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_server/base/config"
	"go_server/base/core"
	"go_server/global"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// 重定向配置
var pathRewrite = map[string]string{
	"/health/liveness": "/health/liveness",
}

// 系统重定向

func proxyPathRewrite(rawPath string) string {
	// 根据路径选择代理地址
	originalPath := rawPath // 路径重写
	// 重定向
	for pattern, replacement := range pathRewrite {
		if pattern == originalPath {
			rawPath = replacement
			break
		}
	}
	return strings.ReplaceAll(rawPath, "/admin/api/proxy", "/api")
	//return "/api" + strings.TrimPrefix(rawPath, "/admin/api/proxy")
}

// todo 通过请求路径解析出代理URL
// biz 配置代理请求地址
// /admin/api/proxy 默认请求地址前缀
// biz 配置 proxy-alias 根据代理别名获取配置的URL
// 如业务服务路径未保持一致 则使用默认代理地址

func getTargetByPath(rewPath string) string {
	aliasSplit := strings.Split(strings.TrimPrefix(rewPath, "/admin/api/proxy"), "/")
	for _, alias := range aliasSplit {
		if alias != "" {
			return proxyAddrByAlias(alias)
		}
	}
	// 未找到则使用默认代理地址
	return config.AppConf().ProxyUrl
}

func proxyAddrByAlias(dbAlias string) string {
	// dbname转换为alias
	for k, dbName := range global.AMS_BIZ_ALIAS_DB_MAP {
		if dbName == dbAlias {
			dbAlias = k
		}
	}
	v, ok := global.AMS_BIZ_ALIAS_PROXY_MAP[dbAlias]
	if !ok {
		// 兼容直接通过数据库找到对应的连接
		return config.AppConf().ProxyUrl
	}
	return v
}

// 使用中间件代理转发 target string, pathRewrite map[string]string

func createReverseProxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		core.Log.Infof("管理后台请求路径:%s", reqPath)
		target := getTargetByPath(reqPath)
		core.Log.Infof("代理URL:%s", target)
		remote, err := url.Parse(target)
		if err != nil {
			panic(err)
		}
		// 解析path 选择目标服务
		authID := c.GetInt64("userId")
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			sign, _ := core.BuildSignMessage()
			req.Header.Set(core.SignKey, sign)
			core.Log.Infof("sign:%s", sign)
			req.Header.Set("authId", fmt.Sprintf("%d", authID))
			// 路径重定向
			req.URL.Path = proxyPathRewrite(req.URL.Path)
			core.Log.Infof("代理请求路径:%s%s", remote, req.URL.Path)
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
