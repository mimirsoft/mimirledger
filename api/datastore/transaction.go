package datastore

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
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
func (store TransactionStore) Store(trn *Transaction) error {
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
		return fmt.Errorf(" store.Client.PrepareNamed(query):%w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(trn).StructScan(trn)
	if err != nil {
		return fmt.Errorf("stmt.QueryRow(trn).StructScan(trn):%w", err)
	}

	return nil
}

func (store TransactionStore) Update(trn *Transaction) error {
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
		return fmt.Errorf(" store.Client.PrepareNamed(query):%w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(trn).StructScan(trn)
	if err != nil {
		return fmt.Errorf("stmt.QueryRow(trn).StructScan(trn):%w", err)
	}

	return nil
}

func (store TransactionStore) GetByID(id uint64) (*Transaction, error) {
	query := `select * from transaction_main where transaction_id = $1`
	row := store.Client.QueryRowx(query, id)

	var myTransaction Transaction

	if err := row.StructScan(&myTransaction); err != nil {
		return nil, fmt.Errorf("row.StructScan(&tn):%w", err)
	}

	return &myTransaction, nil
}

// Delete a transaction
func (store TransactionStore) Delete(trn *Transaction) error {
	query := `Delete FROM transaction_main 
		         where transaction_id = $1`

	_, err := store.Client.Exec(query, trn.TransactionID)
	if err != nil {
		return fmt.Errorf("store.Client.Exec:%w", err)
	}

	return nil
}

// Set Transaction Reconsciled
func (store TransactionStore) SetIsReconciled(trn *Transaction) error {
	query := `UPDATE  transaction_main 
		    SET is_reconciled = $2
		    where transaction_id = $1`

	_, err := store.Client.Exec(query, trn.TransactionID, trn.IsReconciled)
	if err != nil {
		return fmt.Errorf("store.Client.Exec:%w", err)
	}

	return nil
}

// Set TransactionReconcileDate
func (store TransactionStore) SetTransactionReconcileDate(trn *Transaction) error {
	query := `UPDATE  transaction_main 
		    SET transaction_reconcile_date = $2
		    where transaction_id = $1`

	_, err := store.Client.Exec(query, trn.TransactionID, trn.TransactionReconcileDate)
	if err != nil {
		return fmt.Errorf("store.Client.Exec:%w", err)
	}

	return nil
}

type TransactionReconciliation struct {
	TransactionID            uint64       `db:"transaction_id"`
	AccountID                uint64       `db:"account_id"`
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

func (store TransactionStore) GetUnreconciledTransactionsOnAccountForDate(accountLeft, accountRight uint64,
	searchLimitDate time.Time, reconciledCutoffDate time.Time) ([]*TransactionReconciliation, error) {
	query := `SELECT workingtdc.debit_or_credit, 
               workingtdc.transaction_id, 
               workingtdc.transaction_dc_amount, 
               workingtdc.account_id, 
               tm.transaction_date, 
               tm.transaction_reference, 
               tm.transaction_comment, 
               tm.is_reconciled, 
               tm.transaction_reconcile_date, 
               string_agg(odc.account_id::text, ',') AS split
          FROM transaction_debit_credit AS workingtdc
     LEFT JOIN transaction_debit_credit AS odc
            ON workingtdc.transaction_id=odc.transaction_id
    INNER JOIN transaction_main AS tm
            ON tm.transaction_id=workingtdc.transaction_id
         WHERE (workingtdc.account_id 
			IN (SELECT account_id FROM transaction_accounts WHERE account_left BETWEEN $3 AND $4)  
		   AND odc.account_id 
		NOT IN (SELECT account_id FROM transaction_accounts WHERE account_left BETWEEN $3 AND $4) )
           AND ((tm.is_reconciled IS FALSE 
                 AND  EXTRACT(EPOCH FROM tm.transaction_date) <= EXTRACT(EPOCH FROM $1::timestamp) )
               OR
               (tm.is_reconciled IS TRUE 
                 AND EXTRACT(EPOCH FROM tm.transaction_reconcile_date) > EXTRACT(EPOCH FROM  $2::timestamp)
                 AND EXTRACT(EPOCH FROM tm.transaction_reconcile_date) <= EXTRACT(EPOCH FROM  $1::timestamp))
               )
			 GROUP BY  workingtdc.transaction_dc_amount, 
					  workingtdc.debit_or_credit, 
					  workingtdc.transaction_id, 
					  workingtdc.account_id, 
					  tm.transaction_date, 
					  tm.transaction_reconcile_date, 
					  tm.transaction_reference, 
					  tm.transaction_comment, 
					  tm.is_reconciled
      ORDER BY is_reconciled DESC, transaction_reconcile_date ASC, transaction_date ASC, transaction_reference ASC`

	rows, err := store.Client.Queryx(query, searchLimitDate, reconciledCutoffDate, accountLeft, accountRight)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()

	var txnSet []*TransactionReconciliation

	for rows.Next() {
		var txn TransactionReconciliation
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

func (store TransactionStore) GetTransactionsForAccount(accountID uint64) ([]*TransactionLedger, error) {
	query := `SELECT workingDC.transaction_dc_amount, 
    							  workingDC.debit_or_credit, 
    							  tm.transaction_id, 
    							  tm.transaction_reference, 
    							  tm.transaction_date, 
    							  tm.transaction_reconcile_date, 
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
    							  tm.transaction_reconcile_date, 
    							  tm.transaction_comment, 
    							  tm.is_reconciled
                         ORDER BY tm.transaction_date, tm.transaction_id`

	rows, err := store.Client.Queryx(query, accountID)
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

func (store TransactionStore) RetrieveTransactionsNetForDates(accountIDSet []uint64,
	startDate time.Time, endDate time.Time) ([]*TransactionLedger, error) {
	query := `SELECT workingDC.transaction_dc_amount, 
    							  workingDC.debit_or_credit, 
    							  tm.transaction_id, 
    							  tm.transaction_reference, 
    							  tm.transaction_date, 
    							  tm.transaction_reconcile_date, 
    							  tm.transaction_comment, 
    							  tm.is_reconciled, 
    							  string_agg(odc.account_id::text, ',') AS split
                             FROM transaction_debit_credit AS workingDC
                        LEFT JOIN transaction_debit_credit AS odc
                               ON workingDC.transaction_id=odc.transaction_id 
    				   INNER JOIN transaction_main AS tm
                               ON tm.transaction_id=workingDC.transaction_id
                            WHERE workingDC.account_id = ANY($1::int[])
                              AND odc.account_id != ANY($1::int[])
						  AND  EXTRACT(EPOCH FROM tm.transaction_date) >= EXTRACT(EPOCH FROM $2::timestamp) 
               	 		  AND  EXTRACT(EPOCH FROM tm.transaction_date) <= EXTRACT(EPOCH FROM $3::timestamp) 
                         GROUP BY  workingDC.transaction_dc_amount, 
    							  workingDC.debit_or_credit, 
    							  tm.transaction_id, 
    							  tm.transaction_reference, 
    							  tm.transaction_date, 
    							  tm.transaction_reconcile_date, 
    							  tm.transaction_comment, 
    							  tm.is_reconciled
                         ORDER BY tm.transaction_date, tm.transaction_id`

	rows, err := store.Client.Queryx(query, accountIDSet, startDate, endDate)
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

func (store TransactionStore) GetDebitsForAccounts(accountIDs []uint64) (int64, error) {
	return 0, nil
}
func (store TransactionStore) GetDebitsForAccountsFiltered(accountID []uint64, filteredAccounts []uint64) (int64, error) {
	return 0, nil
}

func (store TransactionStore) GetCreditsForAccounts(accountID []uint64) (int64, error) {
	return 0, nil
}
func (store TransactionStore) GetCreditsForAccountsFiltered(accountID []uint64, filteredAccounts []uint64) (int64, error) {
	return 0, nil
}
