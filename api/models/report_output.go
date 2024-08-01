package models

import (
	"time"
)

// the output of a report
type ReportOutput struct {
	ReportName    string
	StartDate     time.Time
	EndDate       time.Time
	ReportDataSet []*ReportData
}

type ReportData struct {
	Expense int64
	Income  int64
}
