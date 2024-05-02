package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mimirsoft/mimirledger/api/web/request"
	"github.com/mimirsoft/mimirledger/api/web/response"
	"net/http"
	"strconv"
)

// GET /accounttypes
func GetAccountTypes(acctController *AccountsController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountsTypes, err := acctController.AccountTypeList(r.Context())
		if err != nil {
			return NewRequestError(http.StatusServiceUnavailable, err)
		}
		return RespondOK(w, accountsTypes)
	}
}

// GET /accounts
func GetAccounts(acctController *AccountsController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		accounts, err := acctController.AccountList(r.Context())
		if err != nil {
			return NewRequestError(http.StatusServiceUnavailable, err)
		}
		jsonResponse := response.ConvertAccountsToRespAccountSet(accounts)
		return RespondOK(w, jsonResponse)
	}
}

var ErrInvalidAccountID = errors.New("invalid accountID request parameter")

// GET /accounts/{accountID}
func GetAccount(acctController *AccountsController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountIDStr := chi.URLParam(r, "accountID")
		accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		if accountID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}
		account, err := acctController.AccountGetByID(r.Context(), accountID)
		if err != nil {
			return NewRequestError(http.StatusNotFound, err)
		}
		jsonResponse := response.AccountToRespAccount(account)
		return RespondOK(w, jsonResponse)
	}
}

var ErrNoRequestBody = errors.New("missing request body")

// POST /accounts
func PostAccounts(acctController *AccountsController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		var acct request.Account
		if r.Body == nil {
			return NewRequestError(http.StatusBadRequest, ErrNoRequestBody)
		}
		err := json.NewDecoder(r.Body).Decode(&acct)
		if err != nil {
			return fmt.Errorf("json.NewDecoder(r.Body).Decode:%w", err)
		}
		mdlAccount := request.ReqAccountToAccount(&acct)
		account, err := acctController.CreateAccount(r.Context(), mdlAccount)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		jsonResponse := response.AccountToRespAccount(account)
		return RespondOK(w, jsonResponse)
	}
}

// PUT /accounts/{accountID}
func PutAccountUpdate(acctController *AccountsController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountIDStr := chi.URLParam(r, "accountID")
		accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		if accountID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}
		var acct request.Account
		if r.Body == nil {
			return NewRequestError(http.StatusBadRequest, ErrNoRequestBody)
		}
		err = json.NewDecoder(r.Body).Decode(&acct)
		if err != nil {
			return fmt.Errorf("json.NewDecoder(r.Body).Decode:%w", err)
		}
		mdlAccount := request.ReqAccountToAccount(&acct)
		mdlAccount.AccountID = accountID
		account, err := acctController.UpdateAccount(r.Context(), mdlAccount)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		jsonResponse := response.AccountToRespAccount(account)
		return RespondOK(w, jsonResponse)
	}
}
