package web

import (
	"github.com/mimirsoft/mimirledger/api/models"
	"github.com/mimirsoft/mimirledger/api/web/response"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"net/http"
	"testing"
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

	NewRouterTableTest([]RouterTest{
		{
			Request: Request{
				Method:     http.MethodGet,
				Router:     TestRouter,
				RequestURL: "/tranasctions/5555",
			},
			GomegaWithT: g,
			Code:        http.StatusNotFound, RespBody: models.ErrTransactionNotFound.Error(),
		},
		{
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/tranasctions",
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: ErrNoRequestBody.Error(),
		},
		{
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/tranasctions",
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
				RequestURL: "/tranasctions",
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
				RequestURL: "/tranasctions",
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
				RequestURL: "/tranasctions",
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
				RequestURL: "/tranasctions",
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
				RequestURL: "/tranasctions",
				Payload: map[string]interface{}{
					"transactionComment": "getting paid",
					"transactionAmount":  10000,
					"debitCreditSet":     mapSlice4,
				},
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: models.ErrTransactionDebitCreditsNotBalanced.Error(),
		},
	}).Exec()
}

func TestTransaction_PostNewTransaction(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	mapSlice4 := []map[string]interface{}{}
	map4a := map[string]interface{}{"transactionDCAmount": 9999, "accountID": 1, "debitOrCredit": "DEBIT"}
	map4b := map[string]interface{}{"transactionDCAmount": 9999, "accountID": 2, "debitOrCredit": "CREDIT"}
	mapSlice4 = append(mapSlice4, map4a, map4b)

	acctReq := map[string]interface{}{
		"transactionComment": "getting paid",
		"transactionAmount":  10000,
		"debitCreditSet":     mapSlice4,
	}
	var test = RouterTest{Request: Request{
		Method:     http.MethodPost,
		Router:     TestRouter,
		RequestURL: "/tranasctions",
		Payload:    acctReq,
	}, GomegaWithT: g, Code: http.StatusOK}

	var res response.Transaction
	test.ExecWithUnmarshal(&res)
	g.Expect(res.TransactionComment).To(gomega.Equal("getting paid"))
	g.Expect(res.TransactionAmount).To(gomega.Equal(uint64(9999)))
	g.Expect(res.DebitCreditSet).To(gomega.HaveLen(2))
}
