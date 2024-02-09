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

// UserNotificationStatus is an enum for UserNotification statuses
type AccountSign string

// UserNotificationType is an enum for UserNotification type
type AccountType string

type Account struct {
	AccountID            uint64          `db:"account_id,omitempty"`
	AccountParent        uint64          `db:"account_parent"`
	AccountName          string          `db:"account_name"`
	AccountFullName      string          `db:"account_full_name"`
	AccountMemo          string          `db:"account_memo"`
	AccountCurrent       bool            `db:"account_current"`
	AccountLeft          uint64          `db:"account_left"`
	AccountRight         uint64          `db:"account_right"`
	AccountBalance       sql.NullFloat64 `db:"account_balance"`
	AccountSubtotal      sql.NullFloat64 `db:"account_subtotal"`
	AccountReconcileDate sql.NullTime    `db:"account_reconcile_date"`
	AccountFlagged       bool            `db:"account_flagged"`
	AccountLocked        bool            `db:"account_locked"`
	AccountOpenDate      time.Time       `db:"account_open_date"`
	AccountCloseDate     sql.NullTime    `db:"account_close_date"`
	AccountCode          sql.NullString  `db:"account_code"`
	AccountSign          AccountSign     `db:"account_sign"`
	AccountType          AccountType     `db:"account_type"`
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
