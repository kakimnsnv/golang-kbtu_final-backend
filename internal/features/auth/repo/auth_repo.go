package auth_repo

import (
	"context"
	auth_entities "final/internal/features/auth/entities"
	auth_interface "final/internal/features/auth/interface"
	"final/pkg/auth"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepoImpl struct {
	logger *zap.Logger
	db     *sqlx.DB
}

var _ auth_interface.AuthRepo = (*AuthRepoImpl)(nil)

func New(logger *zap.Logger, db *sqlx.DB) *AuthRepoImpl {
	return &AuthRepoImpl{
		logger: logger,
		db:     db,
	}
}

func (r *AuthRepoImpl) GetUserByEmail(ctx context.Context, email string) (auth_entities.User, error) {
	var dbUser auth_entities.User
	const q = `SELECT * FROM users WHERE email = $1`

	if err := r.db.GetContext(ctx, &dbUser, q, email); err != nil {
		r.logger.Error("failed to get user by email", zap.Error(err))
		return auth_entities.User{}, err
	}

	return dbUser, nil
}

func (r *AuthRepoImpl) CreateUser(ctx context.Context, email, password string, role auth.Role) (auth_entities.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error("failed to hash password", zap.Error(err))
		return auth_entities.User{}, err
	}

	const q = `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING *`
	var dbUser auth_entities.User

	if err := r.db.GetContext(ctx, &dbUser, q, email, hashedPassword, role); err != nil {
		r.logger.Error("failed to create user", zap.Error(err))
		return auth_entities.User{}, err
	}
	return dbUser, nil
}
