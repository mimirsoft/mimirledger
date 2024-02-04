package main

import (
	"github.com/mimirsoft/mimirledger/api/cfg"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

func main() {
	loggerOutput := zerolog.ConsoleWriter{Out: os.Stderr}
	logger := zerolog.New(loggerOutput)

	appConfig := LoadConfig()
	logger.Info().Msg("####Starting MimirLedger API Server###")
	myClient, err := datastore.NewClient(&appConfig.Postgres)
	if err != nil {
		log.Error().Err(err).Msg("godotenv.Load")
	}
	err = myClient.Ping()
	if err != nil {
		log.Error().Err(err).Msg("myCrlient.Ping()")
	}
	ds := datastore.NewDatastores(myClient)

	r := web.NewRouter(ds, logger)
	err = http.ListenAndServe(":3010", r)
	if err != nil {
		log.Error().Err(err).Msg("http.ListenAndServe")
	}

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
