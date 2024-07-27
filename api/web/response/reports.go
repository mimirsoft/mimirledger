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
	SourceAccountSetType          datastore.ReportAccountSetType `json:"sourceAccountSetType"`
	SourceAccountGroup            *datastore.AccountType         `json:"sourceAccountGroup,omitempty"`
	SourcePredefinedAccounts      []uint64                       `json:"sourcePredefinedAccounts"`
	SourceRecurseSubAccounts      bool                           `json:"sourceRecurseSubAccounts"`
	SourceRecurseSubAccountsDepth int                            `json:"sourceRecurseSubAccountsDepth"`
	FilterAccountSetType          datastore.ReportAccountSetType `json:"filterAccountSetType"`
	FilterAccountGroup            *datastore.AccountType         `json:"filterAccountGroup,omitempty"`
	FilterPredefinedAccounts      []uint64                       `json:"filterPredefinedAccounts"`
	FilterRecurseSubAccounts      bool                           `json:"filterRecurseSubAccounts"`
	FilterRecurseSubAccountsDepth int                            `json:"filterRecurseSubAccountsDepth"`
	DataSetType                   datastore.ReportDataSetType    `json:"dataSetType"`
}

func ReportToRespReport(rpt *models.Report) *Report {
	var accountGroup *datastore.AccountType
	if rpt.ReportBody.SourceAccountGroup != "" {
		accountGroup = &rpt.ReportBody.SourceAccountGroup
	}

	var filterAccountGroup *datastore.AccountType
	if rpt.ReportBody.FilterAccountGroup != "" {
		filterAccountGroup = &rpt.ReportBody.FilterAccountGroup
	}

	myReport := &Report{
		ReportID:   rpt.ReportID,
		ReportName: rpt.ReportName,
		ReportBody: ReportBody{
			SourceAccountSetType:          rpt.ReportBody.SourceAccountSetType,
			SourceAccountGroup:            accountGroup,
			SourcePredefinedAccounts:      rpt.ReportBody.SourcePredefinedAccounts,
			SourceRecurseSubAccounts:      rpt.ReportBody.SourceRecurseSubAccounts,
			SourceRecurseSubAccountsDepth: rpt.ReportBody.SourceRecurseSubAccountsDepth,
			FilterAccountSetType:          rpt.ReportBody.FilterAccountSetType,
			FilterAccountGroup:            filterAccountGroup,
			FilterPredefinedAccounts:      rpt.ReportBody.FilterPredefinedAccounts,
			FilterRecurseSubAccounts:      rpt.ReportBody.FilterRecurseSubAccounts,
			FilterRecurseSubAccountsDepth: rpt.ReportBody.FilterRecurseSubAccountsDepth,
			DataSetType:                   rpt.ReportBody.DataSetType,
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
