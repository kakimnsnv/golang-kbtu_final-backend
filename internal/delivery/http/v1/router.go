package http_v1

import (
	"final/internal/delivery/http/v1/routes"
	auth_interface "final/internal/features/auth/interface"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger is a middleware that logs HTTP requests using zap.
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		c.Next()
	}
}

// RegisterRoutes registers the API routes.
func NewRouter(logger *zap.Logger, authUsecase auth_interface.AuthUsecase) *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // Set mode to "release" for production
	router := gin.Default()

	api := router.Group("/api/v1")

	// Example: Health check route
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	auth := api.Group("/auth")
	{
		routes.NewAuthRoute(auth, logger, authUsecase)
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
