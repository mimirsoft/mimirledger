package cfg

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mimirsoft/mimirledger/api/datastore"
)

const (
	pathSeparator = "/"
)

type Config struct {
	Postgres datastore.PostgresConfig
}

func LoadEnv() error {
	appEnv := os.Getenv("APP_ENV")
	appRoot := os.Getenv("APP_ROOT")

	err := godotenv.Load(
		envFile(appRoot, ".env."+appEnv), // APP_ENV specific file
	)
	if err != nil {
		return fmt.Errorf("godotenv.Load err: %w", err)
	}

	return nil
}

func envFile(appRoot string, file string) string {
	return strings.Join([]string{appRoot, file}, pathSeparator)
}
