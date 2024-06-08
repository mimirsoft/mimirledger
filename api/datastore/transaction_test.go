package datastore

import (
	"database/sql"
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
