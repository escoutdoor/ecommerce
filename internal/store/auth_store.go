package store

import (
	"database/sql"
	"errors"

	"github.com/escoutdoor/ecommerce/internal/models"
	"github.com/escoutdoor/ecommerce/pkg/password"
)

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrEmailAlreadyExists     = errors.New("customer with this email address already exist")
)

type AuthStorer interface {
	Login(data models.LoginReq) (*models.Customer, error)
	Register(data models.RegisterReq) (*models.Customer, error)
}

type AuthStore struct {
	db            *sql.DB
	customerStore CustomerStore
}

func NewAuthStore(db *sql.DB) *AuthStore {
	return &AuthStore{
		db:            db,
		customerStore: CustomerStore{db: db},
	}
}

func (s *AuthStore) Login(data models.LoginReq) (*models.Customer, error) {
	customer, err := s.customerStore.GetByEmail(data.Email)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			return nil, ErrInvalidEmailOrPassword
		}

		return nil, err
	}

	if !password.ComparePasswords(customer.Password, data.Password) {
		return nil, ErrInvalidEmailOrPassword
	}

	return customer, nil
}

func (s *AuthStore) Register(data models.RegisterReq) (*models.Customer, error) {
	c, err := s.customerStore.GetByEmail(data.Email)
	if err != nil {
		if !errors.Is(err, ErrCustomerNotFound) {
			return nil, err
		}
	}
	if c != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPass, err := password.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	stmt, err := s.db.Prepare(`
		INSERT INTO CUSTOMERS(EMAIL, FIRST_NAME, LAST_NAME, PASSWORD) 
		VALUES($1, $2, $3, $4) RETURNING * 
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(data.Email, data.FirstName, data.LastName, hashedPass)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoCustomer(rows)
	}

	return nil, err
}
