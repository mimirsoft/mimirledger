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

// Run executes a report and generates an output
func (c *Report) Run(dStores *datastore.Datastores, startDate time.Time,
	endDate time.Time, runTimeTargetAccounts []uint64) (*ReportOutput, error) {
	myReportOutput := ReportOutput{ReportName: c.ReportName,
		StartDate: startDate,
		EndDate:   endDate}
	// check type of account set
	// build the set of accountIDs to process
	accountSet := make(map[uint64]bool)
	switch c.ReportBody.SourceAccountSetType {
	case datastore.ReportAccountSetGroup:
		// get all accounts in group

	case datastore.ReportAccountSetPredefined:

	case datastore.ReportAccountSetUserSupplied:
		// build total account set
		// get userSuppliedAccountIDs
		// include sub-accounts
		if c.ReportBody.SourceRecurseSubAccounts {
			for _, account := range runTimeTargetAccounts {
				accountAndChildren, err := dStores.AccountStore().GetAccountWithChildrenByLevel(account)
				if err != nil {
					return nil, fmt.Errorf("dStores.AccountStore().GetAccountWithChildrenByLevel:%w", err)
				}
				// how many levels to recurse
				if c.ReportBody.SourceRecurseSubAccountsDepth > 0 {
					for idx := range accountAndChildren {
						if accountAndChildren[idx].Level < c.ReportBody.SourceRecurseSubAccountsDepth {
							accountSet[accountAndChildren[idx].AccountID] = true
						}
					}
				} else { // otherwise no limit to depth add all accounts to the accountSet
					for idx := range accountAndChildren {
						accountSet[accountAndChildren[idx].AccountID] = true
					}
				}
			}
		} else {
			// otherwise we just use the runTimeTargetAccounts
			for idx := range runTimeTargetAccounts {
				accountSet[runTimeTargetAccounts[idx]] = true
			}
		}
	}

	switch c.ReportBody.DataSetType {
	case datastore.ReportDataSetTypeBalance:
	case datastore.ReportDataSetTypeLedger:
	case datastore.ReportDataSetTypeIncome:
	case datastore.ReportDataSetTypeExpense:
	}
	// build data set from account set
	return &myReportOutput, nil
}

func buildDataSetExpense(dStores *datastore.Datastores, accountSet map[uint64]bool) ([]*Transaction, error) {
	var dataSet []*Transaction

	for idx := range accountSet {
		expenseTransactions, err := dStores.TransactionStore().GetExpensesForAccount(idx)
		if err != nil {
			return nil, fmt.Errorf("dStores.TransactionStore().GetExpensesForAccount:%w", err)
		}

		for kdx := range expenseTransactions {
			reportTxn := entTransactionToTransaction(expenseTransactions[kdx])
			dataSet = append(dataSet, reportTxn)
		}
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
