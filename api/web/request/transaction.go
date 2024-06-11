package request

import (
	"database/sql"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
	"time"
)

type Transaction struct {
	TransactionID            uint64                    `json:"transactionID,omitempty"`
	TransactionDate          time.Time                 `json:"transactionDate"`
	TransactionReconcileDate *time.Time                `json:"transactionReconcileDate,omitempty"`
	TransactionComment       string                    `json:"transactionComment"`
	TransactionAmount        uint64                    `json:"transactionAmount"`
	TransactionReference     string                    `json:"transactionReference"` // this could be a check number, batch ,etc
	IsReconciled             bool                      `json:"isReconciled"`
	IsSplit                  bool                      `json:"isSplit"`
	DebitCreditSet           []*TransactionDebitCredit `json:"debitCreditSet"`
}

type TransactionDebitCredit struct {
	TransactionDCID     uint64                `json:"transactionDCID,omitempty"`
	TransactionID       uint64                `json:"transactionID"`
	AccountID           uint64                `json:"accountID"`
	TransactionDCAmount uint64                `json:"transactionDCAmount"`
	DebitOrCredit       datastore.AccountSign `json:"debitOrCredit"`
}

func ReqTransactionToTransaction(rTrans *Transaction) *models.Transaction {
	mytime := sql.NullTime{Time: time.Time{}, Valid: false}
	if rTrans.TransactionReconcileDate != nil {
		mytime = sql.NullTime{Time: *rTrans.TransactionReconcileDate, Valid: true}
	}
	myTransCore := models.TransactionCore{
		TransactionID:            rTrans.TransactionID,
		TransactionDate:          rTrans.TransactionDate,
		TransactionReconcileDate: mytime,
		TransactionComment:       rTrans.TransactionComment,
		TransactionAmount:        rTrans.TransactionAmount,
		TransactionReference:     rTrans.TransactionReference,
		IsReconciled:             rTrans.IsReconciled,
		IsSplit:                  rTrans.IsSplit,
	}
	myDCSet := ConvertReqDebitCreditsToDebitCreditSet(rTrans.DebitCreditSet)
	myTrans := models.Transaction{TransactionCore: myTransCore, DebitCreditSet: myDCSet}
	return &myTrans
}

// ConvertReqDebitCreditsToDebitCreditSet converts []TransactionDebitCreditt []*models.TransactionDebitCredit
func ConvertReqDebitCreditsToDebitCreditSet(dcSet []*TransactionDebitCredit) []*models.TransactionDebitCredit {
	var mset = make([]*models.TransactionDebitCredit, len(dcSet))
	for idx := range dcSet {
		myDS := models.TransactionDebitCredit(*dcSet[idx])
		mset[idx] = &myDS
	}
	return mset
}
