package web

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
)

// ReportsController is the controller struct for reports
type ReportsController struct {
	DataStores *datastore.Datastores
}

// NewReportsController instantiates a new ReportsController struct
func NewReportsController(ds *datastore.Datastores) *ReportsController {
	return &ReportsController{
		DataStores: ds,
	}
}

// GET /reports
func (rc *ReportsController) ReportList(_ context.Context) ([]*models.Report, error) {
	reports, err := models.RetrieveReports(rc.DataStores)
	if err != nil {
		if errors.Is(err, models.ErrNoReports) {
			return nil, nil
		}

		return nil, fmt.Errorf("models.RetrieveReports:%w", err)
	}

	return reports, nil
}

// GET /reports/:reportID
func (rc *ReportsController) GetReportByID(_ context.Context, reportID uint64) (*models.Report, error) {
	report, err := models.RetrieveReportByID(rc.DataStores, reportID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveReportByID:%w", err)
	}

	return report, nil
}

// GET /reports/:reportID/run
func (rc *ReportsController) RunReport(_ context.Context, reportID uint64,
	startDate time.Time, endDate time.Time, accounts []uint64) (*models.ReportOutput, error) {
	report, err := models.RetrieveReportByID(rc.DataStores, reportID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveReportByID:%w", err)
	}

	reportOutPut, err := report.Run(rc.DataStores, startDate, endDate, accounts)
	if err != nil {
		return nil, fmt.Errorf("report.Run:%w", err)
	}

	return reportOutPut, nil
}

// POST /reports
func (rc *ReportsController) CreateReport(_ context.Context, report *models.Report) (*models.Report, error) {
	err := report.Store(rc.DataStores)
	if err != nil {
		return nil, fmt.Errorf("report.Store:%w", err)
	}

	return report, nil
}

// PUT /reports/{reportID}
func (rc *ReportsController) UpdateReport(_ context.Context, report *models.Report) (*models.Report, error) {
	err := report.Update(rc.DataStores)
	if err != nil {
		return nil, fmt.Errorf("report.Update:%w", err)
	}

	return report, nil
}

// DELETE /reports/{reportID}
func (rc *ReportsController) DeleteReport(_ context.Context, reportID uint64) (*models.Report,
	error) {
	myReport, err := models.RetrieveReportByID(rc.DataStores, reportID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveReportByID:%w", err)
	}

	err = myReport.Delete(rc.DataStores)
	if err != nil {
		return nil, fmt.Errorf("myReport.Delete:%w", err)
	}

	return myReport, nil
}
