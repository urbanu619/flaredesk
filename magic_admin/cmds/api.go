package cmds

import (
	"github.com/spf13/cobra"
	"go_server/base/engine"
)

// 定时任务

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "启动Api服务",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
		engine.Run()
	},
}
