package store

import (
	"database/sql"
	"errors"

	"github.com/escoutdoor/ecommerce/internal/models"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderStorer interface {
	Create(id int, data *models.OrderReq) (*models.Order, error)
	GetById(id int) (*models.Order, error)
	Update(id int, data models.OrderReq) (*models.Order, error)
	Delete(id int) error
}

type OrderStore struct {
	db *sql.DB
}

func NewOrderStore(db *sql.DB) *OrderStore {
	return &OrderStore{
		db: db,
	}
}

func (s *OrderStore) Create(id int, data *models.OrderReq) (*models.Order, error) {
	stmt, err := s.db.Prepare(`
		INSERT INTO ORDERS(total, customer_id) VALUES
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(data.Total, data.CustomerId)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoOrder(rows)
	}

	return nil, err
}

func (s *OrderStore) GetById(id int) (*models.Order, error) {
	stmt, err := s.db.Prepare(`
		SELECT * FROM ORDERS WHERE ID = $1
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoOrder(rows)
	}

	return nil, err
}

func (s *OrderStore) Update(id int, data models.OrderReq) (*models.Order, error) {
	stmt, err := s.db.Prepare(`
		UPDATE ORDERS SET
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoOrder(rows)
	}

	return nil, err
}

func (s *OrderStore) Delete(id int) error {
	stmt, err := s.db.Prepare(`
		DELETE FROM ORDERS WHERE ID = $1
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Query(id)
	if err != nil {
		return err
	}

	return err
}

func scanIntoOrder(rows *sql.Rows) (*models.Order, error) {
	order := &models.Order{}
	err := rows.Scan(
		&order.ID,
		&order.Total,
		&order.CustomerId,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	return order, err
}
