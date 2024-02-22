package models

import (
	"database/sql"
	"errors"
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

// Store inserts a URLComment
func (c *Account) Store(ds *datastore.Datastores) (err error) { //nolint:gocyclo
	eAcct := datastore.Account(*c)
	err = ds.AccountStore().Store(&eAcct)
	if err != nil {
		return fmt.Errorf("ds.AccountStore().Store:%w", err)
	}
	*c = Account(eAcct)
	return nil
}

// RetrieveAccounts retrieves accounts
func RetrieveAccounts(ds *datastore.Datastores) ([]*Account, error) {
	as := ds.AccountStore()
	actSet, err := as.GetAccounts()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("AccountStore().GetAccounts:%w", err)
	}
	accts := entAccountToAccounts(actSet)
	return accts, nil
}

func entAccountToAccounts(eAccts []datastore.Account) (ua []*Account) {
	ua = make([]*Account, len(eAccts))
	for idx := range eAccts {
		act := Account(eAccts[idx])
		ua[idx] = &act
	}
	return
}