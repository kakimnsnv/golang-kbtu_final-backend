package http_v1

import (
	"final/internal/delivery/http/v1/middlewares"
	"final/internal/delivery/http/v1/routes"
	auth_interface "final/internal/features/auth/interface"
	cart_interface "final/internal/features/cart/interface"
	product_interface "final/internal/features/product/interface"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// RegisterRoutes registers the API routes.
func NewRouter(logger *zap.Logger, authUsecase auth_interface.AuthUsecase, productUsecase product_interface.ProductUseCase, cartUsecase cart_interface.CartUsecase) *gin.Engine {
	// MARK: create Router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// MARK: Register metrics
	middlewares.RegisterMetrics()

	// MARK: Register Prometheus middleware
	router.Use(middlewares.PrometheusMiddleware())

	// MARK: Register Prometheus metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := router.Group("/api/v1")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	auth := api.Group("/auth")
	{ // MARK: Auth routes
		routes.NewAuthRoute(auth, logger, authUsecase)
	}

	{ // MARK: Product routes
		routes.NewProductRoute(api, logger, productUsecase)
	}

	{ // MARK: Cart routes //TODO: add middleware for auth and role
		routes.NewCartRoute(api, logger, cartUsecase)
	}

	// authorizedAccess := api.Group("/", middlewares.AuthMiddleware(logger))
	// {
	// 	authorizedAccess.GET("/admin", middlewares.RoleMiddleware(logger, auth.RoleAdmin), func(c *gin.Context) {
	// 		c.JSON(200, gin.H{"message": "Admin access granted"})
	// 	})

	// 	authorizedAccess.GET("/dashboard", middlewares.RoleMiddleware(logger, auth.RoleUser|auth.RoleAdmin), func(c *gin.Context) {
	// 		c.JSON(200, gin.H{"message": "Welcom User or Admin!"})
	// 	})
	// }
	return router
}
