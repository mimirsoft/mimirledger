package web

import (
	"fmt"
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
	"github.com/mimirsoft/mimirledger/api/web/response"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"net/http"
	"testing"
	"time"
)

// M is an alias for map[string]interface{}
type M map[string]interface{}

func TestTransaction_Invalid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	mapSlice := []map[string]interface{}{}
	map1 := map[string]interface{}{"transactionDCAmount": 9999}
	mapSlice = append(mapSlice, map1)

	mapSlice2 := []map[string]interface{}{}
	map2 := map[string]interface{}{"transactionDCAmount": 9999, "accountID": 2}
	mapSlice2 = append(mapSlice2, map2)

	mapSlice3 := []map[string]interface{}{}
	map3 := map[string]interface{}{"transactionDCAmount": 9999, "accountID": 2, "debitOrCredit": "DEBIT"}
	mapSlice3 = append(mapSlice3, map3)

	mapSlice4 := []map[string]interface{}{}
	map4 := map[string]interface{}{"transactionDCAmount": 9999, "accountID": 2, "debitOrCredit": "DEBIT"}
	map4b := map[string]interface{}{"transactionDCAmount": 1999, "accountID": 2, "debitOrCredit": "CREDIT"}
	mapSlice4 = append(mapSlice4, map4, map4b)

	mapSlice5 := []map[string]interface{}{}
	map5a := map[string]interface{}{"transactionDCAmount": 3000, "accountID": 5555, "debitOrCredit": "DEBIT"}
	map5b := map[string]interface{}{"transactionDCAmount": 3000, "accountID": 5555, "debitOrCredit": "CREDIT"}
	mapSlice5 = append(mapSlice5, map5a, map5b)

	NewRouterTableTest([]RouterTest{
		{
			Request: Request{
				Method:     http.MethodGet,
				Router:     TestRouter,
				RequestURL: "/transactions/5555",
			},
			GomegaWithT: g,
			Code:        http.StatusNotFound, RespBody: models.ErrTransactionNotFound.Error(),
		},
		{
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/transactions",
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: ErrNoRequestBody.Error(),
		},
		{
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/transactions",
				Payload: map[string]interface{}{
					"woot": "stuff",
				},
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: models.ErrTransactionNoComment.Error(),
		},
		{
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/transactions",
				Payload: map[string]interface{}{
					"transactionComment": "getting paid",
				},
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: models.ErrTransactionNoDebitsCredits.Error(),
		},
		{
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/transactions",
				Payload: map[string]interface{}{
					"transactionComment": "getting paid",
					"debitCreditSet":     mapSlice,
				},
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: models.ErrTransactionDebitCreditAccountInvalid.Error(),
		}, {
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/transactions",
				Payload: map[string]interface{}{
					"transactionComment": "getting paid",
					"transactionAmount":  10000,
					"debitCreditSet":     mapSlice2,
				},
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: models.ErrTransactionDebitCreditsIsNeither.Error(),
		},
		{
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/transactions",
				Payload: map[string]interface{}{
					"transactionComment": "getting paid",
					"transactionAmount":  10000,
					"debitCreditSet":     mapSlice3,
				},
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: models.ErrTransactionDebitCreditsZero.Error(),
		},
		{
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/transactions",
				Payload: map[string]interface{}{
					"transactionComment": "getting paid",
					"transactionAmount":  10000,
					"debitCreditSet":     mapSlice4,
				},
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: models.ErrTransactionDebitCreditsNotBalanced.Error(),
		},
		{
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/transactions",
				Payload: map[string]interface{}{
					"transactionComment": "getting paid",
					"transactionAmount":  10000,
					"debitCreditSet":     mapSlice5,
				},
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: "myTxn.Store:ds.TransactionDCStore().Store:ERROR: insert or update on table \\\"transaction_debit_credit\\\"",
		},
	}).Exec()
}

func TestTransaction_PostNewTransaction(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	// create an account first
	a1 := models.Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a2 := models.Account{AccountName: "Income", AccountSign: datastore.AccountSignCredit, AccountType: datastore.AccountTypeIncome}
	err = a2.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	mapSlice4 := []map[string]interface{}{}
	map4a := map[string]interface{}{"transactionDCAmount": 9999, "accountID": a1.AccountID, "debitOrCredit": "DEBIT"}
	map4b := map[string]interface{}{"transactionDCAmount": 9999, "accountID": a2.AccountID, "debitOrCredit": "CREDIT"}
	mapSlice4 = append(mapSlice4, map4a, map4b)

	acctReq := map[string]interface{}{
		"transactionComment": "getting paid",
		"transactionAmount":  10000,
		"debitCreditSet":     mapSlice4,
	}
	var test = RouterTest{Request: Request{
		Method:     http.MethodPost,
		Router:     TestRouter,
		RequestURL: "/transactions",
		Payload:    acctReq,
	}, GomegaWithT: g, Code: http.StatusOK}

	var res response.Transaction
	test.ExecWithUnmarshal(&res)
	g.Expect(res.TransactionComment).To(gomega.Equal("getting paid"))
	g.Expect(res.TransactionAmount).To(gomega.Equal(uint64(9999)))
	g.Expect(res.DebitCreditSet).To(gomega.HaveLen(2))
	g.Expect(res.TransactionDate).To(gomega.BeTemporally("~", time.Now(), time.Second))
}

func TestTransaction_PutTransactionUpdate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	// create accounts first
	a1 := models.Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a2 := models.Account{AccountName: "Income", AccountSign: datastore.AccountSignCredit, AccountType: datastore.AccountTypeIncome}
	err = a2.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a3 := models.Account{AccountName: "OtherIncome", AccountSign: datastore.AccountSignCredit, AccountType: datastore.AccountTypeIncome}
	err = a3.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	txn := models.Transaction{TransactionCore: models.TransactionCore{TransactionComment: "woot"},
		DebitCreditSet: []*models.TransactionDebitCredit{
			&models.TransactionDebitCredit{AccountID: a2.AccountID,
				DebitOrCredit:       datastore.AccountSignCredit,
				TransactionDCAmount: 10000},
			&models.TransactionDebitCredit{AccountID: a1.AccountID,
				DebitOrCredit:       datastore.AccountSignDebit,
				TransactionDCAmount: 10000},
		},
	}
	err = txn.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	mapSlice4 := []map[string]interface{}{}
	map4a := map[string]interface{}{"transactionDCAmount": 34000, "accountID": a1.AccountID, "debitOrCredit": "DEBIT"}
	map4b := map[string]interface{}{"transactionDCAmount": 34000, "accountID": a3.AccountID, "debitOrCredit": "CREDIT"}
	mapSlice4 = append(mapSlice4, map4b, map4a)

	oldDate, err := time.Parse("2006-01-02", "2016-07-08")
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acctReq := map[string]interface{}{
		"transactionDate":    oldDate.Format(time.RFC3339),
		"transactionComment": "getting paid from other person",
		"transactionAmount":  10000,
		"debitCreditSet":     mapSlice4,
	}
	var test = RouterTest{Request: Request{
		Method:     http.MethodPut,
		Router:     TestRouter,
		RequestURL: fmt.Sprintf("/transactions/%d", txn.TransactionID),
		Payload:    acctReq,
	}, GomegaWithT: g, Code: http.StatusOK}

	var res response.Transaction
	test.ExecWithUnmarshal(&res)
	g.Expect(res.TransactionComment).To(gomega.Equal("getting paid from other person"))
	g.Expect(res.TransactionAmount).To(gomega.Equal(uint64(34000)))
	g.Expect(res.DebitCreditSet).To(gomega.HaveLen(2))
	g.Expect(res.DebitCreditSet).To(gomega.HaveLen(2))
	g.Expect(res.DebitCreditSet[0].TransactionID).To(gomega.Equal(txn.TransactionID))
	g.Expect(res.DebitCreditSet[0].AccountID).To(gomega.Equal(a3.AccountID))
	g.Expect(res.DebitCreditSet[0].DebitOrCredit).To(gomega.Equal(datastore.AccountSignCredit))
	g.Expect(res.DebitCreditSet[0].TransactionDCAmount).To(gomega.Equal(uint64(34000)))
	g.Expect(res.DebitCreditSet[1].TransactionID).To(gomega.Equal(txn.TransactionID))
	g.Expect(res.DebitCreditSet[1].AccountID).To(gomega.Equal(a1.AccountID))
	g.Expect(res.DebitCreditSet[1].DebitOrCredit).To(gomega.Equal(datastore.AccountSignDebit))
	g.Expect(res.DebitCreditSet[1].TransactionDCAmount).To(gomega.Equal(uint64(34000)))
	g.Expect(res.TransactionDate).To(gomega.BeTemporally("~", oldDate, time.Second))
}

func TestTransaction_GetTransactionsOnAccountEmpty(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	// create accounts first
	a1 := models.Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	var test = RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: fmt.Sprintf("/transactions/account/%d", a1.AccountID),
	}, GomegaWithT: g, Code: http.StatusOK}

	var res response.TransactionLedgerSet
	test.ExecWithUnmarshal(&res)
	g.Expect(res.Transactions).To(gomega.HaveLen(0))
	g.Expect(res.AccountID).To(gomega.Equal(a1.AccountID))
	g.Expect(res.AccountFullName).To(gomega.Equal(a1.AccountFullName))
	g.Expect(res.AccountName).To(gomega.Equal(string(a1.AccountName)))
	g.Expect(res.AccountSign).To(gomega.Equal(string(a1.AccountSign)))
}

func TestTransaction_GetTransactionsOnAccountValid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	// create accounts first
	a1 := models.Account{AccountName: "MyBank", AccountSign: datastore.AccountSignDebit, AccountType: datastore.AccountTypeAsset}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	a2 := models.Account{AccountName: "Income", AccountSign: datastore.AccountSignCredit, AccountType: datastore.AccountTypeIncome}
	err = a2.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	txn := models.Transaction{TransactionCore: models.TransactionCore{TransactionComment: "woot"},
		DebitCreditSet: []*models.TransactionDebitCredit{
			&models.TransactionDebitCredit{AccountID: a2.AccountID,
				DebitOrCredit:       datastore.AccountSignCredit,
				TransactionDCAmount: 10000},
			&models.TransactionDebitCredit{AccountID: a1.AccountID,
				DebitOrCredit:       datastore.AccountSignDebit,
				TransactionDCAmount: 10000},
		},
	}
	err = txn.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	txn2 := models.Transaction{TransactionCore: models.TransactionCore{TransactionComment: "woot2"},
		DebitCreditSet: []*models.TransactionDebitCredit{
			&models.TransactionDebitCredit{AccountID: a2.AccountID,
				DebitOrCredit:       datastore.AccountSignCredit,
				TransactionDCAmount: 30000},
			&models.TransactionDebitCredit{AccountID: a1.AccountID,
				DebitOrCredit:       datastore.AccountSignDebit,
				TransactionDCAmount: 30000},
		},
	}
	err = txn2.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	var test = RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: fmt.Sprintf("/transactions/account/%d", a2.AccountID),
	}, GomegaWithT: g, Code: http.StatusOK}

	var res response.TransactionLedgerSet
	test.ExecWithUnmarshal(&res)
	g.Expect(res.Transactions).To(gomega.HaveLen(2))
	g.Expect(res.AccountID).To(gomega.Equal(a2.AccountID))
	g.Expect(res.AccountSign).To(gomega.Equal(string(a2.AccountSign)))
	g.Expect(res.AccountFullName).To(gomega.Equal(a2.AccountFullName))
	g.Expect(res.AccountName).To(gomega.Equal(string(a2.AccountName)))
	g.Expect(res.Transactions[0].TransactionComment).To(gomega.Equal("woot"))
	g.Expect(res.Transactions[0].TransactionDCAmount).To(gomega.Equal(uint64(10000)))
	g.Expect(res.Transactions[0].TransactionDate).To(gomega.BeTemporally("~", time.Now(), time.Second))
	g.Expect(res.Transactions[1].TransactionComment).To(gomega.Equal("woot2"))
	g.Expect(res.Transactions[1].TransactionDCAmount).To(gomega.Equal(uint64(30000)))
	g.Expect(res.Transactions[1].TransactionDate).To(gomega.BeTemporally("~", time.Now(), time.Second))

}
