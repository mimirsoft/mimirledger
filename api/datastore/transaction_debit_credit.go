package datastore

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type TransactionDebitCreditStore struct {
	Client *sqlx.DB
}

type TransactionDebitCredit struct {
	TransactionDCID     uint64      `db:"transaction_dc_id,omitempty"`
	TransactionID       uint64      `db:"transaction_id"`
	AccountID           uint64      `db:"account_id"`
	TransactionDCAmount uint64      `db:"transaction_dc_amount"`
	DebitOrCredit       AccountSign `db:"debit_or_credit"`
}

// Store inserts a UserNotification into postgres
func (store TransactionDebitCreditStore) Store(trn *TransactionDebitCredit) error {
	query := `    INSERT INTO transaction_debit_credit 
		           (transaction_id,
	account_id,
	transaction_dc_amount,
	debit_or_credit)
		    VALUES (:transaction_id,
	:account_id,
	:transaction_dc_amount,
	:debit_or_credit)
		 RETURNING *`

	stmt, err := store.Client.PrepareNamed(query)
	if err != nil {
		return fmt.Errorf("error preparing transaction debit credit: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(trn).StructScan(trn)
	if err != nil {
		return fmt.Errorf("stmt.QueryRow(trn).StructScan(trn):%w", err)
	}

	return nil
}

func (store TransactionDebitCreditStore) GetDCForTransactionID(id uint64) ([]*TransactionDebitCredit, error) {
	query := `select * from transaction_debit_credit where transaction_id = $1`

	rows, err := store.Client.Queryx(query, id)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()

	var txnSet []*TransactionDebitCredit

	for rows.Next() {
		var txn TransactionDebitCredit
		if err = rows.StructScan(&txn); err != nil {
			return nil, fmt.Errorf("rows.StructScan:%w", err)
		}

		txnSet = append(txnSet, &txn)
	}

	if len(txnSet) == 0 {
		return nil, sql.ErrNoRows
	}

	return txnSet, nil
}

func (store TransactionDebitCreditStore) DeleteForTransactionID(id uint64) ([]*TransactionDebitCredit, error) {
	query := `DELETE from transaction_debit_credit where transaction_id = $1
	RETURNING *`

	rows, err := store.Client.Queryx(query, id)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()

	var txnSet []*TransactionDebitCredit

	for rows.Next() {
		var txn TransactionDebitCredit
		if err = rows.StructScan(&txn); err != nil {
			return nil, fmt.Errorf("rows.StructScan:%w", err)
		}

		txnSet = append(txnSet, &txn)
	}

	if len(txnSet) == 0 {
		return nil, sql.ErrNoRows
	}

	return txnSet, nil
}

type AccountSubtotal struct {
	Subtotal      uint64      `db:"subtotal"`
	DebitOrCredit AccountSign `db:"debit_or_credit"`
}

// Gets one account by account ID
func (store TransactionDebitCreditStore) GetSubtotals(accountID uint64) ([]*AccountSubtotal, error) {
	query := `select SUM(transaction_dc_amount) as subtotal,debit_or_credit
	from transaction_debit_credit 
	where account_id = $1
	GROUP BY  debit_or_credit`

	rows, err := store.Client.Queryx(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()

	var txnSet []*AccountSubtotal

	for rows.Next() {
		var txn AccountSubtotal
		if err = rows.StructScan(&txn); err != nil {
			return nil, fmt.Errorf("rows.StructScan:%w", err)
		}

		txnSet = append(txnSet, &txn)
	}

	return txnSet, nil
}

func (store TransactionDebitCreditStore) GetReconciledSubtotals(accountLeft, accountRight uint64,
	reconciledCutoffDate time.Time) ([]*AccountSubtotal, error) {
	query := `SELECT SUM(z.transaction_dc_amount) AS subtotal, z.debit_or_credit
					FROM (SELECT  workingtdc.transaction_dc_amount,
						  workingtdc.debit_or_credit
					FROM  transaction_debit_credit AS workingtdc
					LEFT JOIN transaction_debit_credit AS odc
					ON workingtdc.transaction_id=odc.transaction_id
					INNER JOIN transaction_main AS tm
					ON tm.transaction_id=workingtdc.transaction_id
					WHERE (workingtdc.account_id 
						IN (SELECT account_id FROM transaction_accounts WHERE account_left BETWEEN $2 AND $3)  
						AND odc.account_id 
						NOT IN (SELECT account_id FROM transaction_accounts WHERE account_left BETWEEN $2 AND $3) )
					  
					AND tm.is_reconciled IS TRUE 
					AND EXTRACT(EPOCH FROM tm.transaction_reconcile_date) <= EXTRACT(EPOCH FROM  $1::timestamp)
					  )
				 AS z    
	  GROUP BY z.debit_or_credit`

	rows, err := store.Client.Queryx(query, reconciledCutoffDate, accountLeft, accountRight)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()

	var txnSet []*AccountSubtotal

	for rows.Next() {
		var txn AccountSubtotal
		if err = rows.StructScan(&txn); err != nil {
			return nil, fmt.Errorf("rows.StructScan:%w", err)
		}

		txnSet = append(txnSet, &txn)
	}

	return txnSet, nil
}
