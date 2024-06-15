package datastore

import (
	"database/sql"
	"fmt"
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
func (store TransactionDebitCreditStore) Store(trn *TransactionDebitCredit) (err error) {
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
		return
	}
	defer stmt.Close()
	return stmt.QueryRow(trn).StructScan(trn)
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
