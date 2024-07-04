package web

import (
	"encoding/json"
	"fmt"
	"net/http"

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
