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
func (ac *AccountsController) AccountTypeList(_ context.Context) (*AccountTypeSet, error) {
	at := make([]AccountType, 0, len(datastore.AccountTypeToSign))
	for key, value := range datastore.AccountTypeToSign {
		at = append(at, AccountType{Name: string(key), Sign: string(value)})
	}

	return &AccountTypeSet{AccountTypes: at}, nil
}

// GET /accounts
func (ac *AccountsController) AccountList(_ context.Context) ([]*models.Account, error) {
	accounts, err := models.RetrieveAccounts(ac.DataStores)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveAccounts:%w", err)
	}

	return accounts, nil
}

// GET /accounts/{accountID}
func (ac *AccountsController) AccountGetByID(_ context.Context, accountID uint64) (*models.Account, error) {
	account, err := models.RetrieveAccountByID(ac.DataStores, accountID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveAccountByID:%w", err)
	}

	return account, nil
}

// POST /accounts
func (ac *AccountsController) CreateAccount(_ context.Context, account *models.Account) (*models.Account, error) {
	err := account.Store(ac.DataStores)
	if err != nil {
		return nil, fmt.Errorf("account.Store:%w", err)
	}

	return account, nil
}

// PUT /accounts/{accountID}
func (ac *AccountsController) UpdateAccount(_ context.Context, account *models.Account) (*models.Account, error) {
	err := account.Update(ac.DataStores)
	if err != nil {
		return nil, fmt.Errorf("account.Update:%w", err)
	}

	return account, nil
}

// PUT /accounts/{accountID}/reconciled
func (ac *AccountsController) UpdateAccountReconciledDate(_ context.Context, account *models.Account) (*models.Account, error) {
	myAccount, err := models.RetrieveAccountByID(ac.DataStores, account.AccountID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveAccountByID:%w", err)
	}
	myAccount.AccountReconcileDate = account.AccountReconcileDate
	err = myAccount.UpdateReconciledDate(ac.DataStores)
	if err != nil {
		return nil, fmt.Errorf("account.Update:%w", err)
	}

	return myAccount, nil
}
