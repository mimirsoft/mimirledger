package response

import (
	"database/sql"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
	"time"
)

// TransactionSet is for use in transaction controller responses
type TransactionSet struct {
	Transactions []*Transaction `json:"transactions"`
}

type Transaction struct {
	TransactionID            uint64                    `json:"transaction_id,omitempty"`
	TransactionDate          time.Time                 `json:"transaction_date"`
	TransactionReconcileDate sql.NullTime              `json:"transaction_reconcile_date"`
	TransactionComment       string                    `json:"transaction_comment"`
	TransactionAmount        uint64                    `json:"transaction_amount"`
	TransactionReference     string                    `json:"transaction_reference"` // this could be a check number, batch ,etc
	IsReconciled             bool                      `json:"is_reconciled"`
	IsSplit                  bool                      `json:"is_split"`
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
		TransactionReconcileDate: trans.TransactionReconcileDate,
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
