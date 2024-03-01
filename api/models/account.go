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
	AccountBalance       uint64
	AccountSubtotal      uint64
	AccountDecimals      uint64
	AccountReconcileDate sql.NullTime
	AccountFlagged       bool
	AccountLocked        bool
	AccountOpenDate      time.Time
	AccountCloseDate     sql.NullTime
	AccountCode          sql.NullString
	AccountSign          datastore.AccountSign
	AccountType          datastore.AccountType
}

var errParentAccountNotFound = errors.New("cannot find parent account with ID")

// Store inserts a URLComment
func (c *Account) Store(ds *datastore.Datastores) (err error) { //nolint:gocyclo
	//check if this is top level.  if it is not, the type must be the parent type
	var parentAccount *Account
	if c.AccountParent != 0 {
		parentAccount, err := getAccountByID(ds, c.AccountParent)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("%w, [parent AccountID:%d]", errParentAccountNotFound, c.AccountParent)
			}
			return fmt.Errorf("getAccountByID:%w", err)
		}

		c.AccountType = parentAccount.AccountType
		c.AccountSign = parentAccount.AccountSign
	}
	//Find the new spot in the tree.
	afterLeft, afterRight, err := findSpotInTree(ds, parentAccount, c.AccountName)
	if err != nil {
		return fmt.Errorf("getAccountByID:%w", err)
	}
	err = openSpotInTree(ds, afterLeft, afterRight)
	if err != nil {
		return fmt.Errorf("openSpotInTree:%w", err)
	}
	c.AccountLeft = afterRight + 1
	c.AccountRight = afterRight + 2
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

func findSpotInTree(ds *datastore.Datastores, parent *Account, name string) (uint64, uint64, error) {
	var parentAccountID uint64 = 0
	if parent != nil {
		parentAccountID = parent.AccountID
	}
	var afterThisAccount *Account
	children, err := findDirectChildren(ds, parentAccountID)
	if err != nil {
		return 0, 0, fmt.Errorf("findDirectChildren:%w", err)
	}
	for idx := range children {
		// > means earlier in alphabet
		if children[idx].AccountName < name {
			afterThisAccount = children[idx]
		}
	}
	if afterThisAccount == nil && parent != nil {
		/*		 If this is the first child of the parent, we want to spread the right from the left.
		we do this by setting the right equal to the left, effectively making it one lower
		and then running those values through the openspotintree function.

		*/
		return parent.AccountLeft, parent.AccountLeft, nil
	}
	// if this is the first account
	if afterThisAccount == nil {
		return 0, 0, nil
	}
	return afterThisAccount.AccountLeft, afterThisAccount.AccountRight, nil
}
func openSpotInTree(ds *datastore.Datastores, left uint64, right uint64) error {
	as := ds.AccountStore()
	err := as.OpenSpotInTree(left, right)
	if err != nil {
		return fmt.Errorf("as.OpenSpotInTree:%w", err)
	}
	return nil
}

// getAccountByID retrieves a specificAcocuunt
func getAccountByID(ds *datastore.Datastores, accountID uint64) (*Account, error) {
	as := ds.AccountStore()
	eAccount, err := as.GetAccountByID(accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("AccountStore().GetAccountByID:%w", err)
	}
	myAccount := Account(*eAccount)
	return &myAccount, nil
}

func findDirectChildren(ds *datastore.Datastores, accountID uint64) ([]*Account, error) {
	as := ds.AccountStore()
	actSet, err := as.GetDirectChildren(accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("AccountStore().GetDirectChildren:%w", err)
	}
	accounts := entAccountToAccounts(actSet)
	return accounts, nil
}

func entAccountToAccounts(eAccts []datastore.Account) (ua []*Account) {
	ua = make([]*Account, len(eAccts))
	for idx := range eAccts {
		act := Account(eAccts[idx])
		ua[idx] = &act
	}
	return
}
