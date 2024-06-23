package app

import (
	"go-transactions-test/config"
	"net/http"
	"sync"
)

type ITransactionSvcApp interface {
	Init(string2 string)
	Start() error
	GetWaitGroupVar() *sync.WaitGroup
	GetMainConfig()
}

type transactionSvcApp struct {
	mainConfig *config.MainConfig
	httpServer *http.Server
	router     router.IRouter
	wg         sync.WaitGroup
}

func New(filepath string) ITransactionSvcApp {
	app := new(transactionSvcApp)
	app.mainConfig = config.LoadMainConfig(filepath)
	return app
}

func (app *transactionSvcApp) Init(filepath string) {
	//Initialize DB connection

	//Initialize Router

}

func (app *transactionSvcApp) GetWaitGroupVar() *sync.WaitGroup  { return &app.wg }
func (app *transactionSvcApp) GetMainConfig() *config.MainConfig { return app.mainConfig }
