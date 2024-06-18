package response

import (
	"database/sql"
	"time"

	"github.com/mimirsoft/mimirledger/api/models"
)

// AccountSet is for use in accounts controller responses
type AccountSet struct {
	Accounts []*Account `json:"accounts"`
}

// Account is for use in accounts controller responses
type Account struct {
	AccountID            uint64         `json:"accountID"`
	AccountParent        uint64         `json:"accountParent"`
	AccountName          string         `json:"accountName"`
	AccountFullName      string         `json:"accountFullName"`
	AccountMemo          string         `json:"accountMemo"`
	AccountCurrent       bool           `json:"accountCurrent"`
	AccountLeft          uint64         `json:"accountLeft"`
	AccountRight         uint64         `json:"accountRight"`
	AccountBalance       int64          `json:"accountBalance"`
	AccountSubtotal      int64          `json:"accountSubtotal"`
	AccountDecimals      uint64         `json:"accountDecimals"`
	AccountReconcileDate sql.NullTime   `json:"accountReconcileDate"`
	AccountFlagged       bool           `json:"accountFlagged"`
	AccountLocked        bool           `json:"accountLocked"`
	AccountOpenDate      time.Time      `json:"accountOpenDate"`
	AccountCloseDate     sql.NullTime   `json:"accountCloseDate"`
	AccountCode          sql.NullString `json:"accountCode"`
	AccountSign          string         `json:"accountSign"`
	AccountType          string         `json:"accountType"`
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
