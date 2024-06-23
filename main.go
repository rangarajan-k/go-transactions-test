package main

import (
	"flag"
	"go-transactions-test/app"
	"log"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config.json", "path to configuration file")

	//start app with required config
	application := app.New(configFilePath)

	application.Init(configFilePath)

	err := application.Start()

	if err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}
