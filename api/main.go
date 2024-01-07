package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mimirsoft/mimirledger/api/cfg"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	appConfig := LoadConfig()
	fmt.Printf("appConfig:%v \n", appConfig)
	fmt.Println("Hello, world.")
	myClient, err := datastore.NewClient(&appConfig.Postgres)
	if err != nil {
		log.Error().Err(err).Msg("godotenv.Load")
	}
	err = myClient.Ping()
	if err != nil {
		log.Error().Err(err).Msg("myClient.Ping()")
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok2"))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthOK"))
	})
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	r.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"accounts":[{"name":"blah1"},{"name":"blah2"}]}`))
	})
	r.Get("/accounttypes", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"accountTypes":[{"name":"asset"},{"name":"liability"}]}`))
	})
	http.ListenAndServe(":3010", r)
}

type Config struct {
	Postgres datastore.PostgresConfig
}

func LoadConfig() Config {
	err := cfg.LoadEnv()
	if err != nil {
		log.Error().Err(err).Msg("cfg.LoadEnv()")
	}
	postgresCfg := datastore.LoadPostgresConfigFromEnv()

	myCfg := Config{Postgres: postgresCfg}
	return myCfg
}
