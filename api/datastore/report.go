package datastore

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ReportStore struct {
	Client *sqlx.DB
}

type Report struct {
	ReportID   uint64     `db:"report_id,omitempty"`
	ReportName string     `db:"report_name"`
	ReportBody ReportBody `db:"report_body"`
}
type ReportBody struct {
	AccountSetType     ReportAccountSetType `json:"accountSetType"`
	PredefinedAccounts []uint64             `json:"predefinedAccounts"`
	RecurseSubAccounts int                  `json:"recurseSubAccounts"` // how many layers deep to recurse
}

// AccountType is an enum for account type.
type ReportAccountSetType string

const (
	ReportAccountSetGroup        = ReportAccountSetType("GROUP")
	ReportAccountSetPredefined   = ReportAccountSetType("PREDEFINED")
	ReportAccountSetUserSupplied = ReportAccountSetType("USER_SUPPLIED")
)

// Make the struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (rb *ReportBody) Value() (driver.Value, error) {
	return json.Marshal(rb) //nolint:wrapcheck
}

var errReportBodyScanFailed = errors.New("failed to scan report body:type assertion failed")

// Make the struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (rb *ReportBody) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errReportBodyScanFailed
	}

	return json.Unmarshal(b, &rb) //nolint:wrapcheck
}

// Store inserts a UserNotification into postgres, we do not include :report_id in our insert
func (store ReportStore) Store(myReport *Report) error {
	query := `    INSERT INTO reports 
		           (
	report_name,
	report_body)
		    VALUES (
	:report_name,
	:report_body)
		 RETURNING *`

	stmt, err := store.Client.PrepareNamed(query)
	if err != nil {
		return fmt.Errorf("error preparing report insert: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(myReport).StructScan(myReport) //nolint:musttag
	if err != nil {
		return fmt.Errorf("stmt.QueryRow().StructScan():%w", err)
	}

	return nil
}

// Store inserts a UserNotification into postgres, we do not include :report_id in our insert
func (store ReportStore) Update(myReport *Report) error {
	query := `UPDATE  reports 
		   SET (report_name,
				report_body) 
		       = (:report_name,
				:report_body)
		   WHERE report_id = :report_id
		 RETURNING *`

	stmt, err := store.Client.PrepareNamed(query)
	if err != nil {
		return fmt.Errorf("error preparing report update: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(myReport).StructScan(myReport) //nolint:musttag
	if err != nil {
		return fmt.Errorf("stmt.QueryRow().StructScan():%w", err)
	}

	return nil
}

func (store ReportStore) RetrieveByID(id uint64) (*Report, error) {
	query := `select * from reports where report_id = $1`

	row := store.Client.QueryRowx(query, id)

	var myReport Report

	if err := row.StructScan(&myReport); err != nil { //nolint:musttag
		return nil, fmt.Errorf("row.StructScan(&tn):%w", err)
	}

	return &myReport, nil
}

// Gets All Reports.
func (store ReportStore) Retrieve() ([]*Report, error) {
	query := `select * from reports order by report_name`

	rows, err := store.Client.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("store.Client.Queryx:%w", err)
	}
	defer rows.Close()

	var set []*Report

	for rows.Next() {
		var report Report
		if err = rows.StructScan(&report); err != nil { //nolint:musttag
			return nil, fmt.Errorf("rows.StructScan:%w", err)
		}

		set = append(set, &report)
	}

	if len(set) == 0 {
		return nil, sql.ErrNoRows
	}

	return set, nil
}