package middlewares

import (
	"final/common/consts"
	"final/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RoleMiddleware(logger *zap.Logger, allowedRoles auth.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetInt(consts.ContextUserRole)

		if userRole&int(allowedRoles) == 1 {
			c.Next()
			return
		}

		logger.Error("User does not have the required role")
		c.JSON(http.StatusForbidden, gin.H{"error": consts.ErrForbidden})
		c.Abort()
	}

}
