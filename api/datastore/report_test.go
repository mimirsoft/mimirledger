package datastore

import (
	"database/sql"
	"errors"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func createReportStore() ReportStore {
	return ReportStore{
		Client: TestPostgresClient,
	}
}

func TestReportStore_StoreInvalidEmpty(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	store := createReportStore()

	// failed due to no ReportName
	a1 := Report{}
	err := store.Store(&a1)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring(`new row for relation "reports" violates` +
		` check constraint "reports_report_name_check"`))

	// passes with only a ReportName
	a2 := Report{
		ReportName: "report_name",
	}
	err = store.Store(&a2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// try to save another with same name should throw error
	a3 := Report{
		ReportName: "report_name",
	}
	err = store.Store(&a3)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring(`duplicate key value violates unique constraint ` +
		`"reports_report_name_key"`))
}

func TestReportStore_StoreValid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	store := createReportStore()

	a1 := Report{
		ReportName: "test",
		ReportBody: ReportBody{
			SourceAccountSetType:          ReportAccountSetGroup,
			SourcePredefinedAccounts:      []uint64{1, 2, 3},
			SourceRecurseSubAccounts:      false,
			SourceRecurseSubAccountsDepth: 0,
		},
	}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a1.ReportBody.SourceAccountGroup).To(gomega.BeEmpty())
}

func TestReportStore_StoreAndRetrieveByID(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	store := createReportStore()

	a1 := Report{
		ReportName: "test",
		ReportBody: ReportBody{
			SourceAccountSetType:          ReportAccountSetGroup,
			SourcePredefinedAccounts:      []uint64{1, 2, 3},
			SourceRecurseSubAccounts:      false,
			SourceRecurseSubAccountsDepth: 0,
		},
	}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myReport, err := store.RetrieveByID(a1.ReportID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myReport.ReportName).To(gomega.Equal("test"))
	g.Expect(myReport.ReportBody.SourceAccountSetType).To(gomega.Equal(ReportAccountSetGroup))
	g.Expect(myReport.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(myReport.ReportBody.SourceRecurseSubAccounts).To(gomega.BeFalse())
	g.Expect(myReport.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(0))
}

func TestReportStore_StoreAndRetrieve(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	store := createReportStore()

	a1 := Report{
		ReportName: "test",
		ReportBody: ReportBody{
			SourceAccountSetType:          ReportAccountSetGroup,
			SourceAccountGroup:            AccountTypeIncome,
			SourcePredefinedAccounts:      []uint64{1, 2, 3},
			SourceRecurseSubAccounts:      false,
			SourceRecurseSubAccountsDepth: 0,
			DataSetType:                   ReportDataSetTypeIncome,
		},
	}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a2 := Report{
		ReportName: "test2",
		ReportBody: ReportBody{
			SourceAccountSetType:          ReportAccountSetGroup,
			SourceAccountGroup:            AccountTypeExpense,
			SourcePredefinedAccounts:      []uint64{1, 2, 3},
			SourceRecurseSubAccounts:      false,
			SourceRecurseSubAccountsDepth: 0,
			DataSetType:                   ReportDataSetTypeExpense,
		},
	}
	err = store.Store(&a2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myReports, err := store.Retrieve()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myReports).To(gomega.HaveLen(2))
	g.Expect(myReports[0].ReportName).To(gomega.Equal("test"))
	g.Expect(myReports[0].ReportBody.SourceAccountSetType).To(gomega.Equal(ReportAccountSetGroup))
	g.Expect(myReports[0].ReportBody.SourceAccountGroup).To(gomega.Equal(AccountTypeIncome))
	g.Expect(myReports[0].ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(myReports[0].ReportBody.SourceRecurseSubAccounts).To(gomega.BeFalse())
	g.Expect(myReports[0].ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(0))
	g.Expect(myReports[0].ReportBody.DataSetType).To(gomega.Equal(ReportDataSetTypeIncome))
	g.Expect(myReports[1].ReportName).To(gomega.Equal("test2"))
}

func TestReportStore_StoreAndUpdate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	store := createReportStore()

	a1 := Report{
		ReportName: "test",
		ReportBody: ReportBody{
			SourceAccountSetType:          ReportAccountSetGroup,
			SourceAccountGroup:            AccountTypeExpense,
			SourcePredefinedAccounts:      []uint64{1, 2, 3},
			SourceRecurseSubAccounts:      false,
			SourceRecurseSubAccountsDepth: 0,
			DataSetType:                   ReportDataSetTypeExpense,
		},
	}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myReport, err := store.RetrieveByID(a1.ReportID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myReport.ReportName).To(gomega.Equal("test"))
	g.Expect(myReport.ReportBody.SourceAccountSetType).To(gomega.Equal(ReportAccountSetGroup))
	g.Expect(myReport.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(myReport.ReportBody.SourceRecurseSubAccounts).To(gomega.BeFalse())
	g.Expect(myReport.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(0))
	g.Expect(myReport.ReportBody.SourceAccountGroup).To(gomega.Equal(AccountTypeExpense))
	g.Expect(myReport.ReportBody.DataSetType).To(gomega.Equal(ReportDataSetTypeExpense))

	myReport.ReportName = "updatedName"
	myReport.ReportBody.SourceAccountSetType = ReportAccountSetPredefined
	myReport.ReportBody.SourcePredefinedAccounts = []uint64{2, 3, 4, 5}
	myReport.ReportBody.SourceRecurseSubAccounts = true
	myReport.ReportBody.SourceRecurseSubAccountsDepth = 2
	myReport.ReportBody.SourceAccountGroup = ""
	myReport.ReportBody.DataSetType = ReportDataSetTypeIncome
	err = store.Update(myReport)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myReport.ReportName).To(gomega.Equal("updatedName"))
	g.Expect(myReport.ReportBody.SourceAccountSetType).To(gomega.Equal(ReportAccountSetPredefined))
	g.Expect(myReport.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{2, 3, 4, 5}))
	g.Expect(myReport.ReportBody.SourceRecurseSubAccounts).To(gomega.BeTrue())
	g.Expect(myReport.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(2))
	g.Expect(myReport.ReportBody.SourceAccountGroup).To(gomega.Equal(AccountType("")))
	g.Expect(myReport.ReportBody.DataSetType).To(gomega.Equal(ReportDataSetTypeIncome))

	updatedRetrieve, err := store.RetrieveByID(a1.ReportID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updatedRetrieve.ReportName).To(gomega.Equal("updatedName"))
	g.Expect(updatedRetrieve.ReportBody.SourceAccountSetType).To(gomega.Equal(ReportAccountSetPredefined))
	g.Expect(updatedRetrieve.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{2, 3, 4, 5}))
	g.Expect(myReport.ReportBody.SourceRecurseSubAccounts).To(gomega.BeTrue())
	g.Expect(myReport.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(2))
	g.Expect(updatedRetrieve.ReportBody.SourceAccountGroup).To(gomega.Equal(AccountType("")))
	g.Expect(updatedRetrieve.ReportBody.DataSetType).To(gomega.Equal(ReportDataSetTypeIncome))
}

func TestReportStore_StoreAndDelete(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	store := createReportStore()

	a1 := Report{
		ReportName: "test",
		ReportBody: ReportBody{
			SourceAccountSetType:          ReportAccountSetGroup,
			SourcePredefinedAccounts:      []uint64{1, 2, 3},
			SourceRecurseSubAccounts:      false,
			SourceRecurseSubAccountsDepth: 0,
		},
	}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myReport, err := store.RetrieveByID(a1.ReportID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myReport.ReportName).To(gomega.Equal("test"))
	g.Expect(myReport.ReportBody.SourceAccountSetType).To(gomega.Equal(ReportAccountSetGroup))
	g.Expect(myReport.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(myReport.ReportBody.SourceRecurseSubAccounts).To(gomega.BeFalse())
	g.Expect(myReport.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(0))

	err = store.Delete(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myReport, err = store.RetrieveByID(a1.ReportID)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, sql.ErrNoRows)).To(gomega.BeTrue())
	g.Expect(myReport).To(gomega.BeNil())
}

func TestReportStore_StoreDefault(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	store := createReportStore()

	ledgerReport := Report{
		ReportID:   1,
		ReportName: "Ledger Report",
		ReportBody: ReportBody{
			SourceAccountSetType:          ReportAccountSetUserSupplied,
			SourcePredefinedAccounts:      []uint64{},
			SourceRecurseSubAccounts:      false,
			SourceRecurseSubAccountsDepth: 0,
			DataSetType:                   ReportDataSetTypeLedger,
		},
	}
	err := store.StoreOrUpdate(&ledgerReport)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myReport, err := store.RetrieveByID(ledgerReport.ReportID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myReport.ReportName).To(gomega.Equal("Ledger Report"))
	g.Expect(myReport.ReportBody.SourceAccountSetType).To(gomega.Equal(ReportAccountSetUserSupplied))
	g.Expect(myReport.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{}))
	g.Expect(myReport.ReportBody.SourceRecurseSubAccounts).To(gomega.BeFalse())
	g.Expect(myReport.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(0))
	g.Expect(myReport.ReportBody.DataSetType).To(gomega.Equal(ReportDataSetTypeLedger))

	ledgerReportUpdate := Report{
		ReportID:   myReport.ReportID,
		ReportName: "Ledger Report",
		ReportBody: ReportBody{
			SourceAccountSetType:          ReportAccountSetGroup,
			SourceAccountGroup:            AccountTypeExpense,
			SourcePredefinedAccounts:      []uint64{},
			SourceRecurseSubAccounts:      false,
			SourceRecurseSubAccountsDepth: 0,
			DataSetType:                   ReportDataSetTypeExpense,
		},
	}
	err = store.Update(&ledgerReportUpdate)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(ledgerReport.ReportID).To(gomega.Equal(myReport.ReportID))

	myReport, err = store.RetrieveByID(myReport.ReportID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myReport.ReportName).To(gomega.Equal("Ledger Report"))
	g.Expect(myReport.ReportBody.SourceAccountSetType).To(gomega.Equal(ReportAccountSetGroup))
	g.Expect(myReport.ReportBody.SourceAccountGroup).To(gomega.Equal(AccountTypeExpense))
	g.Expect(myReport.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{}))
	g.Expect(myReport.ReportBody.SourceRecurseSubAccounts).To(gomega.BeFalse())
	g.Expect(myReport.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(0))
	g.Expect(myReport.ReportBody.DataSetType).To(gomega.Equal(ReportDataSetTypeExpense))

	err = store.StoreOrUpdate(&ledgerReport)
	myReport, err = store.RetrieveByID(ledgerReport.ReportID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myReport.ReportName).To(gomega.Equal("Ledger Report"))
	g.Expect(myReport.ReportBody.SourceAccountSetType).To(gomega.Equal(ReportAccountSetUserSupplied))
	g.Expect(myReport.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{}))
	g.Expect(myReport.ReportBody.SourceRecurseSubAccounts).To(gomega.BeFalse())
	g.Expect(myReport.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(0))
	g.Expect(myReport.ReportBody.DataSetType).To(gomega.Equal(ReportDataSetTypeLedger))
}
