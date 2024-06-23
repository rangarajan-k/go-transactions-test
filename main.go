package main

import (
	"flag"
	"go-transactions-test/app"
	"log"
)

func main() {
	var configFilePath, serverPort string
	flag.StringVar(&configFilePath, "config", "config.json", "path to configuration file")
	flag.StringVar(&serverPort, "port", "8080", "port to listen on")

	//start app with required config
	application := app.New(configFilePath)

	application.Init(configFilePath)

	err := application.Start()

	if err != nil {
		log.Fatalf("Error starting application: %v", err)
	}
}
