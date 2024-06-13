package response

import (
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
	"time"
)

// TransactionSet is for use in transaction controller responses
type TransactionSet struct {
	Transactions []*Transaction `json:"transactions"`
}

type Transaction struct {
	TransactionID            uint64                    `json:"transactionID,omitempty"`
	TransactionDate          time.Time                 `json:"transactionDate"`
	TransactionReconcileDate time.Time                 `json:"transactionReconcileDate"`
	TransactionComment       string                    `json:"transactionComment"`
	TransactionAmount        uint64                    `json:"transactionAmount"`
	TransactionReference     string                    `json:"transactionReference"` // this could be a check number, batch ,etc
	IsReconciled             bool                      `json:"isReconciled"`
	IsSplit                  bool                      `json:"isSplit"`
	DebitCreditSet           []*TransactionDebitCredit `json:"debitCreditSet"`
}
type TransactionDebitCredit struct {
	TransactionDCID     uint64                `json:"transactionDCID"`
	TransactionID       uint64                `json:"transactionID"`
	AccountID           uint64                `json:"accountID"`
	TransactionDCAmount uint64                `json:"transactionDCAmount"`
	DebitOrCredit       datastore.AccountSign `json:"debitOrCredit"`
}

func TransactionToRespTransaction(trans *models.Transaction) *Transaction {
	myDCSet := ConvertDebitCreditsToReDebitCreditSet(trans.DebitCreditSet)
	myTrans := Transaction{
		TransactionID:            trans.TransactionID,
		TransactionDate:          trans.TransactionDate,
		TransactionReconcileDate: trans.TransactionReconcileDate.Time,
		TransactionComment:       trans.TransactionComment,
		TransactionAmount:        trans.TransactionAmount,
		TransactionReference:     trans.TransactionReference,
		IsReconciled:             trans.IsReconciled,
		IsSplit:                  trans.IsSplit,
		DebitCreditSet:           myDCSet,
	}
	return &myTrans
}

// ConvertReqDebitCreditsToDebitCreditSet converts []*models.TransactionDebitCredit to TransactionDebitCredit
func ConvertDebitCreditsToReDebitCreditSet(dcSet []*models.TransactionDebitCredit) []*TransactionDebitCredit {
	var mset = make([]*TransactionDebitCredit, len(dcSet))
	for idx := range dcSet {
		myDS := TransactionDebitCredit(*dcSet[idx])
		mset[idx] = &myDS
	}
	return mset
}

// ConvertAccountsToRespAccounts converts []models.Account to AccountSet
func ConvertTransactionsToRespTransactions(txns []*models.Transaction) *TransactionSet {
	var tas = make([]*Transaction, len(txns))
	for idx := range txns {
		tas[idx] = TransactionToRespTransaction(txns[idx])
	}
	return &TransactionSet{Transactions: tas}
}

// TransactionLedgerSet is for use in transaction controller responses for a single acount
type TransactionLedgerSet struct {
	AccountID       uint64               `json:"accountID"`
	AccountSign     string               `json:"accountSign"`
	AccountName     string               `json:"accountName"`
	AccountFullName string               `json:"accountFullName"`
	Transactions    []*TransactionLedger `json:"transactions"`
}

type TransactionLedger struct {
	TransactionID            uint64                `json:"transactionID"`
	TransactionDate          time.Time             `json:"transactionDate"`
	TransactionReconcileDate time.Time             `json:"transactionReconcileDate"`
	TransactionComment       string                `json:"transactionComment"`
	TransactionReference     string                `json:"transactionReference"` // this could be a check number, batch ,etc
	IsReconciled             bool                  `json:"isReconciled"`
	IsSplit                  bool                  `json:"isSplit"`
	TransactionDCAmount      uint64                `json:"transactionDCAmount"`
	DebitOrCredit            datastore.AccountSign `json:"debitOrCredit"`
	Split                    string                `json:"split"` // this could be a check number, batch ,etc
}

// ConvertTransactionLedgerToRespTransactionLedger converts []models.TransactionLedger to TransactionLedger
func ConvertTransactionLedgerToRespTransactionLedger(act *models.Account, txns []*models.TransactionLedger) *TransactionLedgerSet {
	var tas = make([]*TransactionLedger, len(txns))
	for idx := range txns {
		tas[idx] = ConvertTransactionLedgerToRespTransactionLeger(txns[idx])
	}
	return &TransactionLedgerSet{
		AccountID:       act.AccountID,
		AccountName:     act.AccountName,
		AccountFullName: act.AccountFullName,
		AccountSign:     string(act.AccountSign),
		Transactions:    tas}
}

func ConvertTransactionLedgerToRespTransactionLeger(trans *models.TransactionLedger) *TransactionLedger {
	respTransLedger := TransactionLedger{
		TransactionID:            trans.TransactionID,
		TransactionDate:          trans.TransactionDate,
		TransactionReconcileDate: trans.TransactionReconcileDate.Time,
		TransactionComment:       trans.TransactionComment,
		TransactionDCAmount:      trans.TransactionDCAmount,
		TransactionReference:     trans.TransactionReference,
		IsReconciled:             trans.IsReconciled,
		IsSplit:                  trans.IsSplit,
		Split:                    trans.Split,
		DebitOrCredit:            trans.DebitOrCredit,
	}
	return &respTransLedger
}

type TransactionReconciliation struct {
	TransactionID            uint64                `json:"transactionID"`
	TransactionDate          time.Time             `json:"transactionDate"`
	TransactionReconcileDate time.Time             `json:"transactionReconcileDate"`
	TransactionComment       string                `json:"transactionComment"`
	TransactionReference     string                `json:"transactionReference"` // this could be a check number, batch ,etc
	IsReconciled             bool                  `json:"isReconciled"`
	IsSplit                  bool                  `json:"isSplit"`
	TransactionDCAmount      uint64                `json:"transactionDCAmount"`
	DebitOrCredit            datastore.AccountSign `json:"debitOrCredit"`
	Split                    string                `json:"split"` // this could be a check number, batch ,etc
}

// ConvertTransactionRecSetToRespTransactionRecSet converts []models.TransactionReconciliation to []*TransactionReconciliation
func ConvertTransactionRecSetToRespTransactionRecSet(txns []*models.TransactionReconciliation) []*TransactionReconciliation {
	var tas = make([]*TransactionReconciliation, len(txns))
	for idx := range txns {
		tas[idx] = ConvertTransactionReconcileToRespTransactionReconcile(txns[idx])
	}
	return tas
}
func ConvertTransactionReconcileToRespTransactionReconcile(trans *models.TransactionReconciliation) *TransactionReconciliation {
	respTransLedger := TransactionReconciliation{
		TransactionID:            trans.TransactionID,
		TransactionDate:          trans.TransactionDate,
		TransactionReconcileDate: trans.TransactionReconcileDate.Time,
		TransactionComment:       trans.TransactionComment,
		TransactionDCAmount:      trans.TransactionDCAmount,
		TransactionReference:     trans.TransactionReference,
		IsReconciled:             trans.IsReconciled,
		IsSplit:                  trans.IsSplit,
		Split:                    trans.Split,
		DebitOrCredit:            trans.DebitOrCredit,
	}
	return &respTransLedger
}
