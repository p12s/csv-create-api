package rest

import (
	"context"
	"net/http"
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
		c, err := r.Cookie(COOKIE_NAME)
		if err != nil {
			ErrorCookie(w, err)
			return
		}

		userId, err := h.services.Sessioner.GetSession(context.TODO(), c.Value)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserID, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
