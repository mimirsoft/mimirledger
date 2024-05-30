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
	// set c.TransactionCore
	c.TransactionCore = TransactionCore(eTxn)
	// store the debit/credits
	affectedSubTotalAccountIDs := make(map[uint64]bool)
	affectedBalanceAccountIDs := make(map[uint64]bool)
	for idx := range c.DebitCreditSet {
		// get all the parents of this AccountID
		parentIds, err := getParentsAccountIDs(ds, c.DebitCreditSet[idx].AccountID)
		if err != nil {
			return fmt.Errorf("getParentsAccountIDs:%w", err)
		}
		//fmt.Printf("parentIDs:%v accountID:%d \n", parentIds, c.DebitCreditSet[idx].AccountID)
		for jdx := range parentIds {
			affectedBalanceAccountIDs[parentIds[jdx]] = true
		}
		c.DebitCreditSet[idx].TransactionID = eTxn.TransactionID
		entDC := transactionDCToEntTransactionDC(c.DebitCreditSet[idx])
		err = ds.TransactionDebitCreditStore().Store(&entDC)
		if err != nil {
			return fmt.Errorf("ds.TransactionDCStore().Store:%w [transaction:%+v]", err, eTxn)
		}
		myTransDC := TransactionDebitCredit(entDC)
		c.DebitCreditSet[idx] = &myTransDC
		// if it is a debit / credit, update both Balance and Subtotal
		affectedSubTotalAccountIDs[c.DebitCreditSet[idx].AccountID] = true
		affectedBalanceAccountIDs[c.DebitCreditSet[idx].AccountID] = true
	}
	//fmt.Printf("affectedSubTotalAccountIDs:%v \n", affectedSubTotalAccountIDs)
	//fmt.Printf("affectedBalanceAccountIDs:%v \n", affectedBalanceAccountIDs)

	for idx := range affectedSubTotalAccountIDs {
		err := UpdateSubtotalForAccountID(ds, idx)
		if err != nil {
			return fmt.Errorf("UpdateSubtotalForAccountID:%w [accountID:%d]", err, idx)
		}

	}

	for idx := range affectedBalanceAccountIDs {
		err := UpdateBalanceForAccountID(ds, idx)
		if err != nil {
			return fmt.Errorf("UpdateBalanceForAccountID:%w [accountID:%d]", err, idx)
		}
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
	// set c.TransactionCore
	c.TransactionCore = TransactionCore(eTxn)
	// delete the existing TransactionDebitCredits
	deletedDCs, err := ds.TransactionDebitCreditStore().DeleteForTransactionID(c.TransactionID)
	if err != nil {
		return fmt.Errorf("ds.TransactionDebitCreditStore().DeleteForTransactionID:%w [transaction:%+v]", err, c)
	}
	// store the debit/credits
	affectedSubTotalAccountIDs := make(map[uint64]bool)
	affectedBalanceAccountIDs := make(map[uint64]bool)
	for idx := range deletedDCs {
		// parents of the accounts used in the DC records that were deleted
		parentIds, err := getParentsAccountIDs(ds, deletedDCs[idx].AccountID)
		if err != nil {
			return fmt.Errorf("getParentsAccountIDs:%w", err)
		}
		for jdx := range parentIds {
			affectedBalanceAccountIDs[parentIds[jdx]] = true
		}
		affectedSubTotalAccountIDs[deletedDCs[idx].AccountID] = true
		affectedBalanceAccountIDs[deletedDCs[idx].AccountID] = true
	}
	// store the new debit/credits
	for idx := range c.DebitCreditSet {
		// get all the parents of this AccountID
		parentIds, err := getParentsAccountIDs(ds, c.DebitCreditSet[idx].AccountID)
		if err != nil {
			return fmt.Errorf("getParentsAccountIDs:%w", err)
		}
		for jdx := range parentIds {
			affectedBalanceAccountIDs[parentIds[jdx]] = true
		}
		c.DebitCreditSet[idx].TransactionID = eTxn.TransactionID
		entDC := transactionDCToEntTransactionDC(c.DebitCreditSet[idx])
		err = ds.TransactionDebitCreditStore().Store(&entDC)
		if err != nil {
			return fmt.Errorf("ds.TransactionDCStore().Store:%w [transaction:%+v]", err, eTxn)
		}
		myTransDC := TransactionDebitCredit(entDC)
		c.DebitCreditSet[idx] = &myTransDC
		affectedSubTotalAccountIDs[c.DebitCreditSet[idx].AccountID] = true
		affectedBalanceAccountIDs[c.DebitCreditSet[idx].AccountID] = true
	}
	// update all account subtotals and balances
	for idx := range affectedSubTotalAccountIDs {
		err := UpdateSubtotalForAccountID(ds, idx)
		if err != nil {
			return fmt.Errorf("UpdateSubtotalForAccountID:%w [accountID:%d]", err, idx)
		}

	}

	for idx := range affectedBalanceAccountIDs {
		err := UpdateBalanceForAccountID(ds, idx)
		if err != nil {
			return fmt.Errorf("UpdateBalanceForAccountID:%w [accountID:%d]", err, idx)
		}
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
	myDCSet, err := ds.TransactionDebitCreditStore().GetDCForTransactionID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("TransactionDebitCreditStore().GetDCForTransactionID:%w", err)
	}
	myTrans.DebitCreditSet = entTransactionsDCToTransactionsDC(myDCSet)
	return &myTrans, nil
}

func entTransactionsToTransactions(eTxn []*datastore.Transaction) (txnSet []*Transaction) { //nolint:unused
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

func entTransactionsDCToTransactionsDC(eTxn []*datastore.TransactionDebitCredit) (tdcSet []*TransactionDebitCredit) {
	tdcSet = make([]*TransactionDebitCredit, len(eTxn))
	for idx := range eTxn {
		myTDC := TransactionDebitCredit(*eTxn[idx])
		tdcSet[idx] = &myTDC
	}
	return
}

func transactionDCToEntTransactionDC(txn *TransactionDebitCredit) datastore.TransactionDebitCredit {
	eTxn := datastore.TransactionDebitCredit(*txn)
	return eTxn
}

type TransactionLedger struct {
	TransactionID            uint64
	TransactionDate          time.Time
	TransactionReconcileDate sql.NullTime
	TransactionComment       string
	TransactionReference     string
	IsReconciled             bool
	IsSplit                  bool
	TransactionDCAmount      uint64
	DebitOrCredit            datastore.AccountSign
	// split is a generated field, a comma separated list of the other d/c
	Split string
}

// RetrieveTransactionLedgerForAccountID retrieves all transactions ledger records in an account
func RetrieveTransactionLedgerForAccountID(ds *datastore.Datastores, transactionID uint64) ([]*TransactionLedger, error) {
	ts := ds.TransactionStore()
	eTransSet, err := ts.GetTransactionsForAccount(transactionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("TransactionStore().GetTransactionsForAccount:%w", err)
	}
	transSet := entTransactionsLedgerToTransactionsLedger(eTransSet)
	return transSet, nil
}

func entTransactionsLedgerToTransactionsLedger(eTxn []*datastore.TransactionLedger) (tdcSet []*TransactionLedger) {
	tdcSet = make([]*TransactionLedger, len(eTxn))
	for idx := range eTxn {
		myTDC := TransactionLedger(*eTxn[idx])
		tdcSet[idx] = &myTDC
	}
	return
}