package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,containsany=!@#?*"`
}

type RegisterReq struct {
	Email       string `json:"email" validate:"required,email"`
	FirstName   string `json:"first_name" validate:"required,min=2"`
	LastName    string `json:"last_name" validate:"omitempty,min=2"`
	Password    string `json:"password" validate:"required,min=6,containsany=!@#?*"`
	DateOfBirth string `json:"date_of_birth" validate:"omitempty"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	ID string `json:"id"`
}

type AuthResponse struct {
	*User
	Token string `json:"token"`
}
