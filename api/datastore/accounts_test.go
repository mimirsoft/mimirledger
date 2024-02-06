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
	a3 := Account{AccountName: "my bank", AccountSign: AccountSignDebit, AccountType: AccountTypeAsset}
	err := aStore.Store(&a3)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acctSet, err := aStore.GetAccounts()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acctSet).To(gomega.HaveLen(1))
	g.Expect(acctSet[0].AccountName).To(gomega.Equal(a3.AccountName))
}
