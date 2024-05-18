package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/escoutdoor/ecommerce/internal/models"
	"github.com/escoutdoor/ecommerce/internal/store"
	"github.com/escoutdoor/ecommerce/internal/utils/respond"
	"github.com/escoutdoor/ecommerce/pkg/tokens"
)

type AuthHandler struct {
	store store.AuthStorer
}

func NewAuthHandler(s store.AuthStorer) *AuthHandler {
	return &AuthHandler{
		store: s,
	}
}

func (h *AuthHandler) handleLoginCustomer(w http.ResponseWriter, r *http.Request) {
	var req models.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	customer, err := h.store.Login(req)
	if err != nil {
		if errors.Is(err, store.ErrInvalidEmailOrPassword) {
			respond.Error(w, http.StatusUnauthorized, err)
			return
		}

		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	token, err := tokens.CreateJWT(customer.ID)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	response := models.AuthResponse{
		Customer: customer,
		Token:    token,
	}
	respond.JSON(w, http.StatusOK, response)
}

func (h *AuthHandler) handleRegisterCustomer(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	customer, err := h.store.Register(req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	token, err := tokens.CreateJWT(customer.ID)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	response := models.AuthResponse{
		Customer: customer,
		Token:    token,
	}
	respond.JSON(w, http.StatusOK, response)
}
