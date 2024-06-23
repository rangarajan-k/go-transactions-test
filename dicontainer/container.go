package dicontainer

import (
	"github.com/go-pg/pg/v10"
	"go-transactions-test/config"
	"go-transactions-test/controller"
)

type IDiContainer interface {
	StartDependenciesInjection(db *pg.DB)
	GetDiContainer() *DiContainer
	GetDbClient() *pg.DB
}

type DiContainer struct {
	Config                       *config.MainConfig
	PgClient                     *pg.DB
	TransactionServiceController controller.ITransactionServiceController
}

func NewDiContainer(config *config.MainConfig) IDiContainer { return &DiContainer{Config: config} }

func (di *DiContainer) StartDependenciesInjection(pgClient *pg.DB) {
	//initialise pgClient here
	di.PgClient = pgClient
	di.TransactionServiceController = controller.NewTransactionServiceController(di.Config, di.PgClient)
}

func (di *DiContainer) GetDiContainer() *DiContainer {
	return di
}
func (di *DiContainer) GetDbClient() *pg.DB { return di.PgClient }
