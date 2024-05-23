package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/escoutdoor/ecommerce/internal/models"
	"github.com/lib/pq"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductStorer interface {
	Create(data models.ProductReq) (*models.Product, error)
	GetByID(id int) (*models.Product, error)
	Delete(id int) error
	Update(id int, data models.ProductReq) (*models.Product, error)
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

	rows, err := stmt.Query(data.Name, data.Description, data.Price, data.CategoryID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23503" {
			return nil, ErrCategoryNotFound
		}

		return nil, err
	}

	if rows.Next() {
		return scanIntoProduct(rows)
	}

	return nil, err
}

func (s *ProductStore) GetByID(id int) (*models.Product, error) {
	stmt, err := s.db.Prepare(`SELECT * FROM PRODUCTS WHERE ID = $1`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoProduct(rows)
	}

	return nil, err
}

func (s *ProductStore) Delete(id int) error {
	stmt, err := s.db.Prepare("DELETE FROM PRODUCTS WHERE ID = $1")
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
		return fmt.Errorf("product cannot be deleted because it doesn't exist")
	}

	return err
}

func (s *ProductStore) Update(id int, data models.ProductReq) (*models.Product, error) {
	stmt, err := s.db.Prepare(`
		UPDATE PRODUCTS SET
			NAME = $1,
			DESCRIPTION = $2,
			PRICE = $3,
			CATEGORY_ID = $4
		WHERE ID = $5
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(data.Name, data.Description, data.Price, data.CategoryID, id)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23503" {
			return nil, ErrCategoryNotFound
		}

		return nil, err
	}

	if rows.Next() {
		return scanIntoProduct(rows)
	}

	return nil, err
}

func (s *ProductStore) GetByIDs(ids ...int) (map[int]models.Product, error) {
	params := make([]string, len(ids))
	for i := range ids {
		params[i] = fmt.Sprintf("$%d", i+1)
	}
	query := fmt.Sprintf("SELECT * FROM PRODUCTS WHERE ID IN (%s)", strings.Join(params, ","))

	values := make([]interface{}, len(ids))
	for i, v := range ids {
		values[i] = v
	}

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, err
	}

	products := make(map[int]models.Product)
	for rows.Next() {
		p, err := scanIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products[p.ID] = *p
	}

	return products, nil
}

func scanIntoProduct(rows *sql.Rows) (*models.Product, error) {
	product := &models.Product{}
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CategoryID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	return product, err
}
