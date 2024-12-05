package auth_repo

import (
	"context"
	auth_entities "final/internal/features/auth/entities"
	auth_interface "final/internal/features/auth/interface"
	"final/internal/infrastructure/db_gen"
	"final/pkg/auth"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepoImpl struct {
	logger *zap.Logger
	db     *db_gen.Queries
}

var _ auth_interface.AuthRepo = (*AuthRepoImpl)(nil)

func New(logger *zap.Logger, db *db_gen.Queries) *AuthRepoImpl {
	return &AuthRepoImpl{
		logger: logger,
		db:     db,
	}
}

func (r *AuthRepoImpl) GetUserByEmail(ctx context.Context, email string) (auth_entities.User, error) {
	dbUser, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		r.logger.Error("failed to get user by email", zap.Error(err))
		return auth_entities.User{}, err
	}

	return auth_entities.FromDBUser(dbUser), nil
}

func (r *AuthRepoImpl) CreateUser(ctx context.Context, email, password string, role auth.Role) (auth_entities.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error("failed to hash password", zap.Error(err))
		return auth_entities.User{}, err
	}
	user, err := r.db.CreateUser(ctx, db_gen.CreateUserParams{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         int32(role),
	})
	if err != nil {
		r.logger.Error("failed to create user", zap.Error(err))
		return auth_entities.User{}, err
	}
	return auth_entities.FromDBUser(user), nil
}
