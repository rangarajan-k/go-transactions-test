package datastore

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"go-transactions-test/config"
	"log"
	"strconv"
)

func NewPgClient(config config.DBConfig) *pg.DB {
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
	err := CreateSchema(db)
	if err != nil {
		log.Printf("Error creating schema: %v", err)
	}

	return db
}

func CreateSchema(db *pg.DB) error {
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
