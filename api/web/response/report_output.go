package response

import (
	"time"

	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
)

type ReportOutput struct {
	ReportID    uint64                      `json:"reportID"`
	ReportName  string                      `json:"reportName"`
	StartDate   time.Time                   `json:"startDate"`
	EndDate     time.Time                   `json:"endDate"`
	DataSetType datastore.ReportDataSetType `json:"dataSetType"`
	ReportData  []*ReportOutputData         `json:"reportDataSet"`
}

type ReportOutputData struct {
	Expense         int64                `json:"expense,omitempty"`
	Income          int64                `json:"income,omitempty"`
	NetTransactions []*TransactionLedger `json:"netTransactions,omitempty"`
}

func ReportOutputToRespReportOutput(rpt *models.ReportOutput) *ReportOutput {
	reportDataSet := ConvertReportDataToRespReportData(rpt.ReportData)
	myReport := &ReportOutput{
		ReportID:    rpt.ReportID,
		ReportName:  rpt.ReportName,
		StartDate:   rpt.StartDate,
		EndDate:     rpt.EndDate,
		DataSetType: rpt.DataSetType,
		ReportData:  reportDataSet}

	return myReport
}

// ConvertReportDataToRespReportData converts []*models.ReportData to ReportData
func ConvertReportDataToRespReportData(dcSet []*models.ReportOutputData) []*ReportOutputData {
	var mset = make([]*ReportOutputData, len(dcSet))

	for idx := range dcSet {

		var netTrans = make([]*TransactionLedger, len(dcSet[idx].NetTransactions))
		for kdx := range dcSet[idx].NetTransactions {
			netTrans[kdx] = ConvertTransactionLedgerToRespTransactionLeger(dcSet[idx].NetTransactions[kdx])
		}
		myDS := ReportOutputData{
			Expense:         dcSet[idx].Expense,
			Income:          dcSet[idx].Income,
			NetTransactions: netTrans,
		}
		mset[idx] = &myDS
	}

	return mset
}
