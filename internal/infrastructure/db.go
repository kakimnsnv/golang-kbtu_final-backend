package infrastructure

import (
	"context"
	"final/common/consts"
	"os"

	"github.com/jackc/pgx/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewDB(lc fx.Lifecycle, logger *zap.Logger) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv(consts.DB_DSN))
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close(context.Background())
		},
	})

	return conn
}
