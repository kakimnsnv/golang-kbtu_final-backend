package routes

import (
	"final/common/consts"
	pg_errors "final/common/kerrors"
	"final/internal/delivery/http/v1/middlewares"
	cart_interface "final/internal/features/cart/interface"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CartRouteService struct {
	logger  *zap.Logger
	usecase cart_interface.CartUsecase
}

func NewCartRoute(router *gin.RouterGroup, logger *zap.Logger, usecase cart_interface.CartUsecase) {
	rs := &CartRouteService{
		logger:  logger,
		usecase: usecase,
	}

	router.GET("/cart", middlewares.AuthMiddleware(logger), rs.GetCart)
	router.DELETE("/cart", middlewares.AuthMiddleware(logger), rs.DeleteCart)

	router.POST("/cart/product/:id", middlewares.AuthMiddleware(logger), rs.AddToCart)
	router.DELETE("/cart/product/:id", middlewares.AuthMiddleware(logger), rs.RemoveFromCart)
	router.PUT("/cart/product/:id", middlewares.AuthMiddleware(logger), rs.UpdateCart)
}

func (rs *CartRouteService) GetCart(c *gin.Context) {
	userID := c.GetString(consts.ContextUserID)
	cart, err := rs.usecase.GetCart(c.Request.Context(), userID)
	if err != nil {
		rs.logger.Error("failed to get cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, cart)
}

func (rs *CartRouteService) AddToCart(c *gin.Context) {
	userID := c.GetString(consts.ContextUserID)
	productID := c.Param("id")
	quantity := c.Query("quantity")
	if quantity == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity query param is required"})
		return
	}

	q, err := strconv.Atoi(quantity)
	if err != nil {
		rs.logger.Error("failed to convert quantity to int", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be a number"})
		return
	}

	if err := rs.usecase.AddToCart(c.Request.Context(), userID, productID, q); err != nil {
		rs.logger.Error("failed to add product to cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "product added to cart"})
}

func (rs *CartRouteService) RemoveFromCart(c *gin.Context) {
	userID := c.GetString(consts.ContextUserID)
	productID := c.Param("id")
	if err := rs.usecase.RemoveFromCart(c.Request.Context(), userID, productID); err != nil {
		rs.logger.Error("failed to remove product from cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "product removed from cart"})
}

func (rs *CartRouteService) UpdateCart(c *gin.Context) {
	userID := c.GetString(consts.ContextUserID)
	productID := c.Param("id")
	quantity := c.Query("quantity")
	if quantity == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity query param is required"})
		return
	}

	q, err := strconv.Atoi(quantity)
	if err != nil {
		rs.logger.Error("failed to convert quantity to int", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be a number"})
		return
	}

	if err := rs.usecase.UpdateCart(c.Request.Context(), userID, productID, q); err != nil {
		if val, ok := err.(*pg_errors.PgError); ok && val.ErrorType == pg_errors.PgErrorNoRowsAffected {
			rs.logger.Info(err.Error())
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		rs.logger.Error("failed to update cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "cart updated"})
}

func (rs *CartRouteService) DeleteCart(c *gin.Context) {
	userID := c.GetString(consts.ContextUserID)
	if err := rs.usecase.DeleteCart(c.Request.Context(), userID); err != nil {
		rs.logger.Error("failed to delete cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "cart deleted"})
}
