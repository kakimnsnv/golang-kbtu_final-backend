package auth_interface

import (
	"context"
	auth_entities "final/internal/features/auth/entities"
	"final/pkg/auth"
)

type (
	AuthUsecase interface {
		Login(ctx context.Context, loginRequest auth_entities.LoginRequest) (auth_entities.Token, error)
		Register(ctx context.Context, registerRequest auth_entities.RegisterRequest, isAdmin bool) (auth_entities.Token, error)
	}
	AuthRepo interface {
		GetUserByEmail(ctx context.Context, email string) (auth_entities.User, error)
		CreateUser(ctx context.Context, email, password string, role auth.Role) (auth_entities.User, error)
	}
)
