package web

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mimirsoft/mimirledger/api/web/request"
	"github.com/mimirsoft/mimirledger/api/web/response"
	"net/http"
	"strconv"
)

// GET /tranasctions/{transactionID}
func GetTransaction(contoller *TransactionsController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		idStr := chi.URLParam(r, "transactionID")
		transactionID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		if transactionID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}
		transaction, err := contoller.GetTransactionByID(r.Context(), transactionID)
		if err != nil {
			return NewRequestError(http.StatusNotFound, err)
		}
		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(w, jsonResponse)
	}
}

// POST /tranasctions
func PostTransactions(contoller *TransactionsController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		var reqTransaction request.Transaction
		if r.Body == nil {
			return NewRequestError(http.StatusBadRequest, ErrNoRequestBody)
		}
		err := json.NewDecoder(r.Body).Decode(&reqTransaction)
		if err != nil {
			return fmt.Errorf("json.NewDecoder(r.Body).Decode:%w", err)
		}
		mdlTransaction := request.ReqTransactionToTransaction(&reqTransaction)
		transaction, err := contoller.CreateTransaction(r.Context(), mdlTransaction)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(w, jsonResponse)
	}
}

// GET /tranasctions/account/{accountID}
func GetTransactionsOnAccount(contoller *TransactionsController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountIDStr := chi.URLParam(r, "accountID")
		accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		if accountID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}
		transactions, err := contoller.GetTransactionsForAccount(r.Context(), accountID)
		if err != nil {
			return NewRequestError(http.StatusNotFound, err)
		}
		jsonResponse := response.ConvertTransactionLedgerToRespTransactionLedger(transactions)
		return RespondOK(w, jsonResponse)
	}
}

// PUT /tranasctions/{transactionID}
func PutTransactionUpdate(contoller *TransactionsController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		idStr := chi.URLParam(r, "transactionID")
		transactionID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		if transactionID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}
		var reqTransaction request.Transaction
		if r.Body == nil {
			return NewRequestError(http.StatusBadRequest, ErrNoRequestBody)
		}
		err = json.NewDecoder(r.Body).Decode(&reqTransaction)
		if err != nil {
			return fmt.Errorf("json.NewDecoder(r.Body).Decode:%w", err)
		}
		mdlTransaction := request.ReqTransactionToTransaction(&reqTransaction)
		mdlTransaction.TransactionID = transactionID
		transaction, err := contoller.UpdateTransaction(r.Context(), mdlTransaction)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(w, jsonResponse)
	}
}
