package request

import (
	"database/sql"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
	"time"
)

// Account is for use in accounts controller responses
type Account struct {
	AccountID            uint64                `json:"account_id"`
	AccountParent        uint64                `json:"account_parent"`
	AccountName          string                `json:"account_name"`
	AccountFullName      string                `json:"account_fullname"`
	AccountMemo          string                `json:"account_memo"`
	AccountCurrent       bool                  `json:"account_current"`
	AccountBalance       uint64                `json:"account_balance"`
	AccountSubtotal      uint64                `json:"account_subtotal"`
	AccountDecimals      uint64                `json:"account_decimals"`
	AccountReconcileDate sql.NullTime          `json:"account_reconcile_date"`
	AccountFlagged       bool                  `json:"account_flagged"`
	AccountLocked        bool                  `json:"account_locked"`
	AccountOpenDate      time.Time             `json:"account_open_date"`
	AccountCloseDate     sql.NullTime          `json:"account_close_date"`
	AccountCode          sql.NullString        `json:"account_code"`
	AccountSign          datastore.AccountSign `json:"account_sign"`
	AccountType          datastore.AccountType `json:"account_type"`
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
