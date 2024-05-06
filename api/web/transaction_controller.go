package web

import (
	"context"
	"fmt"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
)

// TransactionsController is the controller struct for transactions
type TransactionsController struct {
	DataStores *datastore.Datastores
}

// NewTransactionsController instantiates a new AccountsController struct
func NewTransactionsController(ds *datastore.Datastores) *TransactionsController {
	return &TransactionsController{
		DataStores: ds,
	}
}

// POST /tranasctions
func (tc *TransactionsController) CreateTransaction(ctx context.Context, myTxn *models.Transaction) (*models.Transaction,
	error) {
	err := myTxn.Store(tc.DataStores)
	if err != nil {
		return nil, fmt.Errorf("myTxn.Store:%w", err)
	}
	return myTxn, nil
}

// GET /tranasctions/account/{accountID}
func (tc *TransactionsController) GetTransactionsForAccount(ctx context.Context, transactionID uint64) ([]*models.TransactionLedger,
	error) {
	myTxn, err := models.RetrieveTransactionLedgerForAccountID(tc.DataStores, transactionID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveTransactionsForAccountID:%w", err)
	}
	return myTxn, nil
}

// GET /tranasctions/{transactionID}
func (tc *TransactionsController) GetTransactionByID(ctx context.Context, transactionID uint64) (*models.Transaction, error) {
	myTxn, err := models.RetrieveTransactionByID(tc.DataStores, transactionID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveTransactionByID:%w", err)
	}
	return myTxn, nil
}

// PUT /tranasctions/{transactionID}
func (tc *TransactionsController) UpdateTransaction(ctx context.Context, myTxn *models.Transaction) (*models.Transaction,
	error) {
	err := myTxn.Update(tc.DataStores)
	if err != nil {
		return nil, fmt.Errorf("myTxn.Update:%w", err)
	}
	return myTxn, nil
}
