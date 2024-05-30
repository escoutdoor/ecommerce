package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/escoutdoor/ecommerce/internal/models"
	"github.com/escoutdoor/ecommerce/pkg/password"
)

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrEmailAlreadyExists     = errors.New("user with this email address already exist")
)

type AuthStorer interface {
	Login(models.LoginReq) (*models.User, error)
	Register(models.RegisterReq) (*models.User, error)
}

type AuthStore struct {
	db        *sql.DB
	userStore UserStore
}

func NewAuthStore(db *sql.DB) *AuthStore {
	return &AuthStore{
		db:        db,
		userStore: UserStore{db: db},
	}
}

func (s *AuthStore) Login(data models.LoginReq) (*models.User, error) {
	user, err := s.userStore.GetByEmail(data.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidEmailOrPassword
		}

		return nil, err
	}

	if !password.ComparePasswords(user.Password, data.Password) {
		return nil, ErrInvalidEmailOrPassword
	}

	return user, nil
}

func (s *AuthStore) Register(data models.RegisterReq) (*models.User, error) {
	u, err := s.userStore.GetByEmail(data.Email)
	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			return nil, err
		}
	}
	if u != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPass, err := password.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	stmt, err := s.db.Prepare(`
		INSERT INTO USERS(EMAIL, FIRST_NAME, LAST_NAME, DATE_OF_BIRTH, PASSWORD) 
		VALUES($1, $2, $3, $4, $5) RETURNING * 
	`)
	if err != nil {
		return nil, err
	}

	var birthdate *time.Time
	if data.DateOfBirth != "" {
		pb, err := time.Parse("2006-01-02", data.DateOfBirth)
		if err != nil {
			return nil, err
		}

		birthdate = &pb
	}

	rows, err := stmt.Query(data.Email, data.FirstName, data.LastName, birthdate, hashedPass)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, err
}
