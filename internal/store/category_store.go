package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/escoutdoor/ecommerce/internal/models"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryStorer interface {
	GetByID(id int) (*models.Category, error)
	Create(data models.CategoryReq) (*models.Category, error)
	Delete(id int) error
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
	stmt, err := s.db.Prepare("SELECT * FROM CATEGORIES WHERE ID = $1")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoCategory(rows)
	}

	return nil, err
}

func (s *CategoryStore) Delete(id int) error {
	stmt, err := s.db.Prepare("DELETE FROM CATEGORIES WHERE ID = $1")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category cannot be deleted because it doesn't exist")
	}

	return err
}

func scanIntoCategory(rows *sql.Rows) (*models.Category, error) {
	var category models.Category
	err := rows.Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	return &category, err
}
