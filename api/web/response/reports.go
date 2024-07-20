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
	AccountSetType          datastore.ReportAccountSetType `json:"accountSetType"`
	AccountGroup            *datastore.AccountType         `json:"accountGroup,omitempty"`
	PredefinedAccounts      []uint64                       `json:"predefinedAccounts,omitempty"`
	RecurseSubAccounts      bool                           `json:"recurseSubAccounts"`
	RecurseSubAccountsDepth int                            `json:"recurseSubAccountsDepth"`
	DataSetType             datastore.ReportDataSetType    `json:"dataSetType"`
}

func ReportToRespReport(rpt *models.Report) *Report {
	var accountGroup *datastore.AccountType
	if rpt.ReportBody.AccountGroup != "" {
		accountGroup = &rpt.ReportBody.AccountGroup
	}

	myReport := &Report{
		ReportID:   rpt.ReportID,
		ReportName: rpt.ReportName,
		ReportBody: ReportBody{
			AccountSetType:          rpt.ReportBody.AccountSetType,
			AccountGroup:            accountGroup,
			PredefinedAccounts:      rpt.ReportBody.PredefinedAccounts,
			RecurseSubAccounts:      rpt.ReportBody.RecurseSubAccounts,
			RecurseSubAccountsDepth: rpt.ReportBody.RecurseSubAccountsDepth,
			DataSetType:             rpt.ReportBody.DataSetType,
		},
	}

	return myReport
}

// ConvertReportsToRespReportsSet converts []models.Report to ReportSet
func ConvertReportsToRespReportsSet(rpts []*models.Report) *ReportSet {
	var ras = make([]*Report, len(rpts))
	for idx := range rpts {
		ras[idx] = ReportToRespReport(rpts[idx])
	}

	return &ReportSet{Reports: ras}
}
