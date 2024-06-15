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
func (c *Transaction) Store(dStores *datastore.Datastores) error {
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

	err = dStores.TransactionStore().Store(&eTxn)
	if err != nil {
		return fmt.Errorf("ds.TransactionStore().Store:%w [transaction:%+v]", err, eTxn)
	}
	// set c.TransactionCore
	c.TransactionCore = TransactionCore(eTxn)
	// store the debit/credits
	affectedSubTotalAccountIDs := make(map[uint64]bool)
	affectedBalanceAccountIDs := make(map[uint64]bool)

	err = c.handleDCSetStore(dStores, affectedSubTotalAccountIDs, affectedBalanceAccountIDs)
	if err != nil {
		return fmt.Errorf("c.handleDCSetStore:%w", err)
	}

	err = updateSubtotalsAndBalances(dStores, affectedSubTotalAccountIDs, affectedBalanceAccountIDs)
	if err != nil {
		return fmt.Errorf("updateSubtotalsAndBalances:%w", err)
	}
	return nil
}

func (c *Transaction) handleDCSetStore(dStores *datastore.Datastores, affectedSubTotalAccountIDs map[uint64]bool,
	affectedBalanceAccountIDs map[uint64]bool) error {
	for idx := range c.DebitCreditSet {
		// get all the parents of this AccountID
		parentIds, err := getParentsAccountIDs(dStores, c.DebitCreditSet[idx].AccountID)
		if err != nil {
			return fmt.Errorf("getParentsAccountIDs:%w", err)
		}

		for jdx := range parentIds {
			affectedBalanceAccountIDs[parentIds[jdx]] = true
		}

		c.DebitCreditSet[idx].TransactionID = c.TransactionID
		entDC := transactionDCToEntTransactionDC(c.DebitCreditSet[idx])

		err = dStores.TransactionDebitCreditStore().Store(&entDC)
		if err != nil {
			return fmt.Errorf("ds.TransactionDCStore().Store:%w [transaction:%+v]", err, c)
		}

		myTransDC := TransactionDebitCredit(entDC)
		c.DebitCreditSet[idx] = &myTransDC
		// if it is a debit / credit, update both Balance and Subtotal
		affectedSubTotalAccountIDs[c.DebitCreditSet[idx].AccountID] = true
		affectedBalanceAccountIDs[c.DebitCreditSet[idx].AccountID] = true
	}
	return nil
}

func updateSubtotalsAndBalances(dStores *datastore.Datastores, affectedSubTotalAccountIDs map[uint64]bool,
	affectedBalanceAccountIDs map[uint64]bool) error {
	for idx := range affectedSubTotalAccountIDs {
		err := UpdateSubtotalForAccountID(dStores, idx)
		if err != nil {
			return fmt.Errorf("UpdateSubtotalForAccountID:%w [accountID:%d]", err, idx)
		}

	}

	for idx := range affectedBalanceAccountIDs {
		err := UpdateBalanceForAccountID(dStores, idx)
		if err != nil {
			return fmt.Errorf("UpdateBalanceForAccountID:%w [accountID:%d]", err, idx)
		}
	}
	return nil
}

// This whole function should be a transaction for safety
func (c *Transaction) Update(dStores *datastore.Datastores) error {
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

	err = dStores.TransactionStore().Update(&eTxn)
	if err != nil {
		return fmt.Errorf("ds.TransactionStore().Update:%w [transaction:%+v]", err, eTxn)
	}
	// set c.TransactionCore
	c.TransactionCore = TransactionCore(eTxn)
	affectedSubTotalAccountIDs := make(map[uint64]bool)
	affectedBalanceAccountIDs := make(map[uint64]bool)
	// delete the existing DC for transaction
	err = c.handleDeletedDCs(dStores, affectedSubTotalAccountIDs, affectedBalanceAccountIDs)
	if err != nil {
		return fmt.Errorf("c.handleDeletedDCs:%w", err)
	}
	// store the new debit/credits
	err = c.handleDCSetStore(dStores, affectedSubTotalAccountIDs, affectedBalanceAccountIDs)
	if err != nil {
		return fmt.Errorf("c.handleDCSetStore:%w", err)
	}
	// update all account subtotals and balances
	err = updateSubtotalsAndBalances(dStores, affectedSubTotalAccountIDs, affectedBalanceAccountIDs)
	if err != nil {
		return fmt.Errorf("updateSubtotalsAndBalances:%w", err)
	}
	return nil
}
func (c *Transaction) handleDeletedDCs(dStores *datastore.Datastores, affectedSubTotalAccountIDs map[uint64]bool,
	affectedBalanceAccountIDs map[uint64]bool) error {
	deletedDCs, err := dStores.TransactionDebitCreditStore().DeleteForTransactionID(c.TransactionID)
	if err != nil {
		return fmt.Errorf("ds.TransactionDebitCreditStore().DeleteForTransactionID:%w [transaction:%+v]", err, c)
	}

	for idx := range deletedDCs {
		// parents of the accounts used in the DC records that were deleted
		parentIds, err := getParentsAccountIDs(dStores, deletedDCs[idx].AccountID)
		if err != nil {
			return fmt.Errorf("getParentsAccountIDs:%w", err)
		}

		for jdx := range parentIds {
			affectedBalanceAccountIDs[parentIds[jdx]] = true
		}

		affectedSubTotalAccountIDs[deletedDCs[idx].AccountID] = true
		affectedBalanceAccountIDs[deletedDCs[idx].AccountID] = true
	}
	return nil
}

// This whole function should be a transaction for safety, delete the DC for a transaction first
// record the affected accounts, then perform updates
func (c *Transaction) Delete(dStores *datastore.Datastores) error {
	eTxn := transactionToEntTransaction(c)
	// store the data about the affected account for updating at end
	affectedSubTotalAccountIDs := make(map[uint64]bool)
	affectedBalanceAccountIDs := make(map[uint64]bool)
	// delete the existing DC for transaction
	err := c.handleDeletedDCs(dStores, affectedSubTotalAccountIDs, affectedBalanceAccountIDs)
	if err != nil {
		return fmt.Errorf("c.handleDCSetStore:%w", err)
	}
	// delete the existing Transaction
	err = dStores.TransactionStore().Delete(&eTxn)
	if err != nil {
		return fmt.Errorf("ds.TransactionStore().Delete:%w [transaction:%+v]", err, c)
	}
	// update all account subtotals and balances
	err = updateSubtotalsAndBalances(dStores, affectedSubTotalAccountIDs, affectedBalanceAccountIDs)
	if err != nil {
		return fmt.Errorf("updateSubtotalsAndBalances:%w", err)
	}
	return nil
}

func (c *Transaction) validate() (err error) {
	if c.TransactionComment == "" {
		return ErrTransactionNoComment
	}

	if len(c.DebitCreditSet) == 0 {
		return ErrTransactionNoDebitsCredits
	}

	var (
		debitTotal  uint64
		creditTotal uint64
	)

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

func (c *Transaction) transactionTotal() (uint64, error) {
	var (
		debitTotal  uint64
		creditTotal uint64
	)

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

var ErrReconciledDateInvalid = errors.New("ReconciledDateInvalid")

// This whole function should be a transaction for safety
func (c *Transaction) UpdateReconciled(dStores *datastore.Datastores) error {
	if !c.TransactionReconcileDate.Valid {
		return ErrReconciledDateInvalid
	}

	eTxn := transactionToEntTransaction(c)

	err := dStores.TransactionStore().SetIsReconciled(&eTxn)
	if err != nil {
		return fmt.Errorf("ds.TransactionStore().SetIsReconciled:%w [transaction:%+v]", err, eTxn)
	}

	err = dStores.TransactionStore().SetTransactionReconcileDate(&eTxn)
	if err != nil {
		return fmt.Errorf("ds.TransactionStore().SetTransactionReconcileDate:%w [transaction:%+v]", err, eTxn)
	}
	// set c.TransactionCore
	c.TransactionCore = TransactionCore(eTxn)
	return nil
}

// This whole function should be a transaction for safety
func (c *Transaction) UpdateUnreconciled(dStores *datastore.Datastores) error {
	eTxn := transactionToEntTransaction(c)

	err := dStores.TransactionStore().SetIsReconciled(&eTxn)
	if err != nil {
		return fmt.Errorf("ds.TransactionStore().SetIsReconciled:%w [transaction:%+v]", err, eTxn)
	}
	// set c.TransactionCore
	c.TransactionCore = TransactionCore(eTxn)
	return nil
}

// RetrieveTransactionByID retrieves a specific transactions
func RetrieveTransactionByID(dStores *datastore.Datastores, transactionID uint64) (*Transaction, error) {
	ts := dStores.TransactionStore()

	eTxn, err := ts.GetByID(transactionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("TransactionStore().GetByID:%w", err)
	}

	myTransCore := TransactionCore(*eTxn)
	myTrans := Transaction{TransactionCore: myTransCore}

	myDCSet, err := dStores.TransactionDebitCreditStore().GetDCForTransactionID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("TransactionDebitCreditStore().GetDCForTransactionID:%w", err)
	}

	myTrans.DebitCreditSet = entTransactionsDCToTransactionsDC(myDCSet)
	return &myTrans, nil
}

func entTransactionsToTransactions(eTxn []*datastore.Transaction) []*Transaction { //nolint:unused
	txnSet := make([]*Transaction, len(eTxn))

	for idx := range eTxn {
		myTransCore := TransactionCore(*eTxn[idx])
		txnSet[idx] = &Transaction{TransactionCore: myTransCore}
	}
	return txnSet
}

func transactionToEntTransaction(txn *Transaction) datastore.Transaction {
	etxn := datastore.Transaction(txn.TransactionCore)
	return etxn
}

func entTransactionsDCToTransactionsDC(eTxn []*datastore.TransactionDebitCredit) []*TransactionDebitCredit {
	tdcSet := make([]*TransactionDebitCredit, len(eTxn))

	for idx := range eTxn {
		myTDC := TransactionDebitCredit(*eTxn[idx])
		tdcSet[idx] = &myTDC
	}
	return tdcSet
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
func RetrieveTransactionLedgerForAccountID(dStores *datastore.Datastores, transactionID uint64) ([]*TransactionLedger, error) {
	ts := dStores.TransactionStore()

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

type TransactionReconciliation struct {
	TransactionID            uint64
	AccountID                uint64
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

// RetrieveUnreconciledTransactionsForDate retrieves unreconciled transactions for a date for an account
func RetrieveUnreconciledTransactionsForDate(dStores *datastore.Datastores, accountLeft, accountRight uint64,
	searchLimitDate time.Time, reconciledCutoffDate time.Time) ([]*TransactionReconciliation, error) {
	ts := dStores.TransactionStore()

	eTransSet, err := ts.GetUnreconciledTransactionsOnAccountForDate(accountLeft, accountRight, searchLimitDate, reconciledCutoffDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("TransactionStore().GetTransactionsForAccount:%w", err)
	}

	transSet := entTransactionsRecToTransactionsRec(eTransSet)
	return transSet, nil
}

func entTransactionsRecToTransactionsRec(eTxn []*datastore.TransactionReconciliation) []*TransactionReconciliation {
	tdcSet := make([]*TransactionReconciliation, len(eTxn))

	for idx := range eTxn {
		myTDC := TransactionReconciliation(*eTxn[idx])
		tdcSet[idx] = &myTDC
	}
	return tdcSet
}

func entTransactionsLedgerToTransactionsLedger(eTxn []*datastore.TransactionLedger) []*TransactionLedger {
	tdcSet := make([]*TransactionLedger, len(eTxn))

	for idx := range eTxn {
		myTDC := TransactionLedger(*eTxn[idx])
		tdcSet[idx] = &myTDC
	}
	return tdcSet
}
