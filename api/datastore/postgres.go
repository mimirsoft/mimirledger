package datastore

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx.
	"github.com/jmoiron/sqlx"
)

const (
	defaultSSLMode = "disable"
)

type Datastores struct {
	postgresClient     *sqlx.DB
	accountStore       AccountStore
	transactionStore   TransactionStore
	transactionDCStore TransactionDebitCreditStore
	reportStore        ReportStore
}

// AccountStore is the way to access the AccountStore.
func (ds *Datastores) AccountStore() AccountStore {
	return ds.accountStore
}

// TransactionStore is the way to access the TransactionStore.
func (ds *Datastores) TransactionStore() TransactionStore {
	return ds.transactionStore
}

// TransactionDebitCreditStore is the way to access the TransactionDebitCreditStore.
func (ds *Datastores) TransactionDebitCreditStore() TransactionDebitCreditStore {
	return ds.transactionDCStore
}

// ReportStore is the way to access the ReportStore.
func (ds *Datastores) ReportStore() ReportStore {
	return ds.reportStore
}

// PGClient is the way to access the Postgres Client
func (ds *Datastores) PGClient() *sqlx.DB {
	return ds.postgresClient
}

func NewDatastores(conn *sqlx.DB) *Datastores {
	return &Datastores{postgresClient: conn,
		accountStore: AccountStore{
			Client: conn,
		},
		reportStore: ReportStore{
			Client: conn,
		},
		transactionStore: TransactionStore{
			Client: conn,
		},
		transactionDCStore: TransactionDebitCreditStore{
			Client: conn,
		},
	}
}

type PostgresConfig struct {
	Host, Username, Password, DBName string
	Port, MaxConnLifetime            int
	DisableSSL                       bool
}

func NewClient(config *PostgresConfig) (*sqlx.DB, error) {
	sslMode := defaultSSLMode
	if config.DisableSSL {
		sslMode = defaultSSLMode
	}

	conn, err := sqlx.Open(
		"pgx",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", //nolint:nosprintfhostport
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.DBName,
			sslMode,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres: %w", err)
	}

	if config.MaxConnLifetime > 0 {
		conn.SetConnMaxLifetime(time.Second * time.Duration(config.MaxConnLifetime))
	}

	return conn, nil
}

func LoadPostgresConfigFromEnv() PostgresConfig {
	pgConfig := PostgresConfig{} //nolint:exhaustruct
	// Postgres
	if pgDBHost := os.Getenv("PG_DB_HOST"); pgDBHost != "" {
		pgConfig.Host = pgDBHost
	}

	if pgDBPort := os.Getenv("PG_DB_PORT"); pgDBPort != "" {
		portInt, _ := strconv.Atoi(pgDBPort)
		pgConfig.Port = portInt
	}

	if pgDBUser := os.Getenv("PG_DB_USER"); pgDBUser != "" {
		pgConfig.Username = pgDBUser
	}

	if pgDBPassword := os.Getenv("PG_DB_PASSWORD"); pgDBPassword != "" {
		pgConfig.Password = pgDBPassword
	}

	if pgDBName := os.Getenv("PG_DB_NAME"); pgDBName != "" {
		pgConfig.DBName = pgDBName
	}

	if pgMaxConnLifetimeSecs := os.Getenv("PG_MAX_CONN_LIFETIME"); pgMaxConnLifetimeSecs != "" {
		ltInt, _ := strconv.Atoi(pgMaxConnLifetimeSecs)
		pgConfig.MaxConnLifetime = ltInt
	}

	if pgDisableSSL := os.Getenv("PG_DISABLE_SSL"); pgDisableSSL != "" {
		pgConfig.DisableSSL, _ = strconv.ParseBool(pgDisableSSL)
	}

	return pgConfig
}
