package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "certificateMonitor",
	Short: "Certificate Monitor module cli",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// InitCmd initiates module cli commands
func InitCmd(rCmd *cobra.Command) {
	rCmd.AddCommand(rootCmd)
}
