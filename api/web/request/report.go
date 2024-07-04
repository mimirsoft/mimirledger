package request

import (
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
)

type Report struct {
	ReportName string     `json:"reportName"`
	ReportBody ReportBody `json:"reportBody"`
}
type ReportBody struct {
	AccountSetType     datastore.ReportAccountSetType `json:"accountSetType"`
	PredefinedAccounts []uint64                       `json:"predefinedAccounts"`
	RecurseSubAccounts int                            `json:"recurseSubAccounts"` // how many layers deep to recurse
}

func ReqReportToReport(rpt *Report) *models.Report {
	return &models.Report{ //nolint:exhaustruct
		ReportName: rpt.ReportName,
		ReportBody: models.ReportBody{
			AccountSetType:     rpt.ReportBody.AccountSetType,
			PredefinedAccounts: rpt.ReportBody.PredefinedAccounts,
			RecurseSubAccounts: rpt.ReportBody.RecurseSubAccounts,
		},
	}
}
