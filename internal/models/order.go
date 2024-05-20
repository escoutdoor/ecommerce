package models

import "time"

type Order struct {
	ID         int         `json:"id"`
	Total      float64     `json:"total"`
	CustomerID int         `json:"customer_id"`
	OrderItems []OrderItem `json:"order_items"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderReq struct {
	OrderItems []CreateOrderItemReq `json:"order_items"`
}

type OrderItem struct {
	ID                int    `json:"id"`
	Status            string `json:"status"`
	ProductID         int    `json:"product_id"`
	OrderID           int    `json:"order_id"`
	ShippingDetailsID int    `json:"shipping_details_id"`
	Quantity          int    `json:"quantity"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateOrderItemReq struct {
	ProductID          int `json:"product_id"`
	ShippingDetailsReq `json:"shipping_details"`
	Quantity           int `json:"quantity"`
}

type UpdateOrderItemReq struct {
	Status   string `json:"status"`
	Quantity int    `json:"quantity"`
}

type ShippingDetails struct {
	ID           int    `json:"id"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	PostalCode   string `json:"postal_code"`
	City         string `json:"city"`
	Country      string `json:"country"`
	Notes        string `json:"notes"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShippingDetailsReq struct {
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	PostalCode   string `json:"postal_code"`
	City         string `json:"city"`
	Country      string `json:"country"`
	Notes        string `json:"notes"`
}
