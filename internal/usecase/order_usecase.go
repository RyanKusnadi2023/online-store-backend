package usecase

import (
	"errors"
	"online-store-backend/internal/domain/models"
	"online-store-backend/internal/domain/repository"
)

type OrderUsecase interface {
	Checkout(customerID int) (*models.Order, []models.OrderItem, error)
}

type orderUsecase struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	orderRepo   repository.OrderRepository
}

func NewOrderUsecase(cartRepo repository.CartRepository, productRepo repository.ProductRepository, orderRepo repository.OrderRepository) OrderUsecase {
	return &orderUsecase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
	}
}

func (u *orderUsecase) Checkout(customerID int) (*models.Order, []models.OrderItem, error) {
	// Get cart items
	cartItems, err := u.cartRepo.GetCartItems(customerID)
	if err != nil {
		return nil, nil, err
	}

	if len(cartItems) == 0 {
		return nil, nil, errors.New("cart is empty, cannot checkout")
	}

	var total float64
	var orderItems []models.OrderItem

	// Check each cart item
	for _, cItem := range cartItems {
		product, err := u.productRepo.GetProductByID(cItem.ProductID)
		if err != nil {
			return nil, nil, err
		}
		if product == nil {
			return nil, nil, errors.New("product not found")
		}

		if product.Stock < cItem.Quantity {
			return nil, nil, errors.New("insufficient stock for product: " + product.Name)
		}

		lineTotal := product.Price * float64(cItem.Quantity)
		total += lineTotal

		orderItem := models.OrderItem{
			ProductID: cItem.ProductID,
			Quantity:  cItem.Quantity,
			Price:     product.Price,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Create order
	order, err := u.orderRepo.CreateOrder(customerID, total, "paid") // "paid" as placeholder
	if err != nil {
		return nil, nil, err
	}

	// Create order items
	err = u.orderRepo.CreateOrderItems(order.ID, orderItems)
	if err != nil {
		return nil, nil, err
	}

	// Update stock
	for i, cItem := range cartItems {
		product, err := u.productRepo.GetProductByID(cItem.ProductID)
		if err != nil {
			return nil, nil, err
		}
		newStock := product.Stock - cItem.Quantity
		err = u.orderRepo.UpdateProductStock(product.ID, newStock)
		if err != nil {
			return nil, nil, err
		}
		// Update the orderItems array with the now known order_id
		orderItems[i].OrderID = order.ID
	}

	// Clear cart
	err = u.orderRepo.ClearCart(customerID)
	if err != nil {
		return nil, nil, err
	}

	return order, orderItems, nil
}
