package datastore

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx
)

const (
	defaultSSLMode = "disable"
)

type Datastores struct {
	postgresClient *sqlx.DB
}

// NewDatastores creates a struct of MixStores
func NewDatastores(conn *sqlx.DB) *Datastores {
	return &Datastores{postgresClient: conn}
}

type PostgresConfig struct {
	Host, Username, Password, DBName string
	Port, MaxConnLifetime            int
	DisableSSL                       bool
}

func NewClient(config *PostgresConfig) (conn *sqlx.DB, err error) {
	sslMode := defaultSSLMode
	if config.DisableSSL {
		sslMode = defaultSSLMode
	}
	conn, err = sqlx.Open(
		"pgx",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.DBName,
			sslMode,
		),
	)
	if config.MaxConnLifetime > 0 {
		conn.SetConnMaxLifetime(time.Second * time.Duration(config.MaxConnLifetime))
	}
	return
}

func LoadPostgresConfigFromEnv() PostgresConfig {
	c := PostgresConfig{}
	// Postgres
	if pgDBHost := os.Getenv("PG_DB_HOST"); pgDBHost != "" {
		c.Host = pgDBHost
	}
	if pgDBPort := os.Getenv("PG_DB_PORT"); pgDBPort != "" {
		portInt, _ := strconv.Atoi(pgDBPort)
		c.Port = portInt
	}
	if pgDBUser := os.Getenv("PG_DB_USER"); pgDBUser != "" {
		c.Username = pgDBUser
	}
	if pgDBPassword := os.Getenv("PG_DB_PASSWORD"); pgDBPassword != "" {
		c.Password = pgDBPassword
	}
	if pgDBName := os.Getenv("PG_DB_NAME"); pgDBName != "" {
		c.DBName = pgDBName
	}
	if pgMaxConnLifetimeSecs := os.Getenv("PG_MAX_CONN_LIFETIME"); pgMaxConnLifetimeSecs != "" {
		ltInt, _ := strconv.Atoi(pgMaxConnLifetimeSecs)
		c.MaxConnLifetime = ltInt
	}
	if pgDisableSSL := os.Getenv("PG_DISABLE_SSL"); pgDisableSSL != "" {
		c.DisableSSL, _ = strconv.ParseBool(pgDisableSSL)
	}
	return c
}
