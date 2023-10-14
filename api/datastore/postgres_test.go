package datastore

import (
	"github.com/mimirsoft/mimirledger/api/cfg"
	"os"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var TestPostgresConfig PostgresConfig

func TestMain(m *testing.M) {
	cfg.LoadEnv()
	myConfig := LoadPostgresConfigFromEnv()
	TestPostgresConfig = myConfig
	result := m.Run()
	os.Exit(result)
}

func TestCommentCreateReply(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	myClient, err := NewClient(&TestPostgresConfig)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	err = myClient.Ping()
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
