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
	TransactionDate          time.Time    `db:"transaction_date,omitempty"`
	TransactionReconcileDate sql.NullTime `db:"transaction_reconcile_date"`
	TransactionComment       string       `db:"transaction_comment"`
	TransactionAmount        uint64       `db:"transaction_amount"`
	TransactionReference     string       `db:"transaction_reference"` // this could be a check number, batch ,etc
	IsReconciled             bool         `db:"is_reconciled"`
	IsSplit                  bool         `db:"is_split"`
}

// Store inserts a UserNotification into postgres
func (store TransactionStore) Store(trn *Transaction) (err error) {
	if trn.TransactionDate.IsZero() {
		trn.TransactionDate = time.Now()
	}
	query := `INSERT INTO transaction_main 
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

// Delete a transaction
func (store TransactionStore) Delete(trn *Transaction) (err error) {
	query := `Delete FROM transaction_main 
		         where transaction_id = $1`
	_, err = store.Client.Exec(query, trn.TransactionID)
	if err != nil {
		return fmt.Errorf("store.Client.Exec:%w", err)
	}
	return nil
}

// Set Transaction Reconsciled
func (store TransactionStore) SetIsReconciled(trn *Transaction) (err error) {
	query := `UPDATE  transaction_main 
		    SET is_reconciled = $2
		    where transaction_id = $1`
	_, err = store.Client.Exec(query, trn.TransactionID, trn.IsReconciled)
	if err != nil {
		return fmt.Errorf("store.Client.Exec:%w", err)
	}
	return nil
}

// Set TransactionReconcileDate
func (store TransactionStore) SetTransactionReconcileDate(trn *Transaction) (err error) {
	query := `UPDATE  transaction_main 
		    SET transaction_reconcile_date = $2
		    where transaction_id = $1`
	_, err = store.Client.Exec(query, trn.TransactionID, trn.TransactionReconcileDate)
	if err != nil {
		return fmt.Errorf("store.Client.Exec:%w", err)
	}
	return nil
}

type TransactionLedger struct {
	TransactionID            uint64       `db:"transaction_id"`
	TransactionDate          time.Time    `db:"transaction_date"`
	TransactionReconcileDate sql.NullTime `db:"transaction_reconcile_date"`
	TransactionComment       string       `db:"transaction_comment"`
	TransactionReference     string       `db:"transaction_reference"` // this could be a check number, batch ,etc
	IsReconciled             bool         `db:"is_reconciled"`
	IsSplit                  bool         `db:"is_split"`
	TransactionDCAmount      uint64       `db:"transaction_dc_amount"`
	DebitOrCredit            AccountSign  `db:"debit_or_credit"`
	// split is a generated field, a comma separated list of the other d/c
	Split string `db:"split"`
}

func (store TransactionStore) GetTransactionsForAccount(id uint64) ([]*TransactionLedger, error) {
	query := `SELECT workingDC.transaction_dc_amount, 
    							  workingDC.debit_or_credit, 
    							  tm.transaction_id, 
    							  tm.transaction_reference, 
    							  tm.transaction_date, 
    							  tm.transaction_comment, 
    							  tm.is_reconciled, 
    							  string_agg(odc.account_id::text, ',') AS split
                             FROM transaction_debit_credit AS workingDC
                        LEFT JOIN transaction_debit_credit AS odc
                               ON workingDC.transaction_id=odc.transaction_id AND odc.account_id != $1
    				   INNER JOIN transaction_main AS tm
                               ON tm.transaction_id=workingDC.transaction_id
                            WHERE workingDC.account_id = $1
                         GROUP BY  workingDC.transaction_dc_amount, 
    							  workingDC.debit_or_credit, 
    							  tm.transaction_id, 
    							  tm.transaction_reference, 
    							  tm.transaction_date, 
    							  tm.transaction_comment, 
    							  tm.is_reconciled
                         ORDER BY tm.transaction_date, tm.transaction_id`
	rows, err := store.Client.Queryx(query, id)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()
	var txnSet []*TransactionLedger
	for rows.Next() {
		var txn TransactionLedger
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
