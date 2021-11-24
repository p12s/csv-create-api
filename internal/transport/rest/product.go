package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/jszwec/csvutil"
	"github.com/p12s/csv-create-api/internal/domain"
)

// @Summary Create product
// @Security ApiKeyAuth
// @Tags Product
// @Description Create product
// @ID createProduct
// @Accept  json
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

	err = h.services.CreateProduct(context.TODO(), product)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Delete
// @Security ApiKeyAuth
// @Tags Product
// @Description Deleting product by {id}
// @ID deleteProduct
// @Param id path integer true "Product id"
// @Success 200
// @Router /products/{id} [delete]
func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.DeleteProductById(context.TODO(), id)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Get all
// @Security ApiKeyAuth
// @Tags Product
// @Description Getting all products
// @ID getAllProducts
// @Produce  json
// @Success 200
// @Router /products/ [get]
func (h *Handler) getAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.services.GetAllProducts(context.TODO())
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
// @Security ApiKeyAuth
// @Tags Product
// @Description Update product by {id}
// @ID updateProduct
// @Accept  json
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

	err = h.services.UpdateProductById(context.TODO(), id, input)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
