package web

import (
	"fmt"
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
			SourceAccountSetType:          datastore.ReportAccountSetGroup,
			SourceAccountGroup:            datastore.AccountTypeExpense,
			SourcePredefinedAccounts:      []uint64{1, 2, 3},
			SourceRecurseSubAccounts:      true,
			SourceRecurseSubAccountsDepth: 1,
			DataSetType:                   datastore.ReportDataSetTypeExpense,
		}}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())                                                // reset datastore
	g.Expect(a1.ReportBody.SourceAccountGroup).To(gomega.Equal(datastore.AccountTypeExpense)) // reset datastore
	g.Expect(a1.ReportBody.SourceRecurseSubAccounts).To(gomega.BeTrue())
	g.Expect(a1.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(1))

	test = RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: "/reports",
	}, GomegaWithT: g, Code: http.StatusOK}

	test.ExecWithUnmarshal(&reportSet)
	g.Expect(reportSet.Reports).To(gomega.HaveLen(1))
	g.Expect(reportSet.Reports[0].ReportID).To(gomega.Equal(a1.ReportID))
	g.Expect(reportSet.Reports[0].ReportName).To(gomega.Equal("testName"))
	g.Expect(reportSet.Reports[0].ReportBody.SourceAccountSetType).To(gomega.Equal(datastore.ReportAccountSetGroup))
	g.Expect(*reportSet.Reports[0].ReportBody.SourceAccountGroup).To(gomega.Equal(datastore.AccountTypeExpense))
	g.Expect(reportSet.Reports[0].ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(reportSet.Reports[0].ReportBody.SourceRecurseSubAccounts).To(gomega.BeTrue())
	g.Expect(reportSet.Reports[0].ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(1))
	g.Expect(reportSet.Reports[0].ReportBody.DataSetType).To(gomega.Equal(datastore.ReportDataSetTypeExpense))
}
func TestReports_GetReport(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	a1 := models.Report{ReportName: "testName",
		ReportBody: models.ReportBody{
			SourceAccountSetType:          datastore.ReportAccountSetGroup,
			SourceAccountGroup:            datastore.AccountTypeExpense,
			SourcePredefinedAccounts:      []uint64{1, 2, 3},
			SourceRecurseSubAccounts:      true,
			SourceRecurseSubAccountsDepth: 1,
			DataSetType:                   datastore.ReportDataSetTypeExpense,
		}}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())                                                // reset datastore
	g.Expect(a1.ReportBody.SourceAccountGroup).To(gomega.Equal(datastore.AccountTypeExpense)) // reset datastore
	g.Expect(a1.ReportBody.SourceRecurseSubAccounts).To(gomega.BeTrue())
	g.Expect(a1.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(1))

	test := RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: fmt.Sprintf("/reports/%d", a1.ReportID),
	}, GomegaWithT: g, Code: http.StatusOK}
	var respReport response.Report

	test.ExecWithUnmarshal(&respReport)
	g.Expect(respReport.ReportID).To(gomega.Equal(a1.ReportID))
	g.Expect(respReport.ReportName).To(gomega.Equal("testName"))
	g.Expect(respReport.ReportBody.SourceAccountSetType).To(gomega.Equal(datastore.ReportAccountSetGroup))
	g.Expect(*respReport.ReportBody.SourceAccountGroup).To(gomega.Equal(datastore.AccountTypeExpense))
	g.Expect(respReport.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(respReport.ReportBody.SourceRecurseSubAccounts).To(gomega.BeTrue())
	g.Expect(respReport.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(1))
	g.Expect(respReport.ReportBody.DataSetType).To(gomega.Equal(datastore.ReportDataSetTypeExpense))
}
func TestReports_DeleteReport(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	a1 := models.Report{ReportName: "testName",
		ReportBody: models.ReportBody{
			SourceAccountSetType:          datastore.ReportAccountSetGroup,
			SourceAccountGroup:            datastore.AccountTypeExpense,
			SourcePredefinedAccounts:      []uint64{1, 2, 3},
			SourceRecurseSubAccounts:      true,
			SourceRecurseSubAccountsDepth: 1,
			DataSetType:                   datastore.ReportDataSetTypeExpense,
		}}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())                                                // reset datastore
	g.Expect(a1.ReportBody.SourceAccountGroup).To(gomega.Equal(datastore.AccountTypeExpense)) // reset datastore
	g.Expect(a1.ReportBody.SourceRecurseSubAccounts).To(gomega.BeTrue())
	g.Expect(a1.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(1))

	test := RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: fmt.Sprintf("/reports/%d", a1.ReportID),
	}, GomegaWithT: g, Code: http.StatusOK}
	var respReport response.Report

	test.ExecWithUnmarshal(&respReport)
	g.Expect(respReport.ReportID).To(gomega.Equal(a1.ReportID))
	g.Expect(respReport.ReportName).To(gomega.Equal("testName"))
	g.Expect(respReport.ReportBody.SourceAccountSetType).To(gomega.Equal(datastore.ReportAccountSetGroup))
	g.Expect(*respReport.ReportBody.SourceAccountGroup).To(gomega.Equal(datastore.AccountTypeExpense))
	g.Expect(respReport.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(respReport.ReportBody.SourceRecurseSubAccounts).To(gomega.BeTrue())
	g.Expect(respReport.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(1))
	g.Expect(respReport.ReportBody.DataSetType).To(gomega.Equal(datastore.ReportDataSetTypeExpense))

	test = RouterTest{Request: Request{
		Method:     http.MethodDelete,
		Router:     TestRouter,
		RequestURL: fmt.Sprintf("/reports/%d", a1.ReportID),
	}, GomegaWithT: g, Code: http.StatusOK}

	test.ExecWithUnmarshal(&respReport)
	g.Expect(respReport.ReportID).To(gomega.Equal(a1.ReportID))
	g.Expect(respReport.ReportName).To(gomega.Equal("testName"))
	g.Expect(respReport.ReportBody.SourceAccountSetType).To(gomega.Equal(datastore.ReportAccountSetGroup))
	g.Expect(*respReport.ReportBody.SourceAccountGroup).To(gomega.Equal(datastore.AccountTypeExpense))
	g.Expect(respReport.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(respReport.ReportBody.SourceRecurseSubAccounts).To(gomega.BeTrue())
	g.Expect(respReport.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(1))
	g.Expect(respReport.ReportBody.DataSetType).To(gomega.Equal(datastore.ReportDataSetTypeExpense))

	test = RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: fmt.Sprintf("/reports/%d", a1.ReportID),
	}, GomegaWithT: g, Code: http.StatusNotFound}
	test.Exec()
}

func TestReports_PostReport(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	reqBody := map[string]interface{}{
		"reportName": "testName",
		"reportBody": map[string]interface{}{
			"sourceAccountSetType":          datastore.ReportAccountSetGroup,
			"sourceAccountGroup":            datastore.AccountTypeExpense,
			"sourcePredefinedAccounts":      []uint64{1, 2, 3},
			"sourceRecurseSubAccounts":      true,
			"sourceRecurseSubAccountsDepth": 1,
			"dataSetType":                   datastore.ReportDataSetTypeExpense,
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
	g.Expect(reportRes.ReportBody.SourceAccountSetType).To(gomega.Equal(datastore.ReportAccountSetGroup))
	g.Expect(*reportRes.ReportBody.SourceAccountGroup).To(gomega.Equal(datastore.AccountTypeExpense))
	g.Expect(reportRes.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{1, 2, 3}))
	g.Expect(reportRes.ReportBody.SourceRecurseSubAccounts).To(gomega.BeTrue())
	g.Expect(reportRes.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(1))
	g.Expect(reportRes.ReportBody.DataSetType).To(gomega.Equal(datastore.ReportDataSetTypeExpense))
}

func TestReports_PostReportUserSupplied(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	reqBody := map[string]interface{}{
		"reportName": "userSuppliedLedger",
		"reportBody": map[string]interface{}{
			"sourceAccountSetType":          datastore.ReportAccountSetUserSupplied,
			"sourcePredefinedAccounts":      []uint64{},
			"sourceRecurseSubAccounts":      true,
			"sourceRecurseSubAccountsDepth": 1,
			"dataSetType":                   datastore.ReportDataSetTypeLedger,
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
	g.Expect(reportRes.ReportName).To(gomega.Equal("userSuppliedLedger"))
	g.Expect(reportRes.ReportBody.SourceAccountSetType).To(gomega.Equal(datastore.ReportAccountSetUserSupplied))
	g.Expect(reportRes.ReportBody.SourceAccountGroup).To(gomega.BeNil())
	g.Expect(reportRes.ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{}))
	g.Expect(reportRes.ReportBody.SourceRecurseSubAccounts).To(gomega.BeTrue())
	g.Expect(reportRes.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(1))
	g.Expect(reportRes.ReportBody.DataSetType).To(gomega.Equal(datastore.ReportDataSetTypeLedger))
}

func TestReports_GetReportOutput(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	a1 := models.Report{ReportName: "testName",
		ReportBody: models.ReportBody{
			SourceAccountSetType:          datastore.ReportAccountSetGroup,
			SourceAccountGroup:            datastore.AccountTypeExpense,
			SourcePredefinedAccounts:      []uint64{1, 2, 3},
			SourceRecurseSubAccounts:      true,
			SourceRecurseSubAccountsDepth: 1,
			DataSetType:                   datastore.ReportDataSetTypeExpense,
		}}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())                                                // reset datastore
	g.Expect(a1.ReportBody.SourceAccountGroup).To(gomega.Equal(datastore.AccountTypeExpense)) // reset datastore
	g.Expect(a1.ReportBody.SourceRecurseSubAccounts).To(gomega.BeTrue())
	g.Expect(a1.ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(1))

	test := RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: fmt.Sprintf("/reports/%d/output?startDate=2020-01-01&endDate=2020-01-31", a1.ReportID),
	}, GomegaWithT: g, Code: http.StatusOK}
	var respReportOutput response.ReportOutput

	test.ExecWithUnmarshal(&respReportOutput)
	g.Expect(respReportOutput.ReportID).To(gomega.Equal(a1.ReportID))
	g.Expect(respReportOutput.ReportName).To(gomega.Equal("testName"))
	g.Expect(respReportOutput.ReportData).To(gomega.HaveLen(1))
}

func TestReports_PostRestoreDefault(t *testing.T) {
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

	var test2 = RouterTest{Request: Request{
		Method:     http.MethodPost,
		Router:     TestRouter,
		RequestURL: "/reports/restore",
	}, GomegaWithT: g, Code: http.StatusOK}

	test2.ExecWithUnmarshal(&reportSet)
	g.Expect(reportSet.Reports).To(gomega.HaveLen(3))
	g.Expect(reportSet.Reports[0].ReportName).To(gomega.Equal("LedgerReport"))
	g.Expect(reportSet.Reports[0].ReportBody.SourceAccountSetType).To(gomega.Equal(datastore.ReportAccountSetUserSupplied))
	g.Expect(reportSet.Reports[0].ReportBody.SourceAccountGroup).To(gomega.BeNil())
	g.Expect(reportSet.Reports[0].ReportBody.SourcePredefinedAccounts).To(gomega.ConsistOf([]uint64{}))
	g.Expect(reportSet.Reports[0].ReportBody.SourceRecurseSubAccounts).To(gomega.BeFalse())
	g.Expect(reportSet.Reports[0].ReportBody.SourceRecurseSubAccountsDepth).To(gomega.Equal(0))
	g.Expect(reportSet.Reports[0].ReportBody.DataSetType).To(gomega.Equal(datastore.ReportDataSetTypeLedger))
}
