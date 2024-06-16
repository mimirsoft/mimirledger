package web

import (
	"context"
	"github.com/mimirsoft/mimirledger/api/datastore"
)

// HealthController is the controller struct for the health check endpoint
type HealthController struct {
	DataStores *datastore.Datastores
}

// GET /api/health
// HEAD /api/health
func (healthController *HealthController) HealthCheck(_ context.Context) error {
	return nil
}

// NewHealthController instantiates a new HealthController struct
func NewHealthController(ds *datastore.Datastores) *HealthController {
	return &HealthController{
		DataStores: ds,
	}
}
