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

type OrderHandler struct {
	store store.OrderStorer
}

func NewOrderHandler(s store.OrderStorer) *OrderHandler {
	return &OrderHandler{
		store: s,
	}
}

func (h *OrderHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := getCustomerIDCtx(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	var req models.OrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := validator.New().Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		respond.Error(w, http.StatusBadRequest, respond.ValidationError(errs))
		return
	}

	order, err := h.store.Create(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, store.ErrProductNotFound) {
			respond.Error(w, http.StatusNotFound, err)
			return
		}

		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) handleGetOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	order, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrOrderNotFound) {
			respond.Error(w, http.StatusNotFound, err)
			return
		}

		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, order)
}

func (h *OrderHandler) handleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.Delete(id); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, "order successfully deleted")
}
