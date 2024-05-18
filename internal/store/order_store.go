package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/escoutdoor/ecommerce/internal/models"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderStorer interface {
	Create(ctx context.Context, customerId int, data models.OrderReq) (*models.Order, error)
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

func (s *OrderStore) Create(ctx context.Context, customerId int, data models.OrderReq) (*models.Order, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	total := 10.10

	order, err := s.createOrder(ctx, tx, customerId, total)
	if err != nil {
		return nil, err
	}

	for _, item := range data.OrderItems {
		shippingDetails, err := s.createShippingDetails(ctx, tx, item.ShippingDetailsReq)
		if err != nil {
			return nil, err
		}

		orderItem, err := s.createOrderItem(ctx, tx, item.ProductId, order.ID, shippingDetails.ID)
		if err != nil {
			return nil, err
		}

		order.OrderItems = append(order.OrderItems, *orderItem)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return order, err
}

func (s *OrderStore) createOrder(ctx context.Context, tx *sql.Tx, customerId int, total float64) (*models.Order, error) {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO ORDERS(TOTAL, CUSTOMER_ID) VALUES ($1, $2) 
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, total, customerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanIntoOrder(rows)
	}

	return nil, err
}

func (s *OrderStore) createShippingDetails(ctx context.Context, tx *sql.Tx, data models.ShippingDetailsReq) (*models.ShippingDetails, error) {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO SHIPPING_DETAILS(ADDRESS_LINE1, ADDRESS_LINE2, POSTAL_CODE, CITY, COUNTRY, NOTES) 
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(
		ctx,
		data.AddressLine1,
		data.AddressLine2,
		data.PostalCode,
		data.City,
		data.Country,
		data.Notes,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanIntoShippingDetails(rows)
	}

	return nil, err
}

func (s *OrderStore) createOrderItem(ctx context.Context, tx *sql.Tx, productId int, orderId int, shippingDetailsId int) (*models.OrderItem, error) {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO ORDER_ITEMS(PRODUCT_ID, ORDER_ID, SHIPPING_DETAILS_ID)
		VALUES($1, $2, $3)
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(
		ctx,
		productId,
		orderId,
		shippingDetailsId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanIntoOrderItem(rows)
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

func scanIntoOrderItem(rows *sql.Rows) (*models.OrderItem, error) {
	orderItem := &models.OrderItem{}
	err := rows.Scan(
		&orderItem.ID,
		&orderItem.Status,
		&orderItem.ProductId,
		&orderItem.OrderId,
		&orderItem.ShippingDetailsId,
		&orderItem.CreatedAt,
		&orderItem.UpdatedAt,
	)

	return orderItem, err
}

func scanIntoShippingDetails(rows *sql.Rows) (*models.ShippingDetails, error) {
	shippingDetails := &models.ShippingDetails{}
	err := rows.Scan(
		&shippingDetails.ID,
		&shippingDetails.AddressLine1,
		&shippingDetails.AddressLine2,
		&shippingDetails.PostalCode,
		&shippingDetails.City,
		&shippingDetails.Country,
		&shippingDetails.Notes,
		&shippingDetails.CreatedAt,
		&shippingDetails.UpdatedAt,
	)

	return shippingDetails, err
}
