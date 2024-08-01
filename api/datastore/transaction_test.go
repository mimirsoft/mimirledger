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

func TestTransactionStore_GetUnreconciledTransactionsOnAccountForDate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	aStore := createAccountStore()
	myAcct := Account{AccountName: "myBank", AccountFullName: "BankAccounts:myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 1, AccountRight: 2}
	err := aStore.Store(&myAcct)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myAcct2 := Account{AccountName: "revenue", AccountFullName: "IncomeToBank",
		AccountSign: AccountSignCredit, AccountType: AccountTypeIncome,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 3, AccountRight: 4}
	err = aStore.Store(&myAcct2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	oldDate1, err := time.Parse("2006-01-02", "2016-07-08")
	g.Expect(err).NotTo(gomega.HaveOccurred())

	transStore := createTransactionStore()
	myTrans := Transaction{TransactionDate: oldDate1, TransactionComment: "added_for_reconciliation_test", TransactionAmount: 1000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	dcStore := createTransactionDCStore()
	// succeed
	myDCa := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 10000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDCa)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCb := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 20000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct2.AccountID}
	err = dcStore.Store(&myDCb)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCSet, err := dcStore.GetDCForTransactionID(myTrans.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(2))

	testSearchLimitDateoldDate1, err := time.Parse("2006-01-02", "2015-07-08")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	// should return nothing, as we cut off the search before the date on the transaction
	unreconciledTransaction, err := transStore.GetUnreconciledTransactionsOnAccountForDate(myAcct.AccountLeft, myAcct.AccountRight,
		testSearchLimitDateoldDate1, time.Time{})
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, sql.ErrNoRows)).To(gomega.BeTrue())
	g.Expect(unreconciledTransaction).To(gomega.HaveLen(0))

	unreconciledTransaction, err = transStore.GetUnreconciledTransactionsOnAccountForDate(myAcct.AccountLeft, myAcct.AccountRight,
		time.Now().Add(time.Hour), time.Time{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(unreconciledTransaction).To(gomega.HaveLen(1))
	g.Expect(unreconciledTransaction[0].TransactionID).To(gomega.Equal(myTrans.TransactionID))

	// set is_reconciled and the reconciled_date on myTrans1
	myTrans.IsReconciled = true
	err = transStore.SetIsReconciled(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	reconciledDate, err := time.Parse("2006-01-02", "2016-07-11")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myTrans.TransactionReconcileDate = sql.NullTime{Time: reconciledDate, Valid: true}
	err = transStore.SetTransactionReconcileDate(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// should still return nothing, as the reconciled date greater than the search limit date
	unreconciledTransaction, err = transStore.GetUnreconciledTransactionsOnAccountForDate(myAcct.AccountLeft, myAcct.AccountRight,
		testSearchLimitDateoldDate1, time.Time{})
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, sql.ErrNoRows)).To(gomega.BeTrue())
	g.Expect(unreconciledTransaction).To(gomega.HaveLen(0))

	// should return 1 record, as this transaction is reconciled, but the reconciled date is after the cutoffdate
	reconciledDateCutoff1, err := time.Parse("2006-01-02", "2016-06-30")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	reconciledDateCutoff2, err := time.Parse("2006-01-02", "2016-07-15")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	testSearchLimitDate2, err := time.Parse("2006-01-02", "2016-07-31")
	g.Expect(err).NotTo(gomega.HaveOccurred())

	unreconciledTransaction, err = transStore.GetUnreconciledTransactionsOnAccountForDate(myAcct.AccountLeft, myAcct.AccountRight,
		testSearchLimitDate2, reconciledDateCutoff1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(unreconciledTransaction).To(gomega.HaveLen(1))

	unreconciledTransaction, err = transStore.GetUnreconciledTransactionsOnAccountForDate(myAcct.AccountLeft, myAcct.AccountRight,
		testSearchLimitDate2, reconciledDateCutoff2)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, sql.ErrNoRows)).To(gomega.BeTrue())
	g.Expect(unreconciledTransaction).To(gomega.HaveLen(0))

	// add another transaction
	oldDate2, err := time.Parse("2006-01-02", "2017-09-08")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myTrans2 := Transaction{TransactionDate: oldDate2, TransactionComment: "added_another_reconciliation_test", TransactionAmount: 1000}
	err = transStore.Store(&myTrans2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// succeed
	myDC2a := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 10000,
		TransactionID: myTrans2.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC2a)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDC2b := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 20000,
		TransactionID: myTrans2.TransactionID, AccountID: myAcct2.AccountID}
	err = dcStore.Store(&myDC2b)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// get unreconciled transactions , there are now two of them
	unreconciledTransaction2, err := transStore.GetUnreconciledTransactionsOnAccountForDate(myAcct.AccountLeft, myAcct.AccountRight,
		time.Now().Add(time.Hour), time.Time{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(unreconciledTransaction2).To(gomega.HaveLen(2))
	g.Expect(unreconciledTransaction2[0].TransactionID).To(gomega.Equal(myTrans.TransactionID))
	g.Expect(unreconciledTransaction2[1].TransactionID).To(gomega.Equal(myTrans2.TransactionID))

	// set is_reconciled and the reconciled_date on one of the transactions
	myTrans2.IsReconciled = true
	err = transStore.SetIsReconciled(&myTrans2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	reconciledDate2, err := time.Parse("2006-01-02", "2023-09-08")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myTrans2.TransactionReconcileDate = sql.NullTime{Time: reconciledDate2, Valid: true}
	err = transStore.SetTransactionReconcileDate(&myTrans2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// get transactions , there are still two of them, one of them unreconciled, and the other reconciled, but after the
	// late reconciled dated on the account(which is zero)
	unreconciledTransaction2, err = transStore.GetUnreconciledTransactionsOnAccountForDate(myAcct.AccountLeft, myAcct.AccountRight,
		time.Now().Add(48*time.Hour), time.Time{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(unreconciledTransaction2).To(gomega.HaveLen(2))
	// both transactions are reconciled, myTrans2 has a later date so it should be last
	g.Expect(unreconciledTransaction2[1].TransactionID).To(gomega.Equal(myTrans2.TransactionID))
	g.Expect(unreconciledTransaction2[1].IsReconciled).To(gomega.BeTrue())
	g.Expect(unreconciledTransaction2[0].TransactionID).To(gomega.Equal(myTrans.TransactionID))
	g.Expect(unreconciledTransaction2[0].IsReconciled).To(gomega.BeTrue())

	accountReconciledDateCutoff2, err := time.Parse("2006-01-02", "2024-01-31")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	// both transactions are reconciled, and both are reconciled prior to the newer cutoff date
	unreconciledTransaction3, err := transStore.GetUnreconciledTransactionsOnAccountForDate(myAcct.AccountLeft, myAcct.AccountRight,
		time.Now(), accountReconciledDateCutoff2)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, sql.ErrNoRows)).To(gomega.BeTrue())
	g.Expect(unreconciledTransaction3).To(gomega.HaveLen(0))
}

func TestTransactionStore_RetrieveTransactionsNetForDates(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	aStore := createAccountStore()
	myAcct := Account{AccountName: "myBank", AccountFullName: "BankAccounts:myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 1, AccountRight: 2}
	err := aStore.Store(&myAcct)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myAcct2 := Account{AccountName: "revenue", AccountFullName: "IncomeToBank",
		AccountSign: AccountSignCredit, AccountType: AccountTypeIncome,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 3, AccountRight: 4}
	err = aStore.Store(&myAcct2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	oldDate1, err := time.Parse("2006-01-02", "2016-07-08")
	g.Expect(err).NotTo(gomega.HaveOccurred())

	transStore := createTransactionStore()
	myTrans := Transaction{TransactionDate: oldDate1, TransactionComment: "added_for_reconciliation_test", TransactionAmount: 1000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	dcStore := createTransactionDCStore()
	// succeed
	myDCa := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 10000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDCa)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCb := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 20000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct2.AccountID}
	err = dcStore.Store(&myDCb)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCSet, err := dcStore.GetDCForTransactionID(myTrans.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(2))

	testSearchLimitDate1, err := time.Parse("2006-01-02", "2015-07-08")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	// should return nothing, as we cut off the search before the date on the transaction
	unreconciledTransaction, err := transStore.RetrieveTransactionsNetForDates([]uint64{myAcct.AccountID},
		time.Time{}, testSearchLimitDate1)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, sql.ErrNoRows)).To(gomega.BeTrue())
	g.Expect(unreconciledTransaction).To(gomega.HaveLen(0))

	unreconciledTransaction, err = transStore.RetrieveTransactionsNetForDates([]uint64{myAcct.AccountID},
		time.Time{}, time.Now().Add(time.Hour))
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(unreconciledTransaction).To(gomega.HaveLen(1))
	g.Expect(unreconciledTransaction[0].TransactionID).To(gomega.Equal(myTrans.TransactionID))

	// should  return nothing, because end date is later than start date
	unreconciledTransaction, err = transStore.RetrieveTransactionsNetForDates([]uint64{myAcct.AccountID},
		testSearchLimitDate1, time.Time{})
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, sql.ErrNoRows)).To(gomega.BeTrue())
	g.Expect(unreconciledTransaction).To(gomega.HaveLen(0))

	// add another transaction
	oldDate2, err := time.Parse("2006-01-02", "2017-09-08")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myTrans2 := Transaction{TransactionDate: oldDate2, TransactionComment: "added_another_reconciliation_test", TransactionAmount: 1000}
	err = transStore.Store(&myTrans2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// succeed
	myDC2a := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 10000,
		TransactionID: myTrans2.TransactionID, AccountID: myAcct2.AccountID}
	err = dcStore.Store(&myDC2a)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDC2b := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 20000,
		TransactionID: myTrans2.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC2b)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// get unreconciled transactions , there are now two of them
	unreconciledTransaction2, err := transStore.RetrieveTransactionsNetForDates([]uint64{myAcct.AccountID},
		time.Time{}, time.Now().Add(time.Hour))
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(unreconciledTransaction2).To(gomega.HaveLen(2))
	g.Expect(unreconciledTransaction2[0].TransactionID).To(gomega.Equal(myTrans.TransactionID))
	g.Expect(unreconciledTransaction2[1].TransactionID).To(gomega.Equal(myTrans2.TransactionID))

	testSearchLimitDate2, err := time.Parse("2006-01-02", "2017-08-31")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	// get unreconciled transactions , there are now two of them
	unreconciledTransaction2, err = transStore.RetrieveTransactionsNetForDates([]uint64{myAcct.AccountID},
		testSearchLimitDate2, time.Now().Add(time.Hour))
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(unreconciledTransaction2).To(gomega.HaveLen(1))
	g.Expect(unreconciledTransaction2[0].TransactionID).To(gomega.Equal(myTrans2.TransactionID))
}
