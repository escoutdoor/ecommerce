package server

import (
	"encoding/json"
	"net/http"

	"github.com/escoutdoor/ecommerce/internal/models"
	"github.com/escoutdoor/ecommerce/internal/store"
	"github.com/escoutdoor/ecommerce/internal/utils/respond"
)

type CategoryHandler struct {
	store store.CategoryStorer
}

func NewCategoryHandler(s store.CategoryStorer) *CategoryHandler {
	return &CategoryHandler{
		store: s,
	}
}

func (h *CategoryHandler) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	category, err := h.store.Create(req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	respond.JSON(w, http.StatusOK, category)
}
