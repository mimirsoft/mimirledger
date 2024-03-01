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

func TestAccountStoreOpenSpotInTree(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	aStore := createAccountStore()
	err := aStore.OpenSpotInTree(1, 1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
