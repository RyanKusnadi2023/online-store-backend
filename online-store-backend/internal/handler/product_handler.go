package handler

import (
	"net/http"
	"online-store-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

func NewProductHandler(u usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: u}
}

func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")

	products, err := h.productUsecase.FetchProductsByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
