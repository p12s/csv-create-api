package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jszwec/csvutil"
	"github.com/p12s/csv-create-api/internal/domain"
)

type Producter interface {
	Create(ctx context.Context, product domain.Product) error
	UpdateById(ctx context.Context, id int, input domain.UpdateProductInput) error
	DeleteById(ctx context.Context, id int) error
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
}

type Handler struct {
	productService Producter
}

func NewHandler(product Producter) *Handler {
	return &Handler{
		productService: product,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	products := r.PathPrefix("/products").Subrouter()
	{
		products.HandleFunc("/", h.createProduct).Methods(http.MethodPost)
		products.HandleFunc("/{id:[0-9]+}", h.updateProduct).Methods(http.MethodPut)
		products.HandleFunc("/{id:[0-9]+}", h.deleteProduct).Methods(http.MethodDelete)
		products.HandleFunc("/", h.getAllProducts).Methods(http.MethodGet)
	}

	return r
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var product domain.Product
	if err = json.Unmarshal(reqBytes, &product); err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.productService.Create(context.TODO(), product)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.productService.DeleteById(context.TODO(), id)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetAllProducts(context.TODO())
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response, err := csvutil.Marshal(products)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	CsvResponse(w, os.Getenv("DEFAULT_CSV_FILE_NAME"), response)
}

func (h *Handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var input domain.UpdateProductInput
	if err = json.Unmarshal(reqBytes, &input); err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.productService.UpdateById(context.TODO(), id, input)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
