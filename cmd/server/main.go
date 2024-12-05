package main

import (
	"context"
	http_v1 "final/internal/delivery/http/v1"
	auth_interface "final/internal/features/auth/interface"
	auth_repo "final/internal/features/auth/repo"
	auth_usecase "final/internal/features/auth/usecase"
	product_interface "final/internal/features/product/interface"
	product_repo "final/internal/features/product/repo"
	product_usecase "final/internal/features/product/usecase"
	"final/internal/infrastructure"
	"final/internal/infrastructure/db_gen"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fx.New(
		fx.Provide(
			zap.NewProduction,
			fx.Annotate(
				auth_repo.New,
				fx.As(new(auth_interface.AuthRepo)),
			),
			fx.Annotate(
				auth_usecase.New,
				fx.As(new(auth_interface.AuthUsecase)),
			),
			fx.Annotate(
				product_repo.New,
				fx.As(new(product_interface.ProductRepo)),
			),
			fx.Annotate(
				product_usecase.New,
				fx.As(new(product_interface.ProductUseCase)),
			),
			http_v1.NewRouter,
			fx.Annotate(
				infrastructure.NewDB,
				fx.As(new(db_gen.DBTX)),
			),
			db_gen.New,
		),
		fx.Invoke(StartServer),
	)

	app.Run()
}

// StartServer starts the HTTP server and gracefully shuts it down.
func StartServer(lifecycle fx.Lifecycle, logger *zap.Logger, router *gin.Engine) {
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("Starting HTTP server", zap.String("addr", server.Addr))
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("ListenAndServe: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server...")
			timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return server.Shutdown(timeoutCtx)
		},
	})
}
