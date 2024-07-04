package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mimirsoft/mimirledger/api/datastore"
)

// define set of accounts
type Report struct {
	ReportID   uint64
	ReportName string
	ReportBody ReportBody
}
type ReportBody struct {
	AccountSetType     datastore.ReportAccountSetType
	PredefinedAccounts []uint64
	RecurseSubAccounts int
}

// Store inserts a Report
func (c *Report) Store(dStores *datastore.Datastores) error {
	eReport := reportToEntReport(c)

	err := dStores.ReportStore().Store(eReport)
	if err != nil {
		return fmt.Errorf("ds.ReportStore().Store:%w [Report:%+v]", err, eReport)
	}

	myReport := entReportToReport(eReport)
	*c = *myReport

	return nil
}

// Update updates a Report
func (c *Report) Update(dStores *datastore.Datastores) error {
	eReport := reportToEntReport(c)

	err := dStores.ReportStore().Update(eReport)
	if err != nil {
		return fmt.Errorf("ds.ReportStore().Update:%w [Report:%+v]", err, eReport)
	}

	myReport := entReportToReport(eReport)
	*c = *myReport

	return nil
}

var ErrReportNotFound = errors.New("report not found")

// RetrieveReportByID retrieves a specific report
func RetrieveReportByID(dStores *datastore.Datastores, reportID uint64) (*Report, error) {
	store := dStores.ReportStore()

	eReport, err := store.RetrieveByID(reportID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReportNotFound
		}

		return nil, fmt.Errorf("ReportStore().RetrieveByID:%w", err)
	}

	myReport := entReportToReport(eReport)

	return myReport, nil
}

func reportToEntReport(myReport *Report) *datastore.Report {
	eReport := datastore.Report{
		ReportID:   myReport.ReportID,
		ReportName: myReport.ReportName,
		ReportBody: datastore.ReportBody{
			AccountSetType:     myReport.ReportBody.AccountSetType,
			PredefinedAccounts: myReport.ReportBody.PredefinedAccounts,
			RecurseSubAccounts: myReport.ReportBody.RecurseSubAccounts,
		},
	}

	return &eReport
}

func entReportToReport(entReport *datastore.Report) *Report {
	myReport := Report{
		ReportID:   entReport.ReportID,
		ReportName: entReport.ReportName,
		ReportBody: ReportBody{
			AccountSetType:     entReport.ReportBody.AccountSetType,
			PredefinedAccounts: entReport.ReportBody.PredefinedAccounts,
			RecurseSubAccounts: entReport.ReportBody.RecurseSubAccounts,
		},
	}

	return &myReport
}
