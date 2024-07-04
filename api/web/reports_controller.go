package web

import (
	"context"
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
		return nil, fmt.Errorf("models.RetrieveReports:%w", err)
	}

	return reports, nil
}

// POST /reports
func (ac *ReportsController) CreateReport(_ context.Context, report *models.Report) (*models.Report, error) {
	err := report.Store(ac.DataStores)
	if err != nil {
		return nil, fmt.Errorf("report.Store:%w", err)
	}

	return report, nil
}
