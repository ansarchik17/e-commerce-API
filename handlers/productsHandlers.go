package handlers

import (
	"e-commerce-API/models"
	"e-commerce-API/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productRepo *repositories.ProductsRepository
}

func NewProductHandler(productRepo *repositories.ProductsRepository) *ProductHandler {
	return &ProductHandler{productRepo: productRepo}
}

func (handler *ProductHandler) AddProduct(c *gin.Context) {
	var request models.Product
	err := c.BindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Could not bind json body"))
		return
	}
	product := models.Product{
		ID:    request.ID,
		Name:  request.Name,
		Price: request.Price,
	}
	id, err := handler.productRepo.Create(c, product)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Could not add a product"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
