package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/escoutdoor/ecommerce/internal/store"
	"github.com/escoutdoor/ecommerce/internal/utils/respond"
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

}

func (h *OrderHandler) handleGetOrderById(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	order, err := h.store.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond.Error(w, http.StatusNotFound, store.ErrOrderNotFound)
			return
		}

		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	respond.JSON(w, http.StatusOK, order)
}

func (h *OrderHandler) handleUpdateOrder(w http.ResponseWriter, r *http.Request) {

}

func (h *OrderHandler) handleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.Delete(id); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	respond.JSON(w, http.StatusOK, "order successfully deleted")
}
