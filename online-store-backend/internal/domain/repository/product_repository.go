package repository

import (
	"database/sql"
	"errors"
	"online-store-backend/internal/domain/models"
)

type ProductRepository interface {
	GetProductsByCategory(category string) ([]models.Product, error)
	GetProductByID(productID int) (*models.Product, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetProductsByCategory(category string) ([]models.Product, error) {
	var rows *sql.Rows
	var err error

	if category == "" {
		rows, err = r.db.Query(`SELECT id, name, description, price, category, stock FROM products`)
	} else {
		rows, err = r.db.Query(`SELECT id, name, description, price, category, stock FROM products WHERE category = $1`, category)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Category, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *productRepository) GetProductByID(productID int) (*models.Product, error) {
	query := `SELECT id, name, description, price, category, stock FROM products WHERE id=$1`
	row := r.db.QueryRow(query, productID)
	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Category, &p.Stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
