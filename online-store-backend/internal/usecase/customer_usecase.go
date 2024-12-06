package usecase

import (
	"errors"
	"online-store-backend/internal/domain/models"
	"online-store-backend/internal/domain/repository"
	"online-store-backend/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type CustomerUsecase interface {
	Register(name, email, password string) (*models.Customer, error)
	Login(email, password string) (string, error)
}

type customerUsecase struct {
	customerRepo repository.CustomerRepository
	jwtSecret    string
}

func NewCustomerUsecase(repo repository.CustomerRepository, jwtSecret string) CustomerUsecase {
	return &customerUsecase{
		customerRepo: repo,
		jwtSecret:    jwtSecret,
	}
}

func (u *customerUsecase) Register(name, email, password string) (*models.Customer, error) {
	// Check if customer already exists
	existing, err := u.customerRepo.GetCustomerByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	cust := &models.Customer{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}

	_, err = u.customerRepo.CreateCustomer(cust)
	if err != nil {
		return nil, err
	}

	// Remove password before returning
	cust.Password = ""
	return cust, nil
}

func (u *customerUsecase) Login(email, password string) (string, error) {
	cust, err := u.customerRepo.GetCustomerByEmail(email)
	if err != nil {
		return "", err
	}
	if cust == nil {
		return "", errors.New("invalid credentials")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(cust.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(u.jwtSecret, cust.ID, cust.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
