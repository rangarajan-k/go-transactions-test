package dicontainer

import {"go-transactions-test/config"}

type IDiContainer interface {
	StartDependenciesInjection()
	GetDiContainer() *DiContainer
	GetMainConfig() *config.MainConfig
}

type DiContainer struct {
	Config *config.MainConfig
	Post
}