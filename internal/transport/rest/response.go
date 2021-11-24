package rest

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	SESSION_HOUR_TTL = 12
	COOKIE_NAME      = "session_token"
)

// ErrorResponse
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(message))
}

// ErrorCookie
func ErrorCookie(w http.ResponseWriter, err error) {
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// OkResponse
func OkResponse(w http.ResponseWriter, message []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(message)
}

// CsvResponse
func CsvResponse(w http.ResponseWriter, fileName string, message []byte) {
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "text/csv")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(message)
}
