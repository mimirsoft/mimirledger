package request

import (
	"database/sql"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
	"time"
)

// Account is for use in accounts controller responses
type Account struct {
	AccountID            uint64                `json:"accountID"`
	AccountParent        uint64                `json:"accountParent"`
	AccountName          string                `json:"accountName"`
	AccountFullName      string                `json:"accountFullname"`
	AccountMemo          string                `json:"accountMemo"`
	AccountCurrent       bool                  `json:"accountCurrent"`
	AccountBalance       uint64                `json:"accountBalance"`
	AccountSubtotal      uint64                `json:"accountSubtotal"`
	AccountDecimals      uint64                `json:"accountDecimals"`
	AccountReconcileDate sql.NullTime          `json:"accountReconcileDate"`
	AccountFlagged       bool                  `json:"accountFlagged"`
	AccountLocked        bool                  `json:"accountLocked"`
	AccountOpenDate      time.Time             `json:"accountOpenDate"`
	AccountCloseDate     sql.NullTime          `json:"accountCloseDate"`
	AccountCode          sql.NullString        `json:"accountCode"`
	AccountSign          datastore.AccountSign `json:"accountSign"`
	AccountType          datastore.AccountType `json:"accountType"`
}

func ReqAccountToAccount(act *Account) *models.Account {
	return &models.Account{
		AccountID:            act.AccountID,
		AccountParent:        act.AccountParent,
		AccountName:          act.AccountName,
		AccountFullName:      act.AccountFullName,
		AccountMemo:          act.AccountMemo,
		AccountCurrent:       act.AccountCurrent,
		AccountBalance:       act.AccountBalance,
		AccountSubtotal:      act.AccountSubtotal,
		AccountDecimals:      act.AccountDecimals,
		AccountReconcileDate: act.AccountReconcileDate,
		AccountFlagged:       act.AccountFlagged,
		AccountLocked:        act.AccountLocked,
		AccountOpenDate:      act.AccountOpenDate,
		AccountCloseDate:     act.AccountCloseDate,
		AccountCode:          act.AccountCode,
		AccountSign:          act.AccountSign,
		AccountType:          act.AccountType,
	}
}
