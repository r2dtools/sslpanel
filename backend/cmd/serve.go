package cmd

import (
	"backend/server"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves the application",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := server.Init()

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
