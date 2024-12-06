package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"online-store-backend/internal/config"
	"online-store-backend/internal/domain/repository"
	"online-store-backend/internal/handler"
	"online-store-backend/internal/usecase"
	"online-store-backend/pkg/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()
	db := database.NewPostgresConnection()
	defer db.Close()

	customerRepo := repository.NewCustomerRepository(db)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo, cfg.JWT.Secret)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	router := gin.Default()

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Auth routes
	router.POST("/register", customerHandler.Register)
	router.POST("/login", customerHandler.Login)

	// Later, for protected routes, we will add middleware, e.g:
	// auth := middleware.AuthMiddleware(cfg.JWT.Secret)
	// protected := router.Group("/")
	// protected.Use(auth)
	// protected.GET("/profile", someProfileHandler)

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
