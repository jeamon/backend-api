package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // used to register a source for migration files.
	healthPgx4 "github.com/hellofresh/health-go/v4/checks/pgx4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

// Handler holds two connections to the database. sqlx connection and original pgx connection.
// PGx is primarily used for PostgreSQL and StdDB for migration and integration testing.
type Handler struct {
	StdDB *sql.DB
	PGx   *pgxpool.Pool
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

// ConnectAndMigrate establishes a connection to the database based on the configuration provided.
// Additionally it performs a data migration to ensure that the schema is on the latest version.
func (c *Config) ConnectAndMigrate(logger *zap.Logger) (*Handler, error) {
	logger.Info("Connecting to postgres database...")

	pgxConn, stdConn, err := Connect(c)
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %w", err)
	}

	logger.Info("Migrate db schema...")

	err = Migrate(stdConn, c.MigrationPath, "up")
	if err != nil {
		logger.Info("could not migrate schema up", zap.Error(err))
		return nil, fmt.Errorf("could not migrate schema up: %w", err)
	}

	return &Handler{
		PGx:   pgxConn,
		StdDB: stdConn,
	}, nil
}

// Shutdown closes all connections.
func (h *Handler) Shutdown(ctx context.Context) error {
	err := h.StdDB.Close()
	if err != nil {
		return fmt.Errorf("close sqlx failed: %w", err)
	}

	h.PGx.Close()

	return nil
}

func (c *Config) ToDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Database)
}

func Connect(c *Config) (*pgxpool.Pool, *sql.DB, error) {
	connURL := c.ToDSN()

	conf, err := pgxpool.ParseConfig(connURL)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse config string: %w", err)
	}

	pgxConn, err := pgxpool.ConnectConfig(context.Background(), conf)
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect via pgx: %w", err)
	}

	connStr := stdlib.RegisterConnConfig(conf.ConnConfig)

	stdConn, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect via std library: %w", err)
	}

	return pgxConn, stdConn, nil
}

// Check allows to probe the liveness of postgres database.
func (c *Config) Check() func(ctx context.Context) error {
	return healthPgx4.New(healthPgx4.Config{DSN: c.ToDSN()})
}

func Migrate(db *sql.DB, path, action string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+path,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create new migrate instance: %w", err)
	}

	switch action {
	case "up":
		err = m.Up()
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	case "down":
		err = m.Down()
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
