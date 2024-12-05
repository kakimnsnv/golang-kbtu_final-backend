package infrastructure

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// LoggerModule exports the logger for fx.
var LoggerModule = fx.Options(
	fx.Provide(NewLogger),
)

// NewLogger initializes a zap logger.
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction() // Replace with zap.NewDevelopment() for local development
}
