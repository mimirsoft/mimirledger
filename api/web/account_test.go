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

func TestAccount_GetAccountsIDInvalid(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	NewRouterTableTest([]RouterTest{
		{
			Request: Request{
				Method:     http.MethodGet,
				Router:     TestRouter,
				RequestURL: "/accounts/5555",
			},
			GomegaWithT: g,
			Code:        http.StatusNotFound, RespBody: models.ErrAccountNotFound.Error(),
		},
		{
			Request: Request{
				Method:     http.MethodGet,
				Router:     TestRouter,
				RequestURL: "/accounts/asdf",
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: "strconv.ParseUint: parsing",
		},
	}).Exec()
}
func TestAccountGetAll(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	var test = RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: "/accounts",
	}, GomegaWithT: g, Code: http.StatusOK}

	var accountSetRes response.AccountSet
	test.ExecWithUnmarshal(&accountSetRes)
	g.Expect(accountSetRes.Accounts).To(gomega.HaveLen(0))

	// store 2 accounts manually
	a1 := models.Account{AccountName: "MY BANK", AccountSign: datastore.AccountSignDebit,
		AccountType: datastore.AccountTypeAsset}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred()) // reset datastore
	a2 := models.Account{AccountName: "ZZZ BANK", AccountSign: datastore.AccountSignDebit,
		AccountType: datastore.AccountTypeAsset}
	err = a2.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred()) // reset datastore

	test = RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: "/accounts",
	}, GomegaWithT: g, Code: http.StatusOK}

	test.ExecWithUnmarshal(&accountSetRes)
	g.Expect(accountSetRes.Accounts).To(gomega.HaveLen(2))
	g.Expect(accountSetRes.Accounts[0].AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(accountSetRes.Accounts[0].AccountRight).To(gomega.Equal(uint64(2)))
	g.Expect(accountSetRes.Accounts[0].AccountName).To(gomega.Equal("MY BANK"))
	g.Expect(accountSetRes.Accounts[0].AccountType).To(gomega.Equal("ASSET"))
	g.Expect(accountSetRes.Accounts[0].AccountSign).To(gomega.Equal("DEBIT"))

	g.Expect(accountSetRes.Accounts[1].AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(accountSetRes.Accounts[1].AccountRight).To(gomega.Equal(uint64(4)))
	g.Expect(accountSetRes.Accounts[1].AccountName).To(gomega.Equal("ZZZ BANK"))
	g.Expect(accountSetRes.Accounts[1].AccountType).To(gomega.Equal("ASSET"))
	g.Expect(accountSetRes.Accounts[1].AccountSign).To(gomega.Equal("DEBIT"))

	// add a third account
	a3 := models.Account{AccountName: "AAA BANK", AccountSign: datastore.AccountSignDebit,
		AccountType: datastore.AccountTypeAsset}
	err = a3.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred()) // reset datastore

	test = RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: "/accounts",
	}, GomegaWithT: g, Code: http.StatusOK}

	test.ExecWithUnmarshal(&accountSetRes)
	g.Expect(accountSetRes.Accounts).To(gomega.HaveLen(3))
	g.Expect(accountSetRes.Accounts[0].AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(accountSetRes.Accounts[0].AccountRight).To(gomega.Equal(uint64(2)))
	g.Expect(accountSetRes.Accounts[0].AccountName).To(gomega.Equal("AAA BANK"))
	g.Expect(accountSetRes.Accounts[0].AccountType).To(gomega.Equal("ASSET"))
	g.Expect(accountSetRes.Accounts[0].AccountSign).To(gomega.Equal("DEBIT"))

	g.Expect(accountSetRes.Accounts[1].AccountName).To(gomega.Equal("MY BANK"))
	g.Expect(accountSetRes.Accounts[1].AccountType).To(gomega.Equal("ASSET"))
	g.Expect(accountSetRes.Accounts[1].AccountSign).To(gomega.Equal("DEBIT"))
	g.Expect(accountSetRes.Accounts[1].AccountLeft).To(gomega.Equal(uint64(3)))
	g.Expect(accountSetRes.Accounts[1].AccountRight).To(gomega.Equal(uint64(4)))

	g.Expect(accountSetRes.Accounts[2].AccountName).To(gomega.Equal("ZZZ BANK"))
	g.Expect(accountSetRes.Accounts[2].AccountType).To(gomega.Equal("ASSET"))
	g.Expect(accountSetRes.Accounts[2].AccountSign).To(gomega.Equal("DEBIT"))
	g.Expect(accountSetRes.Accounts[2].AccountLeft).To(gomega.Equal(uint64(5)))
	g.Expect(accountSetRes.Accounts[2].AccountRight).To(gomega.Equal(uint64(6)))
}

func TestAccountsController_AccountGetByID(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	// store 1 account manually
	a1 := models.Account{AccountName: "MY BANK", AccountSign: datastore.AccountSignDebit,
		AccountType: datastore.AccountTypeAsset}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred()) // reset datastore

	var test = RouterTest{Request: Request{
		Method:     http.MethodGet,
		Router:     TestRouter,
		RequestURL: fmt.Sprintf("/accounts/%d", a1.AccountID),
	}, GomegaWithT: g, Code: http.StatusOK}

	var myAcct response.Account

	test.ExecWithUnmarshal(&myAcct)
	g.Expect(myAcct.AccountName).To(gomega.Equal("MY BANK"))
	g.Expect(myAcct.AccountType).To(gomega.Equal("ASSET"))
	g.Expect(myAcct.AccountSign).To(gomega.Equal("DEBIT"))
	g.Expect(myAcct.AccountLeft).To(gomega.Equal(uint64(1)))
	g.Expect(myAcct.AccountRight).To(gomega.Equal(uint64(2)))
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
					"accountName": "my bank acct",
				},
			},
			GomegaWithT: g,
			Code:        http.StatusBadRequest, RespBody: "accountType is not valid, cannot determine AccountSign",
		},
	}).Exec()
}

func TestAccountPostAccounts(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	acctReq := map[string]interface{}{
		"accountName": "my bank",
		"accountType": "ASSET",
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

func TestAccountPost_WithParent(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	// store 1 accounts manually
	a1 := models.Account{AccountName: "MY BANK", AccountSign: datastore.AccountSignDebit,
		AccountType: datastore.AccountTypeAsset}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred()) // reset datastore

	acctReq := map[string]interface{}{
		"accountName":   "my_bank_sub",
		"accountParent": a1.AccountID,
	}
	var test = RouterTest{Request: Request{
		Method:     http.MethodPost,
		Router:     TestRouter,
		RequestURL: "/accounts",
		Payload:    acctReq,
	}, GomegaWithT: g, Code: http.StatusOK}

	var accountRes response.Account
	test.ExecWithUnmarshal(&accountRes)
	g.Expect(accountRes.AccountName).To(gomega.Equal("my_bank_sub"))
	g.Expect(accountRes.AccountType).To(gomega.Equal("ASSET"))
	g.Expect(accountRes.AccountSign).To(gomega.Equal("DEBIT"))
	g.Expect(accountRes.AccountParent).To(gomega.Equal(a1.AccountID))
	g.Expect(accountRes.AccountFullName).To(gomega.Equal("MY BANK:my_bank_sub"))
}

func TestAccountPost_Update(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	// store 1 account manually
	a1 := models.Account{AccountName: "MY BANK", AccountSign: datastore.AccountSignDebit,
		AccountType: datastore.AccountTypeAsset}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred()) // reset datastore

	acctReq := map[string]interface{}{
		"accountName": "my bank",
		"accountType": "ASSET",
	}
	// create the account
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

	acctReq = map[string]interface{}{
		"accountName": "my_bank_update",
		"accountType": "ASSET",
	}
	// update the account name
	var test2 = RouterTest{Request: Request{
		Method:     http.MethodPut,
		Router:     TestRouter,
		RequestURL: fmt.Sprintf("/accounts/%d", accountRes.AccountID),
		Payload:    acctReq,
	}, GomegaWithT: g, Code: http.StatusOK}

	var resAcct response.Account
	test2.ExecWithUnmarshal(&resAcct)
	g.Expect(resAcct.AccountName).To(gomega.Equal("my_bank_update"))
	g.Expect(resAcct.AccountType).To(gomega.Equal("ASSET"))
	g.Expect(resAcct.AccountSign).To(gomega.Equal("DEBIT"))
	g.Expect(resAcct.AccountParent).To(gomega.Equal(uint64(0)))
	g.Expect(resAcct.AccountFullName).To(gomega.Equal("my_bank_update"))

	acctReq = map[string]interface{}{
		"accountName":   "my_bank_update_sub",
		"accountParent": a1.AccountID,
		"accountMemo":   "updated_memo",
	}
	var test3 = RouterTest{Request: Request{
		Method:     http.MethodPut,
		Router:     TestRouter,
		RequestURL: fmt.Sprintf("/accounts/%d", accountRes.AccountID),
		Payload:    acctReq,
	}, GomegaWithT: g, Code: http.StatusOK}

	test3.ExecWithUnmarshal(&resAcct)
	g.Expect(resAcct.AccountName).To(gomega.Equal("my_bank_update_sub"))
	g.Expect(resAcct.AccountType).To(gomega.Equal("ASSET"))
	g.Expect(resAcct.AccountSign).To(gomega.Equal("DEBIT"))
	g.Expect(resAcct.AccountParent).To(gomega.Equal(a1.AccountID))
	g.Expect(resAcct.AccountFullName).To(gomega.Equal("MY BANK:my_bank_update_sub"))
}

func TestAccount_PutUpdateAccountReconciledDate(t *testing.T) {
	g := gomega.NewWithT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupDatastores(TestDataStore)

	// store 1 account manually
	a1 := models.Account{AccountName: "MY BANK", AccountSign: datastore.AccountSignDebit,
		AccountType: datastore.AccountTypeAsset}
	err := a1.Store(TestDataStore)
	g.Expect(err).NotTo(gomega.HaveOccurred()) // reset datastore

	acctReq := map[string]interface{}{
		"accountName": "my bank",
		"accountType": "ASSET",
	}
	// create the account
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
	g.Expect(accountRes.AccountReconcileDate).To(gomega.Equal(time.Time{}))

	oldDate, err := time.Parse("2006-01-02", "2016-07-08")
	g.Expect(err).NotTo(gomega.HaveOccurred())

	acctReq = map[string]interface{}{
		"accountReconcileDate": oldDate.Format(time.RFC3339),
	}
	// update the account name
	var test2 = RouterTest{Request: Request{
		Method:     http.MethodPut,
		Router:     TestRouter,
		RequestURL: fmt.Sprintf("/accounts/%d/reconciled", accountRes.AccountID),
		Payload:    acctReq,
	}, GomegaWithT: g, Code: http.StatusOK}

	var resAcct response.Account
	test2.ExecWithUnmarshal(&resAcct)
	g.Expect(resAcct.AccountName).To(gomega.Equal("my bank"))
	g.Expect(resAcct.AccountType).To(gomega.Equal("ASSET"))
	g.Expect(resAcct.AccountSign).To(gomega.Equal("DEBIT"))
	g.Expect(resAcct.AccountParent).To(gomega.Equal(uint64(0)))
	g.Expect(resAcct.AccountFullName).To(gomega.Equal("my bank"))
	g.Expect(resAcct.AccountReconcileDate).To(gomega.Equal(oldDate))
}
