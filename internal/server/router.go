package server

import (
	"github.com/escoutdoor/ecommerce/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddle "github.com/go-chi/chi/v5/middleware"
)

func (s *Server) Router() *chi.Mux {
	router := chi.NewRouter()
	router.Use(chimiddle.Logger)
	router.Use(chimiddle.StripSlashes)

	router.Route("/customers", func(r chi.Router) {
		r.Get("/{id}", s.customer.handleGetCustomerById)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(s.customer.store))

			r.Put("/", s.customer.handleUpdateCustomer)
			r.Delete("/", s.customer.handleDeleteCustomer)
		})
	})

	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", s.auth.handleLoginCustomer)
		r.Post("/register", s.auth.handleRegisterCustomer)
	})

	router.Route("/orders", func(r chi.Router) {
	})

	return router
}
