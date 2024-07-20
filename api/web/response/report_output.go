package response

import (
	"time"

	"github.com/mimirsoft/mimirledger/api/models"
)

type ReportOutput struct {
	ReportName string    `db:"reportName"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
}

func ReportOutputToRespReportOutput(rpt *models.ReportOutput) *ReportOutput {
	myReport := &ReportOutput{ReportName: rpt.ReportName,
		StartDate: rpt.StartDate,
		EndDate:   rpt.EndDate}

	return myReport
}
