package dicontainer

import (
	"github.com/go-pg/pg/v10"
	"go-transactions-test/config"
	"go-transactions-test/controller"
	"go-transactions-test/datastore"
)

type IDiContainer interface {
	StartDependenciesInjection()
	GetDiContainer() *DiContainer
	GetDbClient() *pg.DB
}

type DiContainer struct {
	Config                       *config.MainConfig
	PgClient                     *pg.DB
	TransactionServiceController controller.ITransactionServiceController
}

func NewDiContainer(config *config.MainConfig) IDiContainer { return &DiContainer{Config: config} }

func (di *DiContainer) StartDependenciesInjection() {
	//initialise pgClient here
	di.PgClient = datastore.NewPgClient(di.Config.DBConfig)

	di.TransactionServiceController = controller.NewTransactionServiceController(di.Config, di)
}

func (di *DiContainer) GetDiContainer() *DiContainer {
	return di
}
func (di *DiContainer) GetDbClient() *pg.DB { return di.PgClient }
