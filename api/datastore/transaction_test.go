package datastore

import (
	"database/sql"
	"errors"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
	"time"
)

func createTransactionStore() TransactionStore {
	return TransactionStore{
		Client: TestPostgresClient,
	}
}
func TestTransactionStore_StoreValidEmpty(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	store := createTransactionStore()

	// failed due to no TransactionAmount
	a1 := Transaction{}
	err := store.Store(&a1)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring(`new row for relation "transaction_main" violates check constraint "transaction_main_transaction_amount_check"`))

	// failed due to no TransactionComment
	a2 := Transaction{TransactionAmount: 2000}
	err = store.Store(&a2)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(
		gomega.ContainSubstring(`new row for relation "transaction_main" violates check constraint "transaction_main_transaction_comment_check"`))
}

func TestTransactionStore_StoreValidAndRetrieve(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	store := createTransactionStore()

	a1 := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a1.TransactionID).NotTo(gomega.BeZero())
	g.Expect(a1.IsSplit).To(gomega.BeFalse())
	g.Expect(a1.IsReconciled).To(gomega.BeFalse())

	myTrans, err := store.GetByID(a1.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myTrans).NotTo(gomega.BeNil())
	g.Expect(myTrans.TransactionComment).To(gomega.Equal("woot"))
	g.Expect(myTrans.TransactionAmount).To(gomega.Equal(uint64(1000)))
	g.Expect(myTrans.TransactionDate).To(gomega.BeTemporally("~", time.Now(), time.Second))
	g.Expect(myTrans.TransactionReconcileDate).To(gomega.Equal(sql.NullTime{Time: time.Time{}, Valid: false}))
}

func TestTransactionStore_StoreUpdateAndRetrieve(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	store := createTransactionStore()

	a1 := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	a1.TransactionComment = "updated woot"
	a1.TransactionAmount = 4444
	err = store.Update(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(a1.TransactionComment).To(gomega.Equal("updated woot"))
	g.Expect(a1.TransactionAmount).To(gomega.Equal(uint64(4444)))

	myTrans, err := store.GetByID(a1.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myTrans).NotTo(gomega.BeNil())
	g.Expect(myTrans.TransactionComment).To(gomega.Equal("updated woot"))
	g.Expect(myTrans.TransactionAmount).To(gomega.Equal(uint64(4444)))
	g.Expect(myTrans.TransactionDate).To(gomega.BeTemporally("~", time.Now(), time.Second))
	g.Expect(myTrans.TransactionReconcileDate).To(gomega.Equal(sql.NullTime{Time: time.Time{}, Valid: false}))
}

func TestTransactionStore_StoreAndDelete(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	store := createTransactionStore()

	a1 := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myTrans, err := store.GetByID(a1.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myTrans).NotTo(gomega.BeNil())

	// Delete once
	err = store.Delete(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// try to get - should return nil, err
	myTrans, err = store.GetByID(a1.TransactionID)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, sql.ErrNoRows)).To(gomega.BeTrue())
	g.Expect(myTrans).To(gomega.BeNil())

	// try to delete again, after already deleted
	err = store.Delete(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestTransactionStore_StoreWithDCAndDelete(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	store := createTransactionStore()

	a1 := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	aStore := createAccountStore()
	myAcct := Account{AccountName: "myBank", AccountFullName: "BankAccounts:myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 1, AccountRight: 2}
	err = aStore.Store(&myAcct)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myAcct2 := Account{AccountName: "revenue", AccountFullName: "IncomeToBank",
		AccountSign: AccountSignCredit, AccountType: AccountTypeIncome,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 3, AccountRight: 4}
	err = aStore.Store(&myAcct2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	dcStore := createTransactionDCStore()
	// succeed
	myDC := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 10000,
		TransactionID: a1.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDC2 := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 20000,
		TransactionID: a1.TransactionID, AccountID: myAcct2.AccountID}
	err = dcStore.Store(&myDC2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myTrans, err := store.GetByID(a1.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myTrans).NotTo(gomega.BeNil())

	// Delete once
	err = store.Delete(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// try to get - should return nil, err
	myTrans, err = store.GetByID(a1.TransactionID)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, sql.ErrNoRows)).To(gomega.BeTrue())
	g.Expect(myTrans).To(gomega.BeNil())

	// try to delete again, after already deleted
	err = store.Delete(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestTransactionStore_StoreAndToggleReconciled(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	store := createTransactionStore()

	a1 := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err := store.Store(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myTrans, err := store.GetByID(a1.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myTrans).NotTo(gomega.BeNil())
	g.Expect(myTrans.IsReconciled).To(gomega.BeFalse())
	g.Expect(myTrans.TransactionReconcileDate).To(gomega.Equal(sql.NullTime{Time: time.Time{}, Valid: false}))

	a1.IsReconciled = true
	err = store.SetIsReconciled(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myTrans, err = store.GetByID(a1.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myTrans).NotTo(gomega.BeNil())
	g.Expect(myTrans).NotTo(gomega.BeNil())
	g.Expect(myTrans.IsReconciled).To(gomega.BeTrue())
	g.Expect(myTrans.TransactionReconcileDate).To(gomega.Equal(sql.NullTime{Time: time.Time{}, Valid: false}))

	a1.TransactionReconcileDate = sql.NullTime{Time: time.Now(), Valid: true}
	err = store.SetTransactionReconcileDate(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myTrans, err = store.GetByID(a1.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myTrans).NotTo(gomega.BeNil())
	g.Expect(myTrans).NotTo(gomega.BeNil())
	g.Expect(myTrans.IsReconciled).To(gomega.BeTrue())
	g.Expect(myTrans.TransactionReconcileDate.Time).To(gomega.BeTemporally("~", time.Now(), time.Second))
	g.Expect(myTrans.TransactionReconcileDate.Valid).To(gomega.BeTrue())

	// unset both
	a1.TransactionReconcileDate.Valid = false
	a1.IsReconciled = false
	err = store.SetIsReconciled(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	err = store.SetTransactionReconcileDate(&a1)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myTrans, err = store.GetByID(a1.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myTrans).NotTo(gomega.BeNil())
	g.Expect(myTrans.IsReconciled).To(gomega.BeFalse())
	g.Expect(myTrans.TransactionReconcileDate).To(gomega.Equal(sql.NullTime{Time: time.Time{}, Valid: false}))
}
