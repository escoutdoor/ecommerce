package server

import (
	"encoding/json"
	"net/http"

	"github.com/escoutdoor/ecommerce/internal/models"
	"github.com/escoutdoor/ecommerce/internal/store"
	"github.com/escoutdoor/ecommerce/internal/utils/respond"
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
		respond.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	product, err := h.store.Create(req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	respond.JSON(w, http.StatusOK, product)
}
