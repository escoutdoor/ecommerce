package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/escoutdoor/ecommerce/internal/models"
	"github.com/escoutdoor/ecommerce/internal/store"
	"github.com/escoutdoor/ecommerce/internal/utils/respond"
	"github.com/go-playground/validator/v10"
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
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := validator.New().Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		respond.Error(w, http.StatusBadRequest, respond.ValidationError(errs))
		return
	}

	category, err := h.store.Create(req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusCreated, category)
}

func (h *CategoryHandler) handleGetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	category, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond.Error(w, http.StatusNotFound, store.ErrCategoryNotFound)
			return
		}

		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) handleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.Delete(id)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, "category successfully deleted")
}

func (h *CategoryHandler) handleUpdateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CategoryReq
	_, err := getUserIDCtx(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	categoryID, err := getID(r)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := validator.New().Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		respond.Error(w, http.StatusBadRequest, respond.ValidationError(errs))
		return
	}

	category, err := h.store.Update(categoryID, req)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, category)
}
