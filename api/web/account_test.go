package web

import (
	"github.com/mimirsoft/mimirledger/api/datastore"
	"github.com/mimirsoft/mimirledger/api/models"
	"github.com/mimirsoft/mimirledger/api/web/response"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"net/http"
	"testing"
)

func TestAccountGetAll(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	// store 2 accounts manually
	a1 := models.Account{AccountName: "MY BANK", AccountSign: datastore.AccountSignDebit,
		AccountType: datastore.AccountTypeAsset}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred()) // reset datastore
	a2 := models.Account{AccountName: "MY OTHER BANK", AccountSign: datastore.AccountSignDebit,
		AccountType: datastore.AccountTypeAsset}
	err = a2.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred()) // reset datastore

	var test = RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: "/accounts",
	}, GomegaWithT: g, Code: http.StatusOK}

	var accountSetRes response.AccountSet
	test.ExecWithUnmarshal(&accountSetRes)
	g.Expect(accountSetRes.Accounts).To(gomega.HaveLen(2))
	g.Expect(accountSetRes.Accounts[0].AccountName).To(gomega.Equal("MY BANK"))
	g.Expect(accountSetRes.Accounts[0].AccountType).To(gomega.Equal("ASSET"))
	g.Expect(accountSetRes.Accounts[0].AccountSign).To(gomega.Equal("DEBIT"))

	g.Expect(accountSetRes.Accounts[1].AccountName).To(gomega.Equal("MY OTHER BANK"))
	g.Expect(accountSetRes.Accounts[1].AccountType).To(gomega.Equal("ASSET"))
	g.Expect(accountSetRes.Accounts[1].AccountSign).To(gomega.Equal("DEBIT"))

}

func TestAccountPostAccountsInvalid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)
	NewRouterTableTest([]RouterTest{
		{
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/accounts",
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: ErrNoRequestBody.Error(),
		},
		{
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/accounts",
				Payload: map[string]interface{}{
					"account_name": "my bank acct",
				},
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: "invalid input value for enum transaction_account_sign_type",
		},
		{
			Request: Request{
				Method:     http.MethodPost,
				Router:     TestRouter,
				RequestURL: "/accounts",
				Payload: map[string]interface{}{
					"account_name": "my bank acct",
					"account_sign": "DEBIT",
				},
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: "invalid input value for enum transaction_account_type",
		},
	}).Exec()
}

func TestAccountPostAccounts(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	acctReq := map[string]interface{}{
		"account_name": "my bank",
		"account_type": "ASSET",
		"account_sign": "DEBIT",
	}
	var test = RouterTest{Request: Request{
		Method:     http.MethodPost,
		Router:     TestRouter,
		RequestURL: "/accounts",
		Payload:    acctReq,
	}, GomegaWithT: g, Code: http.StatusOK}

	var accountRes response.Account
	test.ExecWithUnmarshal(&accountRes)
	g.Expect(accountRes.AccountName).To(gomega.Equal("my bank"))
	g.Expect(accountRes.AccountType).To(gomega.Equal("ASSET"))
	g.Expect(accountRes.AccountSign).To(gomega.Equal("DEBIT"))
}
