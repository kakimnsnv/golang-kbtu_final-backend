package auth_usecase

import (
	"context"
	auth_entities "final/internal/features/auth/entities"
	auth_interface "final/internal/features/auth/interface"
	"final/pkg/auth"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecaseImpl struct {
	logger *zap.Logger
	repo   auth_interface.AuthRepo
}

var _ auth_interface.AuthUsecase = (*AuthUsecaseImpl)(nil)

func New(logger *zap.Logger, repo auth_interface.AuthRepo) *AuthUsecaseImpl {
	return &AuthUsecaseImpl{
		logger: logger,
		repo:   repo,
	}
}

func (u *AuthUsecaseImpl) Login(ctx context.Context, loginRequest auth_entities.LoginRequest) (auth_entities.Token, error) {
	user, err := u.repo.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		return "", errors.Wrap(err, "failed to get user by email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)); err != nil {
		return "", errors.Wrap(err, "failed to compare password")
	}

	token, err := auth.CreateJWT(user.ID.String(), user.Role, time.Hour*24)
	if err != nil {
		return "", errors.Wrap(err, "failed to create JWT")
	}

	return auth_entities.Token(token), nil

}
func (u *AuthUsecaseImpl) Register(ctx context.Context, registerRequest auth_entities.RegisterRequest, isAdmin bool) (auth_entities.Token, error) {
	role := auth.RoleUser
	if isAdmin {
		role = auth.RoleAdmin
	}
	user, err := u.repo.CreateUser(ctx, registerRequest.Email, registerRequest.Password, role)
	if err != nil {
		return "", errors.Wrap(err, "failed to create user")
	}

	token, err := auth.CreateJWT(user.ID.String(), user.Role, time.Hour*24)
	if err != nil {
		return "", errors.Wrap(err, "failed to create JWT")
	}

	return auth_entities.Token(token), nil
}
