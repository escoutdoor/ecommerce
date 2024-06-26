package models

import "time"

type User struct {
	ID          int        `json:"id"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Password    string     `json:"-"`
	Role        string     `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserReq struct {
	Email       string `json:"email" validate:"required,email"`
	FirstName   string `json:"first_name" validate:"required,min=2"`
	LastName    string `json:"last_name,omitempty" validate:"omitempty,min=2"`
	DateOfBirth string `json:"date_of_birth" validate:"omitempty"`
}
