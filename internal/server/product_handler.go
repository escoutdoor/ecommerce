package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/escoutdoor/ecommerce/internal/models"
	"github.com/escoutdoor/ecommerce/internal/store"
	"github.com/escoutdoor/ecommerce/internal/utils/respond"
	"github.com/go-playground/validator/v10"
)

type ProductHandler struct {
	store store.ProductStorer
}

func NewProductHandler(s store.ProductStorer) *ProductHandler {
	return &ProductHandler{
		store: s,
	}
}

func (h *ProductHandler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.ProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := validator.New().Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		respond.Error(w, http.StatusBadRequest, respond.ValidationError(errs))
		return
	}

	product, err := h.store.Create(req)
	if err != nil {
		if errors.Is(err, store.ErrCategoryNotFound) {
			respond.Error(w, http.StatusBadRequest, err)
			return
		}

		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) handleGetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	product, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrProductNotFound) {
			respond.Error(w, http.StatusNotFound, err)
			return
		}

		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	respond.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.Delete(id); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, "product successfully deleted")
}

func (h *ProductHandler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.ProductReq
	_, err := getCustomerIDCtx(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	productID, err := getID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := validator.New().Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		respond.Error(w, http.StatusBadRequest, respond.ValidationError(errs))
		return
	}

	product, err := h.store.Update(productID, req)
	if err != nil {
		if errors.Is(err, store.ErrCategoryNotFound) {
			respond.Error(w, http.StatusBadRequest, err)
			return
		}

		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, product)

}
