package middlewares

import (
	"errors"
	"final/common/consts"
	"final/pkg/auth"
	"net/http"
	"strings"

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
		if !strings.HasPrefix(token, "Bearer ") {
			logger.Error("Invalid token format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": consts.ErrUnauthorized})
			return
		}

		userID, userRole, err := decodeTokenAndGetRole(token)
		if err != nil {
			logger.Error("Failed to decode token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": consts.ErrUnauthorized})
			c.Abort()
			return
		}

		c.Set(consts.ContextUserRole, userRole)
		c.Set(consts.ContextUserID, userID)
		c.Next()
	}
}

func decodeTokenAndGetRole(tokenStr string) (string, int, error) {
	token := strings.Split(tokenStr, "Bearer ")[1]

	claims, err := auth.ParseJWT(token)
	if err != nil {
		return "", 0, err
	}
	val, ok := claims[consts.ClaimsRole].(float64)
	if !ok {
		return "", 0, errors.New(consts.ErrInvalidRole)
	}

	userID, ok := claims[consts.ClaimsUserID].(string)
	if !ok {
		return "", 0, errors.New(consts.ErrInvalidUserID)
	}

	role := int(val)

	return userID, role, nil
}
