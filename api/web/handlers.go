package web

import (
	"encoding/json"
	"fmt"
	"github.com/mimirsoft/mimirledger/api/models"
	"github.com/mimirsoft/mimirledger/api/web/request"
	"github.com/mimirsoft/mimirledger/api/web/response"
	"github.com/rs/zerolog"
	"net/http"
)

type UnknownError struct {
	StatusCode int    `json:"status_code,omitempty"`
	Err        string `json:"error,omitempty"`
}

type RequestError struct {
	StatusCode int    `json:"status_code,omitempty"`
	Err        string `json:"error,omitempty"`
}

func (r *RequestError) Error() string {
	return r.Err
}

func NewRequestError(status int, recErr error) *RequestError {
	return &RequestError{StatusCode: status,
		Err: recErr.Error()}
}

// the root handler takes the normal arguments of an http handler
// it allows the logging and error handling to be in one place
// The Handler struct that takes a function matching our useful signature.
type RootHandler struct {
	H func(w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
// The root handler allows unified logging and error handling for the application
func (h RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := zerolog.Ctx(r.Context())
	err := h.H(w, r)

	if err != nil {
		switch e := err.(type) {
		case *RequestError:
			// since this was a structured RequestError, we do not log and just
			// return the response.  Logging middleware will handle the log
			w.WriteHeader(e.StatusCode)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			if err = json.NewEncoder(w).Encode(e); err != nil {
				l.Printf("failed json encode RequestError response %s", err)
			}
			return
		default:
			// since we don't know what the error is, we log the error
			// then we return a generic StatusInternalServerError with a JSON payload
			l.Error().Err(err).Msgf("HTTP %d", http.StatusInternalServerError)
			unknownErr := UnknownError{StatusCode: http.StatusInternalServerError, Err: err.Error()}
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			if err1 := json.NewEncoder(w).Encode(unknownErr); err1 != nil {
				l.Error().Err(err1).Msg("failed json encode UnknownError response")
			}
			return
		}
	}
}

// returns our universal custom root handler
func NewRootHandler(handler func(w http.ResponseWriter, r *http.Request) error) *RootHandler {
	return &RootHandler{H: handler}
}

// GET /api/health
// HEAD /api/health
func HealthCheck(healthController *HealthController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := healthController.HealthCheck(r.Context())
		if err != nil {
			return NewRequestError(http.StatusServiceUnavailable, err)
		}
		return RespondOK(w, "ok")
	}
}

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

// POST /accounts
func PostAccounts(acctController *AccountsController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		var acct request.Account
		err := json.NewDecoder(r.Body).Decode(&acct)
		if err != nil {
			return fmt.Errorf("son.NewDecoder(r.Body).Decode:%w", err)
		}
		mdlAccount := models.Account(acct)
		account, err := acctController.CreateAccount(r.Context(), mdlAccount)
		if err != nil {
			return NewRequestError(http.StatusServiceUnavailable, err)
		}
		jsonResponse := response.Account(*account)
		return RespondOK(w, jsonResponse)
	}
}
