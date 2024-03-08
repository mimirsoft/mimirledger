package datastore

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func createAccountStore() AccountStore {
	return AccountStore{
		Client: TestPostgresClient,
	}
}
func TestAccountStoreInvalid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	aStore := createAccountStore()

	// failed due to no transaction_account_sign_type
	a1 := Account{}
	err := aStore.Store(&a1)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring(`invalid input value for enum transaction_account_sign_type`))

	// failed due to no transaction_account_sign_type
	a2 := Account{AccountSign: AccountSignCredit}
	err = aStore.Store(&a2)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring(`invalid input value for enum transaction_account_type`))
}

func TestAccountStoreValid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	aStore := createAccountStore()
	a3 := Account{AccountName: "my bank", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 1, AccountRight: 2}
	err := aStore.Store(&a3)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a3.AccountID).To(gomega.Equal(uint64(1)))

	acctSet, err := aStore.GetAccounts()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acctSet).To(gomega.HaveLen(1))
	g.Expect(acctSet[0].AccountName).To(gomega.Equal(a3.AccountName))
	g.Expect(acctSet[0].AccountFullName).To(gomega.Equal("BankAccounts:MyBank"))
	g.Expect(acctSet[0].AccountID).To(gomega.Equal(uint64(1)))
	g.Expect(acctSet[0].AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(acctSet[0].AccountRight).To(gomega.Equal(uint64(2)))
	g.Expect(acctSet[0].AccountBalance).To(gomega.Equal(uint64(3000)))
	g.Expect(acctSet[0].AccountSubtotal).To(gomega.Equal(uint64(2000)))
	g.Expect(acctSet[0].AccountDecimals).To(gomega.Equal(uint64(2)))
}

func TestAccountStoreGetByID(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	aStore := createAccountStore()
	a3 := Account{AccountName: "otherBank2", AccountFullName: "BankAccounts:otherBank2",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 3, AccountRight: 4}
	err := aStore.Store(&a3)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a3.AccountID).To(gomega.Equal(uint64(2)))

	acct, err := aStore.GetAccountByID(2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acct.AccountName).To(gomega.Equal("otherBank2"))
	g.Expect(acct.AccountFullName).To(gomega.Equal("BankAccounts:otherBank2"))
	g.Expect(acct.AccountID).To(gomega.Equal(uint64(2)))
	g.Expect(acct.AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(acct.AccountRight).To(gomega.Equal(uint64(4)))
	g.Expect(acct.AccountBalance).To(gomega.Equal(uint64(3000)))
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(uint64(2000)))
	g.Expect(acct.AccountDecimals).To(gomega.Equal(uint64(2)))
}
func TestAccountStore_OpenSpotInTree(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	aStore := createAccountStore()
	a3 := Account{AccountName: "otherBank2", AccountFullName: "BankAccounts:otherBank2",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 1, AccountRight: 2}
	err := aStore.Store(&a3)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = aStore.OpenSpotInTree(1, 2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	acct, err := aStore.GetAccountByID(a3.AccountID)
	g.Expect(acct.AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(acct.AccountRight).To(gomega.Equal(uint64(4)))
}

func TestAccountStoreAndUpdate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	aStore := createAccountStore()
	a1 := Account{AccountName: "my bank", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 1, AccountRight: 2}
	err := aStore.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acct, err := aStore.GetAccountByID(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acct.AccountName = "UPDATED NAME"
	acct.AccountFullName = "BankAccounts:UPDATED NAME"
	acct.AccountLeft = 5
	acct.AccountRight = 6
	acct.AccountBalance = 4000
	acct.AccountSubtotal = 3000
	err = aStore.Update(acct)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acct, err = aStore.GetAccountByID(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acct.AccountName).To(gomega.Equal("UPDATED NAME"))
	g.Expect(acct.AccountFullName).To(gomega.Equal("BankAccounts:UPDATED NAME"))
	g.Expect(acct.AccountID).To(gomega.Equal(a1.AccountID))
	g.Expect(acct.AccountLeft).To(gomega.Equal(uint64(5)))
	g.Expect(acct.AccountRight).To(gomega.Equal(uint64(6)))
	g.Expect(acct.AccountBalance).To(gomega.Equal(uint64(4000)))
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(uint64(3000)))
	g.Expect(acct.AccountDecimals).To(gomega.Equal(uint64(2)))
}

func TestAccountStore_GetDirectChildren(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	aStore := createAccountStore()
	a1 := Account{AccountName: "my bank", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 1, AccountRight: 6}
	err := aStore.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a2 := Account{AccountName: "subact1", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 2, AccountRight: 5, AccountParent: a1.AccountID}
	err = aStore.Store(&a2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a3 := Account{AccountName: "subact2", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 3, AccountRight: 4, AccountParent: a2.AccountID}
	err = aStore.Store(&a3)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	children, err := aStore.GetDirectChildren(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(children).To(gomega.HaveLen(1))
	g.Expect(children[0].AccountName).To(gomega.Equal(a2.AccountName))
	g.Expect(children[0].AccountID).To(gomega.Equal(a2.AccountID))
	g.Expect(children[0].AccountLeft).To(gomega.Equal(uint64(2)))
	g.Expect(children[0].AccountRight).To(gomega.Equal(uint64(5)))

	children, err = aStore.GetDirectChildren(a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(children).To(gomega.HaveLen(1))
	g.Expect(children[0].AccountName).To(gomega.Equal(a3.AccountName))
	g.Expect(children[0].AccountID).To(gomega.Equal(a3.AccountID))
	g.Expect(children[0].AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(children[0].AccountRight).To(gomega.Equal(uint64(4)))
}

func TestAccountStore_GetAllChildren(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	aStore := createAccountStore()
	a1 := Account{AccountName: "my bank", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 1, AccountRight: 6}
	err := aStore.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a2 := Account{AccountName: "subact1", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 2, AccountRight: 5, AccountParent: a1.AccountID}
	err = aStore.Store(&a2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a3 := Account{AccountName: "subact2", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 3, AccountRight: 4, AccountParent: a2.AccountID}
	err = aStore.Store(&a3)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	children, err := aStore.GetAllChildren(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(children).To(gomega.HaveLen(2))
	g.Expect(children[0].AccountName).To(gomega.Equal(a2.AccountName))
	g.Expect(children[0].AccountID).To(gomega.Equal(a2.AccountID))
	g.Expect(children[0].AccountLeft).To(gomega.Equal(uint64(2)))
	g.Expect(children[0].AccountRight).To(gomega.Equal(uint64(5)))
	g.Expect(children[1].AccountName).To(gomega.Equal(a3.AccountName))
	g.Expect(children[1].AccountID).To(gomega.Equal(a3.AccountID))
	g.Expect(children[1].AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(children[1].AccountRight).To(gomega.Equal(uint64(4)))

	children, err = aStore.GetAllChildren(a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(children).To(gomega.HaveLen(1))
	g.Expect(children[0].AccountName).To(gomega.Equal(a3.AccountName))
	g.Expect(children[0].AccountID).To(gomega.Equal(a3.AccountID))
	g.Expect(children[0].AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(children[0].AccountRight).To(gomega.Equal(uint64(4)))
}

func setupDB(g *gomega.WithT) {
	query := `delete  from transaction_accounts `
	_, err := TestPostgresClient.Exec(query)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
