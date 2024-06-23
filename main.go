package main

import (
	"flag"
	"go-transactions-test/app"
	_ "go-transactions-test/docs"
	"log"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config.json", "path to configuration file")

	//associate the required config then init & start application
	application := app.New(configFilePath)

	application.Init(configFilePath)

	err := application.Start()

	if err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}
