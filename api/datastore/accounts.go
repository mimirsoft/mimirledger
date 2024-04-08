package datastore

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type AccountStore struct {
	Client *sqlx.DB
}

const (
	// AccountSignDebit is the AccountSign status for DEBIT Accounts
	AccountSignDebit = AccountSign("DEBIT")
	// AccountSignCredit is the AccountSign status for CREDIT Accounts
	AccountSignCredit = AccountSign("CREDIT")

	AccountTypeAsset     = AccountType("ASSET")
	AccountTypeLiability = AccountType("LIABILITY")
	AccountTypeEquity    = AccountType("EQUITY")
	AccountTypeIncome    = AccountType("INCOME")
	AccountTypeExpense   = AccountType("EXPENSE")
	AccountTypeGain      = AccountType("GAIN")
	AccountTypeLoss      = AccountType("LOSS")
)

var AccountTypeToSign = map[AccountType]AccountSign{
	AccountTypeAsset:     AccountSignDebit,
	AccountTypeLiability: AccountSignCredit,
	AccountTypeEquity:    AccountSignCredit,
	AccountTypeIncome:    AccountSignDebit,
	AccountTypeExpense:   AccountSignCredit,
	AccountTypeGain:      AccountSignCredit,
	AccountTypeLoss:      AccountSignDebit}

// UserNotificationStatus is an enum for UserNotification statuses
type AccountSign string

// UserNotificationType is an enum for UserNotification type
type AccountType string

type Account struct {
	AccountID            uint64         `db:"account_id,omitempty"`
	AccountParent        uint64         `db:"account_parent"`
	AccountName          string         `db:"account_name"`
	AccountFullName      string         `db:"account_full_name"`
	AccountMemo          string         `db:"account_memo"`
	AccountCurrent       bool           `db:"account_current"`
	AccountLeft          uint64         `db:"account_left"`
	AccountRight         uint64         `db:"account_right"`
	AccountBalance       uint64         `db:"account_balance"`
	AccountSubtotal      uint64         `db:"account_subtotal"`
	AccountDecimals      uint64         `db:"account_decimals"`
	AccountReconcileDate sql.NullTime   `db:"account_reconcile_date"`
	AccountFlagged       bool           `db:"account_flagged"`
	AccountLocked        bool           `db:"account_locked"`
	AccountOpenDate      time.Time      `db:"account_open_date"`
	AccountCloseDate     sql.NullTime   `db:"account_close_date"`
	AccountCode          sql.NullString `db:"account_code"`
	AccountSign          AccountSign    `db:"account_sign"`
	AccountType          AccountType    `db:"account_type"`
}

// Store inserts a UserNotification into postgres
func (store AccountStore) Store(acct *Account) (err error) {
	query := `    INSERT INTO transaction_accounts 
		           (account_parent,
	account_name,
	account_full_name,
	account_memo,
	account_current,
	account_left,
	account_right,
	account_balance,
	account_subtotal,
	account_reconcile_date,
	account_flagged,
	account_locked,
	account_open_date,
	account_close_date,
	account_code,
	account_sign,
	account_type)
		    VALUES (:account_parent,
	:account_name,
	:account_full_name,
	:account_memo,
	:account_current,
	:account_left,
	:account_right,
	:account_balance,
	:account_subtotal,
	:account_reconcile_date,
	:account_flagged,
	:account_locked,
	:account_open_date,
	:account_close_date,
	:account_code,
	:account_sign,
	:account_type)
		 RETURNING *`
	stmt, err := store.Client.PrepareNamed(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	return stmt.QueryRow(acct).StructScan(acct)
}

// Store inserts a UserNotification into postgres
func (store AccountStore) Update(acct *Account) (err error) {
	query := `    UPDATE  transaction_accounts 
		    SET       (account_parent,
	account_name,
	account_full_name,
	account_memo,
	account_current,
	account_left,
	account_right,
	account_balance,
	account_subtotal,
	account_reconcile_date,
	account_flagged,
	account_locked,
	account_open_date,
	account_close_date,
	account_code,
	account_sign,
	account_type) = (:account_parent,
	:account_name,
	:account_full_name,
	:account_memo,
	:account_current,
	:account_left,
	:account_right,
	:account_balance,
	:account_subtotal,
	:account_reconcile_date,
	:account_flagged,
	:account_locked,
	:account_open_date,
	:account_close_date,
	:account_code,
	:account_sign,
	:account_type)
		    WHERE account_id = :account_id
		 RETURNING *`
	stmt, err := store.Client.PrepareNamed(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	return stmt.QueryRow(acct).StructScan(acct)
}

// Gets All Accounts
func (store AccountStore) GetAccounts() (as []Account, err error) {
	query := `select * from transaction_accounts order by account_left`
	rows, err := store.Client.Queryx(query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var acct Account
		if err = rows.StructScan(&acct); err != nil {
			return
		}
		as = append(as, acct)
	}
	if len(as) == 0 {
		return nil, sql.ErrNoRows
	}
	return
}

// Gets one account by account ID
func (store AccountStore) GetAccountByID(id uint64) (*Account, error) {
	query := `select * from transaction_accounts where account_id = $1`
	row := store.Client.QueryRowx(query, id)
	var as Account
	if err := row.StructScan(&as); err != nil {
		return nil, err
	}
	return &as, nil
}

// GetDirectChildren gets first level children of an account
func (store AccountStore) GetDirectChildren(id uint64) (as []Account, err error) {
	query := `select * from transaction_accounts WHERE account_parent = $1
	   ORDER BY account_name `
	rows, err := store.Client.Queryx(query, id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var acct Account
		if err = rows.StructScan(&acct); err != nil {
			return
		}
		as = append(as, acct)
	}
	if len(as) == 0 {
		return nil, sql.ErrNoRows
	}
	return
}

// GetDirectChildren gets first level children of an account
func (store AccountStore) GetParents(id uint64) (as []Account, err error) {
	query := `select Parents.* 
         FROM transaction_accounts AS Parents, transaction_accounts AS BaseAccount
WHERE (BaseAccount.account_left BETWEEN Parents.account_left AND Parents.account_right)
AND (BaseAccount.account_id =$1)
ORDER BY Parents.account_left`
	rows, err := store.Client.Queryx(query, id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var acct Account
		if err = rows.StructScan(&acct); err != nil {
			return
		}
		as = append(as, acct)
	}
	if len(as) == 0 {
		return nil, sql.ErrNoRows
	}
	return
}

// Gets All Children of a given account, regardless of dept
func (store AccountStore) GetAllChildren(id uint64) (as []Account, err error) {
	query := `SELECT 
children.*
FROM transaction_accounts AS parents,
transaction_accounts AS children
WHERE children.account_left BETWEEN parents.account_left AND parents.account_right
AND children.account_left <> parents.account_left
AND parents.account_id=$1
ORDER BY account_left`
	rows, err := store.Client.Queryx(query, id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var acct Account
		if err = rows.StructScan(&acct); err != nil {
			return
		}
		as = append(as, acct)
	}
	if len(as) == 0 {
		return nil, sql.ErrNoRows
	}
	return
}

// OpenSpotInTree opens a spot in our nested set
func (store AccountStore) OpenSpotInTree(afterValue, spread uint64) error {
	query := `UPDATE transaction_accounts
	SET account_right=account_right+$2
	WHERE account_right > $1`
	_, err := store.Client.Exec(query, afterValue, spread)
	if err != nil {
		return err
	}
	query = `UPDATE transaction_accounts
	SET account_left=account_left+$2
	WHERE account_left > $1`
	_, err = store.Client.Exec(query, afterValue, spread)
	if err != nil {
		return nil
	}
	return nil
}

// CloseSpotInTree closes a gap in our account tree
func (store AccountStore) CloseSpotInTree(afterValue, spread uint64) error {
	query := `UPDATE transaction_accounts
	SET account_right=account_right-$2
	WHERE account_right > $1`
	_, err := store.Client.Exec(query, afterValue, spread)
	if err != nil {
		return err
	}
	query = `UPDATE transaction_accounts
	SET account_left=account_left-$2
	WHERE account_left > $1`
	_, err = store.Client.Exec(query, afterValue, spread)
	if err != nil {
		return nil
	}
	return nil
}
