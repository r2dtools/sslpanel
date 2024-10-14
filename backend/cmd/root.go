package cmd

import (
	"backend/db"
	"backend/modules"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "r2sm",
	Short: "R2 Server monitoring system",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(func() {
		onError(db.Init())
	})
	modules.InitModulesCli(rootCmd)
}

func onError(err error) {
	if err != nil {
		panic(err)
	}
}
