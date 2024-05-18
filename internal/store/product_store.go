package store

import (
	"database/sql"

	"github.com/escoutdoor/ecommerce/internal/models"
)

type ProductStorer interface {
	Create(data models.ProductReq) (*models.Product, error)
}

type ProductStore struct {
	db *sql.DB
}

func NewProductStore(db *sql.DB) *ProductStore {
	return &ProductStore{db: db}
}

func (s *ProductStore) Create(data models.ProductReq) (*models.Product, error) {
	stmt, err := s.db.Prepare(`
		INSERT INTO PRODUCTS(NAME, DESCRIPTION, PRICE, CATEGORY_ID)
		VALUES($1, $2, $3, $4)
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(data.Name, data.Description, data.Price, data.CategoryId)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoProduct(rows)
	}

	return nil, err
}

func scanIntoProduct(rows *sql.Rows) (*models.Product, error) {
	product := &models.Product{}
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CategoryId,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	return product, err
}
