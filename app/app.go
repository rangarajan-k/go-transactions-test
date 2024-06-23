package app

import (
	"fmt"
	"go-transactions-test/config"
	"go-transactions-test/datastore"
	"go-transactions-test/dicontainer"
	"go-transactions-test/router"
	"log"
	"net/http"
)

type ITransactionSvcApp interface {
	Init(string2 string)
	Start() error
	GetMainConfig() *config.MainConfig
}

type transactionSvcApp struct {
	mainConfig *config.MainConfig
	httpServer *http.Server
	router     router.IRouter
}

func New(filepath string) ITransactionSvcApp {
	app := new(transactionSvcApp)
	app.mainConfig = config.LoadMainConfig(filepath)
	return app
}

func (app *transactionSvcApp) Init(filepath string) {

	//Initialize DB, Controller, Util dependencies
	dbClient := datastore.NewPgClient(app.mainConfig.DBConfig)
	diContainer := dicontainer.NewDiContainer(app.mainConfig)
	diContainer.StartDependenciesInjection(dbClient) //Initialize Router
	app.router = router.NewRouter(app.mainConfig.GinMode)
	app.router.InitRoutes(diContainer)

	//Start Http server here
	app.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", app.mainConfig.Port),
		Handler: app.router.GetEngine(),
	}
}

func (app *transactionSvcApp) Start() error {
	log.Printf("########## Server starting ##########")
	err := app.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("server failed to start %v", err)
		return err
	}
	return nil
}

func (app *transactionSvcApp) GetMainConfig() *config.MainConfig { return app.mainConfig }
