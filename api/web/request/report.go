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
	SourceAccountSetType          datastore.ReportAccountSetType `json:"sourceAccountSetType"`
	SourceAccountGroup            datastore.AccountType          `json:"sourceAccountGroup,omitempty"`
	SourcePredefinedAccounts      []uint64                       `json:"sourcePredefinedAccounts"`
	SourceRecurseSubAccounts      bool                           `json:"sourceRecurseSubAccounts"`
	SourceRecurseSubAccountsDepth int                            `json:"sourceRecurseSubAccountsDepth"`
	FilterAccountSetType          datastore.ReportAccountSetType `json:"filterAccountSetType"`
	FilterAccountGroup            datastore.AccountType          `json:"filterAccountGroup,omitempty"`
	FilterPredefinedAccounts      []uint64                       `json:"filterPredefinedAccounts"`
	FilterRecurseSubAccounts      bool                           `json:"filterRecurseSubAccounts"`
	FilterRecurseSubAccountsDepth int                            `json:"filterRecurseSubAccountsDepth"`
	DataSetType                   datastore.ReportDataSetType    `json:"dataSetType"`
}

func ReqReportToReport(rpt *Report) *models.Report {
	return &models.Report{ //nolint:exhaustruct
		ReportName: rpt.ReportName,
		ReportBody: models.ReportBody{
			SourceAccountSetType:          rpt.ReportBody.SourceAccountSetType,
			SourceAccountGroup:            rpt.ReportBody.SourceAccountGroup,
			SourcePredefinedAccounts:      rpt.ReportBody.SourcePredefinedAccounts,
			SourceRecurseSubAccounts:      rpt.ReportBody.SourceRecurseSubAccounts,
			SourceRecurseSubAccountsDepth: rpt.ReportBody.SourceRecurseSubAccountsDepth,
			FilterAccountSetType:          rpt.ReportBody.FilterAccountSetType,
			FilterAccountGroup:            rpt.ReportBody.FilterAccountGroup,
			FilterPredefinedAccounts:      rpt.ReportBody.FilterPredefinedAccounts,
			FilterRecurseSubAccounts:      rpt.ReportBody.FilterRecurseSubAccounts,
			FilterRecurseSubAccountsDepth: rpt.ReportBody.FilterRecurseSubAccountsDepth,
			DataSetType:                   rpt.ReportBody.DataSetType,
		},
	}
}
