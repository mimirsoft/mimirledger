package web

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog"
	"net/http"
)

type UnknownError struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Err        string `json:"err,omitempty"`
}

type RequestError struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Err        string `json:"err,omitempty"`
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
	H func(res http.ResponseWriter, req *http.Request) error
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
// The root handler allows unified logging and error handling for the application
func (h RootHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log := zerolog.Ctx(req.Context())
	err := h.H(res, req)

	if err != nil {
		var myRequestError *RequestError
		if errors.As(err, &myRequestError) {
			// since this was a structured RequestError, we do not log and just
			// return the response.  Logging middleware will handle the log
			res.WriteHeader(myRequestError.StatusCode)
			res.Header().Set("Content-Type", "application/json; charset=utf-8")

			if err = json.NewEncoder(res).Encode(myRequestError); err != nil {
				log.Printf("failed json encode RequestError response %s", err)
			}

			return
		}
		// since we don't know what the error is, we log the error
		// then we return a generic StatusInternalServerError with a JSON payload
		log.Error().Err(err).Msgf("HTTP %d", http.StatusInternalServerError)
		unknownErr := UnknownError{StatusCode: http.StatusInternalServerError, Err: err.Error()}
		// Any error types we don't specifically look out for default
		// to serving a HTTP 500
		res.WriteHeader(http.StatusInternalServerError)
		res.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err1 := json.NewEncoder(res).Encode(unknownErr); err1 != nil {
			log.Error().Err(err1).Msg("failed json encode UnknownError response")
		}

		return

	}
}

// returns our universal custom root handler
func NewRootHandler(handler func(w http.ResponseWriter, r *http.Request) error) *RootHandler {
	return &RootHandler{H: handler}
}

// GET /api/health
// HEAD /api/health
func HealthCheck(healthController *HealthController) func(w http.ResponseWriter, r *http.Request) error {
	return func(res http.ResponseWriter, req *http.Request) error {
		err := healthController.HealthCheck(req.Context())
		if err != nil {
			return NewRequestError(http.StatusServiceUnavailable, err)
		}

		return RespondOK(res, "ok")
	}
}
