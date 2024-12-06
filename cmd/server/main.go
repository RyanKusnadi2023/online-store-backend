package main

import (
	"log"
	"net/http"
	"online-store-backend/internal/config"
	"online-store-backend/internal/domain/repository"
	"online-store-backend/internal/handler"
	"online-store-backend/internal/usecase"
	"online-store-backend/pkg/database"
	"online-store-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()
	db := database.NewPostgresConnection()
	defer db.Close()

	// Repositories
	customerRepo := repository.NewCustomerRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Usecases
	customerUsecase := usecase.NewCustomerUsecase(customerRepo, cfg.JWT.Secret)
	productUsecase := usecase.NewProductUsecase(productRepo)
	cartUsecase := usecase.NewCartUsecase(cartRepo, productRepo)
	orderUsecase := usecase.NewOrderUsecase(cartRepo, productRepo, orderRepo)

	// Handlers
	customerHandler := handler.NewCustomerHandler(customerUsecase)
	productHandler := handler.NewProductHandler(productUsecase)
	cartHandler := handler.NewCartHandler(cartUsecase)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	router := gin.Default()

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Auth
	router.POST("/register", customerHandler.Register)
	router.POST("/login", customerHandler.Login)

	// Products
	router.GET("/products/category/:category", productHandler.GetProductsByCategory)

	// Cart (Protected)
	auth := middleware.AuthMiddleware(cfg.JWT.Secret)
	cartGroup := router.Group("/cart")
	cartGroup.Use(auth)
	{
		cartGroup.POST("/add", cartHandler.AddToCart)
		cartGroup.GET("/", cartHandler.GetCartItems)
		cartGroup.DELETE("/:id", cartHandler.RemoveFromCart)
	}

	// Checkout (Protected)
	router.POST("/checkout", auth, orderHandler.Checkout)

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
