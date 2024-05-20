package models

import "time"

type Customer struct {
	ID          int        `json:"id"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Password    string     `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateCustomerReq struct {
	Email       string    `json:"email,omitempty"`
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	Password    string    `json:"password,omitempty"`
}
