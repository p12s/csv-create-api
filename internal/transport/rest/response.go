package rest

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// ErrorResponse
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(message))
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
