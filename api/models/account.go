package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"strings"
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
	AccountBalance       int64
	AccountSubtotal      int64
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
var errAccountTypeInvalid = errors.New("accountType is not valid, cannot determine AccountSign")
var ErrAccountNotFound = errors.New("account not found")

const spreadForOneAccount = uint64(2)

func (c *Account) Store(dStores *datastore.Datastores) error { //nolint:gocyclo
	//check if this is top level.  if it is not, the type must be the parent type
	if c.AccountName == "" {
		return errAccountNameEmptyString
	}

	if c.AccountParent != 0 {
		parentAccount, err := RetrieveAccountByID(dStores, c.AccountParent)
		if err != nil {
			return fmt.Errorf("getAccountByID:%w", err)
		}

		c.AccountType = parentAccount.AccountType
		c.AccountSign = parentAccount.AccountSign
	}
	// fill in sign from AccountType
	var (
		ok          bool
		accountSign datastore.AccountSign
	)

	if accountSign, ok = datastore.AccountTypeToSign[c.AccountType]; !ok {
		return fmt.Errorf("%w c.AccountType:%s:", errAccountTypeInvalid, c.AccountType)
	}

	c.AccountSign = accountSign

	//Find the new spot in the tree.
	afterValue, err := findSpotInTree(dStores, c.AccountParent, c.AccountName)
	if err != nil {
		return fmt.Errorf("findSpotInTree:%w", err)
	}

	err = openSpotInTree(dStores, afterValue, spreadForOneAccount)
	if err != nil {
		return fmt.Errorf("openSpotInTree:%w", err)
	}
	// a single account with not children has a right and left that are sequential numbers
	c.AccountLeft = afterValue + 1
	c.AccountRight = afterValue + 2 //nolint:mnd
	eAcct := datastore.Account(*c)

	err = dStores.AccountStore().Store(&eAcct)
	if err != nil {
		return fmt.Errorf("ds.AccountStore().Store:%w [account:%+v]", err, eAcct)
	}

	fullName, err := retrieveAccountFullName(dStores, eAcct.AccountID)
	if err != nil {
		return fmt.Errorf("retrieveAccountFullName:%w [account:%+v]", err, eAcct)
	}

	eAcct.AccountFullName = fullName

	err = dStores.AccountStore().Update(&eAcct)
	if err != nil {
		return fmt.Errorf("ds.AccountStore().Update:%w", err)
	}

	*c = Account(eAcct)

	return nil
}

// This whole function should be a transaction for safety
func (c *Account) Update(dStores *datastore.Datastores) (err error) { //nolint:gocyclo
	// get existing Account record
	acctB4Update, err := RetrieveAccountByID(dStores, c.AccountID)
	if err != nil {
		return fmt.Errorf("getAccountByID: %w [accountID: %d ", err, c.AccountID)
	}

	oldAccountLeft := acctB4Update.AccountLeft
	affectedBalanceAccountIDs := make(map[uint64]bool)
	// if we have a new parent, we must close the old spot in the tree and open a new one
	// before updating
	if acctB4Update.AccountParent != c.AccountParent {
		// if these are not equal, then the account tree is changed and will need to be rebalanced
		oldParentIDs, err := getParentsAccountIDs(dStores, acctB4Update.AccountID)
		if err != nil {
			return fmt.Errorf("getParentsAccountIDs:%w", err)
		}
		//check if this is top level.  if it is not, the type must be the parent type
		var parentAccount *Account
		if c.AccountParent != 0 {
			parentAccount, err = RetrieveAccountByID(dStores, c.AccountParent)
			if err != nil {
				return fmt.Errorf("getAccountByID:%w", err)
			}

			c.AccountType = parentAccount.AccountType
			c.AccountSign = parentAccount.AccountSign
		}

		for idx := range oldParentIDs {
			affectedBalanceAccountIDs[oldParentIDs[idx]] = true
		}
	}

	var (
		ok          bool
		accountSign datastore.AccountSign
	)
	// fill in sign from AccountType
	if accountSign, ok = datastore.AccountTypeToSign[c.AccountType]; !ok {
		return fmt.Errorf("%w c.AccountType:%s:", errAccountTypeInvalid, c.AccountType)
	}

	c.AccountSign = accountSign
	//Find all children
	children, err := findAllChildren(dStores, c.AccountID)
	if err != nil {
		return fmt.Errorf("findAllChildren: %w [ AccountID:%d]", err, c.AccountID)
	}
	// if this account has no children, spread is 2
	// we calculate it this way so that any error in accountRight is fix
	spread := uint64(uint64(len(children))*spreadForOneAccount) + spreadForOneAccount
	// close old spot in tree
	err = closeSpotInTree(dStores, acctB4Update.AccountRight, spread)
	if err != nil {
		return fmt.Errorf("closeSpotInTree:%w", err)
	}
	//Find the new spot in the tree. after value is the value which this account's AccountLeft should be
	afterValue, err := findSpotInTree(dStores, c.AccountParent, c.AccountName)
	if err != nil {
		return fmt.Errorf("findSpotInTree:%w", err)
	}

	err = openSpotInTree(dStores, afterValue, spread)
	if err != nil {
		return fmt.Errorf("openSpotInTree:%w", err)
	}

	c.AccountLeft = afterValue + 1
	c.AccountRight = afterValue + spread
	eAcct := datastore.Account(*c)

	err = dStores.AccountStore().Update(&eAcct)
	if err != nil {
		return fmt.Errorf("ds.AccountStore().Update:%w", err)
	}
	//Update all the children
	shift := c.AccountLeft - oldAccountLeft

	for idx := range children {
		child := children[idx]
		child.AccountLeft += shift
		child.AccountRight += shift
		childEAcct := datastore.Account(*child)

		err = dStores.AccountStore().Update(&childEAcct)
		if err != nil {
			return fmt.Errorf("accountStore().Update:%w", err)
		}
		// retrieve name after updating left and right
		fullName, err := retrieveAccountFullName(dStores, children[idx].AccountID)
		if err != nil {
			return fmt.Errorf("retrieveAccountFullName:%w [account:%+v]", err, eAcct)
		}
		// update the record again, with the full name
		childEAcct.AccountFullName = fullName

		err = dStores.AccountStore().Update(&childEAcct)
		if err != nil {
			return fmt.Errorf("accountStore().Update:%w", err)
		}
	}

	newParentIDs, err := getParentsAccountIDs(dStores, c.AccountID)
	if err != nil {
		return fmt.Errorf("getParentsAccountIDs:%w", err)
	}

	for idx := range newParentIDs {
		affectedBalanceAccountIDs[newParentIDs[idx]] = true
	}

	for idx := range affectedBalanceAccountIDs {
		err := UpdateBalanceForAccountID(dStores, idx)
		if err != nil {
			return fmt.Errorf("UpdateBalanceForAccountID:%w [accountID:%d]", err, idx)
		}
	}

	// now get the fullName, after the Update to all AccountLefts and Account Rights
	fullName, err := retrieveAccountFullName(dStores, eAcct.AccountID)
	if err != nil {
		return fmt.Errorf("retrieveAccountFullName:%w [account:%+v]", err, eAcct)
	}

	eAcct.AccountFullName = fullName

	err = dStores.AccountStore().Update(&eAcct)
	if err != nil {
		return fmt.Errorf("ds.AccountStore().Update:%w", err)
	}

	*c = Account(eAcct)

	return nil
}

// updateSubtotal updates the subtotal on account
func (c *Account) updateSubtotal(dStores *datastore.Datastores) error {
	tcdStore := dStores.TransactionDebitCreditStore()

	subtotals, err := tcdStore.GetSubtotals(c.AccountID)
	if err != nil {
		return fmt.Errorf("tcdStore.GetSubtotals:%w", err)
	}
	var debitSubtotal int64 = 0
	var creditSubtotal int64 = 0

	for idx := range subtotals {
		switch accountSign := subtotals[idx].DebitOrCredit; accountSign {
		case datastore.AccountSignDebit:
			debitSubtotal = int64(subtotals[idx].Subtotal)
		case datastore.AccountSignCredit:
			creditSubtotal = int64(subtotals[idx].Subtotal)
		}
	}

	switch c.AccountSign {
	case datastore.AccountSignDebit:
		c.AccountSubtotal = debitSubtotal - creditSubtotal
	case datastore.AccountSignCredit:
		c.AccountSubtotal = creditSubtotal - debitSubtotal
	}

	eAcct := datastore.Account(*c)

	err = dStores.AccountStore().UpdateSubtotal(&eAcct)
	if err != nil {
		return fmt.Errorf("ds.AccountStore().UpdateSubtotal:%w", err)
	}

	*c = Account(eAcct)

	return nil
}

// updateBalance retrieves a specificAccount
func (c *Account) updateBalance(dStores *datastore.Datastores) error {
	eAcct := datastore.Account(*c)

	balance, err := dStores.AccountStore().GetBalance(c.AccountID)
	if err != nil {
		return fmt.Errorf("ds.AccountStore().GetBalance:%w", err)
	}

	eAcct.AccountBalance = balance

	err = dStores.AccountStore().UpdateBalance(&eAcct)
	if err != nil {
		return fmt.Errorf("ds.AccountStore().UpdateBalance:%w", err)
	}

	*c = Account(eAcct)

	return nil
}

// updateBalance retrieves a specificAccount
func (c *Account) UpdateReconciledDate(dStores *datastore.Datastores) error {
	eAcct := datastore.Account(*c)

	err := dStores.AccountStore().SetAccountReconciledDate(&eAcct)
	if err != nil {
		return fmt.Errorf("ds.AccountStore().GetBalance:%w", err)
	}

	*c = Account(eAcct)

	return nil
}

// UpdateBalanceForAccountID
func UpdateBalanceForAccountID(dStores *datastore.Datastores, accountID uint64) error {
	myAcct, err := RetrieveAccountByID(dStores, accountID)
	if err != nil {
		return fmt.Errorf("RetrieveAccountByID:%w", err)
	}

	err = myAcct.updateBalance(dStores)
	if err != nil {
		return fmt.Errorf("myAcct.updateBalance:%w", err)
	}

	return nil
}

// UpdateSubtotalForAccountID
func UpdateSubtotalForAccountID(dStores *datastore.Datastores, accountID uint64) error {
	myAcct, err := RetrieveAccountByID(dStores, accountID)
	if err != nil {
		return fmt.Errorf("RetrieveAccountByID:%w", err)
	}

	err = myAcct.updateSubtotal(dStores)
	if err != nil {
		return fmt.Errorf("myAcct.updateSubtotal:%w", err)
	}

	return nil
}

// RetrieveAccounts retrieves accounts
func RetrieveAccounts(dStores *datastore.Datastores) ([]*Account, error) {
	as := dStores.AccountStore()
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

// RetrieveAccountByID retrieves a specific account
func RetrieveAccountByID(dStores *datastore.Datastores, accountID uint64) (*Account, error) {
	as := dStores.AccountStore()

	eAcct, err := as.GetAccountByID(accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}

		return nil, fmt.Errorf("AccountStore().GetAccountByID:%w", err)
	}

	acct := Account(*eAcct)

	return &acct, nil
}

func findSpotInTree(dStores *datastore.Datastores, parentAccountID uint64, name string) (uint64, error) {
	var afterThisAccount *Account

	children, err := findDirectChildren(dStores, parentAccountID)
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
		parent, err := RetrieveAccountByID(dStores, parentAccountID)
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

func closeSpotInTree(dStores *datastore.Datastores, afterValue uint64, spread uint64) error {
	as := dStores.AccountStore()

	err := as.CloseSpotInTree(afterValue, spread)
	if err != nil {
		return fmt.Errorf("as.OpenSpotInTree:%w", err)
	}

	return nil
}

func findDirectChildren(dStores *datastore.Datastores, accountID uint64) ([]*Account, error) {
	as := dStores.AccountStore()

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
func findAllChildren(dStores *datastore.Datastores, parentID uint64) ([]*Account, error) {
	as := dStores.AccountStore()

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

func entAccountToAccounts(eAccts []datastore.Account) []*Account {
	accountSet := make([]*Account, len(eAccts))

	for idx := range eAccts {
		act := Account(eAccts[idx])
		accountSet[idx] = &act
	}

	return accountSet
}
func retrieveAccountFullName(dStores *datastore.Datastores, accountID uint64) (string, error) {
	accountFullName := ""
	as := dStores.AccountStore()

	actSet, err := as.GetParents(accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return accountFullName, nil
		}

		return "", fmt.Errorf("AccountStore().GetParents:%w", err)
	}

	for idx := range actSet {
		accountFullName += actSet[idx].AccountName + ":"
	}

	accountFullName = strings.TrimSuffix(accountFullName, ":")

	return accountFullName, nil
}

func getParentsAccountIDs(dStores *datastore.Datastores, accountID uint64) ([]uint64, error) {
	as := dStores.AccountStore()

	actSet, err := as.GetParents(accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("AccountStore().GetParents:%w", err)
	}
	var parentIDs []uint64

	for idx := range actSet {
		parentIDs = append(parentIDs, actSet[idx].AccountID)
	}

	return parentIDs, nil
}
