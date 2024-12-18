package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mimirsoft/mimirledger/api/datastore"
)

// define a report
type Report struct {
	ReportID   uint64
	ReportName string
	ReportBody ReportBody
}
type ReportBody struct {
	SourceAccountSetType          datastore.ReportAccountSetType
	SourceAccountGroup            datastore.AccountType
	SourcePredefinedAccounts      []uint64
	SourceRecurseSubAccounts      bool
	SourceRecurseSubAccountsDepth int
	FilterAccountSetType          datastore.ReportAccountSetType
	FilterAccountGroup            datastore.AccountType
	FilterPredefinedAccounts      []uint64
	FilterRecurseSubAccounts      bool
	FilterRecurseSubAccountsDepth int
	DataSetType                   datastore.ReportDataSetType
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

// StoreOrUpdate stores a report, or updates if a report of the same name already existst
func (c *Report) StoreOrUpdate(dStores *datastore.Datastores) error {
	eReport := reportToEntReport(c)

	err := dStores.ReportStore().StoreOrUpdate(eReport)
	if err != nil {
		return fmt.Errorf("ds.ReportStore().StoreOrUpdate:%w [Report:%+v]", err, eReport)
	}

	myReport := entReportToReport(eReport)
	*c = *myReport

	return nil
}

func (c *Report) Delete(dStores *datastore.Datastores) error {
	eReport := reportToEntReport(c)

	// delete the existing report
	err := dStores.ReportStore().Delete(eReport)
	if err != nil {
		return fmt.Errorf("ds.ReportStore().Delete:%w [Report:%+v]", err, c)
	}

	return nil
}

// Run executes a report and generates an output
func (c *Report) Run(dStores *datastore.Datastores, startDate time.Time,
	endDate time.Time, runTimeTargetAccounts []uint64) (*ReportOutput, error) {
	myReportOutput := ReportOutput{
		ReportID:   c.ReportID,
		ReportName: c.ReportName, //nolint:exhaustruct
		StartDate:  startDate,
		EndDate:    endDate}

	sourceAccountSet, err := c.buildAccountSet(dStores, runTimeTargetAccounts)
	if err != nil {
		return nil, fmt.Errorf("buildDataSetLedger:%w", err)
	}

	var reportDataSet []*ReportOutputData

	switch c.ReportBody.DataSetType {
	case datastore.ReportDataSetTypeBalance:
	case datastore.ReportDataSetTypeLedger:
		var err error

		reportDataSet, err = buildDataSetLedger(dStores, sourceAccountSet, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("buildDataSetLedger:%w", err)
		}
	case datastore.ReportDataSetTypeIncome:
	case datastore.ReportDataSetTypeExpense:
		var err error

		reportDataSet, err = buildDataSetExpense(dStores, sourceAccountSet, nil)
		if err != nil {
			return nil, fmt.Errorf("buildDataSetExpense:%w", err)
		}
	}

	myReportOutput.DataSetType = c.ReportBody.DataSetType
	myReportOutput.ReportData = reportDataSet
	// build data set from account set
	return &myReportOutput, nil
}

// check type of account set
// build the set of accountIDs to process
func (c *Report) buildAccountSet(dStores *datastore.Datastores, runTimeTargetAccounts []uint64) ([]uint64, error) { //nolint:gocognit
	var sourceAccountSet []uint64
	sourceAccountMap := make(map[uint64]bool)

	switch c.ReportBody.SourceAccountSetType {
	case datastore.ReportAccountSetNone:
	case datastore.ReportAccountSetGroup:
		// get all accounts in group

	case datastore.ReportAccountSetPredefined:

	case datastore.ReportAccountSetUserSupplied:
		// build total account set
		// get userSuppliedAccountIDs
		// include sub-accounts
		if c.ReportBody.SourceRecurseSubAccounts { //nolint:nestif
			for _, account := range runTimeTargetAccounts {
				accountAndChildren, err := dStores.AccountStore().GetAccountWithChildrenByLevel(account)
				if err != nil {
					return nil, fmt.Errorf("dStores.AccountStore().GetAccountWithChildrenByLevel:%w", err)
				}
				// how many levels to recurse
				if c.ReportBody.SourceRecurseSubAccountsDepth > 0 {
					for idx := range accountAndChildren {
						if accountAndChildren[idx].Level < c.ReportBody.SourceRecurseSubAccountsDepth {
							sourceAccountMap[accountAndChildren[idx].AccountID] = true
						}
					}
				} else { // otherwise no limit to depth add all accounts to the accountSet
					for idx := range accountAndChildren {
						sourceAccountMap[accountAndChildren[idx].AccountID] = true
					}
				}
			}
		} else {
			// otherwise we just use the runTimeTargetAccounts
			for idx := range runTimeTargetAccounts {
				sourceAccountMap[runTimeTargetAccounts[idx]] = true
			}
		}

		sourceAccountSet = make([]uint64, 0, len(sourceAccountMap))

		for idx := range sourceAccountMap {
			sourceAccountSet = append(sourceAccountSet, idx)
		}
	}
	return sourceAccountSet, nil
}

func buildDataSetLedger(dStores *datastore.Datastores, sourceAccountSet []uint64,
	startDate time.Time, endDate time.Time) ([]*ReportOutputData, error) {
	var dataSet []*ReportOutputData

	entNetTransactions, err := dStores.TransactionStore().RetrieveTransactionsNetForDates(sourceAccountSet,
		startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("dStores.TransactionStore().GetExpensesForAccount:%w", err)
	}
	// convert to model
	netTransactions := entTransactionsLedgerToTransactionsLedger(entNetTransactions)
	dataRaw := ReportOutputData{NetTransactions: netTransactions} //nolint:exhaustruct
	dataSet = append(dataSet, &dataRaw)

	return dataSet, nil
}
func buildDataSetExpense(dStores *datastore.Datastores, sourceAccountSet []uint64,
	filterAccountSet []uint64) ([]*ReportOutputData, error) {
	var dataSet []*ReportOutputData

	// to do, at add switch on accountType "which are expenses, DEBITS or CREDIT"
	if len(filterAccountSet) == 0 {
		expenses, err := dStores.TransactionStore().GetDebitTotalForAccounts(sourceAccountSet)
		if err != nil {
			return nil, fmt.Errorf("dStores.TransactionStore().GetExpensesForAccount:%w", err)
		}

		dataRaw := ReportOutputData{Expense: expenses} //nolint:exhaustruct
		dataSet = append(dataSet, &dataRaw)
	} else {
		expenses, err := dStores.TransactionStore().GetDebitTotalForAccountsFiltered(sourceAccountSet, filterAccountSet)
		if err != nil {
			return nil, fmt.Errorf("dStores.TransactionStore().GetExpensesForAccount:%w", err)
		}

		dataRaw := ReportOutputData{Expense: expenses} //nolint:exhaustruct
		dataSet = append(dataSet, &dataRaw)
	}

	return dataSet, nil
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
			SourceAccountSetType:          myReport.ReportBody.SourceAccountSetType,
			SourceAccountGroup:            myReport.ReportBody.SourceAccountGroup,
			SourcePredefinedAccounts:      myReport.ReportBody.SourcePredefinedAccounts,
			SourceRecurseSubAccounts:      myReport.ReportBody.SourceRecurseSubAccounts,
			SourceRecurseSubAccountsDepth: myReport.ReportBody.SourceRecurseSubAccountsDepth,
			FilterAccountSetType:          myReport.ReportBody.FilterAccountSetType,
			FilterAccountGroup:            myReport.ReportBody.FilterAccountGroup,
			FilterPredefinedAccounts:      myReport.ReportBody.FilterPredefinedAccounts,
			FilterRecurseSubAccounts:      myReport.ReportBody.FilterRecurseSubAccounts,
			FilterRecurseSubAccountsDepth: myReport.ReportBody.FilterRecurseSubAccountsDepth,
			DataSetType:                   myReport.ReportBody.DataSetType,
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
			SourceAccountSetType:          entReport.ReportBody.SourceAccountSetType,
			SourceAccountGroup:            entReport.ReportBody.SourceAccountGroup,
			SourcePredefinedAccounts:      entReport.ReportBody.SourcePredefinedAccounts,
			SourceRecurseSubAccounts:      entReport.ReportBody.SourceRecurseSubAccounts,
			SourceRecurseSubAccountsDepth: entReport.ReportBody.SourceRecurseSubAccountsDepth,
			FilterAccountSetType:          entReport.ReportBody.FilterAccountSetType,
			FilterAccountGroup:            entReport.ReportBody.FilterAccountGroup,
			FilterPredefinedAccounts:      entReport.ReportBody.FilterPredefinedAccounts,
			FilterRecurseSubAccounts:      entReport.ReportBody.FilterRecurseSubAccounts,
			FilterRecurseSubAccountsDepth: entReport.ReportBody.FilterRecurseSubAccountsDepth,
			DataSetType:                   entReport.ReportBody.DataSetType,
		},
	}

	return &myReport
}
