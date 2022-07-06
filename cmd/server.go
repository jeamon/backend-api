package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	health "github.com/hellofresh/health-go/v4"
	"github.com/jeamon/sample-rest-api/pkg/application"
	"github.com/jeamon/sample-rest-api/pkg/domain"
	"github.com/jeamon/sample-rest-api/pkg/infrastructure/config"
	"github.com/jeamon/sample-rest-api/pkg/infrastructure/mockdb"
	"github.com/jeamon/sample-rest-api/pkg/infrastructure/mongo"
	"github.com/jeamon/sample-rest-api/pkg/infrastructure/postgres"
	apisweb "github.com/jeamon/sample-rest-api/pkg/interfaces/public"
	"github.com/jeamon/sample-rest-api/pkg/interfaces/repository"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Execute runs the main service command.
func Execute(gitCommit, gitTag string) {
	var configData *config.Config
	rootCmd := &cobra.Command{
		Use:   "rest-api service",
		Short: "Demo RestFul API Backend",
		Run: func(cmd *cobra.Command, args []string) {
			configData.GitCommit = gitCommit
			configData.GitTag = gitTag
			err := run(configData, gitCommit, gitTag)
			if err != nil {
				log.Fatalf("Error occurred: %v", err)
			}

			fmt.Println("Stopped.")
		},
	}

	var cfgFile string
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./server.config.yaml)")

	cobra.OnInitialize(func() {
		configData = config.InitConfig(cfgFile)
	})

	if err := rootCmd.Execute(); err != nil {
		panic("unable to execute " + err.Error())
	}
}

//nolint
func run(configData *config.Config, gitCommit, gitTag string) error {
	var logger *zap.Logger
	var err error

	// Setup the logger with default fields.
	if configData.IsProduction {
		logger, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	} else {
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	}
	defer logger.Sync()

	logger = logger.With(
		zap.Bool("is_prod", configData.IsProduction),
		zap.String("git_commit", gitCommit),
		zap.String("git_tag", gitTag),
	)

	logger.Info("starting demo-rest-api service...", zap.String("version", gitTag))

	undo := zap.RedirectStdLog(logger)
	defer undo()

	if !configData.GinDisableReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(ginzap.RecoveryWithZap(logger, true))
	router.Use(cors.New(cors.Config{
		AllowWildcard:    true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Accept", "content-type", "User-Agent", "Accept-Language", "Referer", "DNT", "Connection", "Pragma", "Cache-Control", "TE", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	var h *health.Health
	var dbHandler config.DBHandler
	var scanInfosRepo domain.ScanInfosRepository
	switch configData.Database {
	case "postgres":
		postgresDB := postgres.Config{
			Host:          configData.DBPostgresConfig.Host,
			Port:          configData.DBPostgresConfig.Port,
			User:          configData.DBPostgresConfig.User,
			Password:      configData.DBPostgresConfig.Password,
			Database:      configData.DBPostgresConfig.DatabaseName,
			MigrationPath: "pkg/infrastructure/postgres/migrations",
		}
		pgHandler, err := postgresDB.ConnectAndMigrate(logger)
		if err != nil {
			return err
		}
		dbHandler = pgHandler
		scanInfosRepo = repository.NewPostgresScanInfosRepository(logger, pgHandler)

		h, err = health.New(health.WithChecks(
			health.Config{
				Name:    "postgres",
				Timeout: time.Second * 5,
				Check:   postgresDB.Check(),
			},
		))
		if err != nil {
			return fmt.Errorf("unable to initialize postgres database health checks: %v", err)
		}

	case "mongo":
		mongoDB := mongo.Config{
			Host:          configData.DBMongoConfig.Host,
			Port:          configData.DBMongoConfig.Port,
			User:          configData.DBMongoConfig.User,
			Password:      configData.DBMongoConfig.Password,
			Database:      configData.DBMongoConfig.DatabaseName,
			MigrationPath: "pkg/infrastructure/mongo/migrations",
		}
		mgoHandler, err := mongoDB.ConnectAndMigrate(logger)
		if err != nil {
			return err
		}
		dbHandler = mgoHandler
		scanInfosRepo = repository.NewMongoScanInfosRepository(logger, mgoHandler, mongoDB.Database)
		h, err = health.New(health.WithChecks(
			health.Config{
				Name:    "mongo",
				Timeout: time.Second * 5,
				Check:   mongoDB.Check(),
			},
		))
		if err != nil {
			return fmt.Errorf("unable to initialize mongo database health checks: %v", err)
		}

	case "mockdb": // set <database> field into the configuration to this for quick end-to-end test.
		mockDB := mockdb.Config{}
		mockdbHandler, err := mockDB.ConnectAndMigrate(zap.L())
		if err != nil {
			return err
		}
		dbHandler = mockdbHandler
		scanInfosRepo = repository.NewMockDBScanInfosRepository(zap.L(), mockdbHandler, "mockdb")
	}

	scanInfosUc := application.NewScanInfosUsecase(logger, configData, scanInfosRepo)

	// create the web service and setup the api endpoints.
	apisweb.New(logger, scanInfosUc).Router(router)

	// Useful routes to quickly check platform state.
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/status", func(c *gin.Context) {
		healthCheck := h.Measure(c)
		status := http.StatusOK
		if healthCheck.Status == health.StatusOK {
			status = http.StatusInternalServerError
		}

		c.JSON(status, gin.H{
			"app":       "Demo-Rest-API",
			"gitCommit": gitCommit,
			"gitTag":    gitTag,
			"database":  configData.Database,
			"health":    healthCheck,
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", configData.Server.Host, configData.Server.Port),
		Handler: router,
	}

	var g errgroup.Group
	g.Go(func() error {
		logger.Info("https server infos", zap.String("host", configData.Server.Host), zap.String("port", configData.Server.Port))
		return srv.ListenAndServeTLS(configData.Server.CertsFile, configData.Server.KeyFile)
	})
	go shutdownServer(srv, logger, dbHandler)

	logger.Info("starting HTTPS server")
	err = g.Wait()
	return fmt.Errorf("serving server failed: %w", err)
}

// shutdownServer handles a stop signal and teardown the connection to in-use database.
func shutdownServer(srv *http.Server, logger *zap.Logger, dbHandler config.DBHandler) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down all services...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("platform shutdown: error shutting down server", zap.Error(err))
	}

	if err := dbHandler.Shutdown(ctx); err != nil {
		logger.Error("platform shutdown: error disconnecting from database", zap.Error(err))
	}
}
