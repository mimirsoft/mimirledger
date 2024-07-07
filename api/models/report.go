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
	AccountGroup       datastore.AccountType
	PredefinedAccounts []uint64
	RecurseSubAccounts int
	DataSetType        datastore.ReportDataSetType
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

var ErrNoReports = errors.New("no reports found")

// RetrieveReports retrieves all reports
func RetrieveReports(dStores *datastore.Datastores) ([]*Report, error) {
	store := dStores.ReportStore()

	eReport, err := store.Retrieve()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoReports
		}

		return nil, fmt.Errorf("ReportStore().RetrieveReports:%w", err)
	}

	myReportSet := entReportsSetToReportsSet(eReport)

	return myReportSet, nil
}

func reportToEntReport(myReport *Report) *datastore.Report {
	eReport := datastore.Report{
		ReportID:   myReport.ReportID,
		ReportName: myReport.ReportName,
		ReportBody: datastore.ReportBody{
			AccountSetType:     myReport.ReportBody.AccountSetType,
			AccountGroup:       myReport.ReportBody.AccountGroup,
			PredefinedAccounts: myReport.ReportBody.PredefinedAccounts,
			RecurseSubAccounts: myReport.ReportBody.RecurseSubAccounts,
			DataSetType:        myReport.ReportBody.DataSetType,
		},
	}

	return &eReport
}
func entReportsSetToReportsSet(eReportSet []*datastore.Report) []*Report {
	reportSet := make([]*Report, len(eReportSet))

	for idx := range eReportSet {
		myRpt := entReportToReport(eReportSet[idx])
		reportSet[idx] = myRpt
	}

	return reportSet
}

func entReportToReport(entReport *datastore.Report) *Report {
	myReport := Report{
		ReportID:   entReport.ReportID,
		ReportName: entReport.ReportName,
		ReportBody: ReportBody{
			AccountSetType:     entReport.ReportBody.AccountSetType,
			AccountGroup:       entReport.ReportBody.AccountGroup,
			PredefinedAccounts: entReport.ReportBody.PredefinedAccounts,
			RecurseSubAccounts: entReport.ReportBody.RecurseSubAccounts,
			DataSetType:        entReport.ReportBody.DataSetType,
		},
	}

	return &myReport
}
