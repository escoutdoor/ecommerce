package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/escoutdoor/ecommerce/internal/store"
	"github.com/escoutdoor/ecommerce/internal/utils/respond"
	"github.com/escoutdoor/ecommerce/pkg/tokens"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

func JWTAuth(s store.CustomerStorer) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if len(token) == 0 {
				respond.Error(w, http.StatusUnauthorized, ErrUnauthorized)
				return
			}
			token = token[len("Bearer "):]

			customerId, err := tokens.VerifyToken(token)
			if err != nil {
				respond.Error(w, http.StatusUnauthorized, err)
				return
			}

			_, err = s.GetById(customerId)
			if err != nil {
				respond.Error(w, http.StatusUnauthorized, err)
				return
			}

			ctx := context.WithValue(r.Context(), "customer_id", fmt.Sprintf("%d", customerId))
			newReq := r.WithContext(ctx)
			h.ServeHTTP(w, newReq)
		})
	}
}
