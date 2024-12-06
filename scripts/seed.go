package main

import (
	"fmt"
	"log"
	"online-store-backend/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db := database.NewPostgresConnection()
	defer db.Close()

	products := []struct {
		Name        string
		Description string
		Price       float64
		Category    string
		Stock       int
	}{
		{"Apple iPhone 14", "Latest Apple smartphone", 999.99, "electronics", 50},
		{"Samsung Galaxy S23", "New Samsung flagship phone", 899.99, "electronics", 40},
		{"Nike Running Shoes", "Comfortable running shoes", 129.99, "fashion", 100},
		{"Levi's Jeans", "Classic denim jeans", 59.99, "fashion", 80},
		{"Sony Headphones", "Noise cancelling headphones", 199.99, "electronics", 30},
		{"Coffee Beans", "Premium Arabica coffee beans", 29.99, "grocery", 200},
		{"Green Tea", "Organic green tea leaves", 19.99, "grocery", 150},
	}

	for _, p := range products {
		query := `INSERT INTO products (name, description, price, category, stock) VALUES ($1, $2, $3, $4, $5)`
		_, err := db.Exec(query, p.Name, p.Description, p.Price, p.Category, p.Stock)
		if err != nil {
			log.Printf("Error seeding product %s: %v", p.Name, err)
		} else {
			log.Printf("Seeded product: %s", p.Name)
		}
	}

	fmt.Println("Seeding complete.")
}
