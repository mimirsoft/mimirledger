package models

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/mimirsoft/mimirledger/api/cfg"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"os"
	"testing"
)

var dbClient *sqlx.DB
var testDS *datastore.Datastores

func TestMain(m *testing.M) {
	cfg.LoadEnv()
	myConfig := datastore.LoadPostgresConfigFromEnv()
	myClient, err := datastore.NewClient(&myConfig)
	if err != nil {
		panic(err)
	}
	testDS = datastore.NewDatastores(myClient)
	dbClient = myClient

	result := m.Run()
	os.Exit(result)
}

func TestAccount_StoreInvalid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	// failed due to no empty AccountNmae
	a1 := Account{}
	err := a1.Store(testDS)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, errAccountNameEmptyString)).To(gomega.BeTrue())

	// failed due to no transaction_account_sign_type
	a2 := Account{AccountName: "MyBank"}
	err = a2.Store(testDS)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring("accountType is not valid, cannot determine AccountSign"))

}
func TestAccount_Store(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	a1 := Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err := a1.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myAcct, err := getAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct).NotTo(gomega.BeNil())
	g.Expect(myAcct.AccountID).To(gomega.Equal(a1.AccountID))

}

// test make account and children and grand children
func TestAccount_FindSpotInTree(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	acctSet, err := RetrieveAccounts(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acctSet).To(gomega.HaveLen(0))

	a1 := Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a1.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a1.AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(a1.AccountRight).To(gomega.Equal(uint64(2)))
	// child of a1, so should be after a1.AccountLeft
	afterValue, err := findSpotInTree(testDS, a1.AccountID, "MyBank_subacct")
	g.Expect(afterValue).To(gomega.Equal(uint64(1)))

	a2 := Account{AccountName: "OtherBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a2.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a1.AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(a1.AccountRight).To(gomega.Equal(uint64(2)))

	// top level account, so after a1.AccountRight, which is 2
	afterValue, err = findSpotInTree(testDS, 0, "MyBank2")
	g.Expect(afterValue).To(gomega.Equal(uint64(2)))
}

// test make account and children and grand children
func TestAccount_StoreParentAndChildren(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	acctSet, err := RetrieveAccounts(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acctSet).To(gomega.HaveLen(0))

	a1 := Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a1.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a1.AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(a1.AccountRight).To(gomega.Equal(uint64(2)))

	a2 := Account{AccountName: "MyBank_subacct", AccountParent: a1.AccountID,
		AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a2.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myAcct, err := getAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct).NotTo(gomega.BeNil())
	g.Expect(myAcct.AccountID).To(gomega.Equal(a1.AccountID))
	g.Expect(myAcct.AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(myAcct.AccountRight).To(gomega.Equal(uint64(4)))

	g.Expect(a2.AccountLeft).To(gomega.Equal(uint64(2)))
	g.Expect(a2.AccountRight).To(gomega.Equal(uint64(3)))

	a3 := Account{AccountName: "MyBank_sub_subacct", AccountParent: a2.AccountID,
		AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a3.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a3.AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(a3.AccountRight).To(gomega.Equal(uint64(4)))

	children, err := findDirectChildren(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(children).To(gomega.HaveLen(1))
	g.Expect(children[0].AccountID).To(gomega.Equal(a2.AccountID))
	g.Expect(children[0].AccountLeft).To(gomega.Equal(uint64(2)))
	g.Expect(children[0].AccountRight).To(gomega.Equal(uint64(5)))

	children, err = findDirectChildren(testDS, a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(children).To(gomega.HaveLen(1))
	g.Expect(children[0].AccountID).To(gomega.Equal(a3.AccountID))
	g.Expect(children[0].AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(children[0].AccountRight).To(gomega.Equal(uint64(4)))

	children, err = findAllChildren(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(children).To(gomega.HaveLen(2))
	g.Expect(children[0].AccountID).To(gomega.Equal(a2.AccountID))
	g.Expect(children[0].AccountLeft).To(gomega.Equal(uint64(2)))
	g.Expect(children[0].AccountRight).To(gomega.Equal(uint64(5)))
	g.Expect(children[1].AccountID).To(gomega.Equal(a3.AccountID))
	g.Expect(children[1].AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(children[1].AccountRight).To(gomega.Equal(uint64(4)))
}

// test make account and children and grand children
func TestAccount_StoreParentAndChildrenMoved(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	acctSet, err := RetrieveAccounts(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acctSet).To(gomega.HaveLen(0))

	a1 := Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a1.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a2 := Account{AccountName: "MyBank_subacct", AccountParent: a1.AccountID,
		AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a2.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a3 := Account{AccountName: "MyBank_sub_subacct", AccountParent: a2.AccountID,
		AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a3.Store(testDS)
	g.Expect(a3.AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(a3.AccountRight).To(gomega.Equal(uint64(4)))

	a4 := Account{AccountName: "OtherBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a4.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// move a2-a3 under a4
	a2.AccountParent = a4.AccountID
	err = a2.Update(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myAcct, err := getAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct).NotTo(gomega.BeNil())
	g.Expect(myAcct.AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(myAcct.AccountRight).To(gomega.Equal(uint64(2)))

	myAcct, err = getAccountByID(testDS, a4.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct).NotTo(gomega.BeNil())
	g.Expect(myAcct.AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(myAcct.AccountRight).To(gomega.Equal(uint64(8)))

	myAcct, err = getAccountByID(testDS, a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct).NotTo(gomega.BeNil())
	g.Expect(myAcct.AccountLeft).To(gomega.Equal(uint64(4)))
	g.Expect(myAcct.AccountRight).To(gomega.Equal(uint64(7)))

	myAcct, err = getAccountByID(testDS, a3.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct).NotTo(gomega.BeNil())
	g.Expect(myAcct.AccountLeft).To(gomega.Equal(uint64(5)))
	g.Expect(myAcct.AccountRight).To(gomega.Equal(uint64(6)))

}

// test make account and children and grand children and move them around
func setupDB(g *gomega.WithT) {
	query := `delete  from transaction_accounts `
	_, err := dbClient.Exec(query)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}