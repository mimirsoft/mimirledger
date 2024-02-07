package datastore

import (
	"github.com/jmoiron/sqlx"
	"github.com/mimirsoft/mimirledger/api/cfg"
	"os"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var TestPostgresConfig PostgresConfig

var TestPostgresClient *sqlx.DB

func TestMain(m *testing.M) {
	cfg.LoadEnv()
	myConfig := LoadPostgresConfigFromEnv()
	TestPostgresConfig = myConfig
	myClient, err := NewClient(&TestPostgresConfig)
	if err != nil {
		panic(err)
	}
	TestPostgresClient = myClient
	result := m.Run()
	os.Exit(result)
}

func TestPostgresClientTestAndPing(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	myClient, err := NewClient(&TestPostgresConfig)
	g.Expect(TestPostgresConfig.DBName).To(gomega.Equal("mimirledgertest"))
	g.Expect(err).NotTo(gomega.HaveOccurred())
	err = myClient.Ping()
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
