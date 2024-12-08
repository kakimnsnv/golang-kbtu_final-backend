package routes

import (
	"final/common/consts"
	"final/internal/delivery/http/v1/middlewares"
	product_entities "final/internal/features/product/entities"
	product_interface "final/internal/features/product/interface"
	"final/pkg/auth"
	"net/http"

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

	// Collections
	router.GET("/products", rs.GetProducts)
	router.GET("/trending", rs.GetProducts) // TODO: Implement trending products

	router.GET("/products/:id", rs.GetProduct)

	// User
	router.POST("/products/:id/like", middlewares.AuthMiddleware(logger), rs.LikeProduct)
	router.POST("/products/:id/unlike", middlewares.AuthMiddleware(logger), rs.UnlikeProduct)

	// Admin
	router.POST("/products", middlewares.AuthMiddleware(logger), middlewares.RoleMiddleware(logger, auth.RoleAdmin), rs.CreateProduct)
	router.PUT("/products/:id", middlewares.AuthMiddleware(logger), middlewares.RoleMiddleware(logger, auth.RoleAdmin), rs.UpdateProduct)
	router.DELETE("/products/:id", middlewares.AuthMiddleware(logger), middlewares.RoleMiddleware(logger, auth.RoleAdmin), rs.DeleteProduct)
}

func (rs *ProductRouteService) GetProducts(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		products, err := rs.usecase.GetAllProducts(c.Request.Context())
		if err != nil {
			rs.logger.Error("Failed to get products", zap.Error(err))
			c.JSON(500, gin.H{"error": "Failed to get products"})
			return
		}

		c.JSON(200, products)
	} else {
		userID, _, err := middlewares.DecodeTokenAndGetIDAndRole(authHeader)
		if err != nil {
			rs.logger.Error("Failed to decode token", zap.Error(err))
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		products, err := rs.usecase.GetAllProducts(c.Request.Context(), userID)
		if err != nil {
			rs.logger.Error("Failed to get products", zap.Error(err))
			c.JSON(500, gin.H{"error": "Failed to get products"})
			return
		}

		c.JSON(200, products)
	}
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
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		id := c.Param("id")
		product, err := rs.usecase.GetProductByID(c.Request.Context(), id)
		if err != nil {
			rs.logger.Error("Failed to get product", zap.Error(err))
			c.JSON(500, gin.H{"error": "Failed to get product"})
			return
		}

		c.JSON(200, product)
	} else {
		userID, _, err := middlewares.DecodeTokenAndGetIDAndRole(authHeader)
		if err != nil {
			rs.logger.Error("Failed to decode token", zap.Error(err))
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		product, err := rs.usecase.GetProductByID(c.Request.Context(), c.Param("id"), userID)
		if err != nil {
			rs.logger.Error("Failed to get product", zap.Error(err))
			c.JSON(500, gin.H{"error": "Failed to get product"})
			return
		}
		c.JSON(200, product)
	}
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

func (rs *ProductRouteService) LikeProduct(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(consts.ContextUserID)

	if err := rs.usecase.LikeProduct(c.Request.Context(), userID, id); err != nil {
		rs.logger.Error("Failed to like product", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to like product"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Product liked"})
}

func (rs *ProductRouteService) UnlikeProduct(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString(consts.ContextUserID)

	if err := rs.usecase.UnlikeProduct(c.Request.Context(), userID, id); err != nil {
		rs.logger.Error("Failed to unlike product", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to unlike product"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Product unliked"})
}
