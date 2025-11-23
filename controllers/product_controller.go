package controllers

import (
	"net/http"
	"p3-graded-challenge-2-ziancarlos/models"
	"p3-graded-challenge-2-ziancarlos/service"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service service.ProductService
}

func NewProductController(service service.ProductService) *ProductController {
	return &ProductController{
		service: service,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.ProductRequest true "Product Request"
// @Success 201 {object} models.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /products [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req models.ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := c.service.CreateProduct(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get a list of all products
// @Tags products
// @Produce json
// @Success 200 {array} models.ProductResponse
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /products [get]
func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	products, err := c.service.GetAllProducts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /products/{id} [get]
func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := c.service.GetProductByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update product by ID
// @Description Update a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.ProductRequest true "Product Request"
// @Success 200 {object} models.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /products/{id} [put]
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := c.service.UpdateProduct(ctx.Request.Context(), id, &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete product by ID
// @Description Delete a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeleteProduct(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

