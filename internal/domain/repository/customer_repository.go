package repository

import (
	"database/sql"
	"errors"
	"online-store-backend/internal/domain/models"
)

// CustomerRepository interface defines the methods to interact with the customer data
type CustomerRepository interface {
	CreateCustomer(customer *models.Customer) (int, error)
	GetCustomerByEmail(email string) (*models.Customer, error)
}

// customerRepository is the struct that implements the CustomerRepository interface
type customerRepository struct {
	db *sql.DB
}

// NewCustomerRepository creates a new instance of the customerRepository
func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db: db}
}

// CreateCustomer adds a new customer to the database
func (r *customerRepository) CreateCustomer(customer *models.Customer) (int, error) {
	query := `INSERT INTO customers (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, customer.Name, customer.Email, customer.Password).Scan(&customer.ID)
	if err != nil {
		return 0, err
	}
	return customer.ID, nil
}

// GetCustomerByEmail retrieves a customer by their email
func (r *customerRepository) GetCustomerByEmail(email string) (*models.Customer, error) {
	query := `SELECT id, name, email, password FROM customers WHERE email=$1`
	row := r.db.QueryRow(query, email)
	var cust models.Customer
	err := row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &cust, nil
}
