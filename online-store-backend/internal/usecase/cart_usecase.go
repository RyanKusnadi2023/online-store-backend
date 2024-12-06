package usecase

import (
	"errors"
	"online-store-backend/internal/domain/models"
	"online-store-backend/internal/domain/repository"
)

type CartUsecase interface {
	AddToCart(customerID, productID, quantity int) (*models.CartItem, error)
	GetCartItems(customerID int) ([]models.CartItem, error)
	RemoveFromCart(customerID, cartItemID int) error
}

type cartUsecase struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartUsecase(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartUsecase {
	return &cartUsecase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (u *cartUsecase) AddToCart(customerID, productID, quantity int) (*models.CartItem, error) {
	if quantity < 1 {
		return nil, errors.New("quantity must be at least 1")
	}

	product, err := u.productRepo.GetProductByID(productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	if product.Stock < quantity {
		return nil, errors.New("insufficient stock")
	}

	return u.cartRepo.AddToCart(customerID, productID, quantity)
}

func (u *cartUsecase) GetCartItems(customerID int) ([]models.CartItem, error) {
	return u.cartRepo.GetCartItems(customerID)
}

func (u *cartUsecase) RemoveFromCart(customerID, cartItemID int) error {
	return u.cartRepo.RemoveFromCart(customerID, cartItemID)
}
