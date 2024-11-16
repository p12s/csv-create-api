package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/p12s/csv-create-api/internal/service"

	_ "github.com/p12s/csv-create-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Handler
type Handler struct {
	services *service.Service
}

// NewHandler - constructor
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// InitRouter - init routes
func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/health", h.health).Methods(http.MethodGet)
	r.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
	r.HandleFunc("/sign-in", h.signIn).Methods(http.MethodPost)
	r.HandleFunc("/logout", h.logout).Methods(http.MethodGet)

	products := r.PathPrefix("/products").Subrouter()
	{
		products.Use(h.authMiddleware)
		products.HandleFunc("/", h.createProduct).Methods(http.MethodPost)
		products.HandleFunc("/{id:[0-9]+}", h.updateProduct).Methods(http.MethodPut)
		products.HandleFunc("/{id:[0-9]+}", h.deleteProduct).Methods(http.MethodDelete)
		products.HandleFunc("/", h.getAllProducts).Methods(http.MethodGet)
	}

	return r
}

// getIdFromRequest - getting id from request
func getIdFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
