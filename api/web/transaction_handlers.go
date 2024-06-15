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
	"time"
)

// GET /transactions/{transactionID}
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

// POST /transactions
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
			return NewRequestError(http.StatusBadRequest, fmt.Errorf("reqTransaction:%+v %w", reqTransaction, err))
		}
		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(w, jsonResponse)
	}
}

// GET /transactions/account/{accountID}
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
		account, transactions, err := contoller.GetTransactionsForAccount(r.Context(), accountID)
		if err != nil {
			return NewRequestError(http.StatusNotFound, err)
		}
		jsonResponse := response.ConvertTransactionLedgerToRespTransactionLedger(account, transactions)
		return RespondOK(w, jsonResponse)
	}
}

var ErrInvalidReconcileDate = errors.New("invalid reconcile date")

// GET /transactions/account/{accountID}/unreconciled?date=<date>
func GetUnreconciledTransactionsOnAccount(contoller *TransactionsController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		accountIDStr := chi.URLParam(r, "accountID")
		accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		if accountID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}
		dateStr := r.URL.Query().Get("date")
		dateCutoff, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, ErrInvalidReconcileDate)
		}
		transactions, err := contoller.GetUnreconciledTransactionsOnAccount(r.Context(), accountID, dateCutoff)
		if err != nil {
			return NewRequestError(http.StatusNotFound, err)
		}
		jsonResponse := response.ConvertTransactionRecSetToRespTransactionRecSet(transactions)
		return RespondOK(w, jsonResponse)
	}
}

// PUT /transactions/{transactionID}
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

// DELETE /transactions/{transactionID}
func DeleteTransaction(contoller *TransactionsController) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		idStr := chi.URLParam(r, "transactionID")
		transactionID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		if transactionID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}
		transaction, err := contoller.DeleteTransaction(r.Context(), transactionID)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(w, jsonResponse)
	}
}

// PUT /transactions/{transactionID}/reconciled
func PutTransactionReconciled(contoller *TransactionsController) func(w http.ResponseWriter, r *http.Request) error {
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
		transaction, err := contoller.UpdateReconciled(r.Context(), mdlTransaction)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(w, jsonResponse)
	}
}

// PUT /transactions/{transactionID}/unreconciled
func PutTransactionUnreconciled(contoller *TransactionsController) func(w http.ResponseWriter, r *http.Request) error {
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
		transaction, err := contoller.UpdateUnreconciled(r.Context(), mdlTransaction)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}
		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(w, jsonResponse)
	}
}
