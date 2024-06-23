package datastore

import (
	"github.com/go-pg/pg/v10"
)

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

func (a Account) String() string {
	return ""
}

func (t Transaction) String() string {
	return ""
}

func CreateAccountQuery(db *pg.DB, model *Account) error {
	_, err := db.Model(model).Insert()
	if err != nil {
		return err
	}
	return nil
}

func SubmitTransactionQuery(db *pg.DB, model *Transaction) error {
	_, err := db.Model(model).Insert()
	if err != nil {
		return err
	}
	return nil
}

func GetAccountDetailsQuery(db *pg.DB, model *Account) (*Account, error) {
	err := db.Model(model).WherePK().Select()
	if err != nil {
		return &Account{}, err
	}
	return model, nil
}

func GetMultipleAccountDetailsQuery(db *pg.DB, ids []int) ([]*Account, error) {
	var models []*Account
	err := db.Model(&models).Where("account_id in (?)", pg.In(ids)).Select()
	if err != nil {
		return []*Account{}, err
	}
	return models, nil
}

func UpdateAccountBalanceQuery(db *pg.DB, models []*Account) error {
	_, err := db.Model(&models).WherePK().Update()
	if err != nil {
		return err
	}
	return nil
}
