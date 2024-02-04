package web

import (
	"context"
	"fmt"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
)

// HealthController is the controller struct for the health check endpoint
type HealthController struct {
	DataStores *datastore.Datastores
}

// GET /api/health
// HEAD /api/health
func (healthController *HealthController) HealthCheck(ctx context.Context) (err error) {
	return nil
}

// NewHealthController instantiates a new HealthController struct
func NewHealthController(ds *datastore.Datastores) *HealthController {
	return &HealthController{
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

// GET /accounttypes
func (ac *AccountsController) AccountTypeList(ctx context.Context) (*AccountTypeSet, error) {
	at := []AccountType{
		{
			Name: "assets",
		},
		{
			Name: "liability",
		},
		{
			Name: "equity3",
		},
	}
	return &AccountTypeSet{AccountTypes: at}, nil
}

// GET /accounts
func (ac *AccountsController) AccountList(ctx context.Context) ([]models.Account, error) {
	accounts, err := models.RetrieveAccounts(ac.DataStores)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveAccounts:%w", err)
	}
	return accounts, nil
}
