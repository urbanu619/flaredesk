package main

import (
	"go_server/base/config"
	core2 "go_server/base/core"
	"go_server/base/core/redisclient"
	"go_server/cmds"
	"go_server/global"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func init() {
	appConf := config.AppConf() // 配置初始化
	core2.Log.Infof("当前运行模式:%s http服务监听:%d 路由前缀:%s 业务代理默认地址:%s",
		appConf.Mod,
		appConf.Addr,
		appConf.RouterPrefix, appConf.ProxyUrl,
	)
	global.AMS_DB = core2.MainDb()                                                                   // 系统库初始化
	global.AMS_BIZ_DBS, global.AMS_BIZ_ALIAS_DB_MAP, global.AMS_BIZ_ALIAS_PROXY_MAP = core2.BizDbs() // 业务库初始化
	core2.Log.Infof("alias-db 映射:%+v", global.AMS_BIZ_DBS)
	core2.Log.Infof("dbName-alias映射:%+v", global.AMS_BIZ_ALIAS_DB_MAP)
	core2.Log.Infof("alias-proxy映射:%+v", global.AMS_BIZ_ALIAS_PROXY_MAP)
	global.AMS_REDIS = redisclient.DefaultClient()
}

func main() {
	// 启动cmd 命令服务
	cmds.Execute()
}
