package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// CtxValue
type CtxValue int

const (
	ctxUserID CtxValue = iota
)

// loggingMiddleware - logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Printf("%s: [%s] - %s ", time.Now().Format(time.RFC3339), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// authMiddleware - authentication
func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			ErrorResponse(w, http.StatusUnauthorized, "empty or wrong token")
			return
		}

		userId, err := h.services.ParseToken(r.Context(), token)
		if err != nil {
			ErrorResponse(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserID, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// getTokenFromRequest
func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
