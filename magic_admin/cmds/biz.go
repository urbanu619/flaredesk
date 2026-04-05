package cmds

import (
	"github.com/spf13/cobra"
	"go_server/base/core"
	"go_server/service/system"
)

const (
	generateBiz  = "generate"
	rollbackBiz  = "rollback"
	upgradeModel = "model"
)

// 需要配置正确业务库后执行该命令
// 可增加子命令

var bizCmd = &cobra.Command{
	Use:   "biz",
	Short: "根据biz数据库配置注入服务路由模型 可用命令 biz [generate|rollback|model] alias",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
		core.Log.Infof("args:%+v", args)
		if len(args) != 3 {
			core.Log.Infof("args:%+v 必须输入库名+表名  biz [generate|rollback|upgrade] dbAlias tableName[需同步所有请输入all]", args)
			return
		}
		server := &system.AutoService{}
		switch args[0] {
		case generateBiz:
			// 自动生成所有服务 -- 危险操作 demo: go run main.go biz generate app all
			if args[2] == "all" {
				if err := server.AutoServerCodeWithAlias(args[1], ""); err != nil {
					panic(err)
				}
			} else {
				// 指定表名 生成服务demo: go run main.go biz generate app user
				if err := server.AutoServerCodeWithAlias(args[1], args[2]); err != nil {
					panic(err)
				}
			}
			break
		case rollbackBiz:
			// 回滚表所有管理接口、路由、服务 demo: go run main.go biz rollback app all
			if args[2] == "all" {
				if err := server.RollbackWithAlias(args[1], ""); err != nil {
					panic(err)
				}
			} else {
				// 回滚表所有管理接口、路由、服务  -- 危险操作
				// demo: go run main.go biz rollback app user
				if err := server.RollbackWithAlias(args[1], args[2]); err != nil {
					panic(err)
				}
			}
			break
		case upgradeModel:
			// demo: go run main.go biz model app all  // 同步所有表结构到代码
			// demo: go run main.go biz model app sys_config  // 同步指定表结构
			if args[2] == "all" {
				if err := server.ModelAutoForCmd(args[1], ""); err != nil {
					panic(err)
				}
			} else {
				if err := server.ModelAutoForCmd(args[1], args[2]); err != nil {
					panic(err)
				}
			}
			break
		default:
			core.Log.Infof("未知命令：%s", args[0])
			break
		}

		// todo:根据传参数 执行不同命令

	},
}
