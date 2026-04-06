package cmds

import (
	"github.com/spf13/cobra"
	"go_server/base/core"
	"go_server/base/engine"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "启动Api服务",
	Run: func(cmd *cobra.Command, args []string) {
		core.Migrates() // 与独立 migrate 子命令相同：建表 + 基础数据（SQLite/MySQL 均需）
		engine.Run()
	},
}
