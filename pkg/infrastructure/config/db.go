package config

import (
	"context"

	"go.uber.org/zap"
)

type DBConfig struct {
	Host          string
	Port          string
	User          string
	Password      string
	Database      string
	MigrationPath string
}

type DBHandler interface {
	Shutdown(ctx context.Context) error
}
type Database interface {
	ConnectAndMigrate(logger *zap.Logger, c DBConfig) (*DBHandler, error)
}
