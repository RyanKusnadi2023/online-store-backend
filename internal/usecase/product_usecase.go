package usecase

import (
	"online-store-backend/internal/domain/models"
	"online-store-backend/internal/domain/repository"
)

type ProductUsecase interface {
	FetchProductsByCategory(category string) ([]models.Product, error)
	FetchProductByID(productID int) (*models.Product, error)
}

type productUsecase struct {
	productRepo repository.ProductRepository
}

func NewProductUsecase(productRepo repository.ProductRepository) ProductUsecase {
	return &productUsecase{productRepo: productRepo}
}

func (u *productUsecase) FetchProductsByCategory(category string) ([]models.Product, error) {
	return u.productRepo.GetProductsByCategory(category)
}

func (u *productUsecase) FetchProductByID(productID int) (*models.Product, error) {
	return u.productRepo.GetProductByID(productID)
}
