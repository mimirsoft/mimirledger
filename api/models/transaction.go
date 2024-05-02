package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"time"
)

type Transaction struct {
	TransactionCore
	DebitCreditSet []*TransactionDebitCredit
}

type TransactionCore struct {
	TransactionID            uint64
	TransactionDate          time.Time
	TransactionReconcileDate sql.NullTime
	TransactionComment       string
	TransactionAmount        uint64
	TransactionReference     string
	IsReconciled             bool
	IsSplit                  bool
}
type TransactionDebitCredit struct {
	TransactionDCID     uint64
	TransactionID       uint64
	AccountID           uint64
	TransactionDCAmount uint64
	DebitOrCredit       datastore.AccountSign
}

var ErrTransactionNotFound = errors.New("transaction not found")
var ErrTransactionNoComment = errors.New("transaction has no comment")
var ErrTransactionDebitCreditAccountInvalid = errors.New("transaction debit-credit has invalid accountID")
var ErrTransactionNoDebitsCredits = errors.New("transaction has no debits/credits")
var ErrTransactionDebitCreditsNotBalanced = errors.New("transaction debits do not equal credits")
var ErrTransactionDebitCreditsZero = errors.New("transaction debits or credits total zero")
var ErrTransactionDebitCreditsIsNeither = errors.New("transaction debits credit missing types")

// Store inserts a Transaction
func (c *Transaction) Store(ds *datastore.Datastores) error { //nolint:gocyclo
	if err := c.validate(); err != nil {
		return fmt.Errorf("c.validate:%w", err)
	}
	// total up debits / credits
	// set transactionAmount = total
	total, err := c.transactionTotal()
	if err != nil {
		return fmt.Errorf("c.transactionTotal():%w [transaction:%+v]", err, c)
	}
	c.TransactionAmount = total
	eTxn := transactionToEntTransaction(c)
	err = ds.TransactionStore().Store(&eTxn)
	if err != nil {
		return fmt.Errorf("ds.TransactionStore().Store:%w [transaction:%+v]", err, eTxn)
	}

	return nil
}

// This whole function should be a transaction for safety
func (c *Transaction) Update(ds *datastore.Datastores) error { //nolint:gocyclo
	if err := c.validate(); err != nil {
		return fmt.Errorf("c.validate:%w", err)
	}
	// total up debits / credits
	// set transactionAmount = total
	total, err := c.transactionTotal()
	if err != nil {
		return fmt.Errorf("c.transactionTotal():%w [transaction:%+v]", err, c)
	}
	c.TransactionAmount = total
	eTxn := transactionToEntTransaction(c)
	err = ds.TransactionStore().Update(&eTxn)
	if err != nil {
		return fmt.Errorf("ds.TransactionStore().Update:%w [transaction:%+v]", err, eTxn)
	}
	return nil
}

func (c *Transaction) validate() (err error) { //nolint:gocyclo
	if c.TransactionComment == "" {
		return ErrTransactionNoComment
	}
	if len(c.DebitCreditSet) == 0 {
		return ErrTransactionNoDebitsCredits
	}
	var debitTotal uint64
	var creditTotal uint64
	for idx := range c.DebitCreditSet {
		if c.DebitCreditSet[idx].AccountID == 0 {
			return ErrTransactionDebitCreditAccountInvalid
		}
		switch accountSign := c.DebitCreditSet[idx].DebitOrCredit; accountSign {
		case datastore.AccountSignDebit:
			debitTotal += c.DebitCreditSet[idx].TransactionDCAmount
		case datastore.AccountSignCredit:
			creditTotal += c.DebitCreditSet[idx].TransactionDCAmount
		default:
			return ErrTransactionDebitCreditsIsNeither
		}
	}
	if debitTotal == 0 || creditTotal == 0 {
		return ErrTransactionDebitCreditsZero
	}

	if debitTotal != creditTotal {
		return ErrTransactionDebitCreditsNotBalanced
	}
	return nil
}

func (c *Transaction) transactionTotal() (uint64, error) { //nolint:gocyclo
	var debitTotal uint64
	var creditTotal uint64
	for idx := range c.DebitCreditSet {
		switch accountSign := c.DebitCreditSet[idx].DebitOrCredit; accountSign {
		case datastore.AccountSignDebit:
			debitTotal += c.DebitCreditSet[idx].TransactionDCAmount
		case datastore.AccountSignCredit:
			creditTotal += c.DebitCreditSet[idx].TransactionDCAmount
		default:
			return 0, ErrTransactionDebitCreditsIsNeither
		}
	}
	if debitTotal != creditTotal {
		return 0, ErrTransactionDebitCreditsNotBalanced
	}
	return debitTotal, nil
}

// RetrieveTransactionsForAccountID retrieves all transactions on an account
func RetrieveTransactionsForAccountID(ds *datastore.Datastores, transactionID uint64) ([]*Transaction, error) {
	ts := ds.TransactionStore()
	eTransSet, err := ts.GetTransactionsForAccount(transactionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("TransactionStore().GetTransactionsForAccount:%w", err)
	}
	transSet := entTransactionsToTransactions(eTransSet)
	return transSet, nil
}

// RetrieveTransactionByID retrieves a specific transactions
func RetrieveTransactionByID(ds *datastore.Datastores, transactionID uint64) (*Transaction, error) {
	ts := ds.TransactionStore()
	eTxn, err := ts.GetByID(transactionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("TransactionStore().GetByID:%w", err)
	}
	myTransCore := TransactionCore(*eTxn)
	myTrans := Transaction{TransactionCore: myTransCore}
	return &myTrans, nil
}

func entTransactionsToTransactions(eTxn []*datastore.Transaction) (txnSet []*Transaction) {
	txnSet = make([]*Transaction, len(eTxn))
	for idx := range eTxn {
		myTransCore := TransactionCore(*eTxn[idx])
		txnSet[idx] = &Transaction{TransactionCore: myTransCore}
	}
	return
}

func transactionToEntTransaction(txn *Transaction) datastore.Transaction {
	etxn := datastore.Transaction(txn.TransactionCore)
	return etxn
}
