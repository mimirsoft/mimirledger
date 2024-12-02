package models

import (
	"time"

	"github.com/mimirsoft/mimirledger/api/datastore"
)

// the output of a report
type ReportOutput struct {
	ReportID    uint64
	ReportName  string
	StartDate   time.Time
	EndDate     time.Time
	DataSetType datastore.ReportDataSetType
	ReportData  []*ReportOutputData
}

type ReportOutputData struct {
	Expense         int64
	Income          int64
	NetTransactions []*TransactionLedger
}
