package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/escoutdoor/ecommerce/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type Server struct {
	listenAddr string

	customer *CustomerHandler
	auth     *AuthHandler
	order    *OrderHandler
}

func NewServer() *http.Server {
	if err := godotenv.Load(); err != nil {
		log.Fatal("load env error: ", err)
	}

	db, err := store.ConnectToDB()
	if err != nil {
		log.Fatal("new server error: ", err)
	}

	customerStore := store.NewCustomerStore(db)
	customer := NewCustomerHandler(customerStore)

	authStore := store.NewAuthStore(db)
	auth := NewAuthHandler(authStore)

	orderStore := store.NewOrderStore(db)
	order := NewOrderHandler(orderStore)

	s := &Server{
		listenAddr: ":8080",
		customer:   customer,
		auth:       auth,
		order:      order,
	}

	server := &http.Server{
		Addr:         s.listenAddr,
		Handler:      s.Router(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func getID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %s", idStr)
	}

	return id, nil
}

func getIDFromCtx(r *http.Request) (int, error) {
	idStr := r.Context().Value("customer_id").(string)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %s", idStr)
	}

	return id, nil
}
