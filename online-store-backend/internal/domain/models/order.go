package models

type Order struct {
	ID         int     `json:"id" db:"id"`
	CustomerID int     `json:"customer_id" db:"customer_id"`
	Total      float64 `json:"total" db:"total"`
	Status     string  `json:"status" db:"status"`
	CreatedAt  string  `json:"created_at" db:"created_at"`
}

type OrderItem struct {
	ID        int     `json:"id" db:"id"`
	OrderID   int     `json:"order_id" db:"order_id"`
	ProductID int     `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	Price     float64 `json:"price" db:"price"`
}
