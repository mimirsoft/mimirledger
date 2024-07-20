package models

import (
	"time"
)

// the output of a report
type ReportOutput struct {
	ReportName string
	StartDate  time.Time
	EndDate    time.Time
}
