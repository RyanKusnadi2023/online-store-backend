package repository

import (
	"database/sql"
	"errors"
	"online-store-backend/internal/domain/models"
)

type CartRepository interface {
	AddToCart(customerID, productID, quantity int) (*models.CartItem, error)
	GetCartItems(customerID int) ([]models.CartItem, error)
	RemoveFromCart(customerID, cartItemID int) error
	GetCartItemByID(cartItemID int) (*models.CartItem, error)
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddToCart(customerID, productID, quantity int) (*models.CartItem, error) {
	existing, err := r.getCartItemByCustomerAndProduct(customerID, productID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		newQuantity := existing.Quantity + quantity
		updateQuery := `UPDATE cart_items SET quantity=$1 WHERE id=$2 RETURNING id, customer_id, product_id, quantity`
		updated := &models.CartItem{}
		err = r.db.QueryRow(updateQuery, newQuantity, existing.ID).
			Scan(&updated.ID, &updated.CustomerID, &updated.ProductID, &updated.Quantity)
		if err != nil {
			return nil, err
		}
		return updated, nil
	}

	insertQuery := `INSERT INTO cart_items (customer_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id, customer_id, product_id, quantity`
	newItem := &models.CartItem{}
	err = r.db.QueryRow(insertQuery, customerID, productID, quantity).
		Scan(&newItem.ID, &newItem.CustomerID, &newItem.ProductID, &newItem.Quantity)
	if err != nil {
		return nil, err
	}

	return newItem, nil
}

func (r *cartRepository) GetCartItems(customerID int) ([]models.CartItem, error) {
	query := `SELECT id, customer_id, product_id, quantity FROM cart_items WHERE customer_id=$1`
	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(&item.ID, &item.CustomerID, &item.ProductID, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *cartRepository) RemoveFromCart(customerID, cartItemID int) error {
	query := `DELETE FROM cart_items WHERE id=$1 AND customer_id=$2`
	res, err := r.db.Exec(query, cartItemID, customerID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no cart item found or not authorized to remove")
	}

	return nil
}

func (r *cartRepository) GetCartItemByID(cartItemID int) (*models.CartItem, error) {
	query := `SELECT id, customer_id, product_id, quantity FROM cart_items WHERE id=$1`
	row := r.db.QueryRow(query, cartItemID)
	var item models.CartItem
	err := row.Scan(&item.ID, &item.CustomerID, &item.ProductID, &item.Quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) getCartItemByCustomerAndProduct(customerID, productID int) (*models.CartItem, error) {
	query := `SELECT id, customer_id, product_id, quantity FROM cart_items WHERE customer_id=$1 AND product_id=$2`
	row := r.db.QueryRow(query, customerID, productID)
	var item models.CartItem
	err := row.Scan(&item.ID, &item.CustomerID, &item.ProductID, &item.Quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
