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

func TestTransaction_StoreValid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)

	txn := Transaction{TransactionCore: TransactionCore{TransactionComment: "woot"},
		DebitCreditSet: []*TransactionDebitCredit{
			&TransactionDebitCredit{AccountID: 1, DebitOrCredit: datastore.AccountSignCredit, TransactionDCAmount: 10000},
			&TransactionDebitCredit{AccountID: 1, DebitOrCredit: datastore.AccountSignDebit, TransactionDCAmount: 10000},
		},
	}
	err := txn.Store(testDS)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
