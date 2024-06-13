package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/middlewares"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"net/http"
)

func NewRouter(ds *datastore.Datastores, logger *zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.RequestId)
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
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	healthController := NewHealthController(ds)
	acctsController := NewAccountsController(ds)
	transController := NewTransactionsController(ds)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok2"))
		if err != nil {
			log.Error().Err(err).Msg("w.Write failed")
		}
	})
	r.Get("/health", NewRootHandler(HealthCheck(healthController)).ServeHTTP)
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello"))
		if err != nil {
			log.Error().Err(err).Msg("w.Write failed")
		}
	})
	r.Get("/accounts", NewRootHandler(GetAccounts(acctsController)).ServeHTTP)
	r.Post("/accounts", NewRootHandler(PostAccounts(acctsController)).ServeHTTP)
	r.Get("/accounts/{accountID}", NewRootHandler(GetAccount(acctsController)).ServeHTTP)
	r.Put("/accounts/{accountID}", NewRootHandler(PutAccountUpdate(acctsController)).ServeHTTP)
	r.Post("/transactions", NewRootHandler(PostTransactions(transController)).ServeHTTP)
	r.Get("/transactions/account/{accountID}", NewRootHandler(GetTransactionsOnAccount(transController)).ServeHTTP)
	r.Get("/transactions/account/{accountID}/unreconciled", NewRootHandler(GetUnreconciledTransactionsOnAccount(transController)).ServeHTTP)
	r.Get("/transactions/{transactionID}", NewRootHandler(GetTransaction(transController)).ServeHTTP)
	r.Put("/transactions/{transactionID}", NewRootHandler(PutTransactionUpdate(transController)).ServeHTTP)
	r.Put("/transactions/{transactionID}/reconciled", NewRootHandler(PutTransactionReconciled(transController)).ServeHTTP)
	r.Put("/transactions/{transactionID}/unreconciled", NewRootHandler(PutTransactionUnreconciled(transController)).ServeHTTP)
	r.Delete("/transactions/{transactionID}", NewRootHandler(DeleteTransaction(transController)).ServeHTTP)

	r.Get("/accounttypes", NewRootHandler(GetAccountTypes(acctsController)).ServeHTTP)
	return r
}
