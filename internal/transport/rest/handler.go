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

	_ "github.com/p12s/csv-create-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Producter - transport contract (handler)
type Producter interface {
	Create(ctx context.Context, product domain.Product) error
	UpdateById(ctx context.Context, id int, input domain.UpdateProductInput) error
	DeleteById(ctx context.Context, id int) error
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
}

// Handler
type Handler struct {
	productService Producter
}

// NewHandler - constructor
func NewHandler(product Producter) *Handler {
	return &Handler{
		productService: product,
	}
}

// InitRouter - init routes
func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	products := r.PathPrefix("/products").Subrouter()
	{
		products.HandleFunc("/", h.createProduct).Methods(http.MethodPost)
		products.HandleFunc("/{id:[0-9]+}", h.updateProduct).Methods(http.MethodPut)
		products.HandleFunc("/{id:[0-9]+}", h.deleteProduct).Methods(http.MethodDelete)
		products.HandleFunc("/", h.getAllProducts).Methods(http.MethodGet)
	}

	return r
}

// @Summary Create product
// @Tags Product
// @Description Create product
// @ID createProduct
// @Accept  json
// @Produce  json
// @Param input body domain.Product true "Product created info"
// @Success 200
// @Router /products/ [post]
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

// @Summary Delete
// @Tags Product
// @Description Deleting product by {id}
// @ID deleteProduct
// @Accept  json
// @Produce  json
// @Param id path integer true "Product id"
// @Success 200
// @Router /products/{id} [delete]
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

// @Summary Get all
// @Tags Product
// @Description Getting all products
// @ID getAllProducts
// @Accept  json
// @Produce  json
// @Success 200
// @Router /products/ [get]
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

// @Summary Update by id
// @Tags Product
// @Description Update product by {id}
// @ID updateProduct
// @Accept  json
// @Produce  json
// @Param id path integer true "Product id"
// @Param input body domain.UpdateProductInput true "Product updated info"
// @Success 200
// @Router /products/{id} [put]
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
