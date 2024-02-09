package web

import (
	"context"
	"fmt"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
)

// AccountsController is the controller struct for accounts
type AccountsController struct {
	DataStores *datastore.Datastores
}

// NewAccountsController instantiates a new AccountsController struct
func NewAccountsController(ds *datastore.Datastores) *AccountsController {
	return &AccountsController{
		DataStores: ds,
	}
}

// AccountTypes is for use in accounts controller responses
type AccountType struct {
	Name string `json:"name"`
}

// AccountTypeSet is for use in accounts controller responses
type AccountTypeSet struct {
	AccountTypes []AccountType `json:"accountTypes"`
}

// GET /accounttypes
func (ac *AccountsController) AccountTypeList(ctx context.Context) (*AccountTypeSet, error) {
	at := []AccountType{
		{
			Name: "ASSET",
		},
		{
			Name: "LIABILITY",
		},
		{
			Name: "EQUITY",
		},
		{
			Name: "INCOME",
		},
		{
			Name: "EXPENSE",
		},
		{
			Name: "GAIN",
		},
		{
			Name: "LOSS",
		},
	}
	return &AccountTypeSet{AccountTypes: at}, nil
}

// GET /accounts
func (ac *AccountsController) AccountList(ctx context.Context) ([]*models.Account, error) {
	accounts, err := models.RetrieveAccounts(ac.DataStores)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveAccounts:%w", err)
	}
	return accounts, nil
}

// POST /accounts
func (ac *AccountsController) CreateAccount(ctx context.Context, account models.Account) (*models.Account, error) {
	err := account.Store(ac.DataStores)
	if err != nil {
		return nil, fmt.Errorf("account.Store:%w", err)
	}
	return &account, nil
}
