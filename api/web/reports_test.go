package web

import (
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
	"github.com/mimirsoft/mimirledger/api/web/response"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"net/http"
	"testing"
)

func TestReports_GetAll(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	var test = RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: "/reports",
	}, GomegaWithT: g, Code: http.StatusOK}

	var reportSet response.ReportSet
	test.ExecWithUnmarshal(&reportSet)
	g.Expect(reportSet.Reports).To(gomega.HaveLen(0))

	a1 := models.Report{ReportName: "testName",
		ReportBody: models.ReportBody{
			AccountSetType:     datastore.ReportAccountSetGroup,
			PredefinedAccounts: []uint64{1, 2, 3},
			RecurseSubAccounts: 1,
		}}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred()) // reset datastore

	test = RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: "/reports",
	}, GomegaWithT: g, Code: http.StatusOK}

	test.ExecWithUnmarshal(&reportSet)
	g.Expect(reportSet.Reports).To(gomega.HaveLen(1))
	g.Expect(reportSet.Reports[0].ReportName).To(gomega.Equal("testName"))
	g.Expect(reportSet.Reports[0].ReportBody.AccountSetType).To(gomega.Equal(datastore.ReportAccountSetGroup))
	g.Expect(reportSet.Reports[0].ReportBody.PredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(reportSet.Reports[0].ReportBody.RecurseSubAccounts).To(gomega.Equal(1))
}

func TestReports_PostReport(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	reqBody := map[string]interface{}{
		"reportName": "testName",
		"reportBody": map[string]interface{}{
			"accountSetType":     datastore.ReportAccountSetGroup,
			"predefinedAccounts": []uint64{1, 2, 3},
			"recurseSubAccounts": 1,
		},
	}
	var test = RouterTest{Request: Request{
		Method:     http.MethodPost,
		Router:     TestRouter,
		RequestURL: "/reports",
		Payload:    reqBody,
	}, GomegaWithT: g, Code: http.StatusOK}

	var reportRes response.Report
	test.ExecWithUnmarshal(&reportRes)
	g.Expect(reportRes.ReportName).To(gomega.Equal("testName"))
	g.Expect(reportRes.ReportBody.AccountSetType).To(gomega.Equal(datastore.ReportAccountSetGroup))
	g.Expect(reportRes.ReportBody.PredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(reportRes.ReportBody.RecurseSubAccounts).To(gomega.Equal(1))
}
