package cmds

import (
	"github.com/spf13/cobra"
	"go_server/base/core"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "启动Migrate服务",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
		core.Migrates() // 数据迁移
	},
}
