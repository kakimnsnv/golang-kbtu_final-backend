package http_v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var RouterModule = fx.Options(
	fx.Provide(NewRouter),
	fx.Invoke(RegisterRoutes),
)

// NewRouter initializes the Gin engine.
func NewRouter(logger *zap.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // Set mode to "release" for production
	router := gin.New()

	return router
}

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
func RegisterRoutes(router *gin.Engine) {
	// Example: Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})
}
