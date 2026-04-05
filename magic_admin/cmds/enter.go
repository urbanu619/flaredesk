package cmds

import (
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(migrateCmd, apiCmd, bizCmd, configCmd)
}

var rootCmd = &cobra.Command{
	Short: "Cmd命令管理",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
