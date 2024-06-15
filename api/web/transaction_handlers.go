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
func GetTransaction(contoller *TransactionsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		idStr := chi.URLParam(req, "transactionID")

		transactionID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if transactionID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}

		transaction, err := contoller.GetTransactionByID(req.Context(), transactionID)
		if err != nil {
			return NewRequestError(http.StatusNotFound, err)
		}

		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(res, jsonResponse)
	}
}

// POST /transactions
func PostTransactions(contoller *TransactionsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		var reqTransaction request.Transaction

		if req.Body == nil {
			return NewRequestError(http.StatusBadRequest, ErrNoRequestBody)
		}

		err := json.NewDecoder(req.Body).Decode(&reqTransaction)
		if err != nil {
			return fmt.Errorf("json.NewDecoder(r.Body).Decode:%w", err)
		}

		mdlTransaction := request.ReqTransactionToTransaction(&reqTransaction)

		transaction, err := contoller.CreateTransaction(req.Context(), mdlTransaction)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, fmt.Errorf("reqTransaction:%+v %w", reqTransaction, err))
		}

		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(res, jsonResponse)
	}
}

// GET /transactions/account/{accountID}
func GetTransactionsOnAccount(contoller *TransactionsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		accountIDStr := chi.URLParam(req, "accountID")

		accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if accountID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}

		account, transactions, err := contoller.GetTransactionsForAccount(req.Context(), accountID)
		if err != nil {
			return NewRequestError(http.StatusNotFound, err)
		}

		jsonResponse := response.ConvertTransactionLedgerToRespTransactionLedger(account, transactions)
		return RespondOK(res, jsonResponse)
	}
}

var ErrInvalidReconcileDate = errors.New("invalid reconcile date")

// GET /transactions/account/{accountID}/unreconciled?date=<date>
func GetUnreconciledTransactionsOnAccount(contoller *TransactionsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		accountIDStr := chi.URLParam(req, "accountID")

		accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if accountID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}

		dateStr := req.URL.Query().Get("date")

		dateCutoff, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, ErrInvalidReconcileDate)
		}

		transactions, err := contoller.GetUnreconciledTransactionsOnAccount(req.Context(), accountID, dateCutoff)
		if err != nil {
			return NewRequestError(http.StatusNotFound, err)
		}

		jsonResponse := response.ConvertTransactionRecSetToRespTransactionRecSet(transactions)
		return RespondOK(res, jsonResponse)
	}
}

// PUT /transactions/{transactionID}
func PutTransactionUpdate(contoller *TransactionsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		idStr := chi.URLParam(req, "transactionID")

		transactionID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if transactionID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}

		var reqTransaction request.Transaction

		if req.Body == nil {
			return NewRequestError(http.StatusBadRequest, ErrNoRequestBody)
		}

		err = json.NewDecoder(req.Body).Decode(&reqTransaction)
		if err != nil {
			return fmt.Errorf("json.NewDecoder(r.Body).Decode:%w", err)
		}

		mdlTransaction := request.ReqTransactionToTransaction(&reqTransaction)
		mdlTransaction.TransactionID = transactionID

		transaction, err := contoller.UpdateTransaction(req.Context(), mdlTransaction)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(res, jsonResponse)
	}
}

// DELETE /transactions/{transactionID}
func DeleteTransaction(contoller *TransactionsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		idStr := chi.URLParam(req, "transactionID")

		transactionID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if transactionID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}

		transaction, err := contoller.DeleteTransaction(req.Context(), transactionID)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(res, jsonResponse)
	}
}

// PUT /transactions/{transactionID}/reconciled
func PutTransactionReconciled(contoller *TransactionsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		idStr := chi.URLParam(req, "transactionID")

		transactionID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if transactionID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}

		var reqTransaction request.Transaction

		if req.Body == nil {
			return NewRequestError(http.StatusBadRequest, ErrNoRequestBody)
		}

		err = json.NewDecoder(req.Body).Decode(&reqTransaction)
		if err != nil {
			return fmt.Errorf("json.NewDecoder(r.Body).Decode:%w", err)
		}

		mdlTransaction := request.ReqTransactionToTransaction(&reqTransaction)
		mdlTransaction.TransactionID = transactionID

		transaction, err := contoller.UpdateReconciled(req.Context(), mdlTransaction)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(res, jsonResponse)
	}
}

// PUT /transactions/{transactionID}/unreconciled
func PutTransactionUnreconciled(contoller *TransactionsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		idStr := chi.URLParam(req, "transactionID")

		transactionID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if transactionID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}

		var reqTransaction request.Transaction

		if req.Body == nil {
			return NewRequestError(http.StatusBadRequest, ErrNoRequestBody)
		}

		err = json.NewDecoder(req.Body).Decode(&reqTransaction)
		if err != nil {
			return fmt.Errorf("json.NewDecoder(r.Body).Decode:%w", err)
		}

		mdlTransaction := request.ReqTransactionToTransaction(&reqTransaction)
		mdlTransaction.TransactionID = transactionID

		transaction, err := contoller.UpdateUnreconciled(req.Context(), mdlTransaction)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		jsonResponse := response.TransactionToRespTransaction(transaction)
		return RespondOK(res, jsonResponse)
	}
}
