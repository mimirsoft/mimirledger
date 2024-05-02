package datastore

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type TransactionStore struct {
	Client *sqlx.DB
}
type Transaction struct {
	TransactionID            uint64       `db:"transaction_id,omitempty"`
	TransactionDate          time.Time    `db:"transaction_date"`
	TransactionReconcileDate sql.NullTime `db:"transaction_reconcile_date"`
	TransactionComment       string       `db:"transaction_comment"`
	TransactionAmount        uint64       `db:"transaction_amount"`
	TransactionReference     string       `db:"transaction_reference"` // this could be a check number, batch ,etc
	IsReconciled             bool         `db:"is_reconciled"`
	IsSplit                  bool         `db:"is_split"`
}

// Store inserts a UserNotification into postgres
func (store TransactionStore) Store(trn *Transaction) (err error) {
	query := `    INSERT INTO transaction_main 
		           (transaction_date,
	transaction_reconcile_date,
	transaction_comment,
	transaction_amount,
	transaction_reference,
	is_reconciled,
	is_split)
		    VALUES (:transaction_date,
	:transaction_reconcile_date,
	:transaction_comment,
	:transaction_amount,
	:transaction_reference,
	:is_reconciled,
	:is_split)
		 RETURNING *`
	stmt, err := store.Client.PrepareNamed(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	return stmt.QueryRow(trn).StructScan(trn)
}

func (store TransactionStore) Update(trn *Transaction) (err error) {
	query := `    UPDATE  transaction_main 
		        SET  (transaction_date,
	transaction_reconcile_date,
	transaction_comment,
	transaction_amount,
	transaction_reference,
	is_reconciled,
	is_split)
		    = (:transaction_date,
	:transaction_reconcile_date,
	:transaction_comment,
	:transaction_amount,
	:transaction_reference,
	:is_reconciled,
	:is_split)
		     WHERE transaction_id = :transaction_id
		 RETURNING *`
	stmt, err := store.Client.PrepareNamed(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	return stmt.QueryRow(trn).StructScan(trn)
}

func (store TransactionStore) GetByID(id uint64) (*Transaction, error) {
	query := `select * from transaction_main where transaction_id = $1`
	row := store.Client.QueryRowx(query, id)
	var tn Transaction
	if err := row.StructScan(&tn); err != nil {
		return nil, err
	}
	return &tn, nil
}

func (store TransactionStore) GetTransactionsForAccount(id uint64) ([]*Transaction, error) {
	query := `select * from transaction_main where transaction_id = $1`
	rows, err := store.Client.Queryx(query, id)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()
	var txnSet []*Transaction
	for rows.Next() {
		var txn Transaction
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
