package datastore

import "go-transactions-test/dicontainer"

type Account struct {
	AccountId     int     `pg:",pk" json:"account_id"`
	BalanceAmount float32 `pg:",use_zero" json:"balance_amount"`
}

type Transaction struct {
	ID                   int     `pg:",pk"`
	SourceAccountId      int     `json:"source_account_id"`
	DestinationAccountId int     `json:"destination_account_id"`
	SourceAccount        Account `pg:"rel:has-one,join_fk:source_account_id"`
	DestinationAccount   Account `pg:"rel:has-one,join_fk:destination_account_id"`
	Amount               float32 `pg:",use_zero" json:"amount"`
}

func (a Account) String() string {
	return ""
}

func (t Transaction) String() string {
	return ""
}

func CreateAccountQuery(container dicontainer.IDiContainer, model Account) error {
	db := container.GetDbClient()
	_, err := db.Model(model).Insert()
	if err != nil {
		return err
	}
	return nil
}

func SubmitTransactionQuery(container dicontainer.IDiContainer, model Transaction) error {
	db := container.GetDbClient()
	_, err := db.Model(model).Insert()
	if err != nil {
		return err
	}
	return nil
}

func GetAccountDetailsQuery(container dicontainer.IDiContainer, model Account, accountId int) (Account, error) {
	db := container.GetDbClient()
	err := db.Model(model).Where("account_id = ?", accountId).Select()
	if err != nil {
		return Account{}, err
	}
	return model, nil
}
