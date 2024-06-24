package datastore

import (
	"github.com/go-pg/pg/v10"
)

type IPostgressStore interface {
	CreateAccountQuery(model *Account) error
	GetAccountDetailsQuery(model *Account) (*Account, error)
	GetMultipleAccountDetailsQuery(tx *pg.Tx, ids []int) ([]*Account, error)
	UpdateAccountBalanceQuery(tx *pg.Tx, models []*Account) error
	SubmitTransactionQuery(tx *pg.Tx, model *Transaction) error
}

type PostgressStore struct {
	DB *pg.DB
}

func NewPostgressStore(db *pg.DB) IPostgressStore {
	return &PostgressStore{
		DB: db,
	}
}

type Account struct {
	AccountId     int     `pg:",pk" json:"account_id"`
	BalanceAmount float32 `pg:",use_zero" json:"balance"`
}

type Transaction struct {
	Id                   int     `pg:",pk"`
	SourceAccountId      int     `json:"source_account_id" pg:"join_fk:account_id"`
	DestinationAccountId int     `json:"destination_account_id" pg:"join_fk:account_id"`
	Amount               float32 `pg:",use_zero" json:"amount"`
}

// CreateAccountQuery straight forward create account if account id doesnt exist
func (pgStore *PostgressStore) CreateAccountQuery(model *Account) error {
	db := pgStore.DB
	_, err := db.Model(model).Insert()
	if err != nil {
		return err
	}
	return nil
}

// GetAccountDetailsQuery queries and returns single account detail based on account id
func (pgStore *PostgressStore) GetAccountDetailsQuery(model *Account) (*Account, error) {
	db := pgStore.DB
	err := db.Model(model).WherePK().Select()
	if err != nil {
		return &Account{}, err
	}
	return model, nil
}

// GetMultipleAccountDetailsQuery takes multiple account ids as inputs then returns array of structs
func (pgStore *PostgressStore) GetMultipleAccountDetailsQuery(tx *pg.Tx, ids []int) ([]*Account, error) {
	var models []*Account
	err := tx.Model(&models).Where("account_id in (?)", pg.In(ids)).Select()
	if err != nil {
		return []*Account{}, err
	}
	return models, nil
}

// UpdateAccountBalanceQuery Update multiple accounts passed as array of structs based on the primary key
func (pgStore *PostgressStore) UpdateAccountBalanceQuery(tx *pg.Tx, models []*Account) error {
	_, err := tx.Model(&models).WherePK().Update()
	if err != nil {
		return err
	}
	return nil
}

// SubmitTransactionQuery create a transaction audit entry
func (pgStore *PostgressStore) SubmitTransactionQuery(tx *pg.Tx, model *Transaction) error {
	_, err := tx.Model(model).Insert()
	if err != nil {
		return err
	}
	return nil
}
