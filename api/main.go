package main

import (
	"fmt"
	"github.com/mimirsoft/mimirledger/api/cfg"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := logger.Sugar()

	appConfig := LoadConfig()
	fmt.Printf("appConfig:%v \n", appConfig)
	fmt.Println("Hello, world.")
	myClient, err := datastore.NewClient(&appConfig.Postgres)
	if err != nil {
		sugar.Errorf("godotenv.Load err: %v", err)
	}
	err = myClient.Ping()
	if err != nil {
		sugar.Errorf("myClient.Ping() err: %v", err)
	}

}

type Config struct {
	Postgres datastore.PostgresConfig
}

func LoadConfig() Config {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := logger.Sugar()

	err = cfg.LoadEnv()
	if err != nil {
		sugar.Errorf("godotenv.Load err: %v", err)
	}
	postgresCfg := datastore.LoadPostgresConfigFromEnv()

	myCfg := Config{Postgres: postgresCfg}
	return myCfg
}
