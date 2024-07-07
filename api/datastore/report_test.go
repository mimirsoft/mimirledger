package datastore

import (
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
			AccountSetType:     ReportAccountSetGroup,
			PredefinedAccounts: []uint64{1, 2, 3},
			RecurseSubAccounts: 0,
		},
	}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a1.ReportBody.AccountGroup).To(gomega.BeEmpty())
}

func TestReportStore_StoreAndRetrieveByID(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	store := createReportStore()

	a1 := Report{
		ReportName: "test",
		ReportBody: ReportBody{
			AccountSetType:     ReportAccountSetGroup,
			PredefinedAccounts: []uint64{1, 2, 3},
			RecurseSubAccounts: 0,
		},
	}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myReport, err := store.RetrieveByID(a1.ReportID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myReport.ReportName).To(gomega.Equal("test"))
	g.Expect(myReport.ReportBody.AccountSetType).To(gomega.Equal(ReportAccountSetGroup))
	g.Expect(myReport.ReportBody.PredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(myReport.ReportBody.RecurseSubAccounts).To(gomega.Equal(0))
}

func TestReportStore_StoreAndRetrieve(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	store := createReportStore()

	a1 := Report{
		ReportName: "test",
		ReportBody: ReportBody{
			AccountSetType:     ReportAccountSetGroup,
			AccountGroup:       AccountTypeIncome,
			PredefinedAccounts: []uint64{1, 2, 3},
			RecurseSubAccounts: 0,
			DataSetType:        ReportDataSetTypeIncome,
		},
	}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a2 := Report{
		ReportName: "test2",
		ReportBody: ReportBody{
			AccountSetType:     ReportAccountSetGroup,
			AccountGroup:       AccountTypeExpense,
			PredefinedAccounts: []uint64{1, 2, 3},
			RecurseSubAccounts: 0,
			DataSetType:        ReportDataSetTypeExpense,
		},
	}
	err = store.Store(&a2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myReports, err := store.Retrieve()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myReports).To(gomega.HaveLen(2))
	g.Expect(myReports[0].ReportName).To(gomega.Equal("test"))
	g.Expect(myReports[0].ReportBody.AccountSetType).To(gomega.Equal(ReportAccountSetGroup))
	g.Expect(myReports[0].ReportBody.AccountGroup).To(gomega.Equal(AccountTypeIncome))
	g.Expect(myReports[0].ReportBody.PredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(myReports[0].ReportBody.RecurseSubAccounts).To(gomega.Equal(0))
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
			AccountSetType:     ReportAccountSetGroup,
			AccountGroup:       AccountTypeExpense,
			PredefinedAccounts: []uint64{1, 2, 3},
			RecurseSubAccounts: 0,
			DataSetType:        ReportDataSetTypeExpense,
		},
	}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myReport, err := store.RetrieveByID(a1.ReportID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myReport.ReportName).To(gomega.Equal("test"))
	g.Expect(myReport.ReportBody.AccountSetType).To(gomega.Equal(ReportAccountSetGroup))
	g.Expect(myReport.ReportBody.PredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(myReport.ReportBody.RecurseSubAccounts).To(gomega.Equal(0))
	g.Expect(myReport.ReportBody.AccountGroup).To(gomega.Equal(AccountTypeExpense))
	g.Expect(myReport.ReportBody.DataSetType).To(gomega.Equal(ReportDataSetTypeExpense))

	myReport.ReportName = "updatedName"
	myReport.ReportBody.AccountSetType = ReportAccountSetPredefined
	myReport.ReportBody.PredefinedAccounts = []uint64{2, 3, 4, 5}
	myReport.ReportBody.RecurseSubAccounts = 1
	myReport.ReportBody.AccountGroup = ""
	myReport.ReportBody.DataSetType = ReportDataSetTypeIncome
	err = store.Update(myReport)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myReport.ReportName).To(gomega.Equal("updatedName"))
	g.Expect(myReport.ReportBody.AccountSetType).To(gomega.Equal(ReportAccountSetPredefined))
	g.Expect(myReport.ReportBody.PredefinedAccounts).To(gomega.ConsistOf([]uint64{2, 3, 4, 5}))
	g.Expect(myReport.ReportBody.RecurseSubAccounts).To(gomega.Equal(1))
	g.Expect(myReport.ReportBody.AccountGroup).To(gomega.Equal(AccountType("")))
	g.Expect(myReport.ReportBody.DataSetType).To(gomega.Equal(ReportDataSetTypeIncome))

	updatedRetrieve, err := store.RetrieveByID(a1.ReportID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updatedRetrieve.ReportName).To(gomega.Equal("updatedName"))
	g.Expect(updatedRetrieve.ReportBody.AccountSetType).To(gomega.Equal(ReportAccountSetPredefined))
	g.Expect(updatedRetrieve.ReportBody.PredefinedAccounts).To(gomega.ConsistOf([]uint64{2, 3, 4, 5}))
	g.Expect(updatedRetrieve.ReportBody.RecurseSubAccounts).To(gomega.Equal(1))
	g.Expect(updatedRetrieve.ReportBody.AccountGroup).To(gomega.Equal(AccountType("")))
	g.Expect(updatedRetrieve.ReportBody.DataSetType).To(gomega.Equal(ReportDataSetTypeIncome))
}
