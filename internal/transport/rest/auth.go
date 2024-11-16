package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
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

	if err = input.Validate(); err != nil {
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
// @Param input body domain.Credentials true "credentials"
// @Success 200
// @Router /sign-in [post]
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var creds domain.Credentials
	if err = json.Unmarshal(reqBytes, &creds); err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = creds.Validate(); err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := h.services.GetUserByCredentials(r.Context(), creds.Email, creds.Password)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sessionToken, err := uuid.NewV4()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Sessioner.SetExpiredSession(context.TODO(), sessionToken.String(), userId)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    COOKIE_NAME,
		Value:   sessionToken.String(),
		Expires: time.Now().Add(SESSION_HOUR_TTL * time.Hour),
	})
	w.WriteHeader(http.StatusOK)
}

// @Summary Logout
// @Tags Auth
// @Description Logout, expire cookies
// @ID logout
// @Success 200
// @Router /logout [get]
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		ErrorCookie(w, err)
		return
	}

	err = h.services.Sessioner.DeleteSession(context.TODO(), c.Value)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    COOKIE_NAME,
		Value:   "",
		Expires: time.Now().Add(-1 * SESSION_HOUR_TTL * time.Hour),
	})
	w.WriteHeader(http.StatusOK)
}
