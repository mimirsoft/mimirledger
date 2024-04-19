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
	Sign string `json:"sign"`
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
			Sign: "DEBIT",
		},
		{
			Name: "LIABILITY",
			Sign: "CREDIT",
		},
		{
			Name: "EQUITY",
			Sign: "CREDIT",
		},
		{
			Name: "INCOME",
			Sign: "DEBIT",
		},
		{
			Name: "EXPENSE",
			Sign: "CREDIT",
		},
		{
			Name: "GAIN",
			Sign: "CREDIT",
		},
		{
			Name: "LOSS",
			Sign: "DEBIT",
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

// GET /accounts/{accountID}
func (ac *AccountsController) AccountGetByID(ctx context.Context, accountID uint64) (*models.Account, error) {
	account, err := models.RetrieveAccountByID(ac.DataStores, accountID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveAccountByID:%w", err)
	}
	return account, nil
}

// POST /accounts
func (ac *AccountsController) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	err := account.Store(ac.DataStores)
	if err != nil {
		return nil, fmt.Errorf("account.Store:%w", err)
	}
	return account, nil
}

// POST /accounts/{accountID}
func (ac *AccountsController) UpdateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	err := account.Update(ac.DataStores)
	if err != nil {
		return nil, fmt.Errorf("account.Update:%w", err)
	}
	return account, nil
}
