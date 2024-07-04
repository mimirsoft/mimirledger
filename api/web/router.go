package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/middlewares"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const maxAgeSeconds = 300

func NewRouter(dStores *datastore.Datastores, logger *zerolog.Logger) *chi.Mux { //nolint:funlen
	r := chi.NewRouter() //nolint:varnamelen
	r.Use(middlewares.RequestID)

	if logger != nil {
		r.Use(middlewares.Logger(*logger))
	}

	r.Use(middleware.Recoverer)
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   false,
		MaxAge:             maxAgeSeconds, // Maximum value not ignored by any of major browsers
		OptionsPassthrough: false,
		Debug:              false,
		AllowOriginFunc:    nil,
	}))

	healthController := NewHealthController(dStores)
	accountsController := NewAccountsController(dStores)
	reportsController := NewReportsController(dStores)
	transController := NewTransactionsController(dStores)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("{ok}"))
		if err != nil {
			log.Error().Err(err).Msg("w.Write failed")
		}
	})
	r.Get("/health", NewRootHandler(HealthCheck(healthController)).ServeHTTP)
	r.Get("/hello", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("hello"))
		if err != nil {
			log.Error().Err(err).Msg("w.Write failed")
		}
	})
	r.Get("/accounts", NewRootHandler(GetAccounts(accountsController)).ServeHTTP)
	r.Post("/accounts", NewRootHandler(PostAccounts(accountsController)).ServeHTTP)
	r.Get("/accounts/{accountID}", NewRootHandler(GetAccount(accountsController)).ServeHTTP)
	r.Put("/accounts/{accountID}", NewRootHandler(PutAccountUpdate(accountsController)).ServeHTTP)
	r.Put("/accounts/{accountID}/reconciled", NewRootHandler(PutAccountUpdateReconciled(accountsController)).ServeHTTP)
	r.Get("/accounttypes", NewRootHandler(GetAccountTypes(accountsController)).ServeHTTP)
	r.Get("/reports", NewRootHandler(GetReports(reportsController)).ServeHTTP)
	r.Post("/reports", NewRootHandler(PostReports(reportsController)).ServeHTTP)

	r.Post("/transactions", NewRootHandler(PostTransactions(transController)).ServeHTTP)
	r.Get("/transactions/account/{accountID}", NewRootHandler(GetTransactionsOnAccount(transController)).ServeHTTP)
	r.Get("/transactions/account/{accountID}/unreconciled",
		NewRootHandler(GetUnreconciledTransactionsOnAccount(transController)).ServeHTTP)
	r.Get("/transactions/{transactionID}", NewRootHandler(GetTransaction(transController)).ServeHTTP)
	r.Put("/transactions/{transactionID}", NewRootHandler(PutTransactionUpdate(transController)).ServeHTTP)
	r.Put("/transactions/{transactionID}/reconciled",
		NewRootHandler(PutTransactionReconciled(transController)).ServeHTTP)
	r.Put("/transactions/{transactionID}/unreconciled",
		NewRootHandler(PutTransactionUnreconciled(transController)).ServeHTTP)
	r.Delete("/transactions/{transactionID}", NewRootHandler(DeleteTransaction(transController)).ServeHTTP)

	return r
}
