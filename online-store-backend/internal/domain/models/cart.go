package models

type CartItem struct {
	ID         int `json:"id" db:"id"`
	CustomerID int `json:"customer_id" db:"customer_id"`
	ProductID  int `json:"product_id" db:"product_id"`
	Quantity   int `json:"quantity" db:"quantity"`
}
