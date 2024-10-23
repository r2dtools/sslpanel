package main

import (
	"backend/internal/app/panel/adapters/cli"
	"log"
)

func main() {
	app, err := cli.GetApp()

	if err != nil {
		log.Fatalln(err)
	}

	if err := app.Run(); err != nil {
		log.Fatalln(err)
	}
}
