package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/escoutdoor/ecommerce/internal/models"
)

var (
	ErrOrderNotFound          = errors.New("order not found")
	ErrInvalidProductQuantity = errors.New("invalid product quantity")
)

type OrderStorer interface {
	Create(ctx context.Context, customerID int, data models.OrderReq) (*models.Order, error)
	GetByID(id int) (*models.Order, error)
	Delete(id int) error
}

type OrderStore struct {
	db           *sql.DB
	productStore ProductStore
}

func NewOrderStore(db *sql.DB) *OrderStore {
	return &OrderStore{
		db:           db,
		productStore: ProductStore{db: db},
	}
}

func (s *OrderStore) Create(ctx context.Context, customerID int, data models.OrderReq) (*models.Order, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	total := 0.0
	productsIDs := make([]int, len(data.OrderItems))
	for i, v := range data.OrderItems {
		if v.Quantity <= 0 {
			return nil, ErrInvalidProductQuantity
		}

		productsIDs[i] = v.ProductID
	}

	products, err := s.productStore.GetByIDs(productsIDs...)
	if err != nil {
		return nil, err
	}

	for _, v := range data.OrderItems {
		product, ok := products[v.ProductID]
		if !ok {
			return nil, ErrProductNotFound
		}

		total += float64(v.Quantity) * product.Price
	}

	order, err := s.createOrder(ctx, tx, customerID, total)
	if err != nil {
		return nil, err
	}

	for _, item := range data.OrderItems {
		orderItem, err := s.createOrderItem(ctx, tx, order.ID, item)
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

func (s *OrderStore) GetByID(id int) (*models.Order, error) {
	stmt, err := s.db.Prepare(`
		SELECT * FROM ORDERS WHERE ID = $1
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrderNotFound
		}

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

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order cannot be deleted because it doesn't exist")
	}

	return err
}

func (s *OrderStore) createOrder(ctx context.Context, tx *sql.Tx, customerID int, total float64) (*models.Order, error) {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO ORDERS(TOTAL, CUSTOMER_ID) VALUES ($1, $2) 
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, total, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanIntoOrder(rows)
	}

	return nil, err
}

func (s *OrderStore) createOrderItem(ctx context.Context, tx *sql.Tx, orderID int, data models.CreateOrderItemReq) (*models.OrderItem, error) {
	shippingDetails, err := s.createShippingDetails(ctx, tx, data.ShippingDetails)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO ORDER_ITEMS(PRODUCT_ID, ORDER_ID, SHIPPING_DETAILS_ID, QUANTITY)
		VALUES($1, $2, $3, $4)
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(
		ctx,
		data.ProductID,
		orderID,
		shippingDetails.ID,
		data.Quantity,
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

func scanIntoOrder(rows *sql.Rows) (*models.Order, error) {
	order := &models.Order{}
	err := rows.Scan(
		&order.ID,
		&order.Total,
		&order.CustomerID,
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
		&orderItem.ProductID,
		&orderItem.OrderID,
		&orderItem.ShippingDetailsID,
		&orderItem.Quantity,
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
