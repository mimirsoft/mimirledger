package web

import (
	"context"
	"errors"
	"fmt"

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
func (ac *ReportsController) ReportList(_ context.Context) ([]*models.Report, error) {
	reports, err := models.RetrieveReports(ac.DataStores)
	if err != nil {
		if errors.Is(err, models.ErrNoReports) {
			return nil, nil
		}

		return nil, fmt.Errorf("models.RetrieveReports:%w", err)
	}

	return reports, nil
}

// GET /reports/:reportID
func (ac *ReportsController) GetReportByID(_ context.Context, reportID uint64) (*models.Report, error) {
	report, err := models.RetrieveReportByID(ac.DataStores, reportID)
	if err != nil {
		return nil, fmt.Errorf("models.RetrieveReportByID:%w", err)
	}

	return report, nil
}

// POST /reports
func (ac *ReportsController) CreateReport(_ context.Context, report *models.Report) (*models.Report, error) {
	err := report.Store(ac.DataStores)
	if err != nil {
		return nil, fmt.Errorf("report.Store:%w", err)
	}

	return report, nil
}
