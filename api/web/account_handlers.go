package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/mimirsoft/mimirledger/api/web/request"
	"github.com/mimirsoft/mimirledger/api/web/response"
)

// GET /accounttypes
func GetAccountTypes(acctController *AccountsController) func(res http.ResponseWriter, reqr *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		accountsTypes, err := acctController.AccountTypeList(req.Context())
		if err != nil {
			return NewRequestError(http.StatusServiceUnavailable, err)
		}

		return RespondOK(res, accountsTypes)
	}
}

// GET /accounts
func GetAccounts(acctController *AccountsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		accounts, err := acctController.AccountList(req.Context())
		if err != nil {
			return NewRequestError(http.StatusServiceUnavailable, err)
		}

		jsonResponse := response.ConvertAccountsToRespAccountSet(accounts)

		return RespondOK(res, jsonResponse)
	}
}

var ErrInvalidAccountID = errors.New("invalid accountID request parameter")

// GET /accounts/{accountID}
func GetAccount(acctController *AccountsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		accountIDStr := chi.URLParam(req, "accountID")

		accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if accountID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}

		account, err := acctController.AccountGetByID(req.Context(), accountID)
		if err != nil {
			return NewRequestError(http.StatusNotFound, err)
		}

		jsonResponse := response.AccountToRespAccount(account)

		return RespondOK(res, jsonResponse)
	}
}

var ErrNoRequestBody = errors.New("missing request body")

// POST /accounts
func PostAccounts(acctController *AccountsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		var acct request.Account

		if req.Body == nil {
			return NewRequestError(http.StatusBadRequest, ErrNoRequestBody)
		}

		err := json.NewDecoder(req.Body).Decode(&acct)
		if err != nil {
			return fmt.Errorf("json.NewDecoder(r.Body).Decode:%w", err)
		}

		mdlAccount := request.ReqAccountToAccount(&acct)

		account, err := acctController.CreateAccount(req.Context(), mdlAccount)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		jsonResponse := response.AccountToRespAccount(account)

		return RespondOK(res, jsonResponse)
	}
}

// PUT /accounts/{accountID}
func PutAccountUpdate(acctController *AccountsController) func(res http.ResponseWriter, //nolint:dupl
	req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		accountIDStr := chi.URLParam(req, "accountID")

		accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if accountID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}

		var acct request.Account

		if req.Body == nil {
			return NewRequestError(http.StatusBadRequest, ErrNoRequestBody)
		}

		err = json.NewDecoder(req.Body).Decode(&acct)
		if err != nil {
			return fmt.Errorf("json.NewDecoder(r.Body).Decode:%w", err)
		}

		mdlAccount := request.ReqAccountToAccount(&acct)
		mdlAccount.AccountID = accountID

		account, err := acctController.UpdateAccount(req.Context(), mdlAccount)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		jsonResponse := response.AccountToRespAccount(account)

		return RespondOK(res, jsonResponse)
	}
}
