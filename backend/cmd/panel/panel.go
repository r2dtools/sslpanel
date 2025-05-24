package main

import (
	"backend/config"
	"backend/internal/app/panel"
	"backend/internal/pkg/logger"
	"log"
)

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
