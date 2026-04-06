package engine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_server/base/config"
	"go_server/base/core"
	"go_server/base/engine/middleware"
	"go_server/global"
	"go_server/router"
	"go_server/webdist"
	"os"
)

func Run() {
	address := fmt.Sprintf(":%d", config.AppConf().Addr)
	engine := engineInit()
	logRoutes()

	core.Log.Infof("Gin开启http监听: %s", address)
	core.Log.Infof("Gin run error: %+v", engine.Run(address))
}

// gin

func engineInit() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	// 挂在全局中间件
	engine.Use(middleware.Recovery()) // 全局错误恢复中间件
	engine.Use(middleware.Cors())     // 跨域 -- 放行所有
	k8sGroup(engine)                  // k8s 健康检查

	// 挂在静态文件 用于文件服务 代理路径"/static" 存储路径"uploads"
	// 代理转发
	routerPrefix := "admin"
	if config.AppConf().RouterPrefix != "" {
		routerPrefix = config.AppConf().RouterPrefix
	}
	engine.Static(config.EnvConf().File.ProxyPath, config.EnvConf().File.StorePath)

	if _, err := os.Stat(config.EnvConf().File.StorePath); os.IsNotExist(err) {
		err := os.Mkdir(config.EnvConf().File.StorePath, 0755) // 设置文件权限为rwxr-xr-x（默认值）
		if err != nil {
			panic(err)
		}
	}

	proxyPrefix := fmt.Sprintf("%s/api/proxy", routerPrefix)
	proxyGroup := engine.Group(proxyPrefix)
	proxyGroup.Use(middleware.GinLogger())     // 日志处理 -- 自定义日志
	proxyGroup.Use(middleware.JwtMiddleware()) // jwt鉴权
	// todo: 接口权限限制中间件
	proxyGroup.Any("/*path", createReverseProxy())
	// 使用代理中间件
	publicGroup := engine.Group(routerPrefix)
	publicGroup.Use(middleware.GinLogger()) // 日志处理 -- 自定义日志
	privateGroup := engine.Group(routerPrefix)
	privateGroup.Use(middleware.GinLogger())                            // 日志处理 -- 自定义日志
	privateGroup.Use(middleware.JwtMiddleware()).Use(middleware.Auth()) // 挂载中间件 -- token auth
	if config.AppConf().Mod != config.ModEnvDev {
		// 仅本地调试时使用
	}
	global.RegisterRouter(publicGroup, router.PubRouterGroupApp{})
	global.RegisterRouter(privateGroup, router.PriRouterGroupApp{})
	webdist.Mount(engine, config.AppConf().ServeWeb)
	global.GVA_ROUTERS = engine.Routes()
	return engine
}

// 通用k8s健康检查

func k8sGroup(engine *gin.Engine) {
	health := engine.Group("/health")
	{
		health.GET("/liveness", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"status": "alive"})
		})
		health.GET("/readiness", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"status": "ready"})
		})
	}
}

func logRoutes() {
	for _, route := range global.GVA_ROUTERS {
		core.Log.Infof("ROUTE | %s | URI:%s ",
			route.Method, route.Path)
	}
}
