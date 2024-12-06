package repository

import (
	"database/sql"
	"online-store-backend/internal/domain/models"
)

type OrderRepository interface {
	CreateOrder(customerID int, total float64, status string) (*models.Order, error)
	CreateOrderItems(orderID int, items []models.OrderItem) error
	UpdateProductStock(productID, newStock int) error
	ClearCart(customerID int) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(customerID int, total float64, status string) (*models.Order, error) {
	query := `INSERT INTO orders (customer_id, total, status) VALUES ($1, $2, $3) RETURNING id, customer_id, total, status, created_at`
	var o models.Order
	err := r.db.QueryRow(query, customerID, total, status).Scan(&o.ID, &o.CustomerID, &o.Total, &o.Status, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *orderRepository) CreateOrderItems(orderID int, items []models.OrderItem) error {
	query := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`
	for _, item := range items {
		_, err := r.db.Exec(query, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *orderRepository) UpdateProductStock(productID, newStock int) error {
	query := `UPDATE products SET stock=$1 WHERE id=$2`
	_, err := r.db.Exec(query, newStock, productID)
	return err
}

func (r *orderRepository) ClearCart(customerID int) error {
	query := `DELETE FROM cart_items WHERE customer_id=$1`
	_, err := r.db.Exec(query, customerID)
	return err
}
