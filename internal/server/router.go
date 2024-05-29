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
		r.Use(middleware.JWTAuth(s.customer.store))

		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleGuard)

			r.Get("/{id}", s.customer.handleGetCustomerByID)
		})

		r.Put("/", s.customer.handleUpdateCustomer)
		r.Delete("/", s.customer.handleDeleteCustomer)
	})

	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", s.auth.handleLoginCustomer)
		r.Post("/register", s.auth.handleRegisterCustomer)
	})

	router.Route("/categories", func(r chi.Router) {
		r.Get("/{id}", s.category.handleGetCategoryByID)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(s.customer.store))
			r.Use(middleware.RoleGuard)

			r.Post("/", s.category.handleCreateCategory)
			r.Delete("/{id}", s.category.handleDeleteCategory)
			r.Put("/{id}", s.category.handleUpdateCategory)
		})
	})

	router.Route("/products", func(r chi.Router) {
		r.Get("/{id}", s.product.handleGetProductByID)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(s.customer.store))
			r.Use(middleware.RoleGuard)

			r.Post("/", s.product.handleCreateProduct)
			r.Put("/{id}", s.product.handleUpdateProduct)
			r.Delete("/{id}", s.product.handleDeleteProduct)
		})
	})

	router.Route("/orders", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(s.customer.store))

			r.Post("/", s.order.handleCreateOrder)
			r.Delete("/{id}", s.order.handleDeleteOrder)
		})
	})

	return router
}
