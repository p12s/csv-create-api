package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/p12s/csv-create-api/internal/domain"
)

// @Summary Sign up
// @Tags Auth
// @Description Create account
// @ID signUp
// @Accept  json
// @Produce  json
// @Param input body domain.SignUpInput true "credentials"
// @Success 200
// @Router /sign-up [post]
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var input domain.SignUpInput
	if err = json.Unmarshal(reqBytes, &input); err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.CreateUser(r.Context(), input)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Sign in
// @Tags Auth
// @Description Sending data to get authentication with jwt-token
// @ID signIn
// @Accept  json
// @Produce  json
// @Param input body domain.SignInInput true "credentials"
// @Success 200
// @Router /sign-in [post]
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var input domain.SignInInput
	if err = json.Unmarshal(reqBytes, &input); err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.GetUserByCredentials(r.Context(), input.Email, input.Password)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
