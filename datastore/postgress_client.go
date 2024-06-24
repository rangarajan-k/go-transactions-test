package datastore

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"go-transactions-test/config"
	"log"
	"strconv"
)

type IPgClientInterface interface {
	NewPgClient(config config.DBConfig) *pg.DB
	CreateSchema(db *pg.DB) error
	GetDbClient() *pg.DB
}

type PgClient struct {
	Config   config.DBConfig
	DbClient *pg.DB
}

func NewPostgresClient(mainConfig *config.MainConfig) IPgClientInterface {
	return &PgClient{
		Config: mainConfig.DBConfig,
	}
}

func (p *PgClient) GetDbClient() *pg.DB {
	return p.DbClient
}

func (p *PgClient) NewPgClient(config config.DBConfig) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		User:     config.Username,
		Password: config.Password,
		Database: config.DatabaseName,
	})

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	//Create Initial Schema
	err := p.CreateSchema(db)
	if err != nil {
		log.Printf("Error creating schema: %v", err)
	}
	p.DbClient = db
	return db
}

func (p *PgClient) CreateSchema(db *pg.DB) error {
	models := []interface{}{
		(*Account)(nil),
		(*Transaction)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			log.Printf("Error creating schema: %v", err)
		}
	}
	return nil
}
