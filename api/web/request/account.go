package request

import (
	"database/sql"
	"time"

	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
)

// Account is for use in accounts controller responses
type Account struct {
	AccountID            uint64                `json:"accountID,omitempty"`
	AccountParent        uint64                `json:"accountParent"`
	AccountName          string                `json:"accountName"`
	AccountFullName      string                `json:"accountFullName"`
	AccountMemo          string                `json:"accountMemo"`
	AccountCurrent       bool                  `json:"accountCurrent"`
	AccountDecimals      uint64                `json:"accountDecimals"`
	AccountReconcileDate *time.Time            `json:"accountReconcileDate"`
	AccountFlagged       bool                  `json:"accountFlagged"`
	AccountLocked        bool                  `json:"accountLocked"`
	AccountOpenDate      time.Time             `json:"accountOpenDate"`
	AccountCloseDate     sql.NullTime          `json:"accountCloseDate"`
	AccountCode          sql.NullString        `json:"accountCode"`
	AccountType          datastore.AccountType `json:"accountType"`
}

func ReqAccountToAccount(act *Account) *models.Account {
	acctReconcileDate := sql.NullTime{Time: time.Time{}, Valid: false}
	if act.AccountReconcileDate != nil {
		acctReconcileDate = sql.NullTime{Time: *act.AccountReconcileDate, Valid: true}
	}

	return &models.Account{ //nolint:exhaustruct
		AccountID:            act.AccountID,
		AccountParent:        act.AccountParent,
		AccountName:          act.AccountName,
		AccountFullName:      act.AccountFullName,
		AccountMemo:          act.AccountMemo,
		AccountCurrent:       act.AccountCurrent,
		AccountDecimals:      act.AccountDecimals,
		AccountReconcileDate: acctReconcileDate,
		AccountFlagged:       act.AccountFlagged,
		AccountLocked:        act.AccountLocked,
		AccountOpenDate:      act.AccountOpenDate,
		AccountCloseDate:     act.AccountCloseDate,
		AccountCode:          act.AccountCode,
		AccountType:          act.AccountType,
	}
}
