package mockdb

import (
	"context"

	"go.uber.org/zap"
)

type Handler struct{}

// Config holds all values used to configure and connect to a postgres db.
type Config struct{}

func (c *Config) ConnectAndMigrate(logger *zap.Logger) (*Handler, error) {
	logger.Info("Connecting to mock database...")
	return &Handler{}, nil
}

func (h *Handler) Shutdown(ctx context.Context) error {
	return nil
}
