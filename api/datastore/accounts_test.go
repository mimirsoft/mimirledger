package datastore

import (
	"database/sql"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
	"time"
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

func TestAccount_StoreValid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	aStore := createAccountStore()
	a3 := Account{AccountName: "my bank", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 1, AccountRight: 2}
	err := aStore.Store(&a3)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a3.AccountID).NotTo(gomega.BeZero())

	acctSet, err := aStore.GetAccounts()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acctSet).To(gomega.HaveLen(1))
	g.Expect(acctSet[0].AccountName).To(gomega.Equal(a3.AccountName))
	g.Expect(acctSet[0].AccountFullName).To(gomega.Equal("BankAccounts:MyBank"))
	g.Expect(acctSet[0].AccountID).To(gomega.Equal(a3.AccountID))
	g.Expect(acctSet[0].AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(acctSet[0].AccountRight).To(gomega.Equal(uint64(2)))
	// regular store call cannot update AccountBalance or AccountSubtotal
	g.Expect(acctSet[0].AccountBalance).To(gomega.Equal(int64(0)))
	g.Expect(acctSet[0].AccountSubtotal).To(gomega.Equal(int64(0)))
	g.Expect(acctSet[0].AccountDecimals).To(gomega.Equal(uint64(2)))
}

func TestAccount_StoreGetByID(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)
	aStore := createAccountStore()
	a3 := Account{AccountName: "otherBank2", AccountFullName: "BankAccounts:otherBank2",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 3, AccountRight: 4}
	err := aStore.Store(&a3)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acct, err := aStore.GetAccountByID(a3.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acct.AccountName).To(gomega.Equal("otherBank2"))
	g.Expect(acct.AccountFullName).To(gomega.Equal("BankAccounts:otherBank2"))
	g.Expect(acct.AccountID).To(gomega.Equal(a3.AccountID))
	g.Expect(acct.AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(acct.AccountRight).To(gomega.Equal(uint64(4)))
	g.Expect(acct.AccountBalance).To(gomega.Equal(int64(0)))
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(int64(0)))
	g.Expect(acct.AccountDecimals).To(gomega.Equal(uint64(2)))

	a4 := Account{AccountName: "income", AccountFullName: "income",
		AccountSign: AccountSignCredit, AccountType: AccountTypeIncome,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 3, AccountRight: 4}
	err = aStore.Store(&a4)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a4.AccountID).NotTo(gomega.BeZero())
	g.Expect(a4.AccountSign).To(gomega.Equal(AccountSignCredit))
	g.Expect(a4.AccountType).To(gomega.Equal(AccountTypeIncome))

	acct, err = aStore.GetAccountByID(a4.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acct.AccountName).To(gomega.Equal("income"))
	g.Expect(acct.AccountFullName).To(gomega.Equal("income"))
	g.Expect(acct.AccountID).To(gomega.Equal(a4.AccountID))
	g.Expect(acct.AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(acct.AccountRight).To(gomega.Equal(uint64(4)))
	g.Expect(acct.AccountBalance).To(gomega.Equal(int64(0)))
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(int64(0)))
	g.Expect(acct.AccountDecimals).To(gomega.Equal(uint64(2)))
	g.Expect(acct.AccountSign).To(gomega.Equal(AccountSignCredit))
	g.Expect(acct.AccountType).To(gomega.Equal(AccountTypeIncome))

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
	g.Expect(acct.AccountBalance).To(gomega.Equal(int64(0)))
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(int64(0)))
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

func TestAccountStore_GetLevelsOfChildren(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	aStore := createAccountStore()
	a1 := Account{AccountName: "my bank", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 1, AccountRight: 8}
	err := aStore.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a2 := Account{AccountName: "subact1", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 2, AccountRight: 7, AccountParent: a1.AccountID}
	err = aStore.Store(&a2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a3 := Account{AccountName: "subact2", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 3, AccountRight: 6, AccountParent: a2.AccountID}
	err = aStore.Store(&a3)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a4 := Account{AccountName: "sub_subact2", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 4, AccountRight: 5, AccountParent: a3.AccountID}
	err = aStore.Store(&a4)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	children, err := aStore.GetAccountWithChildrenByLevel(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(children).To(gomega.HaveLen(4))
	g.Expect(children[0].AccountName).To(gomega.Equal(a1.AccountName))
	g.Expect(children[0].AccountID).To(gomega.Equal(a1.AccountID))
	g.Expect(children[0].AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(children[0].AccountRight).To(gomega.Equal(uint64(8)))
	g.Expect(children[0].Level).To(gomega.Equal(int(0)))
	g.Expect(children[1].AccountName).To(gomega.Equal(a2.AccountName))
	g.Expect(children[1].AccountID).To(gomega.Equal(a2.AccountID))
	g.Expect(children[1].AccountLeft).To(gomega.Equal(uint64(2)))
	g.Expect(children[1].AccountRight).To(gomega.Equal(uint64(7)))
	g.Expect(children[1].Level).To(gomega.Equal(int(1)))
	g.Expect(children[2].AccountName).To(gomega.Equal(a3.AccountName))
	g.Expect(children[2].AccountID).To(gomega.Equal(a3.AccountID))
	g.Expect(children[2].AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(children[2].AccountRight).To(gomega.Equal(uint64(6)))
	g.Expect(children[2].Level).To(gomega.Equal(int(2)))
	g.Expect(children[3].AccountName).To(gomega.Equal(a4.AccountName))
	g.Expect(children[3].AccountID).To(gomega.Equal(a4.AccountID))
	g.Expect(children[3].AccountLeft).To(gomega.Equal(uint64(4)))
	g.Expect(children[3].AccountRight).To(gomega.Equal(uint64(5)))
	g.Expect(children[3].Level).To(gomega.Equal(int(3)))
}

func TestAccountStore_GetLevelsOfChildrenMultipleChildren(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	aStore := createAccountStore()
	a1 := Account{AccountName: "my bank", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 1, AccountRight: 8}
	err := aStore.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a2 := Account{AccountName: "subact1", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 2, AccountRight: 5, AccountParent: a1.AccountID}
	err = aStore.Store(&a2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a3 := Account{AccountName: "sub_subact1", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 3, AccountRight: 4, AccountParent: a2.AccountID}
	err = aStore.Store(&a3)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	// a4 is sub of a1
	a4 := Account{AccountName: "subact2", AccountFullName: "BankAccounts:MyBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 6, AccountRight: 7, AccountParent: a1.AccountID}
	err = aStore.Store(&a4)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	children, err := aStore.GetAccountWithChildrenByLevel(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(children).To(gomega.HaveLen(4))
	g.Expect(children[0].AccountName).To(gomega.Equal(a1.AccountName))
	g.Expect(children[0].AccountID).To(gomega.Equal(a1.AccountID))
	g.Expect(children[0].AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(children[0].AccountRight).To(gomega.Equal(uint64(8)))
	g.Expect(children[0].Level).To(gomega.Equal(int(0)))
	g.Expect(children[1].AccountName).To(gomega.Equal(a2.AccountName))
	g.Expect(children[1].AccountID).To(gomega.Equal(a2.AccountID))
	g.Expect(children[1].AccountLeft).To(gomega.Equal(uint64(2)))
	g.Expect(children[1].AccountRight).To(gomega.Equal(uint64(5)))
	g.Expect(children[1].Level).To(gomega.Equal(int(1)))
	g.Expect(children[2].AccountName).To(gomega.Equal(a3.AccountName))
	g.Expect(children[2].AccountID).To(gomega.Equal(a3.AccountID))
	g.Expect(children[2].AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(children[2].AccountRight).To(gomega.Equal(uint64(4)))
	g.Expect(children[2].Level).To(gomega.Equal(int(2)))
	g.Expect(children[3].AccountName).To(gomega.Equal(a4.AccountName))
	g.Expect(children[3].AccountID).To(gomega.Equal(a4.AccountID))
	g.Expect(children[3].AccountLeft).To(gomega.Equal(uint64(6)))
	g.Expect(children[3].AccountRight).To(gomega.Equal(uint64(7)))
	g.Expect(children[3].Level).To(gomega.Equal(int(1)))
}

func TestAccountStore_GetParents(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	aStore := createAccountStore()
	a1 := Account{AccountName: "myBank", AccountFullName: "myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 1, AccountRight: 6}
	err := aStore.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a2 := Account{AccountName: "subact2", AccountFullName: "myBank:subact2",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 2, AccountRight: 5, AccountParent: a1.AccountID}
	err = aStore.Store(&a2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a3 := Account{AccountName: "sub_subact3", AccountFullName: "myBank:subact2:sub_subact3",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 3, AccountRight: 4, AccountParent: a2.AccountID}
	err = aStore.Store(&a3)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	children, err := aStore.GetParents(a3.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(children).To(gomega.HaveLen(3))
	g.Expect(children[0].AccountName).To(gomega.Equal(a1.AccountName))
	g.Expect(children[0].AccountID).To(gomega.Equal(a1.AccountID))
	g.Expect(children[0].AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(children[0].AccountRight).To(gomega.Equal(uint64(6)))
	g.Expect(children[1].AccountName).To(gomega.Equal(a2.AccountName))
	g.Expect(children[1].AccountID).To(gomega.Equal(a2.AccountID))
	g.Expect(children[1].AccountLeft).To(gomega.Equal(uint64(2)))
	g.Expect(children[1].AccountRight).To(gomega.Equal(uint64(5)))
	g.Expect(children[2].AccountName).To(gomega.Equal(a3.AccountName))
	g.Expect(children[2].AccountID).To(gomega.Equal(a3.AccountID))
	g.Expect(children[2].AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(children[2].AccountRight).To(gomega.Equal(uint64(4)))
}

func TestAccountStore_UpdateSubtotal(t *testing.T) {
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
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(int64(0)))

	acct.AccountSubtotal = 5555

	err = aStore.UpdateSubtotal(acct)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acct, err = aStore.GetAccountByID(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(int64(5555)))
}

func TestAccountStore_UpdateBalance(t *testing.T) {
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
	g.Expect(acct.AccountSubtotal).To(gomega.Equal(int64(0)))

	acct.AccountBalance = 5555

	err = aStore.UpdateBalance(acct)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acct, err = aStore.GetAccountByID(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(acct.AccountBalance).To(gomega.Equal(int64(5555)))
}

func TestAccountStore_GetBalance(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	aStore := createAccountStore()
	// create two accounts
	a1 := Account{AccountName: "myBank", AccountFullName: "myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 1, AccountRight: 6}
	err := aStore.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a2 := Account{AccountName: "subact2", AccountFullName: "myBank:subact2",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 2, AccountRight: 5, AccountParent: a1.AccountID}
	err = aStore.Store(&a2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// create 2 transactions
	transStore := createTransactionStore()
	myTrans1 := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myTrans2 := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// create debits on both
	dcStore := createTransactionDCStore()
	myDC1 := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 45000,
		TransactionID: myTrans1.TransactionID, AccountID: a1.AccountID}
	err = dcStore.Store(&myDC1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDC2 := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 56000,
		TransactionID: myTrans2.TransactionID, AccountID: a2.AccountID}
	err = dcStore.Store(&myDC2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a1.AccountSubtotal = 45000
	err = aStore.UpdateSubtotal(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a2.AccountSubtotal = 56000
	err = aStore.UpdateSubtotal(&a2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCSet1, err := dcStore.GetSubtotals(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet1).To(gomega.HaveLen(1))
	g.Expect(myDCSet1[0].DebitOrCredit).To(gomega.Equal(AccountSignDebit))
	g.Expect(myDCSet1[0].Subtotal).To(gomega.Equal(uint64(45000)))

	myDCSet2, err := dcStore.GetSubtotals(a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet2).To(gomega.HaveLen(1))
	g.Expect(myDCSet2[0].DebitOrCredit).To(gomega.Equal(AccountSignDebit))
	g.Expect(myDCSet2[0].Subtotal).To(gomega.Equal(uint64(56000)))

	balance, err := aStore.GetBalance(a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(balance).To(gomega.Equal(int64(56000)))

	balance, err = aStore.GetBalance(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(balance).To(gomega.Equal(int64(101000)))
}

func TestAccountStore_GetBalance_FromSubaccounts(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	aStore := createAccountStore()
	// create two accounts
	a1 := Account{AccountName: "myBank", AccountFullName: "myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 1, AccountRight: 6}
	err := aStore.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a2 := Account{AccountName: "subact2", AccountFullName: "myBank:subact2",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 2, AccountRight: 5, AccountParent: a1.AccountID}
	err = aStore.Store(&a2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// create a transactions
	transStore := createTransactionStore()
	myTrans1 := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// create debits on both
	dcStore := createTransactionDCStore()
	myDC1 := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 45000,
		TransactionID: myTrans1.TransactionID, AccountID: a2.AccountID}
	err = dcStore.Store(&myDC1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a2.AccountSubtotal = 45000
	err = aStore.UpdateSubtotal(&a2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	balance, err := aStore.GetBalance(a2.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(balance).To(gomega.Equal(int64(45000)))

	balance, err = aStore.GetBalance(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(balance).To(gomega.Equal(int64(45000)))
}

func TestAccountStore_AccountSetReconcileDate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	aStore := createAccountStore()
	// create two accounts
	a1 := Account{AccountName: "myBank", AccountFullName: "myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 2000,
		AccountLeft: 1, AccountRight: 6}
	err := aStore.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myA1, err := aStore.GetAccountByID(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myA1.AccountReconcileDate).To(gomega.Equal(sql.NullTime{Time: time.Time{}, Valid: false}))

	a1.AccountReconcileDate = sql.NullTime{Time: time.Now(), Valid: true}
	err = aStore.SetAccountReconciledDate(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myA1, err = aStore.GetAccountByID(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myA1.AccountReconcileDate.Time).To(gomega.BeTemporally("~", time.Now(), time.Second))
	g.Expect(myA1.AccountReconcileDate.Valid).To(gomega.BeTrue())

	// unset the AccountReconciled Date
	a1.AccountReconcileDate.Valid = false
	err = aStore.SetAccountReconciledDate(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myA1, err = aStore.GetAccountByID(a1.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myA1.AccountReconcileDate).To(gomega.Equal(sql.NullTime{Time: time.Time{}, Valid: false}))
}

func setupDB(g *gomega.WithT) {
	query := `delete from transaction_debit_credit `
	_, err := TestPostgresClient.Exec(query)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	query = `delete from transaction_accounts `
	_, err = TestPostgresClient.Exec(query)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	query = `delete from reports `
	_, err = TestPostgresClient.Exec(query)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
