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
}

func TestReportStore_StoreValid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

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
}

func TestReportStore_StoreAndRetrieve(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

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
