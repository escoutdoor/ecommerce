package models

import "github.com/golang-jwt/jwt/v5"

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterReq struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	ID string `json:"id"`
}

type AuthResponse struct {
	*Customer
	Token string `json:"token"`
}
