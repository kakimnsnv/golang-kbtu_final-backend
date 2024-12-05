package middlewares

import (
	"errors"
	"final/common/consts"
	"final/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(consts.HTTPAuthorizationHeader)
		if token == "" {
			logger.Error("No token provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": consts.ErrUnauthorized})
			c.Abort()
			return
		}

		userRole, err := decodeTokenAndGetRole(token)
		if err != nil {
			logger.Error("Failed to decode token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": consts.ErrUnauthorized})
			c.Abort()
			return
		}

		c.Set(consts.ContextUserRole, userRole)
		c.Next()
	}
}

func decodeTokenAndGetRole(token string) (auth.Role, error) {
	claims, err := auth.ParseJWT(token)
	if err != nil {
		return auth.RoleNone, err
	}

	role, ok := claims[consts.ClaimsRole].(auth.Role)
	if !ok {
		return auth.RoleNone, errors.New(consts.ErrInvalidRole)
	}

	return role, nil
}
