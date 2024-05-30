package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/escoutdoor/ecommerce/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type Server struct {
	listenAddr string

	user     *UserHandler
	auth     *AuthHandler
	product  *ProductHandler
	order    *OrderHandler
	category *CategoryHandler
}

func NewServer() *http.Server {
	if err := godotenv.Load(); err != nil {
		log.Fatal("load env error: ", err)
	}

	var port = os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	db, err := store.ConnectToDB()
	if err != nil {
		log.Fatal("new server error: ", err)
	}

	userStore := store.NewUserStore(db)
	user := NewUserHandler(userStore)

	authStore := store.NewAuthStore(db)
	auth := NewAuthHandler(authStore)

	orderStore := store.NewOrderStore(db)
	order := NewOrderHandler(orderStore)

	productStore := store.NewProductStore(db)
	product := NewProductHandler(productStore)

	categoryStore := store.NewCategoryStore(db)
	category := NewCategoryHandler(categoryStore)

	s := &Server{
		listenAddr: ":" + port,
		user:       user,
		auth:       auth,
		product:    product,
		order:      order,
		category:   category,
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

func getUserIDCtx(r *http.Request) (int, error) {
	idStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		return 0, fmt.Errorf("user id not found in context")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %s", idStr)
	}

	return id, nil
}
