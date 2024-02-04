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

func NewRouter(ds *datastore.Datastores, logger zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.RequestId)
	r.Use(middlewares.Logger(logger))
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
	r.Get("/accounts", NewRootHandler(Accounts(acctsController)).ServeHTTP)
	r.Get("/accounttypes", NewRootHandler(AccountTypes(acctsController)).ServeHTTP)
	return r
}
