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
var errAccountNameEmptyString = errors.New("account name cannot be empty")

// Store inserts a URLComment
func (c *Account) Store(ds *datastore.Datastores) error { //nolint:gocyclo
	//check if this is top level.  if it is not, the type must be the parent type
	if c.AccountName == "" {
		return errAccountNameEmptyString
	}
	var parentAccount *Account
	var err error
	if c.AccountParent != 0 {
		parentAccount, err = getAccountByID(ds, c.AccountParent)
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
	afterValue, err := findSpotInTree(ds, c.AccountParent, c.AccountName)
	if err != nil {
		return fmt.Errorf("findSpotInTree:%w", err)
	}
	err = openSpotInTree(ds, afterValue, 2)
	if err != nil {
		return fmt.Errorf("openSpotInTree:%w", err)
	}
	c.AccountLeft = afterValue + 1
	c.AccountRight = afterValue + 2
	eAcct := datastore.Account(*c)
	err = ds.AccountStore().Store(&eAcct)
	if err != nil {
		return fmt.Errorf("ds.AccountStore().Store:%w", err)
	}
	*c = Account(eAcct)
	return nil
}

// This whole function should be a transaction for safety
func (c *Account) Update(ds *datastore.Datastores) (err error) { //nolint:gocyclo
	// get existing Account record
	acctB4Update, err := getAccountByID(ds, c.AccountID)
	if err != nil {
		return fmt.Errorf("getAccountByID: %w [accountID: %d ", err, c.AccountID)
	}
	// if we have a new parent, we must close the old spot in the tree and open a new one
	// before updating
	if acctB4Update.AccountParent != c.AccountParent {
		//check if this is top level.  if it is not, the type must be the parent type
		var parentAccount *Account
		if c.AccountParent != 0 {
			parentAccount, err = getAccountByID(ds, c.AccountParent)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("%w, [parent AccountID:%d]", errParentAccountNotFound, c.AccountParent)
				}
				return fmt.Errorf("getAccountByID:%w", err)
			}

			c.AccountType = parentAccount.AccountType
			c.AccountSign = parentAccount.AccountSign
		}
		//Find all children
		children, err := findAllChildren(ds, c.AccountID)
		if err != nil {
			return fmt.Errorf("findAllChildren: %w [ AccountID:%d]", err, c.AccountID)
		}
		// close old spot in tree
		spread := acctB4Update.AccountRight - acctB4Update.AccountLeft + 1
		err = closeSpotInTree(ds, acctB4Update.AccountRight, spread)
		if err != nil {
			return fmt.Errorf("closeSpotInTree:%w", err)
		}
		//Find the new spot in the tree.
		afterValue, err := findSpotInTree(ds, c.AccountParent, c.AccountName)
		if err != nil {
			return fmt.Errorf("findSpotInTree:%w", err)
		}
		err = openSpotInTree(ds, afterValue, spread)
		if err != nil {
			return fmt.Errorf("openSpotInTree:%w", err)
		}
		c.AccountLeft = afterValue + 1
		c.AccountRight = afterValue + spread
		//Update all the children
		shift := c.AccountLeft - acctB4Update.AccountLeft
		for idx := range children {
			child := children[idx]
			child.AccountLeft += shift
			child.AccountRight += shift
			eAcct := datastore.Account(*child)
			err := ds.AccountStore().Update(&eAcct)
			if err != nil {
				return fmt.Errorf("accountStore().Update:%w", err)
			}
		}
	}
	eAcct := datastore.Account(*c)
	err = ds.AccountStore().Update(&eAcct)
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

func findSpotInTree(ds *datastore.Datastores, parentAccountID uint64, name string) (uint64, error) {
	var afterThisAccount *Account
	children, err := findDirectChildren(ds, parentAccountID)
	if err != nil {
		return 0, fmt.Errorf("findDirectChildren:%w", err)
	}
	for idx := range children {
		// > means earlier in alphabet
		if children[idx].AccountName < name {
			afterThisAccount = children[idx]
		}
	}
	if afterThisAccount == nil && parentAccountID != 0 {
		/*		 If this is the first child of the parent, we want to spread the right from the left.
		we do this by setting the right equal to the left, effectively making it one lower
		and then running those values through the openspotintree function.

		*/
		// we do the call here, because this parent's left and right may have changed due to a closeSPontInTree call
		parent, err := getAccountByID(ds, parentAccountID)
		if err != nil {
			return 0, fmt.Errorf("getAccountByID:%w", err)
		}

		return parent.AccountLeft, nil
	}
	// if this is the first account (no sibling account before this and no parent)
	if afterThisAccount == nil {
		return 0, nil
	}
	return afterThisAccount.AccountRight, nil
}
func openSpotInTree(ds *datastore.Datastores, afterValue uint64, spread uint64) error {
	as := ds.AccountStore()
	err := as.OpenSpotInTree(afterValue, spread)
	if err != nil {
		return fmt.Errorf("as.OpenSpotInTree:%w", err)
	}
	return nil
}

func closeSpotInTree(ds *datastore.Datastores, afterValue uint64, spread uint64) error {
	as := ds.AccountStore()
	err := as.CloseSpotInTree(afterValue, spread)
	if err != nil {
		return fmt.Errorf("as.OpenSpotInTree:%w", err)
	}
	return nil
}

// getAccountByID retrieves a specificAccount
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
func findAllChildren(ds *datastore.Datastores, parentID uint64) ([]*Account, error) {
	as := ds.AccountStore()
	actSet, err := as.GetAllChildren(parentID)
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
