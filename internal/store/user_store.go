package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/escoutdoor/ecommerce/internal/models"
	"github.com/lib/pq"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserStorer interface {
	GetByID(id int) (*models.User, error)
	Update(id int, data models.UpdateUserReq) (*models.User, error)
	Delete(id int) error
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) GetByID(id int) (*models.User, error) {
	stmt, err := s.db.Prepare(`
		SELECT * FROM USERS WHERE ID = $1
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, ErrUserNotFound
}

func (s *UserStore) GetByEmail(email string) (*models.User, error) {
	stmt, err := s.db.Prepare(`
		SELECT * FROM USERS WHERE EMAIL = $1
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(email)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, ErrUserNotFound
}

func (s *UserStore) Update(id int, data models.UpdateUserReq) (*models.User, error) {
	stmt, err := s.db.Prepare(`
		UPDATE USERS 
		SET 
			EMAIL = $1,
			FIRST_NAME = $2,
			LAST_NAME = $3,
			DATE_OF_BIRTH = $4
		WHERE ID = $5
		RETURNING *
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

	rows, err := stmt.Query(
		data.Email,
		data.FirstName,
		data.LastName,
		birthdate,
		id,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				constraintName := pqErr.Constraint
				switch constraintName {
				case "users_email_key":
					return nil, ErrEmailAlreadyExists
				default:
					return nil, fmt.Errorf("unique constraint %s violation", constraintName)
				}
			}
		}
		return nil, err
	}

	if rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, err
}

func (s *UserStore) Delete(id int) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}

	stmt, err := s.db.Prepare(`
		DELETE FROM USERS WHERE ID = $1
	`)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user cannot be deleted because it doesn't exist")
	}

	return err
}

func scanIntoUser(rows *sql.Rows) (*models.User, error) {
	u := &models.User{}
	err := rows.Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.DateOfBirth,
		&u.Password,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	return u, err
}
