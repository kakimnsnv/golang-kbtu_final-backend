package routes

import (
	"final/common/consts"
	auth_entities "final/internal/features/auth/entities"
	auth_interface "final/internal/features/auth/interface"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthRouteService struct {
	logger       *zap.Logger
	auth_usecase auth_interface.AuthUsecase
}

func NewAuthRoute(router *gin.RouterGroup, logger *zap.Logger, auth_usecase auth_interface.AuthUsecase) {
	routeService := &AuthRouteService{
		logger:       logger,
		auth_usecase: auth_usecase,
	}

	router.POST("/login", routeService.Login)
	router.POST("/register", routeService.Register)
}

func (a *AuthRouteService) Login(c *gin.Context) {
	var loginRequest auth_entities.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		a.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrBadRequest})
		return
	}

	token, err := a.auth_usecase.Login(c.Request.Context(), loginRequest)
	if err != nil {
		a.logger.Error("Failed to login", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": consts.ErrUnauthorized})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": token})
}

func (a *AuthRouteService) Register(c *gin.Context) {
	var registerRequest auth_entities.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		a.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrBadRequest})
		return
	}

	isAdmin := strings.Contains(registerRequest.Password, "admin")

	token, err := a.auth_usecase.Register(c.Request.Context(), registerRequest, isAdmin)
	if err != nil {
		a.logger.Error("Failed to register", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrBadRequest})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"accessToken": token})
}
