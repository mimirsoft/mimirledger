package datastore

import (
	"database/sql"
	"errors"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
	"time"
)

func createTransactionDCStore() TransactionDebitCreditStore {
	return TransactionDebitCreditStore{
		Client: TestPostgresClient,
	}
}

func TestTransactionDebitCreditStore_StoreInvalid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	store := createTransactionDCStore()

	// failed due to no Account type sign
	myDC := TransactionDebitCredit{}
	err := store.Store(&myDC)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring(`invalid input value for enum transaction_account_sign_type`))

	// failed due to no amount
	myDC = TransactionDebitCredit{DebitOrCredit: AccountSignCredit}
	err = store.Store(&myDC)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring(`new row for relation "transaction_debit_credit" violates check constraint "transaction_debit_credit_transaction_dc_amount_check"`))

	// failed due to no transaction ID
	myDC = TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 10000}
	err = store.Store(&myDC)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring(`insert or update on table "transaction_debit_credit" violates foreign key constraint "transactions_debit_credit_transaction_id_fkey"`))

	transStore := createTransactionStore()
	myTrans := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// failed due to no AccountID
	myDC = TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 10000, TransactionID: myTrans.TransactionID}
	err = store.Store(&myDC)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(err.Error()).To(gomega.ContainSubstring(`insert or update on table "transaction_debit_credit" violates foreign key constraint "transactions_debit_credit_transaction_account_id_fkey`))

	aStore := createAccountStore()
	myAcct := Account{AccountName: "myBank", AccountFullName: "BankAccounts:myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 1, AccountRight: 6}
	err = aStore.Store(&myAcct)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestTransactionDebitCreditStore_StoreValid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	aStore := createAccountStore()
	myAcct := Account{AccountName: "myBank", AccountFullName: "BankAccounts:myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 3000, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 1, AccountRight: 6}
	err := aStore.Store(&myAcct)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	transStore := createTransactionStore()
	myTrans := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	store := createTransactionDCStore()

	// succeed
	myDC := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 10000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct.AccountID}
	err = store.Store(&myDC)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestTransactionDebitCreditStore_RetrieveInvalid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	dcStore := createTransactionDCStore()

	myDCSet, err := dcStore.GetDCForTransactionID(4444)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, sql.ErrNoRows)).To(gomega.BeTrue())
	g.Expect(myDCSet).To(gomega.HaveLen(0))
}

func TestTransactionDebitCreditStore_StoreAndRetrieve(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

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

	transStore := createTransactionStore()
	myTrans := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	dcStore := createTransactionDCStore()

	// succeed
	myDC := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 10000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDC2 := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 20000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct2.AccountID}
	err = dcStore.Store(&myDC2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCSet, err := dcStore.GetDCForTransactionID(myTrans.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(2))
	g.Expect(myDCSet[0].TransactionDCID).To(gomega.Equal(myDC.TransactionDCID))
	g.Expect(myDCSet[0].TransactionID).To(gomega.Equal(myDC.TransactionID))
	g.Expect(myDCSet[0].AccountID).To(gomega.Equal(myAcct.AccountID))
	g.Expect(myDCSet[0].DebitOrCredit).To(gomega.Equal(AccountSignDebit))
	g.Expect(myDCSet[0].TransactionDCAmount).To(gomega.Equal(uint64(10000)))
	g.Expect(myDCSet[1].TransactionDCID).To(gomega.Equal(myDC2.TransactionDCID))
	g.Expect(myDCSet[1].TransactionID).To(gomega.Equal(myDC.TransactionID))
	g.Expect(myDCSet[1].AccountID).To(gomega.Equal(myAcct2.AccountID))
	g.Expect(myDCSet[1].DebitOrCredit).To(gomega.Equal(AccountSignCredit))
	g.Expect(myDCSet[1].TransactionDCAmount).To(gomega.Equal(uint64(20000)))

}

func TestTransactionDebitCreditStore_StoreThenDelete(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

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

	transStore := createTransactionStore()
	myTrans := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	dcStore := createTransactionDCStore()
	// succeed
	myDC := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 10000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDC2 := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 20000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct2.AccountID}
	err = dcStore.Store(&myDC2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCSet, err := dcStore.GetDCForTransactionID(myTrans.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(2))

	deleteDCSet, err := dcStore.DeleteForTransactionID(myTrans.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(deleteDCSet).To(gomega.HaveLen(2))

	myDCSet, err = dcStore.GetDCForTransactionID(myTrans.TransactionID)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, sql.ErrNoRows)).To(gomega.BeTrue())
	g.Expect(myDCSet).To(gomega.HaveLen(0))
}

func TestTransactionDebitCreditStore_GetSubtotals_Empty(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	aStore := createAccountStore()
	myAcct := Account{AccountName: "myBank", AccountFullName: "BankAccounts:myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 1, AccountRight: 2}
	err := aStore.Store(&myAcct)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	dcStore := createTransactionDCStore()
	myDCSet, err := dcStore.GetSubtotals(myAcct.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.BeNil())
	g.Expect(myDCSet).To(gomega.HaveLen(0))
}

func TestTransactionDebitCreditStore_GetSubtotals_OnlyDebits(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	aStore := createAccountStore()
	myAcct := Account{AccountName: "myBank", AccountFullName: "BankAccounts:myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 1, AccountRight: 2}
	err := aStore.Store(&myAcct)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	transStore := createTransactionStore()
	myTrans := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// create only a debit
	dcStore := createTransactionDCStore()
	myDC := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 10000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCSet, err := dcStore.GetSubtotals(myAcct.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(1))
	g.Expect(myDCSet[0].DebitOrCredit).To(gomega.Equal(AccountSignDebit))
	g.Expect(myDCSet[0].Subtotal).To(gomega.Equal(uint64(10000)))

	// add another transaction
	myTrans2 := Transaction{TransactionComment: "woot2", TransactionAmount: 30000}
	err = transStore.Store(&myTrans2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// create another transaction with only a debit
	myDC2 := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 30000,
		TransactionID: myTrans2.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCSet, err = dcStore.GetSubtotals(myAcct.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(1))
	g.Expect(myDCSet[0].DebitOrCredit).To(gomega.Equal(AccountSignDebit))
	g.Expect(myDCSet[0].Subtotal).To(gomega.Equal(uint64(40000)))
}

func TestTransactionDebitCreditStore_GetSubtotals_OnlyCredits(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	aStore := createAccountStore()
	myAcct := Account{AccountName: "myBank", AccountFullName: "BankAccounts:myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 1, AccountRight: 2}
	err := aStore.Store(&myAcct)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	transStore := createTransactionStore()
	myTrans := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// create only a credit
	dcStore := createTransactionDCStore()
	myDC := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 13000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCSet, err := dcStore.GetSubtotals(myAcct.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(1))
	g.Expect(myDCSet[0].DebitOrCredit).To(gomega.Equal(AccountSignCredit))
	g.Expect(myDCSet[0].Subtotal).To(gomega.Equal(uint64(13000)))

	// add another transaction
	myTrans2 := Transaction{TransactionComment: "woot2", TransactionAmount: 30000}
	err = transStore.Store(&myTrans2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// create another transaction with only a credit
	myDC2 := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 34000,
		TransactionID: myTrans2.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCSet, err = dcStore.GetSubtotals(myAcct.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(1))
	g.Expect(myDCSet[0].DebitOrCredit).To(gomega.Equal(AccountSignCredit))
	g.Expect(myDCSet[0].Subtotal).To(gomega.Equal(uint64(47000)))
}

func TestTransactionDebitCreditStore_GetSubtotals_DebitsAndCredits(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	aStore := createAccountStore()
	myAcct := Account{AccountName: "myBank", AccountFullName: "BankAccounts:myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 1, AccountRight: 2}
	err := aStore.Store(&myAcct)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	transStore := createTransactionStore()
	myTrans := Transaction{TransactionComment: "woot", TransactionAmount: 1000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// create only a credit
	dcStore := createTransactionDCStore()
	myDC := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 13000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCb := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 10000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDCb)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCSet, err := dcStore.GetSubtotals(myAcct.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(2))
	g.Expect(myDCSet[0].DebitOrCredit).To(gomega.Equal(AccountSignCredit))
	g.Expect(myDCSet[0].Subtotal).To(gomega.Equal(uint64(13000)))
	g.Expect(myDCSet[1].DebitOrCredit).To(gomega.Equal(AccountSignDebit))
	g.Expect(myDCSet[1].Subtotal).To(gomega.Equal(uint64(10000)))

	// add another transaction
	myTrans2 := Transaction{TransactionComment: "woot2", TransactionAmount: 30000}
	err = transStore.Store(&myTrans2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// create another transaction with only a credit
	myDC2 := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 34000,
		TransactionID: myTrans2.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myDC2b := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 30000,
		TransactionID: myTrans2.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC2b)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCSet, err = dcStore.GetSubtotals(myAcct.AccountID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(2))
	g.Expect(myDCSet[0].DebitOrCredit).To(gomega.Equal(AccountSignCredit))
	g.Expect(myDCSet[0].Subtotal).To(gomega.Equal(uint64(47000)))
	g.Expect(myDCSet[1].DebitOrCredit).To(gomega.Equal(AccountSignDebit))
	g.Expect(myDCSet[1].Subtotal).To(gomega.Equal(uint64(40000)))
}

func TestTransactionDebitCreditStore_GetReconciledSubtotal(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	aStore := createAccountStore()
	myAcct := Account{AccountName: "myBank", AccountFullName: "BankAccounts:myBank",
		AccountSign: AccountSignDebit, AccountType: AccountTypeAsset,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 1, AccountRight: 2}
	err := aStore.Store(&myAcct)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myAcct2 := Account{AccountName: "myBank", AccountFullName: "ExpensesRent",
		AccountSign: AccountSignDebit, AccountType: AccountTypeExpense,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 3, AccountRight: 4}
	err = aStore.Store(&myAcct2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	transStore := createTransactionStore()
	myTrans := Transaction{TransactionComment: "woot", TransactionAmount: 13000}
	err = transStore.Store(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	dcStore := createTransactionDCStore()
	myDC := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 13000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myDCb := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 13000,
		TransactionID: myTrans.TransactionID, AccountID: myAcct2.AccountID}
	err = dcStore.Store(&myDCb)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// set is_reconciled and the reconciled_date on one of the transactions
	myTrans.IsReconciled = true
	err = transStore.SetIsReconciled(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	reconciledDate2, err := time.Parse("2006-01-02", "2023-09-08")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myTrans.TransactionReconcileDate = sql.NullTime{Time: reconciledDate2, Valid: true}
	err = transStore.SetTransactionReconcileDate(&myTrans)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	reconciledDateCutoff0, err := time.Parse("2006-01-02", "2023-08-30")
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// cut off date is far in past, zero transactions in sum
	myDCSet, err := dcStore.GetReconciledSubtotals(1, 2, reconciledDateCutoff0)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(0))

	reconciledDateCutoff1, err := time.Parse("2006-01-02", "2023-09-30")
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myDCSet, err = dcStore.GetReconciledSubtotals(1, 2, reconciledDateCutoff1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(1))
	g.Expect(myDCSet[0].DebitOrCredit).To(gomega.Equal(AccountSignCredit))
	g.Expect(myDCSet[0].Subtotal).To(gomega.Equal(uint64(13000)))

	// add an income account
	myAcct3 := Account{AccountName: "subcriptionFees", AccountFullName: "IncomeSubscriptions",
		AccountSign: AccountSignCredit, AccountType: AccountTypeIncome,
		AccountBalance: 0, AccountDecimals: 2, AccountSubtotal: 0,
		AccountLeft: 5, AccountRight: 6}
	err = aStore.Store(&myAcct3)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// add an income transactions
	myTrans2 := Transaction{TransactionComment: "woot", TransactionAmount: 20000}
	err = transStore.Store(&myTrans2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myDC2 := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 20000,
		TransactionID: myTrans2.TransactionID, AccountID: myAcct3.AccountID}
	err = dcStore.Store(&myDC2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myDC2b := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 20000,
		TransactionID: myTrans2.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC2b)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// set is_reconciled and the reconciled_date on myTrans2
	myTrans2.IsReconciled = true
	err = transStore.SetIsReconciled(&myTrans2)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myTrans2.TransactionReconcileDate = sql.NullTime{Time: reconciledDate2, Valid: true}
	err = transStore.SetTransactionReconcileDate(&myTrans2)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// cut off date is far in past, zero transactions in sum
	myDCSet, err = dcStore.GetReconciledSubtotals(1, 2, reconciledDateCutoff0)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(0))

	// should now be 2 transaction,
	myDCSet, err = dcStore.GetReconciledSubtotals(1, 2, reconciledDateCutoff1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(2))
	g.Expect(myDCSet[0].DebitOrCredit).To(gomega.Equal(AccountSignCredit))
	g.Expect(myDCSet[0].Subtotal).To(gomega.Equal(uint64(13000)))
	g.Expect(myDCSet[1].DebitOrCredit).To(gomega.Equal(AccountSignDebit))
	g.Expect(myDCSet[1].Subtotal).To(gomega.Equal(uint64(20000)))

	// add an 2nd income transactions
	myTrans3 := Transaction{TransactionComment: "woot", TransactionAmount: 45000}
	err = transStore.Store(&myTrans3)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myDC3 := TransactionDebitCredit{DebitOrCredit: AccountSignCredit, TransactionDCAmount: 45000,
		TransactionID: myTrans3.TransactionID, AccountID: myAcct3.AccountID}
	err = dcStore.Store(&myDC3)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myDC3b := TransactionDebitCredit{DebitOrCredit: AccountSignDebit, TransactionDCAmount: 45000,
		TransactionID: myTrans3.TransactionID, AccountID: myAcct.AccountID}
	err = dcStore.Store(&myDC3b)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// set is_reconciled and the reconciled_date on myTrans2
	myTrans3.IsReconciled = true
	err = transStore.SetIsReconciled(&myTrans3)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	myTrans3.TransactionReconcileDate = sql.NullTime{Time: reconciledDate2, Valid: true}
	err = transStore.SetTransactionReconcileDate(&myTrans3)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// should now be 3 transaction, 13000 in credit and 65000 in debit on account1
	myDCSet, err = dcStore.GetReconciledSubtotals(1, 2, reconciledDateCutoff1)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(myDCSet).To(gomega.HaveLen(2))
	g.Expect(myDCSet[0].DebitOrCredit).To(gomega.Equal(AccountSignCredit))
	g.Expect(myDCSet[0].Subtotal).To(gomega.Equal(uint64(13000)))
	g.Expect(myDCSet[1].DebitOrCredit).To(gomega.Equal(AccountSignDebit))
	g.Expect(myDCSet[1].Subtotal).To(gomega.Equal(uint64(65000)))
}
