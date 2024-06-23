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

/* We can choose what dependencies to give to which API controller groups
In this case its single controller only
This can include any external connectors like http client / queue stores etc */

func (di *DiContainer) StartDependenciesInjection(pgClient *pg.DB) {

	di.PgClient = pgClient
	di.TransactionServiceController = controller.NewTransactionServiceController(di.Config, di.PgClient)
}

func (di *DiContainer) GetDiContainer() *DiContainer { return di }
func (di *DiContainer) GetDbClient() *pg.DB          { return di.PgClient }
