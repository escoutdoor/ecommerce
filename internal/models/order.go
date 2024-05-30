package models

import "time"

type Order struct {
	ID         int         `json:"id"`
	Total      float64     `json:"total"`
	UserID     int         `json:"user_id"`
	OrderItems []OrderItem `json:"order_items"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderReq struct {
	OrderItems []CreateOrderItemReq `json:"order_items" validate:"required,min=1,dive"`
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
	ProductID       int                `json:"product_id" validate:"required"`
	ShippingDetails ShippingDetailsReq `json:"shipping_details" validate:"required"`
	Quantity        int                `json:"quantity" validate:"required,min=1"`
}

type UpdateOrderItemReq struct {
	Status   string `json:"status" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
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
	AddressLine1 string `json:"address_line1" validate:"required"`
	AddressLine2 string `json:"address_line2" validate:"omitempty"`
	PostalCode   string `json:"postal_code" validate:"omitempty"`
	City         string `json:"city" validate:"required,min=3"`
	Country      string `json:"country" validate:"required,min=3"`
	Notes        string `json:"notes" validate:"omitempty"`
}
