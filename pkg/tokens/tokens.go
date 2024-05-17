package tokens

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"log"

	"github.com/escoutdoor/ecommerce/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(tokenStr string) (int, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return 0, fmt.Errorf("jwt is empty verifyToken err")
	}

	var claims models.TokenClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return 0, fmt.Errorf("That's not even a token")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return 0, fmt.Errorf("Invalid signature")
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			return 0, fmt.Errorf("Token is either expired or not active yet")
		default:
			return 0, fmt.Errorf("Couldn't handle this token: %s", err.Error())
		}
	}

	if token.Valid {
		log.Print("Token is valid")
	}

	id, err := strconv.Atoi(claims.ID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func CreateJWT(id int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("jwt is empty createToken err")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		ID: fmt.Sprintf("%d", id),
	})

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
