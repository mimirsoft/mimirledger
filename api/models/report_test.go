package models

import (
	"errors"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestReport_StoreAndDelete(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	myReport := Report{ReportName: "MyBank"}
	err := myReport.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	retReport, err := RetrieveReportByID(testDS, myReport.ReportID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(retReport.ReportName).To(gomega.Equal(myReport.ReportName))

	err = myReport.Delete(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	retReport, err = RetrieveReportByID(testDS, myReport.ReportID)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, ErrReportNotFound)).To(gomega.BeTrue())
	g.Expect(retReport).To(gomega.BeNil())
}
