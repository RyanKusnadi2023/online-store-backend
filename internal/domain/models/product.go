package models

type Product struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price"`
	Category    string  `json:"category" db:"category"`
	Stock       int     `json:"stock" db:"stock"`
}
