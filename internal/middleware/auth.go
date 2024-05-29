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
	ErrForbidden    = errors.New("forbidden")
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

			customerID, err := tokens.VerifyToken(token)
			if err != nil {
				respond.Error(w, http.StatusUnauthorized, err)
				return
			}

			customer, err := s.GetByID(customerID)
			if err != nil {
				respond.Error(w, http.StatusUnauthorized, err)
				return
			}

			ctx := context.WithValue(r.Context(), "customer_id", fmt.Sprintf("%d", customerID))
			ctx = context.WithValue(ctx, "role", customer.Role)
			newReq := r.WithContext(ctx)
			h.ServeHTTP(w, newReq)
		})
	}
}

func RoleGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("role").(string)
		if !ok || role != "admin" {
			respond.Error(w, http.StatusForbidden, ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
