package mongo

import (
	"context"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file" // used to register a source for migration files.
	"go.uber.org/zap"

	healthMongo "github.com/hellofresh/health-go/v4/checks/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Handler holds a connection to the nosql database
type Handler struct {
	Client *mongo.Client
}

// Config holds all values used to configure and connect to a postgres db.
type Config struct {
	Host          string
	Port          string
	User          string
	Password      string
	Database      string
	MigrationPath string
}

// Connect establishes a connection to the database based on the configuration provided.
func (c *Config) ConnectAndMigrate(logger *zap.Logger) (*Handler, error) {
	logger.Info("Connecting to mongo database...")

	mongoConn, err := GetClient(c)
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %w", err)
	}

	return &Handler{
		Client: mongoConn,
	}, nil
}

// Shutdown closes all connections.
func (h *Handler) Shutdown(ctx context.Context) error {
	err := h.Client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("failed to close mongo client connection: %w", err)
	}

	return nil
}

func (c *Config) ToDSN() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", c.User, c.Password, c.Host, c.Port)
}

func GetClient(c *Config) (*mongo.Client, error) {
	connURL := c.ToDSN()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connURL))
	if err != nil {
		return nil, fmt.Errorf("could not connect to mongo database: %w", err)
	}

	return client, nil
}

// Check allows to probe the liveness of mongo database.
func (c *Config) Check() func(ctx context.Context) error {
	return healthMongo.New(healthMongo.Config{DSN: c.ToDSN()})
}
