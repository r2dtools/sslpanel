package main

import (
	"backend/config"
	"backend/internal/app/panel"
	"backend/internal/pkg/logger"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "r2panel",
	Short: "R2DTools Control Panel",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func main() {
	config, err := config.GetConfig()

	if err != nil {
		log.Fatalln(err)
	}

	appLogger, err := logger.NewLogger(config)

	if err != nil {
		log.Fatalln(err)
	}

	app, err := panel.GetApp(config, appLogger)

	if err != nil {
		log.Fatalln(err)
	}

	if err := app.Run(); err != nil {
		log.Fatalln(err)
	}
}
