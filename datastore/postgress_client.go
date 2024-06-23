package datastore

import (
	"go-transactions-test/config"
)

type PgDataStore struct {
	Config *config.DBConfig
}

func NewPgDataStore(config *config.DBConfig) (*PgDataStore, error) {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})
	defer db.Close()
}
