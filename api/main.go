package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok2"))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthOK"))
	})
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
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
