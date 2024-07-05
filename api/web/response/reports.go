package response

import (
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
)

// ReportSet is for use in accounts controller responses
type ReportSet struct {
	Reports []*Report `json:"reports"`
}
type Report struct {
	ReportID   uint64     `json:"reportID"`
	ReportName string     `json:"reportName"`
	ReportBody ReportBody `json:"reportBody"`
}
type ReportBody struct {
	AccountSetType     datastore.ReportAccountSetType `json:"accountSetType"`
	PredefinedAccounts []uint64                       `json:"predefinedAccounts"`
	RecurseSubAccounts int                            `json:"recurseSubAccounts"` // how many layers deep to recurse
}

func ReportToRespReport(rpt *models.Report) *Report {
	return &Report{
		ReportID:   rpt.ReportID,
		ReportName: rpt.ReportName,
		ReportBody: ReportBody{
			AccountSetType:     rpt.ReportBody.AccountSetType,
			PredefinedAccounts: rpt.ReportBody.PredefinedAccounts,
			RecurseSubAccounts: rpt.ReportBody.RecurseSubAccounts,
		},
	}
}

// ConvertReportsToRespReportsSet converts []models.Report to ReportSet
func ConvertReportsToRespReportsSet(rpts []*models.Report) *ReportSet {
	var ras = make([]*Report, len(rpts))
	for idx := range rpts {
		ras[idx] = ReportToRespReport(rpts[idx])
	}

	return &ReportSet{Reports: ras}
}
