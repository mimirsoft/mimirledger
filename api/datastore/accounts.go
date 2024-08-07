package datastore

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type AccountStore struct {
	Client *sqlx.DB
}

const (
	// AccountSignDebit is the AccountSign status for DEBIT Accounts.
	AccountSignDebit = AccountSign("DEBIT")
	// AccountSignCredit is the AccountSign status for CREDIT Accounts.
	AccountSignCredit = AccountSign("CREDIT")

	AccountTypeAsset     = AccountType("ASSET")
	AccountTypeLiability = AccountType("LIABILITY")
	AccountTypeEquity    = AccountType("EQUITY")
	AccountTypeIncome    = AccountType("INCOME")
	AccountTypeExpense   = AccountType("EXPENSE")
	AccountTypeGain      = AccountType("GAIN")
	AccountTypeLoss      = AccountType("LOSS")
)

var AccountTypeToSign = map[AccountType]AccountSign{ //nolint:gochecknoglobals
	AccountTypeAsset:     AccountSignDebit,
	AccountTypeLiability: AccountSignCredit,
	AccountTypeEquity:    AccountSignCredit,
	AccountTypeIncome:    AccountSignCredit,
	AccountTypeExpense:   AccountSignDebit,
	AccountTypeGain:      AccountSignCredit,
	AccountTypeLoss:      AccountSignDebit}

// AccountSign is an enum for account signs "DEBIT" or "CREDIT".
type AccountSign string

// AccountType is an enum for account type.
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
	AccountBalance       int64          `db:"account_balance"`
	AccountSubtotal      int64          `db:"account_subtotal"`
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

// Store inserts an Account into postgres.
func (store AccountStore) Store(acct *Account) error {
	query := `    INSERT INTO transaction_accounts 
		           (account_parent,
	account_name,
	account_full_name,
	account_memo,
	account_current,
	account_left,
	account_right,
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
		return fmt.Errorf(" store.Client.PrepareNamed(query):%w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(acct).StructScan(acct)
	if err != nil {
		return fmt.Errorf("stmt.QueryRow(acct).StructScan(acct):%w", err)
	}

	return nil
}

// Update  updates Accounts into postgres.
func (store AccountStore) Update(acct *Account) error {
	query := `    UPDATE  transaction_accounts 
		    SET       (account_parent,
	account_name,
	account_full_name,
	account_memo,
	account_current,
	account_left,
	account_right,
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
		return fmt.Errorf(" store.Client.PrepareNamed(query):%w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(acct).StructScan(acct)
	if err != nil {
		return fmt.Errorf("stmt.QueryRow(acct).StructScan(acct):%w", err)
	}

	return nil
}

// UpdateSubtotal  updates the account_subtotal into postgres.
func (store AccountStore) UpdateSubtotal(acct *Account) error {
	query := `    UPDATE  transaction_accounts 
		    SET  account_subtotal = :account_subtotal
		    WHERE account_id = :account_id
		 RETURNING *`

	stmt, err := store.Client.PrepareNamed(query)
	if err != nil {
		return fmt.Errorf(" store.Client.PrepareNamed:%w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(acct).StructScan(acct)
	if err != nil {
		return fmt.Errorf("stmt.QueryRow(acct).StructScan(acct):%w", err)
	}

	return nil
}

// account_balance  updates the account_balance into postgres.
func (store AccountStore) UpdateBalance(acct *Account) error {
	query := `    UPDATE  transaction_accounts 
		    SET  account_balance = :account_balance
		    WHERE account_id = :account_id
		 RETURNING *`

	stmt, err := store.Client.PrepareNamed(query)
	if err != nil {
		return fmt.Errorf(" store.Client.PrepareNamed:%w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(acct).StructScan(acct)
	if err != nil {
		return fmt.Errorf("stmt.QueryRow(acct).StructScan(acct):%w", err)
	}

	return nil
}

// Set AccountReconcileDate.
func (store AccountStore) SetAccountReconciledDate(acct *Account) error {
	query := `UPDATE  transaction_accounts 
		    SET account_reconcile_date = $2
		    where account_id = $1`

	_, err := store.Client.Exec(query, acct.AccountID, acct.AccountReconcileDate)
	if err != nil {
		return fmt.Errorf("store.Client.Exec:%w", err)
	}

	return nil
}

// GetBalancel  gets the sum of all the subtotals for this accountID and all child accounts.
func (store AccountStore) GetBalance(accountID uint64) (int64, error) {
	query := `    SELECT SUM(subaccount.account_subtotal) AS balance
                            FROM transaction_accounts AS p_account, transaction_accounts AS subaccount
                            WHERE subaccount.account_left BETWEEN p_account.account_left AND p_account.account_right
                            AND p_account.account_id =$1 `
	row := store.Client.QueryRowx(query, accountID)

	var accountBalance int64

	if err := row.Scan(&accountBalance); err != nil {
		return 0, fmt.Errorf("row.Scan(&accountBalance):%w", err)
	}

	return accountBalance, nil
}

// GetBalances  gets the sum of all the subtotals for this accountID and all child accounts.
func (store AccountStore) GetBalances(accountID uint64) ([]Account, error) {
	query := `    SELECT subaccount.*
                            FROM transaction_accounts AS p_account, transaction_accounts AS subaccount
                            WHERE subaccount.account_left BETWEEN p_account.account_left AND p_account.account_right
                            AND p_account.account_id =$1 `

	rows, err := store.Client.Queryx(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()

	var accountSet []Account

	for rows.Next() {
		var acct Account
		if err = rows.StructScan(&acct); err != nil {
			return nil, fmt.Errorf("rows.StructScan:%w", err)
		}

		accountSet = append(accountSet, acct)
	}

	if len(accountSet) == 0 {
		return nil, sql.ErrNoRows
	}

	return accountSet, nil
}

// Gets All Accounts.
func (store AccountStore) GetAccounts() ([]Account, error) {
	query := `select * from transaction_accounts order by account_left`

	rows, err := store.Client.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()

	var accountSet []Account

	for rows.Next() {
		var acct Account
		if err = rows.StructScan(&acct); err != nil {
			return nil, fmt.Errorf("rows.StructScan:%w", err)
		}

		accountSet = append(accountSet, acct)
	}

	if len(accountSet) == 0 {
		return nil, sql.ErrNoRows
	}

	return accountSet, nil
}

// Gets one account by account ID.
func (store AccountStore) GetAccountByID(id uint64) (*Account, error) {
	query := `select * from transaction_accounts where account_id = $1`
	row := store.Client.QueryRowx(query, id)

	var acct Account

	if err := row.StructScan(&acct); err != nil {
		return nil, fmt.Errorf("row.StructScan(&as):%w", err)
	}

	return &acct, nil
}

// GetDirectChildren gets first level children of an account.
func (store AccountStore) GetDirectChildren(acctID uint64) ([]Account, error) {
	query := `select * from transaction_accounts WHERE account_parent = $1
	   ORDER BY account_name `

	rows, err := store.Client.Queryx(query, acctID)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()

	var accountSet []Account

	for rows.Next() {
		var acct Account

		if err = rows.StructScan(&acct); err != nil {
			return nil, fmt.Errorf("rows.StructScan:%w", err)
		}

		accountSet = append(accountSet, acct)
	}

	if len(accountSet) == 0 {
		return nil, sql.ErrNoRows
	}

	return accountSet, nil
}

// GetDirectChildren gets first level children of an account.
func (store AccountStore) GetParents(acctID uint64) ([]Account, error) {
	query := `select Parents.* 
         FROM transaction_accounts AS Parents, transaction_accounts AS BaseAccount
WHERE (BaseAccount.account_left BETWEEN Parents.account_left AND Parents.account_right)
AND (BaseAccount.account_id =$1)
ORDER BY Parents.account_left`

	rows, err := store.Client.Queryx(query, acctID)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()

	var accountSet []Account

	for rows.Next() {
		var acct Account
		if err = rows.StructScan(&acct); err != nil {
			return nil, fmt.Errorf("rows.StructScan:%w", err)
		}

		accountSet = append(accountSet, acct)
	}

	if len(accountSet) == 0 {
		return nil, sql.ErrNoRows
	}

	return accountSet, nil
}

// Gets All Children of a given account, regardless of dept.
func (store AccountStore) GetAllChildren(acctID uint64) ([]Account, error) {
	query := `SELECT 
children.*
FROM transaction_accounts AS parents,
transaction_accounts AS children
WHERE children.account_left BETWEEN parents.account_left AND parents.account_right
AND children.account_left <> parents.account_left
AND parents.account_id=$1
ORDER BY account_left`

	rows, err := store.Client.Queryx(query, acctID)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()

	var accountSet []Account

	for rows.Next() {
		var acct Account
		if err = rows.StructScan(&acct); err != nil {
			return nil, fmt.Errorf("rows.StructScan:%w", err)
		}

		accountSet = append(accountSet, acct)
	}

	if len(accountSet) == 0 {
		return nil, sql.ErrNoRows
	}

	return accountSet, nil
}

type AccountWithLevel struct {
	Account
	Level int `db:"level"`
}

// Gets the Account and Children with Levels
func (store AccountStore) GetAccountWithChildrenByLevel(acctID uint64) ([]AccountWithLevel, error) {
	query := `
		WITH parent AS ( Select * from transaction_accounts where account_id = $1)
		SELECT *
          FROM parent,
               (SELECT A2.*,
                       (Count(A1.account_id))-1 AS level
                   FROM transaction_accounts AS A1,
                   transaction_accounts AS A2
                            WHERE A2.account_left BETWEEN A1.account_left AND A1.account_right
                            GROUP BY A2.account_id
                            ORDER BY A2.account_left)
            AS accounts
         WHERE accounts.account_left
       BETWEEN parent.account_left AND parent.account_right`

	rows, err := store.Client.Queryx(query, acctID)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()

	var accountSet []AccountWithLevel

	for rows.Next() {
		var acct AccountWithLevel
		if err = rows.StructScan(&acct); err != nil {
			return nil, fmt.Errorf("rows.StructScan:%w", err)
		}

		accountSet = append(accountSet, acct)
	}

	if len(accountSet) == 0 {
		return nil, sql.ErrNoRows
	}

	return accountSet, nil
}

// OpenSpotInTree opens a spot in our nested set.
func (store AccountStore) OpenSpotInTree(afterValue, spread uint64) error {
	query := `UPDATE transaction_accounts
	SET account_right=account_right+$2
	WHERE account_right > $1`

	_, err := store.Client.Exec(query, afterValue, spread)
	if err != nil {
		return fmt.Errorf("store.Client.Exec:%w", err)
	}

	query = `UPDATE transaction_accounts
	SET account_left=account_left+$2
	WHERE account_left > $1`

	_, err = store.Client.Exec(query, afterValue, spread)
	if err != nil {
		return fmt.Errorf("store.Client.Exec:%w", err)
	}

	return nil
}

// CloseSpotInTree closes a gap in our account tree.
func (store AccountStore) CloseSpotInTree(afterValue, spread uint64) error {
	query := `UPDATE transaction_accounts
	SET account_right=account_right-$2
	WHERE account_right > $1`

	_, err := store.Client.Exec(query, afterValue, spread)
	if err != nil {
		return fmt.Errorf("store.Client.Exec:%w", err)
	}

	query = `UPDATE transaction_accounts
	SET account_left=account_left-$2
	WHERE account_left > $1`

	_, err = store.Client.Exec(query, afterValue, spread)
	if err != nil {
		return fmt.Errorf("store.Client.Exec:%w", err)
	}

	return nil
}
