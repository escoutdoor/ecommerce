package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/escoutdoor/ecommerce/internal/models"
	"github.com/escoutdoor/ecommerce/internal/store"
	"github.com/escoutdoor/ecommerce/internal/utils/respond"
	"github.com/escoutdoor/ecommerce/pkg/tokens"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	store store.AuthStorer
}

func NewAuthHandler(s store.AuthStorer) *AuthHandler {
	return &AuthHandler{
		store: s,
	}
}

func (h *AuthHandler) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	var req models.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := validator.New().Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		respond.Error(w, http.StatusBadRequest, respond.ValidationError(errs))
		return
	}

	user, err := h.store.Login(req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	token, err := tokens.CreateJWT(user.ID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	response := models.AuthResponse{
		User:  user,
		Token: token,
	}
	respond.JSON(w, http.StatusOK, response)
}

func (h *AuthHandler) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := validator.New().Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		respond.Error(w, http.StatusBadRequest, respond.ValidationError(errs))
		return
	}

	user, err := h.store.Register(req)
	if err != nil {
		if errors.Is(err, store.ErrEmailAlreadyExists) {
			respond.Error(w, http.StatusBadRequest, err)
			return
		}

		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	token, err := tokens.CreateJWT(user.ID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	response := models.AuthResponse{
		User:  user,
		Token: token,
	}
	respond.JSON(w, http.StatusOK, response)
}
