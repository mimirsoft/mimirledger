package models

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/mimirsoft/mimirledger/api/cfg"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"os"
	"testing"
	"time"
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

	a2 := Account{AccountName: "MyBank", AccountSign: datastore.AccountSignCredit, AccountType: datastore.AccountTypeIncome}
	err = a2.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myAcct, err = getAccountByID(testDS, a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct).NotTo(gomega.BeNil())
	g.Expect(myAcct.AccountID).To(gomega.Equal(a2.AccountID))
	g.Expect(myAcct.AccountSign).To(gomega.Equal(datastore.AccountSignCredit))
}

func TestAccount_UpdateBalanceForAccountID(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	a1 := Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err := a1.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	transStore := createTransactionStore()
	myTrans := datastore.Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	dcStore := createTransactionDCStore()
	// succeed
	myDC := datastore.TransactionDebitCredit{DebitOrCredit: datastore.AccountSignDebit, TransactionDCAmount: 54500,
		TransactionID: myTrans.TransactionID, AccountID: a1.AccountID}
	err = dcStore.Store(&myDC)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = UpdateSubtotalForAccountID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acct, err := RetrieveAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(int64(54500)))
	g.Expect(acct.AccountBalance).To(gomega.Equal(int64(0)))

	err = UpdateBalanceForAccountID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acct, err = RetrieveAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(int64(54500)))
	g.Expect(acct.AccountBalance).To(gomega.Equal(int64(54500)))
}
func TestAccount_updateSubtotal(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	a1 := Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err := a1.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	transStore := createTransactionStore()
	myTrans := datastore.Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	dcStore := createTransactionDCStore()
	// succeed
	myDC := datastore.TransactionDebitCredit{DebitOrCredit: datastore.AccountSignDebit, TransactionDCAmount: 54500,
		TransactionID: myTrans.TransactionID, AccountID: a1.AccountID}
	err = dcStore.Store(&myDC)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = a1.updateSubtotal(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a1.AccountSubtotal).To(gomega.Equal(int64(54500)))
	g.Expect(a1.AccountBalance).To(gomega.Equal(int64(0)))

	acct, err := RetrieveAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(int64(54500)))
	g.Expect(acct.AccountBalance).To(gomega.Equal(int64(0)))
}

func TestAccount_UpdateBalance(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	a1 := Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err := a1.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	transStore := createTransactionStore()
	myTrans := datastore.Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	dcStore := createTransactionDCStore()
	// succeed
	myDC := datastore.TransactionDebitCredit{DebitOrCredit: datastore.AccountSignDebit, TransactionDCAmount: 54500,
		TransactionID: myTrans.TransactionID, AccountID: a1.AccountID}
	err = dcStore.Store(&myDC)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = a1.updateSubtotal(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a1.AccountSubtotal).To(gomega.Equal(int64(54500)))
	g.Expect(a1.AccountBalance).To(gomega.Equal(int64(0)))

	err = a1.updateBalance(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a1.AccountSubtotal).To(gomega.Equal(int64(54500)))
	g.Expect(a1.AccountBalance).To(gomega.Equal(int64(54500)))

	acct, err := RetrieveAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(int64(54500)))
	g.Expect(acct.AccountBalance).To(gomega.Equal(int64(54500)))
}

func TestAccount_UpdateReconciledDate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	a1 := Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err := a1.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acct, err := RetrieveAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(int64(0)))
	g.Expect(acct.AccountBalance).To(gomega.Equal(int64(0)))
	g.Expect(acct.AccountReconcileDate).To(gomega.Equal(sql.NullTime{Time: time.Time{}, Valid: false}))

	oldDate1, err := time.Parse("2006-01-02", "2016-07-08")
	a1.AccountReconcileDate = sql.NullTime{Time: oldDate1, Valid: true}

	err = a1.UpdateReconciledDate(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acct, err = RetrieveAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(int64(0)))
	g.Expect(acct.AccountBalance).To(gomega.Equal(int64(0)))
	g.Expect(acct.AccountReconcileDate).To(gomega.Equal(sql.NullTime{Time: oldDate1, Valid: true}))
}

func TestAccount_RetrieveByIDInvalid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	setupDB(g)
	acctSet, err := RetrieveAccountByID(testDS, 5555)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, ErrAccountNotFound)).To(gomega.BeTrue())
	g.Expect(acctSet).To(gomega.BeNil())
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
	g.Expect(a2.AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(a2.AccountRight).To(gomega.Equal(uint64(4)))

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
	g.Expect(a1.AccountName).To(gomega.Equal("MyBank"))
	g.Expect(a1.AccountFullName).To(gomega.Equal("MyBank"))

	a2 := Account{AccountName: "MyBank_subacct", AccountParent: a1.AccountID,
		AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a2.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a2.AccountName).To(gomega.Equal("MyBank_subacct"))
	g.Expect(a2.AccountFullName).To(gomega.Equal("MyBank:MyBank_subacct"))

	myAcct, err := getAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct).NotTo(gomega.BeNil())
	g.Expect(myAcct.AccountID).To(gomega.Equal(a1.AccountID))
	g.Expect(myAcct.AccountName).To(gomega.Equal("MyBank"))
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
	g.Expect(a3.AccountName).To(gomega.Equal("MyBank_sub_subacct"))
	g.Expect(a3.AccountFullName).To(gomega.Equal("MyBank:MyBank_subacct:MyBank_sub_subacct"))

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
func TestAccount_CloseSpotInTree(t *testing.T) {
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
	g.Expect(a1.AccountName).To(gomega.Equal("MyBank"))
	g.Expect(a1.AccountFullName).To(gomega.Equal("MyBank"))

	a2 := Account{AccountName: "MyBank_subacct", AccountParent: a1.AccountID,
		AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a2.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a2.AccountName).To(gomega.Equal("MyBank_subacct"))
	g.Expect(a2.AccountFullName).To(gomega.Equal("MyBank:MyBank_subacct"))

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
	g.Expect(a3.AccountName).To(gomega.Equal("MyBank_sub_subacct"))
	g.Expect(a3.AccountFullName).To(gomega.Equal("MyBank:MyBank_subacct:MyBank_sub_subacct"))

	a4 := Account{AccountName: "OtherBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a4.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

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

	allChildrenA2, err := findAllChildren(testDS, a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(allChildrenA2).To(gomega.HaveLen(1))
	g.Expect(allChildrenA2[0].AccountID).To(gomega.Equal(a3.AccountID))

	updateAcct2, err := getAccountByID(testDS, a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updateAcct2.AccountLeft).To(gomega.Equal(uint64(2)))
	g.Expect(updateAcct2.AccountRight).To(gomega.Equal(uint64(5)))
	oldA2AccountLeft := updateAcct2.AccountLeft

	spread := (updateAcct2.AccountRight - updateAcct2.AccountLeft) + 1
	g.Expect(spread).To(gomega.Equal(uint64(4)))

	err = closeSpotInTree(testDS, a2.AccountRight, spread)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myAcct, err = getAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct).NotTo(gomega.BeNil())
	g.Expect(myAcct.AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(myAcct.AccountRight).To(gomega.Equal(uint64(2)))

	myAcct4, err := getAccountByID(testDS, a4.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct4).NotTo(gomega.BeNil())
	g.Expect(myAcct4.AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(myAcct4.AccountRight).To(gomega.Equal(uint64(4)))

	//Find the new spot in the tree.
	afterValue, err := findSpotInTree(testDS, a4.AccountID, a2.AccountName)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(afterValue).To(gomega.Equal(uint64(3)))

	err = openSpotInTree(testDS, afterValue, spread)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myAcct4, err = getAccountByID(testDS, a4.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct4).NotTo(gomega.BeNil())
	g.Expect(myAcct4.AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(myAcct4.AccountRight).To(gomega.Equal(uint64(8)))

	// aftervalue should be 3, that is the a4 account
	// old a2.AccountLeft should be 2,
	shift := afterValue - oldA2AccountLeft + 1
	// shift should be 2, aft a2.Account left is going from 2 to 4
	g.Expect(shift).To(gomega.Equal(uint64(2)))
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

	myAcct, err = getAccountByID(testDS, a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct).NotTo(gomega.BeNil())
	g.Expect(myAcct.AccountParent).To(gomega.Equal(a4.AccountID))
	g.Expect(myAcct.AccountLeft).To(gomega.Equal(uint64(4)))
	g.Expect(myAcct.AccountRight).To(gomega.Equal(uint64(7)))
	g.Expect(myAcct.AccountName).To(gomega.Equal("MyBank_subacct"))
	g.Expect(myAcct.AccountFullName).To(gomega.Equal("OtherBank:MyBank_subacct"))

	myAcct, err = getAccountByID(testDS, a3.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct).NotTo(gomega.BeNil())
	g.Expect(myAcct.AccountParent).To(gomega.Equal(a2.AccountID))
	g.Expect(myAcct.AccountLeft).To(gomega.Equal(uint64(5)))
	g.Expect(myAcct.AccountRight).To(gomega.Equal(uint64(6)))
	g.Expect(myAcct.AccountName).To(gomega.Equal("MyBank_sub_subacct"))
	g.Expect(myAcct.AccountFullName).To(gomega.Equal("OtherBank:MyBank_subacct:MyBank_sub_subacct"))

	myAcct, err = getAccountByID(testDS, a4.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myAcct).NotTo(gomega.BeNil())
	g.Expect(myAcct.AccountName).To(gomega.Equal("OtherBank"))
	g.Expect(myAcct.AccountFullName).To(gomega.Equal("OtherBank"))
	g.Expect(myAcct.AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(myAcct.AccountRight).To(gomega.Equal(uint64(8)))
}

// test make account and children and grand children
func TestAccount_MoveChildWithBalances(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	acctSet, err := RetrieveAccounts(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acctSet).To(gomega.HaveLen(0))

	incomeAcct := Account{AccountName: "Income",
		AccountSign: datastore.AccountSignCredit, AccountType: datastore.AccountTypeIncome}
	err = incomeAcct.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a1 := Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a1.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a2 := Account{AccountName: "MyBank_subacct", AccountParent: a1.AccountID,
		AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a2.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a4 := Account{AccountName: "OtherBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err = a4.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// add a transaction from income to MyBank_subacct
	txn := Transaction{TransactionCore: TransactionCore{TransactionComment: "woot"},
		DebitCreditSet: []*TransactionDebitCredit{
			&TransactionDebitCredit{AccountID: a2.AccountID,
				DebitOrCredit:       datastore.AccountSignDebit,
				TransactionDCAmount: 10000},
			&TransactionDebitCredit{AccountID: incomeAcct.AccountID,
				DebitOrCredit:       datastore.AccountSignCredit,
				TransactionDCAmount: 10000},
		},
	}
	err = txn.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// retreive accounts and check balances and subtotals
	updatedA2, err := RetrieveAccountByID(testDS, a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updatedA2.AccountSubtotal).To(gomega.Equal(int64(10000)))
	g.Expect(updatedA2.AccountBalance).To(gomega.Equal(int64(10000)))

	updatedA1, err := RetrieveAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updatedA1.AccountSubtotal).To(gomega.Equal(int64(0)))
	g.Expect(updatedA1.AccountBalance).To(gomega.Equal(int64(10000)))

	updatedA4, err := RetrieveAccountByID(testDS, a4.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updatedA4.AccountSubtotal).To(gomega.Equal(int64(0)))
	g.Expect(updatedA4.AccountBalance).To(gomega.Equal(int64(0)))
	updatedIncomeAcct, err := RetrieveAccountByID(testDS, incomeAcct.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updatedIncomeAcct.AccountSubtotal).To(gomega.Equal(int64(10000)))
	g.Expect(updatedIncomeAcct.AccountBalance).To(gomega.Equal(int64(10000)))

	// move a2 under a4
	a2.AccountParent = a4.AccountID
	err = a2.Update(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// retreive accounts and check balances and subtotals
	updatedA1, err = RetrieveAccountByID(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updatedA1.AccountSubtotal).To(gomega.Equal(int64(0)))
	g.Expect(updatedA1.AccountBalance).To(gomega.Equal(int64(0)))
	updatedA2, err = RetrieveAccountByID(testDS, a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updatedA2.AccountSubtotal).To(gomega.Equal(int64(10000)))
	g.Expect(updatedA2.AccountBalance).To(gomega.Equal(int64(10000)))
	updatedA4, err = RetrieveAccountByID(testDS, a4.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updatedA4.AccountSubtotal).To(gomega.Equal(int64(0)))
	g.Expect(updatedA4.AccountBalance).To(gomega.Equal(int64(10000)))
	updatedIncomeAcct, err = RetrieveAccountByID(testDS, incomeAcct.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updatedIncomeAcct.AccountSubtotal).To(gomega.Equal(int64(10000)))
	g.Expect(updatedIncomeAcct.AccountBalance).To(gomega.Equal(int64(10000)))

}

// test make account and children and grand children
func TestAccount_RetrieveAccountFullName(t *testing.T) {
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

	fullName, err := retrieveAccountFullName(testDS, a3.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(fullName).To(gomega.Equal("MyBank:MyBank_subacct:MyBank_sub_subacct"))

	fullName, err = retrieveAccountFullName(testDS, a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(fullName).To(gomega.Equal("MyBank:MyBank_subacct"))

	fullName, err = retrieveAccountFullName(testDS, a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(fullName).To(gomega.Equal("MyBank"))
}

// test make account and children and grand children and move them around
func setupDB(g *gomega.WithT) {
	query := `delete from transaction_debit_credit `
	_, err := dbClient.Exec(query)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	query = `delete from transaction_accounts `
	_, err = dbClient.Exec(query)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func createTransactionStore() datastore.TransactionStore {
	return datastore.TransactionStore{
		Client: dbClient,
	}
}

func createTransactionDCStore() datastore.TransactionDebitCreditStore {
	return datastore.TransactionDebitCreditStore{
		Client: dbClient,
	}
}
