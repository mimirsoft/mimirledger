package web

import (
	"context"
	"fmt"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
	"time"
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

// POST /transactions
func (tc *TransactionsController) CreateTransaction(_ context.Context, myTxn *models.Transaction) (*models.Transaction,
	error) {
	err := myTxn.Store(tc.DataStores)
	if err != nil {
		return nil, fmt.Errorf("myTxn.Store:%w", err)
	}
	// after creating transaction, update balance on all affected accounts
	return myTxn, nil
}

// GET /transactions/account/{accountID}
func (tc *TransactionsController) GetTransactionsForAccount(_ context.Context, accountID uint64) (*models.Account,
	[]*models.TransactionLedger,
	error) {
	account, err := models.RetrieveAccountByID(tc.DataStores, accountID)
	if err != nil {
		return nil, nil, fmt.Errorf("models.RetrieveTransactionsForAccountID:%w", err)
	}

	myTxn, err := models.RetrieveTransactionLedgerForAccountID(tc.DataStores, accountID)
	if err != nil {
		return nil, nil, fmt.Errorf("models.RetrieveTransactionLedgerForAccountID:%w", err)
	}

	return account, myTxn, nil
}

// GET /transactions/account/{accountID}/unreconciled?date=<date>
func (tc *TransactionsController) GetUnreconciledTransactionsOnAccount(_ context.Context, accountID uint64,
	searchDate time.Time) ([]*models.TransactionReconciliation,
	error) {
	account, err := models.RetrieveAccountByID(tc.DataStores, accountID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveAccountByID:%w", err)
	}

	myTxn, err := models.RetrieveUnreconciledTransactionsForDate(tc.DataStores, account.AccountLeft, account.AccountRight,
		searchDate, account.AccountReconcileDate.Time)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveUnreconciledTransactionsForDate:%w", err)
	}

	return myTxn, nil
}

// GET /transactions/{transactionID}
func (tc *TransactionsController) GetTransactionByID(_ context.Context,
	transactionID uint64) (*models.Transaction, error) {
	myTxn, err := models.RetrieveTransactionByID(tc.DataStores, transactionID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveTransactionByID:%w", err)
	}

	return myTxn, nil
}

// PUT /transactions/{transactionID}
func (tc *TransactionsController) UpdateTransaction(_ context.Context, myTxn *models.Transaction) (*models.Transaction,
	error) {
	err := myTxn.Update(tc.DataStores)
	if err != nil {
		return nil, fmt.Errorf("myTxn.Update:%w", err)
	}

	return myTxn, nil
}

// PUT /transactions/{transactionID}/reconciled
func (tc *TransactionsController) UpdateReconciled(_ context.Context, myTxn *models.Transaction) (*models.Transaction,
	error) {
	readTxn, err := models.RetrieveTransactionByID(tc.DataStores, myTxn.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveTransactionByID:%w", err)
	}

	readTxn.IsReconciled = true
	readTxn.TransactionReconcileDate = myTxn.TransactionReconcileDate

	err = readTxn.UpdateReconciled(tc.DataStores)
	if err != nil {
		return nil, fmt.Errorf("readTxn.UpdateReconciled:%w", err)
	}

	return readTxn, nil
}

// PUT /transactions/{transactionID}/unreconciled
func (tc *TransactionsController) UpdateUnreconciled(_ context.Context, myTxn *models.Transaction) (*models.Transaction,
	error) {
	readTxn, err := models.RetrieveTransactionByID(tc.DataStores, myTxn.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveTransactionByID:%w", err)
	}

	readTxn.IsReconciled = false

	err = readTxn.UpdateUnreconciled(tc.DataStores)
	if err != nil {
		return nil, fmt.Errorf("myTxn.UpdateUnreconciled:%w", err)
	}

	return readTxn, nil
}

// DELETE /transactions/{transactionID}
func (tc *TransactionsController) DeleteTransaction(_ context.Context, transactionID uint64) (*models.Transaction,
	error) {
	// retrieve the transaction for checking before deletion
	myTxn, err := models.RetrieveTransactionByID(tc.DataStores, transactionID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveTransactionByID:%w", err)
	}

	err = myTxn.Delete(tc.DataStores)
	if err != nil {
		return nil, fmt.Errorf("myTxn.Delete:%w", err)
	}

	return myTxn, nil
}
