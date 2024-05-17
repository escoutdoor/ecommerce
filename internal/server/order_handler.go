package server

import (
	"net/http"

	"github.com/escoutdoor/ecommerce/internal/store"
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
