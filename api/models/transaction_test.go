package models

import (
	"errors"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestTransaction_StoreInvalid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	// failed due to no  TransactionComment
	txn := Transaction{}
	err := txn.Store(testDS)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, ErrTransactionNoComment)).To(gomega.BeTrue())

	// failed due to no DebitCreditSet
	txn = Transaction{TransactionCore: TransactionCore{TransactionComment: "woot"}}
	err = txn.Store(testDS)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, ErrTransactionNoDebitsCredits)).To(gomega.BeTrue())

	// failed due to DebitCreditSet having no DebitOrCredit AccountSign
	txn = Transaction{TransactionCore: TransactionCore{TransactionComment: "woot"},
		DebitCreditSet: []*TransactionDebitCredit{
			&TransactionDebitCredit{AccountID: 1},
		},
	}
	err = txn.Store(testDS)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, ErrTransactionDebitCreditsIsNeither)).To(gomega.BeTrue())

	// failed due to no zero credit / debit total
	txn = Transaction{TransactionCore: TransactionCore{TransactionComment: "woot"},
		DebitCreditSet: []*TransactionDebitCredit{
			&TransactionDebitCredit{AccountID: 1, DebitOrCredit: datastore.AccountSignCredit},
		},
	}
	err = txn.Store(testDS)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, ErrTransactionDebitCreditsZero)).To(gomega.BeTrue())

	// failed due to no empty TransactionComment
	txn = Transaction{TransactionCore: TransactionCore{TransactionComment: "woot"},
		DebitCreditSet: []*TransactionDebitCredit{
			&TransactionDebitCredit{AccountID: 1, DebitOrCredit: datastore.AccountSignCredit, TransactionDCAmount: 10000},
			&TransactionDebitCredit{AccountID: 1, DebitOrCredit: datastore.AccountSignDebit, TransactionDCAmount: 9900},
		},
	}
	err = txn.Store(testDS)
	g.Expect(err).To(gomega.HaveOccurred())
	g.Expect(errors.Is(err, ErrTransactionDebitCreditsNotBalanced)).To(gomega.BeTrue())

}

func TestTransaction_StoreValidAndRetrieve(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	// create an account first
	a1 := Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err := a1.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a2 := Account{AccountName: "Income", AccountSign: datastore.AccountSignCredit, AccountType: datastore.AccountTypeIncome}
	err = a2.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	txn := Transaction{TransactionCore: TransactionCore{TransactionComment: "woot"},
		DebitCreditSet: []*TransactionDebitCredit{
			&TransactionDebitCredit{AccountID: a2.AccountID,
				DebitOrCredit:       datastore.AccountSignCredit,
				TransactionDCAmount: 10000},
			&TransactionDebitCredit{AccountID: a1.AccountID,
				DebitOrCredit:       datastore.AccountSignDebit,
				TransactionDCAmount: 10000},
		},
	}
	err = txn.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(txn.TransactionComment).To(gomega.Equal("woot"))
	g.Expect(txn.TransactionAmount).To(gomega.Equal(uint64(10000)))
	g.Expect(txn.TransactionID).NotTo(gomega.BeZero())
	g.Expect(txn.DebitCreditSet).To(gomega.HaveLen(2))
	g.Expect(txn.DebitCreditSet[0].TransactionDCID).NotTo(gomega.BeZero())
	g.Expect(txn.DebitCreditSet[0].TransactionID).To(gomega.Equal(txn.TransactionID))
	g.Expect(txn.DebitCreditSet[0].AccountID).To(gomega.Equal(a2.AccountID))
	g.Expect(txn.DebitCreditSet[0].DebitOrCredit).To(gomega.Equal(datastore.AccountSignCredit))
	g.Expect(txn.DebitCreditSet[0].TransactionDCAmount).To(gomega.Equal(uint64(10000)))
	g.Expect(txn.DebitCreditSet[1].TransactionDCID).NotTo(gomega.BeZero())
	g.Expect(txn.DebitCreditSet[1].TransactionID).To(gomega.Equal(txn.TransactionID))
	g.Expect(txn.DebitCreditSet[1].AccountID).To(gomega.Equal(a1.AccountID))
	g.Expect(txn.DebitCreditSet[1].DebitOrCredit).To(gomega.Equal(datastore.AccountSignDebit))
	g.Expect(txn.DebitCreditSet[1].TransactionDCAmount).To(gomega.Equal(uint64(10000)))

	myTxn, err := RetrieveTransactionByID(testDS, txn.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(myTxn.TransactionComment).To(gomega.Equal("woot"))
	g.Expect(myTxn.TransactionAmount).To(gomega.Equal(uint64(10000)))
	g.Expect(myTxn.TransactionID).NotTo(gomega.BeZero())
	g.Expect(myTxn.DebitCreditSet).To(gomega.HaveLen(2))
	g.Expect(myTxn.DebitCreditSet[0].TransactionDCID).NotTo(gomega.BeZero())
	g.Expect(myTxn.DebitCreditSet[0].TransactionID).To(gomega.Equal(txn.TransactionID))
	g.Expect(myTxn.DebitCreditSet[0].AccountID).To(gomega.Equal(a2.AccountID))
	g.Expect(myTxn.DebitCreditSet[0].DebitOrCredit).To(gomega.Equal(datastore.AccountSignCredit))
	g.Expect(myTxn.DebitCreditSet[0].TransactionDCAmount).To(gomega.Equal(uint64(10000)))
	g.Expect(myTxn.DebitCreditSet[1].TransactionDCID).NotTo(gomega.BeZero())
	g.Expect(myTxn.DebitCreditSet[1].TransactionID).To(gomega.Equal(txn.TransactionID))
	g.Expect(myTxn.DebitCreditSet[1].AccountID).To(gomega.Equal(a1.AccountID))
	g.Expect(myTxn.DebitCreditSet[1].DebitOrCredit).To(gomega.Equal(datastore.AccountSignDebit))
	g.Expect(myTxn.DebitCreditSet[1].TransactionDCAmount).To(gomega.Equal(uint64(10000)))
}

func TestTransaction_StoreAndUpdate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDB(g)

	// create an account first
	a1 := Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err := a1.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a2 := Account{AccountName: "Income", AccountSign: datastore.AccountSignCredit, AccountType: datastore.AccountTypeIncome}
	err = a2.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a3 := Account{AccountName: "Other Income", AccountSign: datastore.AccountSignCredit, AccountType: datastore.AccountTypeIncome}
	err = a3.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	txn := Transaction{TransactionCore: TransactionCore{TransactionComment: "woot"},
		DebitCreditSet: []*TransactionDebitCredit{
			&TransactionDebitCredit{AccountID: a2.AccountID,
				DebitOrCredit:       datastore.AccountSignCredit,
				TransactionDCAmount: 10000},
			&TransactionDebitCredit{AccountID: a1.AccountID,
				DebitOrCredit:       datastore.AccountSignDebit,
				TransactionDCAmount: 10000},
		},
	}
	err = txn.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	txn.TransactionComment = "updated woot"
	txn.DebitCreditSet[0].TransactionDCAmount = 33000
	txn.DebitCreditSet[1].TransactionDCAmount = 33000

	txn.DebitCreditSet[0].AccountID = a3.AccountID

	err = txn.Update(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	myTxn, err := RetrieveTransactionByID(testDS, txn.TransactionID)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(myTxn.TransactionComment).To(gomega.Equal("updated woot"))
	g.Expect(myTxn.TransactionAmount).To(gomega.Equal(uint64(33000)))
	g.Expect(myTxn.TransactionID).NotTo(gomega.BeZero())
	g.Expect(myTxn.DebitCreditSet).To(gomega.HaveLen(2))
	g.Expect(myTxn.DebitCreditSet[0].TransactionID).To(gomega.Equal(txn.TransactionID))
	g.Expect(myTxn.DebitCreditSet[0].AccountID).To(gomega.Equal(a3.AccountID))
	g.Expect(myTxn.DebitCreditSet[0].DebitOrCredit).To(gomega.Equal(datastore.AccountSignCredit))
	g.Expect(myTxn.DebitCreditSet[0].TransactionDCAmount).To(gomega.Equal(uint64(33000)))
	g.Expect(myTxn.DebitCreditSet[1].TransactionID).To(gomega.Equal(txn.TransactionID))
	g.Expect(myTxn.DebitCreditSet[1].AccountID).To(gomega.Equal(a1.AccountID))
	g.Expect(myTxn.DebitCreditSet[1].DebitOrCredit).To(gomega.Equal(datastore.AccountSignDebit))
	g.Expect(myTxn.DebitCreditSet[1].TransactionDCAmount).To(gomega.Equal(uint64(33000)))
}
