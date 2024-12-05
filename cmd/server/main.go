package main

import (
	"context"
	http_v1 "final/internal/delivery/http/v1"
	"final/internal/infrastructure"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		http_v1.RouterModule,
		infrastructure.LoggerModule,
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
