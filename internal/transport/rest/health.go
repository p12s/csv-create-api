package rest

import (
	"encoding/json"
	"net/http"
)

// health - app technical service handler
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(map[string]string{
		"status": "Ok",
	})
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
