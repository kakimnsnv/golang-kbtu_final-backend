package routes

import (
	"final/common/consts"
	"final/internal/delivery/http/v1/middlewares"
	order_entities "final/internal/features/order/entities"
	order_interface "final/internal/features/order/interface"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderRouteService struct {
	usecase order_interface.OrderUsecase
	logger  *zap.Logger
}

func NewOrderRoute(router *gin.RouterGroup, usecase order_interface.OrderUsecase, logger *zap.Logger) {
	rs := &OrderRouteService{
		usecase: usecase,
		logger:  logger,
	}

	router.GET("/orders", middlewares.AuthMiddleware(logger), rs.GetOrders)
	router.GET("/orders/:id", middlewares.AuthMiddleware(logger), rs.GetOrder)
	router.POST("/orders", middlewares.AuthMiddleware(logger), rs.CreateOrder)
	router.PUT("/orders/:id", middlewares.AuthMiddleware(logger), rs.UpdateOrder)
	router.DELETE("/orders/:id", middlewares.AuthMiddleware(logger), rs.DeleteOrder)
}

func (rs *OrderRouteService) GetOrders(c *gin.Context) {
	userID := c.GetString(consts.ContextUserID)
	orders, err := rs.usecase.GetOrders(c.Request.Context(), userID)
	if err != nil {
		rs.logger.Error("failed to get orders", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, orders)
}

func (rs *OrderRouteService) GetOrder(c *gin.Context) {
	//! userID := c.GetString(consts.ContextUserID) // TODO: validate that order belongs to user or user is admin. now anyone can get any order
	orderID := c.Param("id")
	order, err := rs.usecase.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		rs.logger.Error("failed to get order", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, order)
}

func (rs *OrderRouteService) CreateOrder(c *gin.Context) {
	userID := c.GetString(consts.ContextUserID)
	orderID, err := rs.usecase.PlaceOrder(c.Request.Context(), userID)
	if err != nil {
		rs.logger.Error("failed to create order", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"order_id": orderID})
}

func (rs *OrderRouteService) UpdateOrder(c *gin.Context) {
	orderID := c.Param("id")
	status := c.Query("status")
	statusInt, err := strconv.ParseInt(status, 10, 64)
	if err != nil {
		rs.logger.Error("failed to parse status", zap.Error(err))
		c.JSON(400, gin.H{"error": "invalid status"})
		return
	}
	err = rs.usecase.ChangeOrderStatus(c.Request.Context(), orderID, order_entities.OrderStatus(statusInt))
	if err != nil {
		rs.logger.Error("failed to update order", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "order updated"})
}

func (rs *OrderRouteService) DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")
	err := rs.usecase.CancelOrder(c.Request.Context(), orderID)
	if err != nil {
		rs.logger.Error("failed to delete order", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "order cancelled"})
}
