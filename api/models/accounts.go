package models

import (
	"database/sql"
	"fmt"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"time"
)

type Account struct {
	AccountID            uint64
	AccountParent        uint64
	AccountName          string
	AccountFullName      string
	AccountMemo          string
	AccountCurrent       bool
	AccountLeft          uint64
	AccountRight         uint64
	AccountBalance       sql.NullFloat64
	AccountSubtotal      sql.NullFloat64
	AccountReconcileDate sql.NullTime
	AccountFlagged       bool
	AccountLocked        bool
	AccountOpenDate      time.Time
	AccountCloseDate     sql.NullTime
	AccountCode          sql.NullString
	AccountSign          datastore.AccountSign
	AccountType          datastore.AccountType
}

// RetrieveAccounts retrieves accounts
func RetrieveAccounts(ds *datastore.Datastores) ([]Account, error) {
	as := ds.AccountStore()
	actSet, err := as.GetAccounts()
	if err != nil {
		return nil, fmt.Errorf("AccountStore().GetAccounts:%w", err)
	}
	accts := entAccountToAccounts(actSet)
	return accts, nil
}

func entAccountToAccounts(eAccts []datastore.Account) (ua []Account) {
	ua = make([]Account, len(eAccts))
	for idx := range eAccts {
		ua[idx] = Account(eAccts[idx])
	}
	return
}
