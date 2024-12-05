package routes

import (
	product_entities "final/internal/features/product/entities"
	product_interface "final/internal/features/product/interface"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProductRouteService struct {
	logger  *zap.Logger
	usecase product_interface.ProductUseCase
}

func NewProductRoute(router *gin.RouterGroup, logger *zap.Logger, usecase product_interface.ProductUseCase) {
	rs := &ProductRouteService{
		logger:  logger,
		usecase: usecase,
	}

	router.GET("/products", rs.GetProducts)
	router.POST("/products", rs.CreateProduct)
	router.GET("/products/:id", rs.GetProduct)
	router.PUT("/products/:id", rs.UpdateProduct)
	router.DELETE("/products/:id", rs.DeleteProduct)
}

func (rs *ProductRouteService) GetProducts(c *gin.Context) {
	products, err := rs.usecase.GetAllProducts(c.Request.Context())
	if err != nil {
		rs.logger.Error("Failed to get products", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to get products"})
		return
	}

	c.JSON(200, products)
}

func (rs *ProductRouteService) CreateProduct(c *gin.Context) {
	var product product_entities.ProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		rs.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(400, gin.H{"error": "Failed to bind JSON"})
		return
	}

	productResponse, err := rs.usecase.CreateProduct(c.Request.Context(), product)
	if err != nil {
		rs.logger.Error("Failed to create product", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(200, productResponse)
}

func (rs *ProductRouteService) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := rs.usecase.GetProductByID(c.Request.Context(), id)
	if err != nil {
		rs.logger.Error("Failed to get product", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to get product"})
		return
	}

	c.JSON(200, product)
}

func (rs *ProductRouteService) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product product_entities.ProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		rs.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(400, gin.H{"error": "Failed to bind JSON"})
		return
	}

	productResponse, err := rs.usecase.UpdateProduct(c.Request.Context(), id, product)
	if err != nil {
		rs.logger.Error("Failed to update product", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(200, productResponse)
}

func (rs *ProductRouteService) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	err := rs.usecase.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		rs.logger.Error("Failed to delete product", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(200, gin.H{"message": "Product deleted"})
}
