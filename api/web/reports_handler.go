package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mimirsoft/mimirledger/api/models"
	"github.com/mimirsoft/mimirledger/api/web/request"
	"github.com/mimirsoft/mimirledger/api/web/response"
)

// GET /reports
func GetReports(reportsCtl *ReportsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		Reports, err := reportsCtl.ReportList(req.Context())
		if err != nil {
			return NewRequestError(http.StatusServiceUnavailable, err)
		}

		jsonResponse := response.ConvertReportsToRespReportsSet(Reports)

		return RespondOK(res, jsonResponse)
	}
}

var ErrInvalidReportID = errors.New("invalid reportID request parameter")

// GET /reports/{reportID}
func GetReport(reportsCtl *ReportsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		reportIDStr := chi.URLParam(req, "reportID")

		reportID, err := strconv.ParseUint(reportIDStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if reportID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}

		report, err := reportsCtl.GetReportByID(req.Context(), reportID)
		if err != nil {
			if errors.Is(err, models.ErrReportNotFound) {
				return NewRequestError(http.StatusNotFound, err)
			}

			return NewRequestError(http.StatusServiceUnavailable, err)
		}

		jsonResponse := response.ReportToRespReport(report)

		return RespondOK(res, jsonResponse)
	}
}

var ErrInvalidStartDate = errors.New("invalid startDate")
var ErrInvalidEndDate = errors.New("invalid endDate")

// GET /reports/{reportID}/output?date
func GetReportOutput(reportsCtl *ReportsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		reportIDStr := chi.URLParam(req, "reportID")

		reportID, err := strconv.ParseUint(reportIDStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if reportID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}

		startDateStr := req.URL.Query().Get("startDate")

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, ErrInvalidStartDate)
		}

		endDateStr := req.URL.Query().Get("endDate")

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, ErrInvalidStartDate)
		}

		var accountSet []uint64

		queryVals, ok := req.URL.Query()["account"]
		if ok {
			for _, str := range queryVals {
				parsedUint, strConvErr := strconv.ParseUint(str, 10, 64)
				if strConvErr != nil {
					return NewRequestError(http.StatusBadRequest, strConvErr)
				}

				accountSet = append(accountSet, parsedUint)
			}
		}

		reportOutput, err := reportsCtl.RunReport(req.Context(), reportID, startDate, endDate, accountSet)
		if err != nil {
			if errors.Is(err, models.ErrReportNotFound) {
				return NewRequestError(http.StatusNotFound, err)
			}

			return NewRequestError(http.StatusServiceUnavailable, err)
		}

		jsonResponse := response.ReportOutputToRespReportOutput(reportOutput)

		return RespondOK(res, jsonResponse)
	}
}

// PUT /reports/{reportID}
func PutReportUpdate(reportsCtl *ReportsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		reportIDStr := chi.URLParam(req, "reportID")

		reportID, err := strconv.ParseUint(reportIDStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if reportID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidAccountID)
		}

		if req.Body == nil {
			return NewRequestError(http.StatusBadRequest, ErrNoRequestBody)
		}

		var report request.Report

		err = json.NewDecoder(req.Body).Decode(&report)
		if err != nil {
			return fmt.Errorf("json.NewDecoder(r.Body).Decode:%w", err)
		}

		mdlReport := request.ReqReportToReport(&report)
		mdlReport.ReportID = reportID

		myReport, err := reportsCtl.UpdateReport(req.Context(), mdlReport)
		if err != nil {
			if errors.Is(err, models.ErrReportNotFound) {
				return NewRequestError(http.StatusNotFound, err)
			}

			return NewRequestError(http.StatusServiceUnavailable, err)
		}

		jsonResponse := response.ReportToRespReport(myReport)

		return RespondOK(res, jsonResponse)
	}
}

// POST /reports
func PostReports(reportsCtl *ReportsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		var report request.Report

		if req.Body == nil {
			return NewRequestError(http.StatusBadRequest, ErrNoRequestBody)
		}

		err := json.NewDecoder(req.Body).Decode(&report)
		if err != nil {
			return fmt.Errorf("json.NewDecoder(r.Body).Decode:%w", err)
		}

		mdlReport := request.ReqReportToReport(&report)

		myReport, err := reportsCtl.CreateReport(req.Context(), mdlReport)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		jsonResponse := response.ReportToRespReport(myReport)

		return RespondOK(res, jsonResponse)
	}
}

// POST /reports/restore
func PostReportsRestore(reportsCtl *ReportsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		Reports, err := reportsCtl.RestoreDefault(req.Context())
		if err != nil {
			return NewRequestError(http.StatusServiceUnavailable, err)
		}
		jsonResponse := response.ConvertReportsToRespReportsSet(Reports)
		return RespondOK(res, jsonResponse)
	}
}

// DELETE /reports/{reportID}
func DeleteReport(reportsCtl *ReportsController) func(res http.ResponseWriter, req *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		idStr := chi.URLParam(req, "reportID")

		reportID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		if reportID == 0 {
			return NewRequestError(http.StatusBadRequest, ErrInvalidTransactionID)
		}

		report, err := reportsCtl.DeleteReport(req.Context(), reportID)
		if err != nil {
			return NewRequestError(http.StatusBadRequest, err)
		}

		jsonResponse := response.ReportToRespReport(report)

		return RespondOK(res, jsonResponse)
	}
}
