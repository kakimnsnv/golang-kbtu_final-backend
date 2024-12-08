package postgres

import (
	"context"
	"final/common/consts"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func ConnectDB(logger *zap.Logger) *sqlx.DB {
	db, err := sqlx.ConnectContext(context.Background(), "postgres", os.Getenv(consts.DB_DSN))
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	logger.Info("connected to database")
	return db
}
