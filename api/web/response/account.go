package response

import (
	"database/sql"
	"github.com/mimirsoft/mimirledger/api/models"
	"time"
)

// AccountSet is for use in accounts controller responses
type AccountSet struct {
	Accounts []*Account `json:"accounts"`
}

// Account is for use in accounts controller responses
type Account struct {
	AccountID            uint64         `json:"account_id"`
	AccountParent        uint64         `json:"account_parent"`
	AccountName          string         `json:"account_name"`
	AccountFullName      string         `json:"account_fullname"`
	AccountMemo          string         `json:"account_memo"`
	AccountCurrent       bool           `json:"account_current"`
	AccountLeft          uint64         `json:"account_left"`
	AccountRight         uint64         `json:"account_right"`
	AccountBalance       uint64         `json:"account_balance"`
	AccountSubtotal      uint64         `json:"account_subtotal"`
	AccountDecimals      uint64         `json:"account_decimals"`
	AccountReconcileDate sql.NullTime   `json:"account_reconcile_date"`
	AccountFlagged       bool           `json:"account_flagged"`
	AccountLocked        bool           `json:"account_locked"`
	AccountOpenDate      time.Time      `json:"account_open_date"`
	AccountCloseDate     sql.NullTime   `json:"account_close_date"`
	AccountCode          sql.NullString `json:"account_code"`
	AccountSign          string         `json:"account_sign"`
	AccountType          string         `json:"account_type"`
}

// ConvertAccountsToRespAccounts converts []models.Account to AccountSet
func ConvertAccountsToRespAccountSet(accts []*models.Account) *AccountSet {
	var ras = make([]*Account, len(accts))
	for idx := range accts {
		ras[idx] = AccountToRespAccount(accts[idx])
	}
	return &AccountSet{Accounts: ras}
}

func AccountToRespAccount(act *models.Account) *Account {
	return &Account{
		AccountID:            act.AccountID,
		AccountParent:        act.AccountParent,
		AccountName:          act.AccountName,
		AccountFullName:      act.AccountFullName,
		AccountMemo:          act.AccountMemo,
		AccountCurrent:       act.AccountCurrent,
		AccountLeft:          act.AccountLeft,
		AccountRight:         act.AccountRight,
		AccountBalance:       act.AccountBalance,
		AccountSubtotal:      act.AccountSubtotal,
		AccountDecimals:      act.AccountDecimals,
		AccountReconcileDate: act.AccountReconcileDate,
		AccountFlagged:       act.AccountFlagged,
		AccountLocked:        act.AccountLocked,
		AccountOpenDate:      act.AccountOpenDate,
		AccountCloseDate:     act.AccountCloseDate,
		AccountCode:          act.AccountCode,
		AccountSign:          string(act.AccountSign),
		AccountType:          string(act.AccountType),
	}
}
