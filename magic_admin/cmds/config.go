package cmds

import (
	"github.com/spf13/cobra"
	"go_server/base/config"
	"go_server/base/core"
	"go_server/global"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "显示配置信息",
	Run: func(cmd *cobra.Command, args []string) {
		appConf := config.AppConf()
		core.Log.Infof("当前运行模式:%s http服务监听:%d 路由前缀:%s 业务代理默认地址:%s",
			appConf.Mod,
			appConf.Addr,
			appConf.RouterPrefix, appConf.ProxyUrl,
		)
		core.Log.Infof("alias-db 映射:%+v", global.AMS_BIZ_DBS)
		core.Log.Infof("dbName-alias映射:%+v", global.AMS_BIZ_ALIAS_DB_MAP)
		core.Log.Infof("alias-proxy映射:%+v", global.AMS_BIZ_ALIAS_PROXY_MAP)
	},
}




