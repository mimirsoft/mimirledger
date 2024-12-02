package main

import (
	"net/http"
	"os"
	"time"

	"github.com/mimirsoft/mimirledger/api/cfg"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const readHeaderTimeout = time.Second * 3

func main() {
	loggerOutput := zerolog.ConsoleWriter{Out: os.Stderr} //nolint:exhaustruct
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
	r := web.NewRouter(ds, &logger)

	server := &http.Server{ //nolint:exhaustruct
		Addr:              ":3010",
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           r,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Error().Err(err).Msg("server.ListenAndServe")
	}
}

func LoadConfig() cfg.Config {
	err := cfg.LoadEnv()
	if err != nil {
		log.Error().Err(err).Msg("cfg.LoadEnv()")
	}

	postgresCfg := datastore.LoadPostgresConfigFromEnv()

	myCfg := cfg.Config{Postgres: postgresCfg}

	return myCfg
}
