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
	ErrCustomerNotFound = errors.New("user not found")
)

type CustomerStorer interface {
	GetById(id int) (*models.Customer, error)
	Update(id int, data models.UpdateCustomerReq) (*models.Customer, error)
	Delete(id int) error
}

type CustomerStore struct {
	db *sql.DB
}

func NewCustomerStore(db *sql.DB) *CustomerStore {
	return &CustomerStore{
		db: db,
	}
}

func (s *CustomerStore) GetById(id int) (*models.Customer, error) {
	stmt, err := s.db.Prepare(`
		SELECT * FROM CUSTOMERS WHERE ID = $1
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoCustomer(rows)
	}

	return nil, ErrCustomerNotFound
}

func (s *CustomerStore) GetByEmail(email string) (*models.Customer, error) {
	stmt, err := s.db.Prepare(`
		SELECT * FROM CUSTOMERS WHERE EMAIL = $1
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(email)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanIntoCustomer(rows)
	}

	return nil, ErrCustomerNotFound
}

func (s *CustomerStore) Update(id int, data models.UpdateCustomerReq) (*models.Customer, error) {
	stmt, err := s.db.Prepare(`
		UPDATE CUSTOMERS 
		SET 
			EMAIL = COALESCE(NULLIF($1, ''), EMAIL),
			FIRST_NAME = COALESCE(NULLIF($2, ''), FIRST_NAME),
			LAST_NAME = COALESCE(NULLIF($3, ''), LAST_NAME),
			DATE_OF_BIRTH = CASE WHEN $4::date IS NOT NULL THEN $4 ELSE NULL END,
			UPDATE_AT=NOW()
		WHERE ID = $5
		RETURNING *
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(
		data.Email,
		data.FirstName,
		data.LastName,
		data.DateOfBirth.Format(time.RFC3339),
		id,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				constraintName := pqErr.Constraint
				switch constraintName {
				case "customers_email_key":
					return nil, ErrEmailAlreadyExists
				default:
					return nil, fmt.Errorf("unique constraint %s violation", constraintName)
				}
			}
		}
		return nil, err
	}

	if rows.Next() {
		return scanIntoCustomer(rows)
	}

	return nil, err
}

func (s *CustomerStore) Delete(id int) error {
	if _, err := s.GetById(id); err != nil {
		return err
	}

	stmt, err := s.db.Prepare(`
		DELETE FROM CUSTOMERS WHERE ID = $1
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Query(id)
	if err != nil {
		return err
	}

	return err
}

func scanIntoCustomer(rows *sql.Rows) (*models.Customer, error) {
	c := &models.Customer{}
	err := rows.Scan(
		&c.ID,
		&c.Email,
		&c.FirstName,
		&c.LastName,
		&c.DateOfBirth,
		&c.Password,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	return c, err
}
