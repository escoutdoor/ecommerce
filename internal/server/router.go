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

	router.Route("/users", func(r chi.Router) {
		r.Use(middleware.JWTAuth(s.user.store))

		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleGuard)

			r.Get("/{id}", s.user.handleGetUserByID)
		})

		r.Put("/", s.user.handleUpdateUser)
		r.Delete("/", s.user.handleDeleteUser)
	})

	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", s.auth.handleLoginUser)
		r.Post("/register", s.auth.handleRegisterUser)
	})

	router.Route("/categories", func(r chi.Router) {
		r.Get("/{id}", s.category.handleGetCategoryByID)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(s.user.store))
			r.Use(middleware.RoleGuard)

			r.Post("/", s.category.handleCreateCategory)
			r.Delete("/{id}", s.category.handleDeleteCategory)
			r.Put("/{id}", s.category.handleUpdateCategory)
		})
	})

	router.Route("/products", func(r chi.Router) {
		r.Get("/{id}", s.product.handleGetProductByID)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(s.user.store))
			r.Use(middleware.RoleGuard)

			r.Post("/", s.product.handleCreateProduct)
			r.Put("/{id}", s.product.handleUpdateProduct)
			r.Delete("/{id}", s.product.handleDeleteProduct)
		})
	})

	router.Route("/orders", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(s.user.store))

			r.Post("/", s.order.handleCreateOrder)
			r.Delete("/{id}", s.order.handleDeleteOrder)
			r.Get("/{id}", s.order.handleGetOrderByID)
		})
	})

	return router
}
