package store

import (
	"database/sql"

	"github.com/escoutdoor/ecommerce/internal/models"
)

type CategoryStorer interface {
	GetByID(id int) (*models.Category, error)
	Create(data models.CategoryReq) (*models.Category, error)
}

type CategoryStore struct {
	db *sql.DB
}

func NewCategoryStore(db *sql.DB) *CategoryStore {
	return &CategoryStore{
		db: db,
	}
}

func (s *CategoryStore) Create(data models.CategoryReq) (*models.Category, error) {
	var category models.Category

	stmt, err := s.db.Prepare("INSERT INTO CATEGORIES(NAME) VALUES($1) RETURNING *")
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(data.Name).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &category, err
}

func (s *CategoryStore) GetByID(id int) (*models.Category, error) {
	return nil, nil
}
