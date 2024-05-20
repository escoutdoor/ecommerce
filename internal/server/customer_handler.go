package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/escoutdoor/ecommerce/internal/models"
	"github.com/escoutdoor/ecommerce/internal/store"
	"github.com/escoutdoor/ecommerce/internal/utils/respond"
)

type CustomerHandler struct {
	store store.CustomerStorer
}

func NewCustomerHandler(s store.CustomerStorer) *CustomerHandler {
	return &CustomerHandler{
		store: s,
	}
}

func (h *CustomerHandler) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	customer, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrCustomerNotFound) {
			respond.Error(w, http.StatusNotFound, err)
			return
		}

		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	respond.JSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) handleUpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromCtx(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	var req models.UpdateCustomerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	customer, err := h.store.Update(id, req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	respond.JSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) handleDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromCtx(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = h.store.Delete(id); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	respond.JSON(w, http.StatusOK, "customer account successfully deleted")
}
